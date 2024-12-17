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
	lblVar                                                *eui.SubjectBase
	btnQuit                                               *eui.Button
	moveTimer, gameTimer                                  *eui.Timer
	gameData                                              *game.GameData
	gameConf                                              game.GameConf
	board                                                 *game.Board
	grid                                                  *eui.GridView
	btnsLayout                                            *eui.BoxLayout
	moveTime, delayTimeShowCell, delayTimeHideCell        int
	posModMove, colModMove, symModMove                    bool
	userMoved, resetOnError, resetOpt, showLbl, checkIn   bool
	clrMoved, clrNeutral, clrCorrect, clrWrong, clrMissed color.Color
	posModalKey, colorModalKey, symbolModalKey            ebiten.Key
	nextLevelDialog                                       *nextLevelDialog
}

func New() *SceneGame {
	s := &SceneGame{}
	s.lblTitle = eui.NewText("nBack ") // (модальность уровень)(ходов осталось)
	s.Add(s.lblTitle)
	s.lblVar = eui.NewSubject()
	s.lblVar.Attach(s.lblTitle)
	s.btnQuit = eui.NewButton("<", func(b *eui.Button) {
		eui.GetUi().Pop()
	})
	s.Add(s.btnQuit)
	s.nextLevelDialog = newNextLevelDialog(5000)
	s.Add(s.nextLevelDialog)
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
	s.btnsLayout.ResetContainerBase()
	for _, v := range s.gameData.Modalities {
		btn := eui.NewButton(string(v.GetSym()), s.buttonsLogic)
		btn.Bg(theme.Get(app.LabelColorDefault))
		btn.Fg(theme.Get(app.GameColorBg))
		s.btnsLayout.Add(btn)
		if v.GetSym() == game.Pos {
			grid := conf.Get(game.GridSize).(int)
			s.grid.Set(float64(grid), float64(grid))
		}
		v.Attach(s)
	}
	s.clrNeutral = theme.Get(app.LabelColorDefault)
	s.clrMoved = theme.Get(app.ColorNeutral)
	s.clrCorrect = theme.Get(app.ColorCorrect)
	s.clrWrong = theme.Get(app.ColorWrong)
	s.clrMissed = theme.Get(app.ColorMissed)
	appConf := eui.GetUi().GetSettings()
	s.posModalKey = appConf.Get(app.PositionKeypress).(ebiten.Key)
	s.colorModalKey = appConf.Get(app.ColorKeypress).(ebiten.Key)
	s.symbolModalKey = appConf.Get(app.SymbolKeypress).(ebiten.Key)
	totalTime := conf.Get(game.TotalTime).(int)
	if totalTime > 0 {
		s.gameTimer = eui.NewTimer(totalTime * 60 * 1000)
		if s.gameData.CheckIn {
			s.checkIn = true
		}
	}
	s.moveTime = int(s.gameData.MoveTime * 1000)
	showCellPercent := conf.Get(game.ShowCellPercent).(float64)
	timeShowCell := int(float64(s.moveTime) * showCellPercent)
	s.delayTimeShowCell = (s.moveTime - timeShowCell) / 2
	s.delayTimeHideCell = s.delayTimeShowCell + timeShowCell
	s.moveTimer = eui.NewTimer(s.delayTimeShowCell) // pause before first move
	s.board.Setup(conf, s.gameData)
	s.showLbl = conf.Get(game.ShowGameLabel).(bool)
	s.lblTitle.Visible(s.showLbl)
	s.lblTitle.Bg(s.clrNeutral)
	s.lblTitle.Fg(theme.Get(app.GameColorBg))
	confApp := eui.GetUi().GetSettings()
	restDuration := confApp.Get(app.RestDuration).(int)
	s.nextLevelDialog.timer.SetDuration(restDuration * 1000)
	log.Printf("init move timer:%v show time:%v delay before show:%v delay hide:%v", s.moveTime, timeShowCell, s.delayTimeShowCell, s.delayTimeHideCell)
}

func (s *SceneGame) Entered() {
	s.Resize()
	eui.GetUi().GetInputKeyboard().Attach(s)
	if s.gameTimer != nil {
		s.gameTimer.On()
	}
	s.moveTimer.On()
	s.moveTimer.SetDuration(s.moveTime)
	s.board.MakeMove()
	s.board.Visible(false)
	log.Println("begin play:00 hide cell", s.board.Move)
}

