package scene_game

import (
	"image/color"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/t0l1k/eui"
	"github.com/t0l1k/nBack/app"
	"github.com/t0l1k/nBack/app/data"
	"github.com/t0l1k/nBack/app/game"
)

type SceneGame struct {
	eui.SceneBase
	lblTitle                                                  *eui.Text
	lblVar                                                    *eui.StringVar
	btnQuit                                                   *eui.Button
	moveTimer                                                 *eui.Timer
	gameData                                                  *data.GameData
	gameConf                                                  data.GameConf
	board                                                     *game.Board
	grid                                                      *eui.GridView
	btnsLayout                                                *eui.BoxLayout
	moveTime, delayTimeShowCell, delayTimeHideCell            int
	posModMove, symModMove, colModMove, ariModMove, userMoved bool
	clrMoved, clrNeutral, clrCorrect, clrWrong, clrMissed     color.Color
}

func New() *SceneGame {
	s := &SceneGame{}
	s.lblTitle = eui.NewText("nBack ") // (модальность уровень) (попыток) (ходов осталось)
	s.Add(s.lblTitle)
	s.lblVar = eui.NewStringVar("")
	s.lblVar.Attach(s.lblTitle)
	s.btnQuit = eui.NewButton("<", func(b *eui.Button) {
		eui.GetUi().Pop()
	})
	s.Add(s.btnQuit)
	s.board = game.New()
	s.Add(s.board)
	s.grid = eui.NewGridView(1, 1)
	s.Add(s.grid)
	s.btnsLayout = eui.NewHLayout()
	s.Add(s.btnsLayout)
	return s
}

func (s *SceneGame) Setup(conf data.GameConf, gd *data.GameData) {
	s.gameConf = conf
	s.gameData = gd
	theme := eui.GetUi().GetTheme()
	s.grid.Bg(theme.Get(app.GameColorBg))
	s.grid.Fg(theme.Get(app.GameColorFgCrosshair))
	s.grid.Visible(conf.Get(data.ShowGrid).(bool))
	s.lblTitle.Fg(theme.Get(app.GameColorBg))
	s.btnsLayout.Container = nil
	for _, v := range s.gameData.Modalities {
		btn := eui.NewButton(v.GetSym(), s.buttonsLogic)
		btn.Bg(theme.Get(app.LabelColorDefault))
		btn.Fg(theme.Get(app.GameColorBg))
		s.btnsLayout.Add(btn)
		if v.GetSym() == data.Pos {
			grid := conf.Get(data.GridSize).(int)
			s.grid.Set(grid, grid)
		}
		v.Attach(s)
	}
	s.clrNeutral = theme.Get(app.LabelColorDefault)
	s.clrMoved = theme.Get(app.ColorNeutral)
	s.clrCorrect = theme.Get(app.ColorCorrect)
	s.clrWrong = theme.Get(app.ColorWrong)
	s.clrMissed = theme.Get(app.ColorMissed)
	s.moveTime = int(s.gameData.MoveTime * 1000)
	appConf := eui.GetUi().GetSettings()
	showCellPercent := appConf.Get(app.ShowCellPercent).(float64)
	timeShowCell := int(float64(s.moveTime) * showCellPercent)
	s.delayTimeShowCell = (s.moveTime - timeShowCell) / 2
	s.delayTimeHideCell = s.delayTimeShowCell + timeShowCell
	s.moveTimer = eui.NewTimer(s.moveTime)
	s.board.Setup(conf, s.gameData)
	s.lblTitle.Bg(s.clrNeutral)
	log.Println("init:", s.moveTime, timeShowCell, s.delayTimeShowCell, s.delayTimeHideCell)
}

func (s *SceneGame) Entered() {
	s.Resize()
	eui.GetUi().GetInputKeyboard().Attach(s)
	time.Sleep(time.Duration(s.delayTimeShowCell))
	s.moveTimer.On()
	s.board.MakeMove()
	s.board.Visible(false)
	log.Println("begin play:00 hide cell", s.board.Move)
}

func (s *SceneGame) Update(dt int) {
	s.moveTimer.Update(dt)
	s.btnQuit.Update(dt)
	for _, v := range s.btnsLayout.Container {
		v.Update(dt)
	}
	if s.moveTimer.TimePassed() > s.delayTimeShowCell && s.moveTimer.TimePassed() < s.delayTimeHideCell && !s.board.IsVisible() {
		s.checkProgress()
		s.board.NextMove()
		log.Println("01 show cell", s.board.Move)
	}
	if s.moveTimer.TimePassed() > s.delayTimeHideCell && s.board.IsVisible() {
		s.board.Visible(false)
		log.Println("02 hide cell", s.board.Move)
	}
	if s.moveTimer.IsDone() {
		s.resetColorsAfterMove()
		if s.board.Move >= s.gameData.TotalMoves {
			log.Println("last move check")
			s.checkProgress()
			s.sendResult()
		} else {
			s.moveTimer.Reset()
			log.Println("reset move timer:", s.board.Move)
		}
	}
	s.updateLbls()
}

func (s *SceneGame) resetColorsAfterMove() {
	s.lblTitle.Bg(s.clrNeutral)
	for _, v := range s.btnsLayout.Container {
		v.(*eui.Button).Bg(s.clrNeutral)
	}
}

