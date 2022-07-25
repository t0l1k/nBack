package game

import (
	"fmt"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/t0l1k/nBack/ui"
)

type SceneGame struct {
	name                                          string
	lblName, lblResult, lblMotiv, lblTimer, lblDt *ui.Label
	btnStart, btnQuit                             *ui.Button
	rect                                          *ui.Rect
	container                                     []ui.Drawable
	stopper, pauseTimer                           int
	board                                         *Board
	count, level, lives                           int
	delayBeginCellShow, delayBeginCellHide        int
	timeToNextCell, timeShowCell                  int
	paused                                        bool
}

func NewSceneGame() *SceneGame {
	s := &SceneGame{
		rect: ui.NewRect([]int{0, 0, 1, 1}),
	}
	s.initUi()
	return s
}

func (s *SceneGame) Entered() {
	s.Resize()
	s.initGame()
	s.initGameTimers()
	log.Printf("Enterd Scene Game")
}

func (s *SceneGame) initGame() {
	s.count = getDb().todayGamesCount
	if s.count > 0 {
		s.level, s.lives, _ = getDb().todayData[s.count].NextLevel()
	} else {
		s.count = 1
		s.level = (*ui.GetPreferences())["default level"].(int)
		s.lives = (*ui.GetPreferences())["threshold fallback sessions"].(int)
	}
	ss := fmt.Sprintf("#%v level:%v", s.count, s.level)
	if (*ui.GetPreferences())["manual mode"].(bool) {
		ss += " Manual game mode."
	} else {
		ss += " Classic game mode."
	}
	s.lblResult.SetText(ss)
}
func (s *SceneGame) initGameTimers() {
	s.timeToNextCell = int((*ui.GetPreferences())["time to next cell"].(float64) * 1000)
	s.timeShowCell = int((*ui.GetPreferences())["time to show cell"].(float64) * 1000)
	s.stopper = 0
	delay := (s.timeToNextCell - s.timeShowCell) / 2
	s.delayBeginCellShow = delay
	s.delayBeginCellHide = delay + s.timeShowCell
}

func (s *SceneGame) initUi() {
	rect := []int{0, 0, 1, 1}
	s.btnStart = ui.NewButton("New Session", rect, (*ui.GetTheme())["correct color"], (*ui.GetTheme())["fg"], func(b *ui.Button) {
		log.Println("Button new session pressed")
		s.paused = false
		s.newSession()
	})
	s.Add(s.btnStart)
	s.btnQuit = ui.NewButton("<", rect, (*ui.GetTheme())["correct color"], (*ui.GetTheme())["fg"], func(b *ui.Button) { ui.GetApp().Pop() })
	s.Add(s.btnQuit)
	s.name = "Game N-Back result"
	s.lblName = ui.NewLabel(s.name, rect, (*ui.GetTheme())["correct color"], (*ui.GetTheme())["fg"])
	s.Add(s.lblName)
	s.board = NewBoard(rect, ui.GetPreferences(), ui.GetTheme())
	s.Add(s.board)
	s.lblResult = ui.NewLabel(" ", rect, (*ui.GetTheme())["correct color"], (*ui.GetTheme())["fg"])
	s.Add(s.lblResult)
	s.lblMotiv = ui.NewLabel("Motivation", rect, (*ui.GetTheme())["correct color"], (*ui.GetTheme())["fg"])
	s.Add(s.lblMotiv)
	s.lblMotiv.Visible = false
	s.lblTimer = ui.NewLabel(s.name, rect, (*ui.GetTheme())["correct color"], (*ui.GetTheme())["fg"])
	s.Add(s.lblTimer)
	s.lblTimer.Visible = false
	s.lblDt = ui.NewLabel("up: 00:00 ", rect, (*ui.GetTheme())["correct color"], (*ui.GetTheme())["fg"])
	s.Add(s.lblDt)
}

