package game

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/t0l1k/eui"
	"github.com/t0l1k/nBack/data"
	"github.com/t0l1k/nBack/game"
	"github.com/t0l1k/nBack/ui/scene/plot"
)

type SceneGame struct {
	name                                                     string
	lblName, lblDt, lblTimer, lblMotiv, lblResult, lblHelper *eui.Label
	movesLine                                                *plot.MovesLine
	btnStart, btnQuit                                        *eui.Button
	rect                                                     *eui.Rect
	stopper, pauseTimer, warningTimer                        int
	board                                                    *game.Board
	count, level, lives                                      int
	delayBeginCellShow, delayBeginCellHide                   int
	timeToNextCell                                           int
	paused                                                   bool
	eui.ContainerDefault
}

func NewSceneGame() *SceneGame {
	s := &SceneGame{
		rect: eui.NewRect([]int{0, 0, 1, 1}),
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
		s.level = eui.GetPreferences().Get("default level").(int)
		s.lives = eui.GetPreferences().Get("threshold fallback sessions").(int)
		count = s.count
	}
	s.parseResult(count)
}

func (s *SceneGame) parseResult(count int) {
	var ss string
	ss += fmt.Sprintf("#%v ", count)
	ss += fmt.Sprintf("%v %v %v.", eui.GetLocale().Get("btnStart"), s.level, eui.GetLocale().Get("wordstepback"))
	res := ""
	tp := eui.GetPreferences().Get("game type").(string)
	switch tp {
	case game.Pos:
		res = eui.GetLocale().Get("optpos")
	case game.Col:
		res = eui.GetLocale().Get("optcol")
	case game.Sym:
		res = eui.GetLocale().Get("optsym")
	case game.Ari:
		res = eui.GetLocale().Get("optari")
	default:
		res = tp
	}
	if eui.GetPreferences().Get("manual mode").(bool) {
		ss += fmt.Sprintf(" %v(%v) %v.", eui.GetLocale().Get("wordGame"), res, eui.GetLocale().Get("wordhand"))
	} else {
		ss += fmt.Sprintf(" %v(%v) %v.", eui.GetLocale().Get("wordGame"), res, eui.GetLocale().Get("wordcclassic"))
	}
	s.lblResult.SetText(ss)
}
func (s *SceneGame) initGameTimers() {
	s.timeToNextCell = int(eui.GetPreferences().Get("time to next cell").(float64) * 1000)
	tm := eui.GetPreferences().Get("time to show cell").(float64)
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
	s.btnStart = eui.NewButton(
		eui.GetLocale().Get("wordnewsess"),
		rect,
		eui.GetTheme().Get("correct color"),
		eui.GetTheme().Get("fg"),
		func(b *eui.Button) {
			log.Println("Button new session pressed")
			s.paused = false
			s.newSession()
		})
	s.Add(s.btnStart)
	s.btnQuit = eui.NewButton("<", rect, eui.GetTheme().Get("correct color"), eui.GetTheme().Get("fg"), func(b *eui.Button) { eui.Pop() })
	s.Add(s.btnQuit)
	s.name = eui.GetLocale().Get("AppName") + " " + eui.GetLocale().Get("btnStart")
	s.lblName = eui.NewLabel(s.name, rect, eui.GetTheme().Get("correct color"), eui.GetTheme().Get("fg"))
	s.Add(s.lblName)
	s.board = game.NewBoard(rect, eui.GetPreferences(), eui.GetTheme())
	s.Add(s.board)
	s.lblResult = eui.NewLabel(" ", rect, eui.GetTheme().Get("correct color"), eui.GetTheme().Get("fg"))
	s.Add(s.lblResult)
	s.lblMotiv = eui.NewLabel(" ", rect, eui.GetTheme().Get("correct color"), eui.GetTheme().Get("fg"))
	s.Add(s.lblMotiv)
	s.lblMotiv.Visible = false
	s.lblTimer = eui.NewLabel(s.name, rect, eui.GetTheme().Get("correct color"), eui.GetTheme().Get("fg"))
	s.Add(s.lblTimer)
	s.lblTimer.Visible = false
	s.lblDt = eui.NewLabel("up: ", rect, eui.GetTheme().Get("correct color"), eui.GetTheme().Get("fg"))
	s.Add(s.lblDt)
	s.movesLine = plot.NewMovesLine(rect)
	s.Add(s.movesLine)
	s.lblHelper = eui.NewLabel(eui.GetLocale().Get("btnHelperInGame"), rect, eui.GetTheme().Get("correct color"), eui.GetTheme().Get("fg"))
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
		curPause := eui.GetPreferences().Get("time to next cell").(float64)
		if curPause < 5 {
			curPause += 0.5
			eui.GetPreferences().Set("time to next cell", curPause)
			s.initGameTimers()
			ss := fmt.Sprintf("%v %v %v %v %v", eui.GetLocale().Get("inc"), eui.GetLocale().Get("opttmnc"), eui.GetLocale().Get("by"), curPause, eui.GetLocale().Get("sec"))
			eui.GetUi().ShowNotification(ss)
			log.Println(ss)
		}
	} else if inpututil.IsKeyJustReleased(ebiten.KeyF6) {
		curPause := eui.GetPreferences().Get("time to next cell").(float64)
		if curPause >= 1.5 {
			curPause -= 0.5
			eui.GetPreferences().Set("time to next cell", curPause)
			s.initGameTimers()
			ss := fmt.Sprintf("%v %v %v %v %v", eui.GetLocale().Get("dec"), eui.GetLocale().Get("opttmnc"), eui.GetLocale().Get("by"), curPause, eui.GetLocale().Get("sec"))
			eui.GetUi().ShowNotification(ss)
			log.Println(ss)
		}
	} else if inpututil.IsKeyJustReleased(ebiten.KeyR) {
		if !s.board.InGame {
			s.level = eui.GetPreferences().Get("default level").(int)
			s.parseResult(s.count)
			ss := fmt.Sprintf("Обнулен уровень нназад на уровень по умолчанию на %v", s.level)
			eui.GetUi().ShowNotification(ss)
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
			eui.GetUi().ShowNotification(ss)
		}
	} else if inpututil.IsKeyJustReleased(ebiten.KeyF2) {
		if !s.board.InGame {
			if s.level > 1 {
				s.level -= 1
				s.parseResult(s.count)
				ss := fmt.Sprintf("Уменьшен уровень нназад на %v", s.level)
				eui.GetUi().ShowNotification(ss)
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
		s.lblDt.SetText(eui.GetUi().UpdateUpTime())
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
			s.lblName.SetBg(eui.GetTheme().Get("correct color"))
			x, y, w, h := int(float64(s.rect.H)*0.05), 0, int(float64(s.rect.W)*0.20), int(float64(s.rect.H)*0.05)
			s.lblName.Resize([]int{x, y, w, h})
			s.lblResult.Visible = true
			s.lblMotiv.Visible = true
			s.lblTimer.Visible = true
			s.btnQuit.Visible = true
			s.lblTimer.SetBg(eui.GetTheme().Get("error color"))
			s.pauseTimer = eui.GetPreferences().Get("pause to rest").(int) * 1000
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
					s.lblTimer.SetBg(eui.GetTheme().Get("correct color"))
				}
			}
		} else if s.pauseTimer <= 0 {
			if s.paused {
				s.paused = false
				s.pauseTimer += eui.GetPreferences().Get("pause to rest").(int) * 1000
				s.lblTimer.SetBg(eui.GetTheme().Get("warning color"))
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
	if eui.GetPreferences().Get("feedback on user move").(bool) {
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
	if eui.GetPreferences().Get("feedback on user move").(bool) {
		switch s.board.MoveStatus {
		case game.Correct:
			s.lblName.SetBg(eui.GetTheme().Get("correct color"))
		case game.Error:
			s.lblName.SetBg(eui.GetTheme().Get("error color"))
		case game.Warning:
			s.lblName.SetBg(eui.GetTheme().Get("warning color"))
		case game.Regular:
			s.lblName.SetBg(eui.GetTheme().Get("regular color"))
		default:
			s.lblName.SetBg(eui.GetTheme().Get("game bg"))
		}
	}
	str1 := ""
	switch eui.GetPreferences().Get("game type").(string) {
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
		GameType:     eui.GetPreferences().Get("game type").(string),
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
		Manual:       eui.GetPreferences().Get("manual mode").(bool),
		Advance:      eui.GetPreferences().Get("threshold advance").(int),
		Fallback:     eui.GetPreferences().Get("threshold fallback").(int),
		Resetonerror: eui.GetPreferences().Get("reset on first wrong").(bool),
		MovesStatus:  s.board.MovesStatus,
	}
	data.GetDb().InsertGame(values)
	log.Println("Game Saved in DB.")
	dur := s.board.DtEnd.Sub(s.board.DtBeg) / 2
	s.warningTimer = int(dur.Seconds()) * 1000
}

func (s *SceneGame) Draw(surface *ebiten.Image) {
	surface.Fill(eui.GetTheme().Get("game bg"))
	for _, value := range s.Container {
		value.Draw(surface)
	}
}

func (s *SceneGame) Resize() {
	w0, h0 := eui.GetUi().GetScreenSize()
	s.rect = eui.NewRect([]int{0, 0, w0, h0})
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
