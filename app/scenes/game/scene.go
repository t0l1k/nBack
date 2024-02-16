package scene_game

import (
	"image/color"
	"log"
	"strconv"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/t0l1k/eui"
	"github.com/t0l1k/nBack/app"
	"github.com/t0l1k/nBack/app/game"
)

type SceneGame struct {
	eui.SceneBase
	lblTitle                                              *eui.Text
	lblVar                                                *eui.StringVar
	btnQuit                                               *eui.Button
	moveTimer                                             *eui.Timer
	gameData                                              *game.GameData
	gameConf                                              game.GameConf
	board                                                 *game.Board
	grid                                                  *eui.GridView
	btnsLayout                                            *eui.BoxLayout
	moveTime, delayTimeShowCell, delayTimeHideCell        int
	posModMove, symModMove, colModMove, ariModMove        bool
	userMoved, resetOnError, resetOpt                     bool
	clrMoved, clrNeutral, clrCorrect, clrWrong, clrMissed color.Color
}

func New() *SceneGame {
	s := &SceneGame{}
	s.lblTitle = eui.NewText("nBack ") // (модальность уровень)(ходов осталось)
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
	return s
}

func (s *SceneGame) Setup(conf game.GameConf, gd *game.GameData) {
	s.gameConf = conf
	s.gameData = gd
	s.resetOpt = conf.Get(game.ResetOnFirstWrong).(bool)
	s.resetOnError = false
	theme := eui.GetUi().GetTheme()
	s.grid.Bg(theme.Get(app.GameColorBg))
	s.grid.Fg(theme.Get(app.GameColorFgCrosshair))
	s.grid.Visible(conf.Get(game.ShowGrid).(bool))
	s.lblTitle.Fg(theme.Get(app.GameColorBg))
	s.btnsLayout.ResetContainerBase()
	for _, v := range s.gameData.Modalities {
		btn := eui.NewButton(string(v.GetSym()), s.buttonsLogic)
		btn.Bg(theme.Get(app.LabelColorDefault))
		btn.Fg(theme.Get(app.GameColorBg))
		s.btnsLayout.Add(btn)
		if v.GetSym() == game.Pos {
			grid := conf.Get(game.GridSize).(int)
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
	showCellPercent := conf.Get(game.ShowCellPercent).(float64)
	timeShowCell := int(float64(s.moveTime) * showCellPercent)
	s.delayTimeShowCell = (s.moveTime - timeShowCell) / 2
	s.delayTimeHideCell = s.delayTimeShowCell + timeShowCell
	s.moveTimer = eui.NewTimer(s.delayTimeShowCell) // pause before first move
	s.board.Setup(conf, s.gameData)
	s.lblTitle.Bg(s.clrNeutral)
	log.Printf("init move timer:%v show time:%v delay before show:%v delay hide:%v", s.moveTime, timeShowCell, s.delayTimeShowCell, s.delayTimeHideCell)
}

func (s *SceneGame) Entered() {
	s.Resize()
	eui.GetUi().GetInputKeyboard().Attach(s)
	s.moveTimer.On()
	s.moveTimer.SetDuration(s.moveTime)
	s.board.MakeMove()
	s.board.Visible(false)
	log.Println("begin play:00 hide cell", s.board.Move)
}

func (s *SceneGame) Update(dt int) {
	s.moveTimer.Update(dt)
	for _, v := range s.GetContainer() {
		v.Update(dt)
	}
	for _, v := range s.btnsLayout.GetContainer() {
		v.Update(dt)
	}
	if s.moveTimer.TimePassed() > s.delayTimeHideCell && s.board.IsVisible() {
		s.board.Visible(false)
		log.Println("02 hide cell", s.board.Move)
	}
	if s.moveTimer.IsDone() {
		s.resetColorsAfterMove()
		s.checkProgress()
		if s.board.Move >= s.gameData.TotalMoves || s.resetOpt && s.resetOnError {
			log.Println("last move check")
			s.sendResult()
		} else {
			s.board.NextMove()
			log.Println("01 show cell", s.board.Move)
			s.moveTimer.Reset()
			log.Println("reset move timer:", s.board.Move)
		}
	}
	s.updateLbls()
}

func (s *SceneGame) resetColorsAfterMove() {
	s.lblTitle.Bg(s.clrNeutral)
	for _, v := range s.btnsLayout.GetContainer() {
		v.(*eui.Button).Bg(s.clrNeutral)
	}
}

func (s *SceneGame) checkProgress() {
	if s.board.Move <= s.gameData.Level {
		if s.board.Move <= 0 {
			return
		}
		for _, v := range s.gameData.Modalities {
			v.SetRegular()
		}
		return
	}
	for _, v := range s.gameData.Modalities {
		if v.GetSym() == game.Pos {
			str := v.CheckMove(s.posModMove, s.board.LastMove, s.board.TestMove)
			s.posModMove = false
			log.Println(str)
		}
		if v.GetSym() == game.Col {
			str := v.CheckMove(s.colModMove, s.board.LastMove, s.board.TestMove)
			s.colModMove = false
			log.Println(str)
		}
		if v.GetSym() == game.Sym {
			str := v.CheckMove(s.symModMove, s.board.LastMove, s.board.TestMove)
			s.symModMove = false
			log.Println(str)
		}
		if v.GetSym() == game.Ari {
			str := v.CheckMove(s.ariModMove, s.board.LastMove, s.board.TestMove)
			s.ariModMove = false
			log.Println(str)
		}
	}
}

func (s *SceneGame) updateLbls() {
	var str strings.Builder
	str.WriteString(string(s.gameData.GameMode()))
	str.WriteString("(")
	str.WriteString(strconv.Itoa(s.gameData.TotalMoves - s.board.Move))
	str.WriteString(")")
	s.lblVar.SetValue(str.String())
	if s.userMoved {
		s.lblTitle.Bg(s.clrMoved)
		for _, v := range s.btnsLayout.GetContainer() {
			if s.posModMove {
				if v.(*eui.Button).GetText() == string(game.Pos) {
					v.(*eui.Button).Bg(s.clrMoved)
				}
			}
			if s.colModMove {
				if v.(*eui.Button).GetText() == string(game.Col) {
					v.(*eui.Button).Bg(s.clrMoved)
				}
			}
			if s.symModMove {
				if v.(*eui.Button).GetText() == string(game.Sym) {
					v.(*eui.Button).Bg(s.clrMoved)
				}
			}
			if s.ariModMove {
				if v.(*eui.Button).GetText() == string(game.Ari) {
					v.(*eui.Button).Bg(s.clrMoved)
				}
			}
		}
		s.userMoved = false
	}
}

func (s *SceneGame) buttonsLogic(b *eui.Button) {
	switch b.GetText() {
	case game.Pos.String():
		s.userMove(game.Pos.String())
		log.Printf("button <%v> pressed", b.GetText())
	case game.Col.String():
		s.userMove(game.Col.String())
		log.Printf("button <%v> pressed", b.GetText())
	case game.Sym.String():
		s.userMove(game.Sym.String())
		log.Printf("button <%v> pressed", b.GetText())
	case game.Ari.String():
		s.userMove(game.Ari.String())
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
				s.userMove(game.Pos.String())
				log.Println("pressed <A>")
			} else if key == ebiten.KeyC {
				s.userMove(game.Col.String())
				log.Println("pressed <C>")
			} else if key == ebiten.KeyS {
				s.userMove(game.Sym.String())
				log.Println("pressed <S>")
			} else if key == ebiten.KeyR {
				s.userMove(game.Ari.String())
				log.Println("pressed <R>")
			}
		}
	}
}