func (s *SceneGame) Update(dt int) {
	for _, value := range s.container {
		value.Update(dt)
	}
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		if s.board.inGame && !s.board.userMoved {
			s.board.CheckUserMove()
		}
	}
	if inpututil.IsKeyJustReleased(ebiten.KeySpace) {
		if s.board.inGame {
			s.board.CheckUserMove()
		} else if !s.paused {
			s.newSession()
		}
	}
	if s.board.inGame {
		s.stopper += dt
		if s.stopper >= s.timeToNextCell {
			s.stopper -= s.timeToNextCell
			s.board.MakeMove()
		} else if !s.board.IsShowActiveCell() && s.delayBeginCellShow < s.stopper && s.stopper < s.delayBeginCellHide {
			s.board.CheckMoveRegular()
			s.board.ShowActiveCell()
		} else if s.board.IsShowActiveCell() && s.stopper > s.delayBeginCellHide {
			s.board.HideActiveCell()
		}
		s.moveStatus()
	} else {
		s.lblDt.SetText(ui.GetApp().UpdateUpTime())
		if !s.lblResult.Visible {
			s.SaveGame()
			var motiv string
			count := getDb().todayGamesCount
			s.level, s.lives, motiv = getDb().todayData[count].NextLevel()
			ss := getDb().todayData[count].String()
			s.lblResult.SetText(ss)
			log.Printf("Game Result is: %v", ss)
			s.count += 1
			s.lblMotiv.SetText(motiv)
			s.lblMotiv.SetBg(getDb().todayData[count].BgColor())
			s.lblName.SetRect(true)
			s.lblName.SetText(s.name)
			s.lblName.SetBg((*ui.GetTheme())["correct color"])
			x, y, w, h := int(float64(s.rect.H)*0.05), 0, int(float64(s.rect.W)*0.20), int(float64(s.rect.H)*0.05)
			s.lblName.Resize([]int{x, y, w, h})
			s.lblResult.Visible = true
			s.lblMotiv.Visible = true
			s.lblTimer.Visible = true
			s.btnQuit.Visible = true
			s.lblTimer.SetBg((*ui.GetTheme())["error color"])
			s.pauseTimer = (*ui.GetPreferences())["pause to rest"].(int) * 1000
			s.paused = true
			s.lblDt.Visible = true
		}
		if s.pauseTimer > 0 {
			if s.paused {
				s.pauseTimer -= dt
				s.lblTimer.SetText(fmt.Sprintf("%v", s.pauseTimer/1000))
			} else {
				s.pauseTimer += dt
				s.lblTimer.SetText(fmt.Sprintf("%02v:%02v", s.pauseTimer/1000/60, s.pauseTimer/1000%60))
				s.btnStart.Visible = true
			}
		} else if s.pauseTimer <= 0 {
			if s.paused {
				s.paused = false
				s.pauseTimer += (*ui.GetPreferences())["pause to rest"].(int) * 1000
				s.lblTimer.SetBg((*ui.GetTheme())["correct color"])
			}
		}
	}
}

func (s *SceneGame) newSession() {
	s.board.Reset(s.count, s.level)
	s.btnStart.Visible = false
	s.lblResult.Visible = false
	s.lblMotiv.Visible = false
	s.lblTimer.Visible = false
	s.btnQuit.Visible = false
	s.lblDt.Visible = false
	if (*ui.GetPreferences())["feedback on user move"].(bool) {
		x, y, w, h := 0, 0, int(float64(s.rect.W)*0.20), int(float64(s.rect.H)*0.05)
		s.lblName.Resize([]int{x, y, w, h})
		s.lblName.SetRect(true)
	} else {
		s.lblName.Visible = false
		sz := s.rect.GetLowestSize()
		cellSize := float64(sz)/3 - float64(sz)*0.02
		marginX := float64(s.rect.W)/2 - cellSize*3/2
		marginY := float64(s.rect.H)/2 - cellSize*3/2
		x, y := int(marginX), int(marginY)
		s.board.Resize([]int{x, y, int(cellSize) * 3, int(cellSize) * 3})
	}
}