func (s *SceneGame) checkProgress() {
	if s.board.Move <= s.gameData.Level {
		return
	}
	for _, v := range s.gameData.Modalities {
		if v.GetSym() == data.Pos {
			str := v.CheckMove(s.posModMove, s.board.LastMove, s.board.TestMove)
			s.posModMove = false
			log.Println(str)
		}
		if v.GetSym() == data.Col {
			str := v.CheckMove(s.colModMove, s.board.LastMove, s.board.TestMove)
			s.colModMove = false
			log.Println(str)
		}
		if v.GetSym() == data.Sym {
			str := v.CheckMove(s.symModMove, s.board.LastMove, s.board.TestMove)
			s.symModMove = false
			log.Println(str)
		}
		if v.GetSym() == data.Ari {
			str := v.CheckMove(s.ariModMove, s.board.LastMove, s.board.TestMove)
			s.ariModMove = false
			log.Println(str)
		}
	}
}

func (s *SceneGame) updateLbls() {
	var str strings.Builder
	str.WriteString(s.gameData.GameMode())
	str.WriteString("(")
	str.WriteString(strconv.Itoa(s.gameData.Lives))
	str.WriteString(")")
	str.WriteString("(")
	str.WriteString(strconv.Itoa(s.gameData.TotalMoves - s.board.Move))
	str.WriteString(")")
	s.lblVar.SetValue(str.String())
	if s.userMoved {
		s.lblTitle.Bg(s.clrMoved)
		for _, v := range s.btnsLayout.Container {
			if s.posModMove {
				if v.(*eui.Button).GetText() == data.Pos {
					v.(*eui.Button).Bg(s.clrMoved)
				}
			}
			if s.colModMove {
				if v.(*eui.Button).GetText() == data.Col {
					v.(*eui.Button).Bg(s.clrMoved)
				}
			}
			if s.symModMove {
				if v.(*eui.Button).GetText() == data.Sym {
					v.(*eui.Button).Bg(s.clrMoved)
				}
			}
			if s.ariModMove {
				if v.(*eui.Button).GetText() == data.Ari {
					v.(*eui.Button).Bg(s.clrMoved)
				}
			}
		}
		s.userMoved = false
	}
}

func (s *SceneGame) buttonsLogic(b *eui.Button) {
	switch b.GetText() {
	case data.Pos:
		s.userMove(data.Pos)
		log.Printf("button <%v> pressed", b.GetText())
	case data.Col:
		s.userMove(data.Col)
		log.Printf("button <%v> pressed", b.GetText())
	case data.Sym:
		s.userMove(data.Sym)
		log.Printf("button <%v> pressed", b.GetText())
	case data.Ari:
		s.userMove(data.Ari)
		log.Printf("button <%v> pressed", b.GetText())
	}
}

func (s *SceneGame) UpdateInput(value interface{}) {
	switch v := value.(type) {
	case eui.KeyboardData:
		for _, key := range v.GetKeys() {
			if key == ebiten.KeySpace {
				log.Println("pressed <space>")
			} else if key == ebiten.KeyA {
				s.userMove(data.Pos)
				log.Println("pressed <A>")
			} else if key == ebiten.KeyC {
				s.userMove(data.Col)
				log.Println("pressed <C>")
			} else if key == ebiten.KeyS {
				s.userMove(data.Sym)
				log.Println("pressed <S>")
			} else if key == ebiten.KeyR {
				s.userMove(data.Ari)
				log.Println("pressed <R>")
			}
		}
	}
}

func (s *SceneGame) userMove(value string) {
	if value == data.Pos {
		s.posModMove = true
	} else if value == data.Col {
		s.colModMove = true
	} else if value == data.Sym {
		s.symModMove = true
	} else if value == data.Ari {
		s.ariModMove = true
	}
	s.userMoved = true
}

func (s *SceneGame) sendResult() {
	s.gameData.SetGameDone(s.board.Move)
	for _, v := range s.gameData.Modalities {
		v.Detach(s)
	}
	log.Println(s.gameData)
	eui.GetUi().GetInputKeyboard().Detach(s)
	eui.GetUi().Pop()
}

func (s *SceneGame) UpdateData(value interface{}) {
	switch v := value.(type) {
	case []string:
		var clr color.Color
		if v[1] == data.AddCorrect {
			clr = s.clrCorrect
		} else if v[1] == data.AddWrong {
			clr = s.clrWrong
		} else if v[1] == data.AddMissed {
			clr = s.clrMissed
		}
		s.lblTitle.Bg(clr)
		for _, btn := range s.btnsLayout.Container {
			if btn.(*eui.Button).GetText() == v[0] {
				btn.(*eui.Button).Bg(clr)
			}
		}
	}
}

func (s *SceneGame) Resize() {
	w0, h0 := eui.GetUi().Size()
	w := int(float64(w0) * 0.268)
	h := int(float64(h0) * 0.05)
	x, y := 0, 0
	s.btnQuit.Resize([]int{x, y, h, h})
	x += h
	s.lblTitle.Resize([]int{x, y, w, h})
	x = h / 2
	y += h + h/2
	s.board.Resize([]int{x, y, w0 - h, h0 - h*4})
	s.grid.Resize([]int{x, y, w0 - h, h0 - h*4})
	y += h0 - h*4 + h/2
	s.btnsLayout.Resize([]int{x, y, w0 - h, h * 2})
}
