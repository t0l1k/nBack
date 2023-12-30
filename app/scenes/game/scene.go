package scene_game

import (
	"fmt"
	"image/color"
	"log"
	"strconv"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/t0l1k/eui"
	"github.com/t0l1k/nBack/app"
	"github.com/t0l1k/nBack/app/data"
	"github.com/t0l1k/nBack/app/game"
)

type SceneGame struct {
	eui.SceneBase
	lblTitle                                              *eui.Text
	lblVar                                                *eui.StringVar
	moveTimer                                             *eui.Timer
	gameData                                              *data.GameData
	game                                                  *game.Board
	grid                                                  *eui.GridView
	btnsLayout                                            *eui.BoxLayout
	moveTime, delayTimeShowCell, delayTimeHideCell        int
	posModMove, symModMove, colModMove, userMoved         bool
	clrMoved, clrNeutral, clrCorrect, clrWrong, clrMissed color.Color
}

func New() *SceneGame {
	s := &SceneGame{}
	s.lblTitle = eui.NewText("nBack ") // (модальность уровень) (попыток) (ход/ходов)
	s.Add(s.lblTitle)
	s.lblVar = eui.NewStringVar("")
	s.lblVar.Attach(s.lblTitle)
	s.game = game.New()
	s.Add(s.game)
	s.grid = eui.NewGridView(1, 1)
	s.grid.Bg(eui.Black)
	s.grid.Fg(eui.Aqua)
	s.Add(s.grid)
	s.btnsLayout = eui.NewHLayout()
	s.Add(s.btnsLayout)
	return s
}

func (s *SceneGame) Setup(gd *data.GameData) {
	s.gameData = gd
	conf := eui.GetUi().GetSettings()
	s.btnsLayout.Container = nil
	for _, v := range s.gameData.Modalities {
		btn := eui.NewButton(v.GetSym(), s.buttonsLogic)
		s.btnsLayout.Add(btn)
		if v.GetSym() == data.Pos {
			grid := conf.Get(app.GridSize).(int)
			s.grid.Set(grid, grid)
		}
		v.Attach(s)
	}
	s.clrNeutral = conf.Get(app.LabelColorDefault).(color.Color)
	s.clrMoved = conf.Get(app.ColorNeutral).(color.Color)
	s.clrCorrect = conf.Get(app.ColorCorrect).(color.Color)
	s.clrWrong = conf.Get(app.ColorWrong).(color.Color)
	s.clrMissed = conf.Get(app.ColorMissed).(color.Color)
	s.moveTime = int(s.gameData.MoveTime * 1000)
	showCellPercent := conf.Get(app.ShowCellPercent).(float64)
	timeShowCell := int(float64(s.moveTime) * showCellPercent)
	s.delayTimeShowCell = (s.moveTime - timeShowCell) / 2
	s.delayTimeHideCell = s.delayTimeShowCell + timeShowCell
	s.moveTimer = eui.NewTimer(s.moveTime + s.delayTimeShowCell)
	s.game.Setup(s.gameData)
	s.lblTitle.Bg(s.clrNeutral)
	log.Println("init:", s.moveTime, timeShowCell, s.delayTimeShowCell, s.delayTimeHideCell)
}

func (s *SceneGame) buttonsLogic(b *eui.Button) {
	switch b.GetText() {
	case data.Pos:
		s.userMove(data.Pos)
		log.Printf("button <%v> pressed", b.GetText())
	case data.Col:
		s.userMove(data.Col)
		log.Printf("button <%v> pressed", b.GetText())
	}
}

func (s *SceneGame) Entered() {
	s.Resize()
	eui.GetUi().GetInputKeyboard().Attach(s)
	s.moveTimer.On()
	s.moveTimer.SetDuration(s.moveTime)
	s.game.MakeMove()
	s.game.Visible(false)
	log.Println("begin play:00 hide cell", s.game.Move, s.moveTimer.TimePassed())
}

func (s *SceneGame) Update(dt int) {
	for _, v := range s.btnsLayout.Container {
		v.Update(dt)
	}
	s.moveTimer.Update(dt)
	if s.moveTimer.TimePassed() > s.delayTimeShowCell && s.moveTimer.TimePassed() < s.delayTimeHideCell && !s.game.IsVisible() {
		s.checkProgress()
		s.game.NextMove()
		s.game.Visible(true)
		log.Println("01 show cell", s.game.Move, s.moveTimer.TimePassed())
	}
	if s.moveTimer.TimePassed() > s.delayTimeHideCell && s.game.IsVisible() {
		s.game.Visible(false)
		log.Println("02 hide cell", s.game.Move, s.moveTimer.TimePassed())
	}
	if s.moveTimer.IsDone() {
		s.resetColors()
		if s.game.Move >= s.gameData.TotalMoves {
			log.Println("last move check")
			s.checkProgress()
			s.sendResult()
		} else {
			s.moveTimer.Reset()
			log.Println("reset move timer:", s.game.Move)
		}
	}
	s.updateLbls()
}

func (s *SceneGame) resetColors() {
	s.lblTitle.Bg(s.clrNeutral)
	for _, v := range s.btnsLayout.Container {
		v.(*eui.Button).Bg(s.clrNeutral)
	}
}