func (s *SceneGame) moveStatus() {
	if (*ui.GetPreferences())["feedback on user move"].(bool) {
		switch s.board.moveStatus {
		case Correct:
			s.lblName.SetBg((*ui.GetTheme())["correct color"])
		case Error:
			s.lblName.SetBg((*ui.GetTheme())["error color"])
		case Warning:
			s.lblName.SetBg((*ui.GetTheme())["warning color"])
		case Regular:
			s.lblName.SetBg((*ui.GetTheme())["regular color"])
		default:
			s.lblName.SetBg((*ui.GetTheme())["game bg"])
		}
	}
	str := fmt.Sprintf("Pos %v (%v) (%v/%v)", s.level, s.lives, s.board.move, s.board.totalMoves)
	s.lblName.SetText(str)
}

func (s *SceneGame) SaveGame() {
	dtBeg := s.board.dtBeg.Format("2006-01-02 15:04:05.000")
	dtEnd := s.board.dtEnd.Format("2006-01-02 15:04:05.000")
	values := &GameData{
		dtBeg:        dtBeg,
		dtEnd:        dtEnd,
		level:        s.level,
		lives:        s.lives,
		percent:      s.board.getPercent(),
		correct:      s.board.countCorrect,
		wrong:        s.board.countWrong,
		moves:        s.board.move,
		totalmoves:   s.board.totalMoves,
		manual:       (*ui.GetPreferences())["manual mode"].(bool),
		advance:      (*ui.GetPreferences())["threshold advance"].(int),
		fallback:     (*ui.GetPreferences())["threshold fallback"].(int),
		resetonerror: (*ui.GetPreferences())["reset on first wrong"].(bool),
	}
	getDb().InsertGame(values)
	log.Println("Game Saved in DB.")
}

func (s *SceneGame) Draw(surface *ebiten.Image) {
	surface.Fill((*ui.GetTheme())["game bg"])
	for _, value := range s.container {
		value.Draw(surface)
	}
}

func (s *SceneGame) Add(item ui.Drawable) {
	s.container = append(s.container, item)
}

func (s *SceneGame) Resize() {
	w, h := ui.GetApp().GetScreenSize()
	s.rect = ui.NewRect([]int{0, 0, w, h})
	x, y, w, h := 0, 0, int(float64(s.rect.H)*0.05), int(float64(s.rect.H)*0.05)
	s.btnQuit.Resize([]int{x, y, w, h})
	x, w = h, int(float64(s.rect.W)*0.20)
	s.lblName.Resize([]int{x, y, w, h})
	x = s.rect.Right() - w
	s.lblDt.Resize([]int{x, y, w, h})
	sz := s.rect.GetLowestSize()
	cellSize := float64(sz)/3 - float64(sz)*0.02
	marginX := float64(s.rect.W)/2 - cellSize*3/2
	marginY := float64(s.rect.H)/2 - cellSize*3/2
	x, y = int(marginX), int(marginY)+h/2
	s.board.Resize([]int{x, y, int(cellSize) * 3, int(cellSize) * 3})
	w, h = int(float64(s.rect.W)*0.5), int(float64(s.rect.H)*0.15)
	x, y = (s.rect.W-w)/2, s.rect.H-int(float64(h)*1)
	s.btnStart.Resize([]int{x, y, w, h})
	w, h = int(float64(s.rect.W)*0.9), int(float64(s.rect.H)*0.08)
	x, y = (s.rect.W-w)/2, s.rect.H-int(float64(h)*3.5)
	s.lblResult.Resize([]int{x, y, w, h})
	w, h = int(float64(s.rect.W)*0.7), int(float64(s.rect.H)*0.08)
	x, y = (s.rect.W-w)/2, s.rect.H-int(float64(h)*5.5)
	s.lblMotiv.Resize([]int{x, y, w, h})
	w, h = int(float64(s.rect.W)*0.5), int(float64(s.rect.H)*0.2)
	x, y = (s.rect.W-w)/2, s.rect.H-int(float64(h)*4)
	s.lblTimer.Resize([]int{x, y, w, h})
}

func (s *SceneGame) Quit() {}
