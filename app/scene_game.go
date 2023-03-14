package app

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	ui "github.com/t0l1k/eui"
	"github.com/t0l1k/nBack/data"
	"github.com/t0l1k/nBack/game"
)

type SceneGame struct {
	name                                                     string
	lblName, lblDt, lblTimer, lblMotiv, lblResult, lblHelper *ui.Label
	movesLine                                                *MovesLine
	btnStart, btnQuit                                        *ui.Button
	rect                                                     *ui.Rect
	stopper, pauseTimer, warningTimer                        int
	board                                                    *game.Board
	count, level, lives                                      int
	delayBeginCellShow, delayBeginCellHide                   int
	timeToNextCell                                           int
	paused                                                   bool
	ui.ContainerDefault
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
	var count int
	if s.count > 0 {
		s.level, s.lives, _ = data.GetDb().TodayData[s.count].NextLevel()
		count = s.count + 1
	} else {
		s.count = 1
		s.level = ui.GetPreferences().Get("default level").(int)
		s.lives = ui.GetPreferences().Get("threshold fallback sessions").(int)
		count = s.count
	}
	s.parseResult(count)
}

func (s *SceneGame) parseResult(count int) {
	var ss string
	ss += fmt.Sprintf("#%v ", count)
	ss += fmt.Sprintf("%v %v %v.", ui.GetLocale().Get("btnStart"), s.level, ui.GetLocale().Get("wordstepback"))
	res := ""
	tp := ui.GetPreferences().Get("game type").(string)
	switch tp {
	case game.Pos:
		res = ui.GetLocale().Get("optpos")
	case game.Col:
		res = ui.GetLocale().Get("optcol")
	case game.Sym:
		res = ui.GetLocale().Get("optsym")
	case game.Ari:
		res = ui.GetLocale().Get("optari")
	default:
		res = tp
	}
	if ui.GetPreferences().Get("manual mode").(bool) {
		ss += fmt.Sprintf(" %v(%v) %v.", ui.GetLocale().Get("wordGame"), res, ui.GetLocale().Get("wordhand"))
	} else {
		ss += fmt.Sprintf(" %v(%v) %v.", ui.GetLocale().Get("wordGame"), res, ui.GetLocale().Get("wordcclassic"))
	}
	s.lblResult.SetText(ss)
}
func (s *SceneGame) initGameTimers() {
	s.timeToNextCell = int(ui.GetPreferences().Get("time to next cell").(float64) * 1000)
	tm := ui.GetPreferences().Get("time to show cell").(float64)
	if tm > 0.9 {
		tm = 0.85
	}
	timeShowCell := int(float64(s.timeToNextCell) * tm)
	s.stopper = 0
	delay := (s.timeToNextCell - timeShowCell) / 2
	s.delayBeginCellShow = delay
	s.delayBeginCellHide = delay + timeShowCell
	log.Println("init board tm:", tm, s.timeToNextCell, timeShowCell, delay, s.delayBeginCellShow, s.delayBeginCellHide)
}

func (s *SceneGame) initUi() {
	rect := []int{0, 0, 1, 1}
	s.btnStart = ui.NewButton(
		ui.GetLocale().Get("wordnewsess"),
		rect,
		ui.GetTheme().Get("correct color"),
		ui.GetTheme().Get("fg"),
		func(b *ui.Button) {
			log.Println("Button new session pressed")
			s.paused = false
			s.newSession()
		})
	s.Add(s.btnStart)
	s.btnQuit = ui.NewButton("<", rect, ui.GetTheme().Get("correct color"), ui.GetTheme().Get("fg"), func(b *ui.Button) { ui.Pop() })
	s.Add(s.btnQuit)
	s.name = ui.GetLocale().Get("AppName") + " " + ui.GetLocale().Get("btnStart")
	s.lblName = ui.NewLabel(s.name, rect, ui.GetTheme().Get("correct color"), ui.GetTheme().Get("fg"))
	s.Add(s.lblName)
	s.board = game.NewBoard(rect, ui.GetPreferences(), ui.GetTheme())
	s.Add(s.board)
	s.lblResult = ui.NewLabel(" ", rect, ui.GetTheme().Get("correct color"), ui.GetTheme().Get("fg"))
	s.Add(s.lblResult)
	s.lblMotiv = ui.NewLabel(" ", rect, ui.GetTheme().Get("correct color"), ui.GetTheme().Get("fg"))
	s.Add(s.lblMotiv)
	s.lblMotiv.Visible = false
	s.lblTimer = ui.NewLabel(s.name, rect, ui.GetTheme().Get("correct color"), ui.GetTheme().Get("fg"))
	s.Add(s.lblTimer)
	s.lblTimer.Visible = false
	s.lblDt = ui.NewLabel("up: ", rect, ui.GetTheme().Get("correct color"), ui.GetTheme().Get("fg"))
	s.Add(s.lblDt)
	s.movesLine = NewMovesLine(rect)
	s.Add(s.movesLine)
	s.lblHelper = ui.NewLabel(ui.GetLocale().Get("btnHelperInGame"), rect, ui.GetTheme().Get("correct color"), ui.GetTheme().Get("fg"))
	s.Add(s.lblHelper)
}