func (s *SceneGame) checkProgress() {
	if s.game.Move <= s.gameData.Level {
		return
	}
	var str string
	for i, v := range s.gameData.Modalities {
		if v.GetSym() == data.Pos {
			lastMove, testMove := v.GetField()[s.game.LastMove], v.GetField()[s.game.TestMove]
			str = fmt.Sprintf("game progress for modal[%v] level(%v) moves[%v-%v] values:[%v-%v] timer:%v ", s.gameData.Modalities[i].GetSym(), s.gameData.Level, s.game.Move, s.game.Move-s.gameData.Level, testMove, lastMove, s.moveTimer.TimePassed())
			if s.posModMove {
				if lastMove == testMove {
					s.gameData.Modalities[i].SetCorrect(1)
					str += "correct answer!"
				} else {
					s.gameData.Modalities[i].SetWrong(1)
					str += "wrong answer!"
				}
				s.posModMove = false
			} else if lastMove == testMove {
				s.gameData.Modalities[i].SetMissed(1)
				str += "missed answer!"
			}
			log.Println(str)
		}

		if v.GetSym() == data.Col {
			lastMove, testMove := v.GetField()[s.game.LastMove], v.GetField()[s.game.TestMove]
			str = fmt.Sprintf("game progress for modal[%v] level(%v) moves[%v-%v] values:[%v-%v] timer:%v ", s.gameData.Modalities[i].GetSym(), s.gameData.Level, s.game.Move, s.game.Move-s.gameData.Level, testMove, lastMove, s.moveTimer.TimePassed())
			if s.colModMove {
				if lastMove == testMove {
					s.gameData.Modalities[i].SetCorrect(1)
					str += "correct answer!"
				} else {
					s.gameData.Modalities[i].SetWrong(1)
					str += "wrong answer!"
				}
				s.colModMove = false
			} else if lastMove == testMove {
				s.gameData.Modalities[i].SetMissed(1)
				str += "missed answer!"
			}
			log.Println(str)
		}

		if v.GetSym() == data.Sym {
			lastMove, testMove := v.GetField()[s.game.LastMove], v.GetField()[s.game.TestMove]
			str = fmt.Sprintf("game progress for modal[%v] level(%v) moves[%v-%v] values:[%v-%v] timer:%v ", s.gameData.Modalities[i].GetSym(), s.gameData.Level, s.game.Move, s.game.Move-s.gameData.Level, testMove, lastMove, s.moveTimer.TimePassed())
			if s.symModMove {
				if lastMove == testMove {
					s.gameData.Modalities[i].SetCorrect(1)
					str += "correct answer!"
				} else {
					s.gameData.Modalities[i].SetWrong(1)
					str += "wrong answer!"
				}
				s.symModMove = false
			} else if lastMove == testMove {
				s.gameData.Modalities[i].SetMissed(1)
				str += "missed answer!"
			}
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
	str.WriteString(strconv.Itoa(s.game.Move))
	str.WriteString("/")
	str.WriteString(strconv.Itoa(s.gameData.TotalMoves))
	str.WriteString(")")
	s.lblVar.SetValue(str.String())

	if s.userMoved {
		s.lblTitle.Bg(s.clrMoved)
		for _, v := range s.btnsLayout.Container {
			if s.posModMove {
				if v.(*eui.Button).GetText() == data.Pos {
					v.(*eui.Button).Bg(s.clrMoved)
				}
			} else if s.colModMove {
				if v.(*eui.Button).GetText() == data.Col {
					v.(*eui.Button).Bg(s.clrMoved)
				}
			} else if s.symModMove {
				if v.(*eui.Button).GetText() == data.Sym {
					v.(*eui.Button).Bg(s.clrMoved)
				}
			}
		}
		s.userMoved = false
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
	}
	s.userMoved = true
}

func (s *SceneGame) sendResult() {
	s.gameData.SetGameDone(s.game.Move)
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
		if v[0] == data.Pos {
			for _, btn := range s.btnsLayout.Container {
				if btn.(*eui.Button).GetText() == data.Pos {
					btn.(*eui.Button).Bg(clr)
				}
			}
		} else if v[0] == data.Col {
			for _, btn := range s.btnsLayout.Container {
				if btn.(*eui.Button).GetText() == data.Col {
					btn.(*eui.Button).Bg(clr)
				}
			}
		} else if v[0] == data.Sym {
			for _, btn := range s.btnsLayout.Container {
				if btn.(*eui.Button).GetText() == data.Sym {
					btn.(*eui.Button).Bg(clr)
				}
			}
		}
	}
}

func (s *SceneGame) Resize() {
	w0, h0 := eui.GetUi().Size()
	w := int(float64(w0) * 0.268)
	h := int(float64(h0) * 0.05)
	x, y := 0, 0
	s.lblTitle.Resize([]int{x, y, w, h})
	x = h / 2
	y += h + h/2
	s.game.Resize([]int{x, y, w0 - h, h0 - h*4})
	s.grid.Resize([]int{x, y, w0 - h, h0 - h*4})
	y += h0 - h*4 + h/2
	s.btnsLayout.Resize([]int{x, y, w0 - h, h * 2})
}