func (s *SceneGame) Update(dt int) {
	for _, v := range s.GetContainer() {
		v.Update(dt)
	}
	if s.nextLevelDialog.IsVisible() {
		return
	}
	if s.gameTimer != nil {
		s.gameTimer.Update(dt)
	}
	s.moveTimer.Update(dt)
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
		if !s.checkIn && s.board.Move >= s.gameData.TotalMoves || (!s.checkIn && s.resetOpt && s.resetOnError) || (s.gameTimer != nil && s.gameTimer.IsDone()) {
			log.Println("last move check")
			s.sendResult()
		} else if s.checkIn && s.board.Move >= s.gameData.TotalMoves || s.checkIn && s.resetOpt && s.resetOnError {
			s.board.Reset()
			s.resetOnError = false
			msg, col := s.gameData.CheckNextLevel(s.gameConf)
			s.nextLevelDialog.show(msg, col)
			s.gameData.FillField(s.gameConf)
			log.Println("begin play:05 hide cell", s.board.Move)
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
	if !s.showLbl {
		return
	}
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
			v.SetRegular(s.gameData.Level)
		}
		return
	}
	for _, v := range s.gameData.Modalities {
		if v.GetSym() == game.Pos {
			str := v.CheckMove(s.posModMove, s.board.LastMove, s.board.TestMove, s.gameData.Level)
			s.posModMove = false
			log.Println(str)
		}
		if v.GetSym() == game.Col {
			str := v.CheckMove(s.colModMove, s.board.LastMove, s.board.TestMove, s.gameData.Level)
			s.colModMove = false
			log.Println(str)
		}
		if v.GetSym() == game.Sym {
			str := v.CheckMove(s.symModMove, s.board.LastMove, s.board.TestMove, s.gameData.Level)
			s.symModMove = false
			log.Println(str)
		} else if v.GetSym() == game.Ari {
			str := v.CheckMove(s.symModMove, s.board.LastMove, s.board.TestMove, s.gameData.Level)
			s.symModMove = false
			log.Println(str)
		}
	}
}

func (s *SceneGame) updateLbls() {
	if !s.showLbl {
		if s.userMoved {
			s.userMoved = false
		}
		return
	}
	var str strings.Builder
	str.WriteString(string(s.gameData.GameMode()))

	if s.gameTimer != nil {
		str.WriteString("(")
		str.WriteString(strconv.Itoa(s.board.Move))
		str.WriteString("/")
		str.WriteString(strconv.Itoa(s.gameData.TotalMoves))
		str.WriteString(")")

		str.WriteString("(")
		str.WriteString(s.gameTimer.String())
		str.WriteString(")")
	} else {
		str.WriteString("(")
		str.WriteString(strconv.Itoa(s.gameData.TotalMoves - s.board.Move))
		str.WriteString(")")

	}

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
				} else if v.(*eui.Button).GetText() == string(game.Ari) {
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
			switch key {
			case s.posModalKey:
				s.userMove(game.Pos.String())
				log.Println("pressed pos modal")
			case s.colorModalKey:
				s.userMove(game.Col.String())
				log.Println("pressed color modal")
			case s.symbolModalKey:
				s.userMove(game.Sym.String())
				log.Println("pressed symbol modal")
			}
		}
	}
}

func (s *SceneGame) userMove(value string) {
	switch value {
	case game.Pos.String():
		s.posModMove = true
	case game.Col.String():
		s.colModMove = true
	case game.Sym.String():
		s.symModMove = true
	case game.Ari.String():
		s.symModMove = true
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
	if !s.showLbl {
		return
	}
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
	if s.nextLevelDialog.IsVisible() {
		s.nextLevelDialog.Draw(surface)
		return
	}
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
	s.grid.Resize(s.board.GetLayRect())
	y += h0 - h*4 + h/2
	s.btnsLayout.Resize([]int{x, y, w0 - h, h * 2})
	w1, h1 := w0/2, h0/2
	x, y = (w0-w1)/2, (h0-h1)/2
	s.nextLevelDialog.Resize([]int{x, y, w1, h1})
}