func (s *SceneGame) Update(dt int) {
	for _, value := range s.Container {
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
			ui.GetUi().ShowNotification(ss)
			log.Println(ss)
		}
	} else if inpututil.IsKeyJustReleased(ebiten.KeyF6) {
		curPause := ui.GetPreferences().Get("time to next cell").(float64)
		if curPause >= 1.5 {
			curPause -= 0.5
			ui.GetPreferences().Set("time to next cell", curPause)
			s.initGameTimers()
			ss := fmt.Sprintf("%v %v %v %v %v", ui.GetLocale().Get("dec"), ui.GetLocale().Get("opttmnc"), ui.GetLocale().Get("by"), curPause, ui.GetLocale().Get("sec"))
			ui.GetUi().ShowNotification(ss)
			log.Println(ss)
		}
	} else if inpututil.IsKeyJustReleased(ebiten.KeyR) {
		if !s.board.InGame {
			s.level = ui.GetPreferences().Get("default level").(int)
			s.parseResult(s.count)
			ss := fmt.Sprintf("Обнулен уровень нназад на уровень по умолчанию на %v", s.level)
			ui.GetUi().ShowNotification(ss)
		}
	} else if inpututil.IsKeyJustReleased(ebiten.KeyF1) {
		if !s.board.InGame {
			values, _ := data.GetDb().ReadAllGamesScore(0, "", "")
			max := values.Max
			if max > s.level {
				s.level += 1
			}
			s.parseResult(s.count)
			ss := fmt.Sprintf("Увеличен уровень нназад на %v", s.level)
			ui.GetUi().ShowNotification(ss)
		}
	} else if inpututil.IsKeyJustReleased(ebiten.KeyF2) {
		if !s.board.InGame {
			if s.level > 1 {
				s.level -= 1
				s.parseResult(s.count)
				ss := fmt.Sprintf("Уменьшен уровень нназад на %v", s.level)
				ui.GetUi().ShowNotification(ss)
			}
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
			ss := "#" + strconv.Itoa(count) + data.GetDb().TodayData[count].String()
			s.lblResult.SetText(ss)
			log.Printf("Game result: %v", ss)
			s.count += 1
			s.lblMotiv.SetText(motiv)
			s.lblMotiv.SetBg(data.GetDb().TodayData[count].BgColor())
			s.lblName.Visible = true
			s.lblName.SetRect(true)
			s.lblName.SetText(s.name)
			s.lblName.SetBg(ui.GetTheme().Get("correct color"))
			x, y, w, h := int(float64(s.rect.H)*0.05), 0, int(float64(s.rect.W)*0.20), int(float64(s.rect.H)*0.05)
			s.lblName.Resize([]int{x, y, w, h})
			s.lblResult.Visible = true
			s.lblMotiv.Visible = true
			s.lblTimer.Visible = true
			s.btnQuit.Visible = true
			s.lblTimer.SetBg(ui.GetTheme().Get("error color"))
			s.pauseTimer = ui.GetPreferences().Get("pause to rest").(int) * 1000
			s.paused = true
			s.lblDt.Visible = true
			s.lblHelper.Visible = true
		}
		if s.pauseTimer > 0 {
			if s.paused {
				s.pauseTimer -= dt
				s.lblTimer.SetText(fmt.Sprintf("%v", s.pauseTimer/1000))
				if s.warningTimer > 0 {
					s.warningTimer -= dt
				}
			} else {
				s.pauseTimer += dt
				s.lblTimer.SetText(fmt.Sprintf("%02v:%02v", s.pauseTimer/1000/60, s.pauseTimer/1000%60))
				s.btnStart.Visible = true
				if s.warningTimer > 0 {
					s.warningTimer -= dt
				} else {
					s.lblTimer.SetBg(ui.GetTheme().Get("correct color"))
				}
			}
		} else if s.pauseTimer <= 0 {
			if s.paused {
				s.paused = false
				s.pauseTimer += ui.GetPreferences().Get("pause to rest").(int) * 1000
				s.lblTimer.SetBg(ui.GetTheme().Get("warning color"))
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
	s.lblHelper.Visible = false
	if ui.GetPreferences().Get("feedback on user move").(bool) {
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
	if ui.GetPreferences().Get("feedback on user move").(bool) {
		switch s.board.MoveStatus {
		case game.Correct:
			s.lblName.SetBg(ui.GetTheme().Get("correct color"))
		case game.Error:
			s.lblName.SetBg(ui.GetTheme().Get("error color"))
		case game.Warning:
			s.lblName.SetBg(ui.GetTheme().Get("warning color"))
		case game.Regular:
			s.lblName.SetBg(ui.GetTheme().Get("regular color"))
		default:
			s.lblName.SetBg(ui.GetTheme().Get("game bg"))
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
	case game.Ari:
		str1 = "Ari"
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
		Manual:       ui.GetPreferences().Get("manual mode").(bool),
		Advance:      ui.GetPreferences().Get("threshold advance").(int),
		Fallback:     ui.GetPreferences().Get("threshold fallback").(int),
		Resetonerror: ui.GetPreferences().Get("reset on first wrong").(bool),
		MovesStatus:  s.board.MovesStatus,
	}
	data.GetDb().InsertGame(values)
	log.Println("Game Saved in DB.")
	dur := s.board.DtEnd.Sub(s.board.DtBeg) / 2
	s.warningTimer = int(dur.Seconds()) * 1000
}

func (s *SceneGame) Draw(surface *ebiten.Image) {
	surface.Fill(ui.GetTheme().Get("game bg"))
	for _, value := range s.Container {
		value.Draw(surface)
	}
}

func (s *SceneGame) Resize() {
	w0, h0 := ui.GetUi().GetScreenSize()
	s.rect = ui.NewRect([]int{0, 0, w0, h0})
	x, y, w, hTop := 0, 0, int(float64(s.rect.H)*0.05), int(float64(s.rect.H)*0.05)
	s.btnQuit.Resize([]int{x, y, w, hTop})
	x, w = hTop, int(float64(s.rect.W)*0.20)
	s.lblName.Resize([]int{x, y, w, hTop})
	x = s.rect.Right() - w
	s.lblDt.Resize([]int{x, y, w, hTop})

	sz := s.rect.GetLowestSize()
	cellSize := float64(sz)/3 - float64(sz)*0.02
	marginX := float64(s.rect.W)/2 - cellSize*3/2
	marginY := float64(s.rect.H)/2 - cellSize*3/2
	x, y = int(marginX), int(marginY)+hTop/2
	s.board.Resize([]int{x, y, int(cellSize) * 3, int(cellSize) * 3})

	wBtn, hBtn := int(float64(s.rect.W)*0.5), int(float64(s.rect.H)*0.15)
	w2 := int(float64(s.rect.W) * 0.87)

	x = (s.rect.W - wBtn) / 2
	y = s.rect.H - hTop*12 - hBtn
	s.lblTimer.Resize([]int{x, y, wBtn, hBtn})

	x = (s.rect.W - w2) / 2
	y = s.rect.H - hTop*8 - hBtn
	s.lblMotiv.Resize([]int{x, y, w2, hTop * 2})

	x = (s.rect.W - w2) / 2
	y = s.rect.H - hTop*6 - hBtn
	s.movesLine.Resize([]int{x, y, w2, hTop})

	x = (s.rect.W - w2) / 2
	y = s.rect.H - hTop*5 - hBtn
	s.lblResult.Resize([]int{x, y, w2, hTop * 2})

	x = (s.rect.W - wBtn) / 2
	y = s.rect.H - hTop*2 - hBtn
	s.btnStart.Resize([]int{x, y, wBtn, hBtn})
	x, y = 0, s.rect.H-hTop
	s.lblHelper.Resize([]int{x, y, w0, hTop})
}

func (s *SceneGame) Close() {
	for _, v := range s.Container {
		v.Close()
	}
}
