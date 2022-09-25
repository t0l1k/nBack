package app

import (
	"fmt"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/t0l1k/nBack/data"
	"github.com/t0l1k/nBack/game"
	"github.com/t0l1k/nBack/ui"
)

type SceneGame struct {
	name                                          string
	lblName, lblResult, lblMotiv, lblTimer, lblDt *ui.Label
	movesLine                                     *MovesLine
	btnStart, btnQuit                             *ui.Button
	rect                                          *ui.Rect
	container                                     []ui.Drawable
	stopper, pauseTimer                           int
	board                                         *game.Board
	count, level, lives                           int
	delayBeginCellShow, delayBeginCellHide        int
	timeToNextCell, timeShowCell                  int
	paused                                        bool
	notify                                        *ui.Notification
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
	s.count = data.GetDb().TodayGamesCount
	if s.count > 0 {
		s.level, s.lives, _ = data.GetDb().TodayData[s.count].NextLevel()
	} else {
		s.count = 1
		s.level = (*ui.GetPreferences())["default level"].(int)
		s.lives = (*ui.GetPreferences())["threshold fallback sessions"].(int)
	}
	ss := fmt.Sprintf("#%v %v %v %v.", s.count, ui.GetLocale().Get("btnStart"), s.level, ui.GetLocale().Get("wordstepback"))
	res := ""
	tp := ui.GetPreferences().Get("game type").(string)
	switch tp {
	case game.Pos:
		res = ui.GetLocale().Get("optpos")
	case game.Col:
		res = ui.GetLocale().Get("optcol")
	case game.Sym:
		res = ui.GetLocale().Get("optsym")
	default:
		res = tp
	}
	if (*ui.GetPreferences())["manual mode"].(bool) {
		ss += fmt.Sprintf(" %v(%v) %v.", ui.GetLocale().Get("wordGame"), res, ui.GetLocale().Get("wordhand"))
	} else {
		ss += fmt.Sprintf(" %v(%v) %v.", ui.GetLocale().Get("wordGame"), res, ui.GetLocale().Get("wordcclassic"))
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
	s.btnStart = ui.NewButton(ui.GetLocale().Get("wordnewsess"), rect, (*ui.GetTheme())["correct color"], (*ui.GetTheme())["fg"], func(b *ui.Button) {
		log.Println("Button new session pressed")
		s.paused = false
		s.newSession()
	})
	s.Add(s.btnStart)
	s.btnQuit = ui.NewButton("<", rect, (*ui.GetTheme())["correct color"], (*ui.GetTheme())["fg"], func(b *ui.Button) { ui.GetUi().Pop() })
	s.Add(s.btnQuit)
	s.name = ui.GetLocale().Get("AppName") + " " + ui.GetLocale().Get("btnStart")
	s.lblName = ui.NewLabel(s.name, rect, (*ui.GetTheme())["correct color"], (*ui.GetTheme())["fg"])
	s.Add(s.lblName)
	s.board = game.NewBoard(rect, ui.GetPreferences(), ui.GetTheme())
	s.Add(s.board)
	s.lblResult = ui.NewLabel(" ", rect, (*ui.GetTheme())["correct color"], (*ui.GetTheme())["fg"])
	s.Add(s.lblResult)
	s.lblMotiv = ui.NewLabel(" ", rect, (*ui.GetTheme())["correct color"], (*ui.GetTheme())["fg"])
	s.Add(s.lblMotiv)
	s.lblMotiv.Visible = false
	s.lblTimer = ui.NewLabel(s.name, rect, (*ui.GetTheme())["correct color"], (*ui.GetTheme())["fg"])
	s.Add(s.lblTimer)
	s.lblTimer.Visible = false
	s.lblDt = ui.NewLabel("up: ", rect, (*ui.GetTheme())["correct color"], (*ui.GetTheme())["fg"])
	s.Add(s.lblDt)
	s.movesLine = NewMovesLine(rect)
	s.Add(s.movesLine)
	s.notify = ui.NewNotification(ui.GetLocale().Get("notifhere"), rect, (*ui.GetTheme())["regular color"], (*ui.GetTheme())["fg"])
	s.Add(s.notify)
}

func (s *SceneGame) Update(dt int) {
	for _, value := range s.container {
		value.Update(dt)
	}
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		if s.board.InGame && !s.board.UserMoved {
			s.board.CheckUserMove()
		}
	}
	if inpututil.IsKeyJustReleased(ebiten.KeySpace) {
		if s.board.InGame {
			s.board.CheckUserMove()
		} else if !s.paused {
			s.newSession()
		}
	} else if inpututil.IsKeyJustReleased(ebiten.KeyF5) {
		curPause := ui.GetPreferences().Get("time to next cell").(float64)
		if curPause < 5 {
			curPause += 0.5
			ui.GetPreferences().Set("time to next cell", curPause)
			s.initGameTimers()
			ss := fmt.Sprintf("%v %v %v %v %v", ui.GetLocale().Get("inc"), ui.GetLocale().Get("opttmnc"), ui.GetLocale().Get("by"), curPause, ui.GetLocale().Get("sec"))
			s.notify.SetText(ss)
			log.Printf(ss, curPause)
		}
	} else if inpututil.IsKeyJustReleased(ebiten.KeyF6) {
		curPause := ui.GetPreferences().Get("time to next cell").(float64)
		if curPause >= 2 {
			curPause -= 0.5
			ui.GetPreferences().Set("time to next cell", curPause)
			s.initGameTimers()
			ss := fmt.Sprintf("%v %v %v %v %v", ui.GetLocale().Get("dec"), ui.GetLocale().Get("opttmnc"), ui.GetLocale().Get("by"), curPause, ui.GetLocale().Get("sec"))
			s.notify.SetText(ss)
			log.Printf(ss, curPause)
		}
	}
	if s.board.InGame {
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
		s.lblDt.SetText(ui.GetUi().UpdateUpTime())
		if !s.lblResult.Visible {
			s.SaveGame()
			s.movesLine.Visible = true
			s.movesLine.Dirty = true
			var motiv string
			count := data.GetDb().TodayGamesCount
			s.level, s.lives, motiv = data.GetDb().TodayData[count].NextLevel()
			ss := data.GetDb().TodayData[count].String()
			s.lblResult.SetText(ss)
			log.Printf("Game result: %v", ss)
			s.count += 1
			s.lblMotiv.SetText(motiv)
			s.lblMotiv.SetBg(data.GetDb().TodayData[count].BgColor())
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
	s.movesLine.Visible = false
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
		switch s.board.MoveStatus {
		case game.Correct:
			s.lblName.SetBg((*ui.GetTheme())["correct color"])
		case game.Error:
			s.lblName.SetBg((*ui.GetTheme())["error color"])
		case game.Warning:
			s.lblName.SetBg((*ui.GetTheme())["warning color"])
		case game.Regular:
			s.lblName.SetBg((*ui.GetTheme())["regular color"])
		default:
			s.lblName.SetBg((*ui.GetTheme())["game bg"])
		}
	}
	str1 := ""
	switch ui.GetPreferences().Get("game type").(string) {
	case game.Pos:
		str1 = "Pos"
	case game.Col:
		str1 = "Col"
	case game.Sym:
		str1 = "Sym"
	}
	str := fmt.Sprintf("%v %v (%v) (%v/%v)", str1, s.level, s.lives, s.board.Move, s.board.TotalMoves)
	s.lblName.SetText(str)
}

func (s *SceneGame) SaveGame() {
	dtBeg := s.board.DtBeg.Format("2006-01-02 15:04:05.000")
	dtEnd := s.board.DtEnd.Format("2006-01-02 15:04:05.000")
	values := &data.GameData{
		GameType:     ui.GetPreferences().Get("game type").(string),
		DtBeg:        dtBeg,
		DtEnd:        dtEnd,
		Level:        s.level,
		Lives:        s.lives,
		Percent:      s.board.GetPercent(),
		Correct:      s.board.CountCorrect,
		Wrong:        s.board.CountWrong,
		Missed:       s.board.CountMissed,
		Moves:        s.board.Move,
		Totalmoves:   s.board.TotalMoves,
		Manual:       (*ui.GetPreferences())["manual mode"].(bool),
		Advance:      (*ui.GetPreferences())["threshold advance"].(int),
		Fallback:     (*ui.GetPreferences())["threshold fallback"].(int),
		Resetonerror: (*ui.GetPreferences())["reset on first wrong"].(bool),
		MovesStatus:  s.board.MovesStatus,
	}
	data.GetDb().InsertGame(values)
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
	w, h := ui.GetUi().GetScreenSize()
	s.rect = ui.NewRect([]int{0, 0, w, h})
	x, y, w, h := 0, 0, int(float64(s.rect.H)*0.05), int(float64(s.rect.H)*0.05)
	s.btnQuit.Resize([]int{x, y, w, h})
	x, w = h, int(float64(s.rect.W)*0.20)
	s.lblName.Resize([]int{x, y, w, h})
	x = s.rect.Right() - w
	s.lblDt.Resize([]int{x, y, w, h})

	w = int(float64(s.rect.W) * 0.50)
	x = s.rect.CenterX() - w/2
	s.notify.Resize([]int{x, y, w, h})

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
	w, h = int(float64(s.rect.W)*0.9), int(float64(s.rect.H)*0.05)
	x, y = (s.rect.W-w)/2, s.rect.H-int(float64(h)*6.7)
	s.movesLine.Resize([]int{x, y, w, h})
}

func (s *SceneGame) Quit() {
	for _, v := range s.container {
		v.Close()
	}
}