func (s *SceneGame) userMove(value string) {
	if value == game.Pos.String() {
		s.posModMove = true
	} else if value == game.Col.String() {
		s.colModMove = true
	} else if value == game.Sym.String() {
		s.symModMove = true
	} else if value == game.Ari.String() {
		s.ariModMove = true
	}
	s.userMoved = true
}

func (s *SceneGame) sendResult() {
	s.gameData.SetGameDone(s.board.Move)
	for _, v := range s.gameData.Modalities {
		v.Detach(s)
	}
	log.Println(s.gameData.LastGameFullResult())
	eui.GetUi().GetInputKeyboard().Detach(s)
	eui.GetUi().Pop()
}

func (s *SceneGame) UpdateData(value interface{}) {
	switch v := value.(type) {
	case map[game.ModalType]game.MoveType:
		for k1, v1 := range v {
			var clr color.Color
			switch v1 {
			case game.AddCorrect:
				clr = s.clrCorrect
			case game.AddWrong:
				clr = s.clrWrong
				s.resetOnError = true
			case game.AddMissed:
				clr = s.clrMissed
				s.resetOnError = true
			case game.AddRegular:
				clr = s.clrNeutral
			}
			s.lblTitle.Bg(clr)
			for _, btn := range s.btnsLayout.GetContainer() {
				if btn.(*eui.Button).GetText() == k1.String() {
					btn.(*eui.Button).Bg(clr)
				}
			}
		}
	}
}

func (s *SceneGame) Draw(surface *ebiten.Image) {
	for _, v := range s.GetContainer() {
		v.Draw(surface)
	}
	for _, v := range s.btnsLayout.GetContainer() {
		v.Draw(surface)
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
