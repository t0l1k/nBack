package main

import (
	"fmt"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/t0l1k/nBack/ui"
)

type SceneGame struct {
	name                                             string
	lblName, lblIntro, lblResult, lblMotiv, lblTimer *ui.Label
	rect                                             *ui.Rect
	container                                        []ui.Drawable
	stopper, pauseTimer                              int
	board                                            *Board
	count, level, lives                              int
	delayBeginCellShow, delayBeginCellHide           int
	timeToNextCell, timeShowCell                     int
	paused                                           bool
}

func NewSceneGame() *SceneGame {
	return &SceneGame{
		rect: getApp().rect}
}

func (s *SceneGame) Entered() {
	s.initUi()
	s.initGame()
	s.initGameTimers()
	log.Printf("Enterd Scene Game")
}

func (s *SceneGame) initGame() {
	s.count = getApp().db.todayGamesCount
	if s.count > 0 {
		s.level, s.lives, _ = getApp().db.todayData[s.count].NextLevel()
	} else {
		s.count = 1
		s.level = 1
		s.lives = 3
	}
	ss := fmt.Sprintf("#%v level:%v", s.count, s.level)
	s.lblResult.SetText(ss)
}
func (s *SceneGame) initGameTimers() {
	s.timeToNextCell = 2000
	s.timeShowCell = 500
	s.stopper = 0
	delay := (s.timeToNextCell - s.timeShowCell) / 2
	s.delayBeginCellShow = delay
	s.delayBeginCellHide = delay + s.timeShowCell
}

func (s *SceneGame) initUi() {
	rect := []int{0, 0, 1, 1}
	s.name = "Game N-Back result"
	s.lblName = ui.NewLabel(s.name, rect, getApp().theme.correct, getApp().theme.fg)
	s.Add(s.lblName)
	s.board = NewBoard(rect)
	s.Add(s.board)
	s.lblIntro = ui.NewLabel("Press the space bar to start the game", rect, getApp().theme.correct, getApp().theme.fg)
	s.Add(s.lblIntro)
	s.lblResult = ui.NewLabel(" ", rect, getApp().theme.correct, getApp().theme.fg)
	s.Add(s.lblResult)
	s.lblMotiv = ui.NewLabel("Motivation", rect, getApp().theme.correct, getApp().theme.fg)
	s.Add(s.lblMotiv)
	s.lblMotiv.Visibe = false
	s.lblTimer = ui.NewLabel(s.name, rect, getApp().theme.correct, getApp().theme.fg)
	s.Add(s.lblTimer)
	s.lblTimer.Visibe = false
	s.Resize()
}

func (s *SceneGame) Update(dt int) {
	for _, value := range s.container {
		value.Update(dt)
	}
	if inpututil.IsKeyJustReleased(ebiten.KeySpace) {
		if s.board.inGame {
			s.board.CheckUserMove()
		} else if !s.paused {
			s.board.Reset(s.count, s.level)
			s.lblIntro.Visibe = false
			s.lblResult.Visibe = false
			s.lblMotiv.Visibe = false
			s.lblName.DrawRect = true
			s.lblTimer.Visibe = false
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
		if !s.lblIntro.Visibe {
			s.SaveGame()
			var motiv string
			count := getApp().db.todayGamesCount
			s.level, s.lives, motiv = getApp().db.todayData[count].NextLevel()
			ss := getApp().db.todayData[count].String()
			s.lblResult.SetText(ss)
			log.Printf("Game Result is: %v", ss)
			s.count += 1
			s.lblMotiv.SetText(motiv)
			s.lblName.DrawRect = false
			s.lblName.SetText(s.name)
			s.lblName.SetBg(getApp().theme.correct)
			s.lblIntro.Visibe = true
			s.lblResult.Visibe = true
			s.lblMotiv.Visibe = true
			s.lblTimer.Visibe = true
			s.lblTimer.SetBg(getApp().theme.error)
			s.pauseTimer = 5000
			s.paused = true
		}
		if s.pauseTimer > 0 {
			if s.paused {
				s.pauseTimer -= dt
				s.lblTimer.SetText(fmt.Sprintf("%v", s.pauseTimer/1000))
			} else {
				s.pauseTimer += dt
				s.lblTimer.SetText(fmt.Sprintf("%02v:%02v", s.pauseTimer/1000/60, s.pauseTimer/1000%60))
			}
		} else if s.pauseTimer <= 0 {
			if s.paused {
				s.paused = false
				s.pauseTimer += 5000
				s.lblTimer.SetBg(getApp().theme.correct)
			}
		}
	}
}

func (s *SceneGame) moveStatus() {
	switch s.board.moveStatus {
	case Correct:
		s.lblName.SetBg(getApp().theme.correct)
	case Error:
		s.lblName.SetBg(getApp().theme.error)
	case Warning:
		s.lblName.SetBg(getApp().theme.warning)
	case Regular:
		s.lblName.SetBg(getApp().theme.regular)
	default:
		s.lblName.SetBg(getApp().theme.bg)
	}
	str := fmt.Sprintf("Pos %v (%v) (%v/%v)", s.level, s.lives, s.board.move, s.board.totalMoves)
	s.lblName.SetText(str)
}

func (s *SceneGame) SaveGame() {
	dtBeg := s.board.dtBeg.Format("2006.01.02 15:04:05.000")
	dtEnd := s.board.dtEnd.Format("2006.01.02 15:04:05.000")
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
		manual:       false,
		advance:      80,
		fallback:     50,
		resetonerror: false,
	}
	getApp().db.Insert(values)
	log.Println("Game Saved in DB.")
}

func (s *SceneGame) Draw(surface *ebiten.Image) {
	for _, value := range s.container {
		value.Draw(surface)
	}
}

func (s *SceneGame) Add(item ui.Drawable) {
	s.container = append(s.container, item)
}

func (s *SceneGame) Resize() {
	s.rect = getApp().rect
	x, y, w, h := 0, 0, int(float64(getApp().rect.W)*0.3), int(float64(getApp().rect.H)*0.05)
	s.lblName.Resize([]int{x, y, w, h})
	sz := s.rect.GetLowestSize()
	cellSize := float64(sz)/3 - float64(sz)*0.02
	marginX := float64(s.rect.W)/2 - cellSize*3/2
	marginY := float64(s.rect.H)/2 - cellSize*3/2
	x, y = int(marginX), int(marginY)+h/2
	s.board.Resize([]int{x, y, int(cellSize) * 3, int(cellSize) * 3})
	w, h = int(float64(getApp().rect.W)*0.8), int(float64(getApp().rect.H)*0.05)
	x, y = (s.rect.W-w)/2, s.rect.H-int(float64(h)*1.5)
	s.lblIntro.Resize([]int{x, y, w, h})
	w, h = int(float64(getApp().rect.W)*0.9), int(float64(getApp().rect.H)*0.08)
	x, y = (s.rect.W-w)/2, s.rect.H-int(float64(h)*3.5)
	s.lblResult.Resize([]int{x, y, w, h})
	w, h = int(float64(getApp().rect.W)*0.7), int(float64(getApp().rect.H)*0.08)
	x, y = (s.rect.W-w)/2, s.rect.H-int(float64(h)*5.5)
	s.lblMotiv.Resize([]int{x, y, w, h})
	w, h = int(float64(getApp().rect.W)*0.5), int(float64(getApp().rect.H)*0.2)
	x, y = (s.rect.W-w)/2, s.rect.H-int(float64(h)*4)
	s.lblTimer.Resize([]int{x, y, w, h})
}

func (s *SceneGame) Quit() {}
