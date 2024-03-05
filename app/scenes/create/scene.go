package create

import (
	"fmt"
	"math"

	"github.com/t0l1k/eui"
	"github.com/t0l1k/nBack/app"
	"github.com/t0l1k/nBack/app/data"
	"github.com/t0l1k/nBack/app/game"
	scene_game "github.com/t0l1k/nBack/app/scenes/game"
)

const (
	bNew   = "Применить"
	bReset = "Обнулить"
	bTest  = "Тестировать"
)

var (
	dtTitleMoves = []string{"Мало", "Нормально", "Много", "Экстремально"}
	movesArr     = [][]int{{10, 1, 1}, {20, 1, 1}, {20, 5, 1}, {20, 1, 2}}
	strMod       = "Модальность "
	symTitle     = []string{strMod + "не использовать", strMod + "Цифры", strMod + "Арифметика"}
	symData      = []interface{}{" ", game.Sym.String(), game.Ari.String()}
)

type SceneCreateGame struct {
	eui.SceneBase
	topBar                                                  *eui.TopBar
	profile                                                 *data.GameProfiles
	inpName                                                 *eui.InputBox
	lblSelectModal, lblSelectMoves, lblSelectThreshold      *eui.Text
	optPos, optCol, optAddSub, optMulDiv, optShowGameLbl    *eui.Checkbox
	optCrossHair, optGrid, optUseCenter, optReset           *eui.Checkbox
	cDefLev, cSym, cGridSz, cMoves, cMoveTime, cShowCellTm  *eui.ComboBox
	cThUp, cThDown, cThAdv, cThFall, cRR, cMaxNum           *eui.ComboBox
	iColors                                                 *eui.Icon
	btnApply, btnReset, btnTest                             *eui.Button
	selPos, selCol, selSym, selAri                          bool
	defLevel, movesConfIndex, rr, gridSz, maxNum            int
	moveTime, showCellTime                                  float64
	thresholdUp, thresholdDown, thresholdAdv, thresholdFall int
	showGrid, showCrossHair, useCenterCell, resetOnWrong    bool
	useMulDiv, useAddSub, showGameLabel                     bool
	inTesting                                               bool
	examples                                                *ExamplesFrame
}

func NewSceneCreateGame(profile *data.GameProfiles) *SceneCreateGame {
	s := &SceneCreateGame{}
	s.profile = profile
	s.topBar = eui.NewTopBar("Создать профиль игры нНазад", nil)
	s.topBar.SetTitleCoverArea(0.6)
	s.Add(s.topBar)

	s.inpName = eui.NewInputBox("Profile name", 20, func(ib *eui.InputBox) {
		fmt.Println("input name", ib.GetText())
	})
	s.Add(s.inpName)

	levels := func() (arr []interface{}) {
		for i := 1; i <= 10; i++ {
			arr = append(arr, i)
		}
		return arr
	}()
	s.defLevel = levels[0].(int)
	s.cDefLev = eui.NewComboBox("Уровень по умолчанию", levels, 0, func(cb *eui.ComboBox) {
		s.defLevel = cb.Value().(int)
	})
	s.Add(s.cDefLev)

	moveTime := func() (arr []interface{}) {
		for i := 1.0; i <= 5; i += 0.5 {
			arr = append(arr, i)
		}
		return arr
	}()
	s.moveTime = moveTime[3].(float64)
	s.cMoveTime = eui.NewComboBox("Время хода", moveTime, 3, func(cb *eui.ComboBox) {
		s.moveTime = cb.Value().(float64)
	})
	s.Add(s.cMoveTime)

	showCellTime := func() (arr []interface{}) {
		for i := 0.5; i <= 0.9; i += 0.15 {
			arr = append(arr, i)
		}
		return arr
	}()
	s.showCellTime = showCellTime[1].(float64)
	s.cShowCellTm = eui.NewComboBox("Процент показа хода от времени хода", showCellTime, 1, func(cb *eui.ComboBox) {
		s.showCellTime = cb.Value().(float64)
	})
	s.Add(s.cShowCellTm)

	s.lblSelectModal = eui.NewText("Выбор и настройка модальностей")
	s.Add(s.lblSelectModal)
	s.optPos = eui.NewCheckbox(strMod+"Позиции", func(c *eui.Checkbox) {
		if c.IsChecked() {
			s.selPos = true
		} else {
			s.selPos = false
		}
	})
	s.Add(s.optPos)
	s.optCol = eui.NewCheckbox(strMod+"Цвета", func(c *eui.Checkbox) {
		if c.IsChecked() {
			s.selCol = true
		} else {
			s.selCol = false
		}
	})
	s.Add(s.optCol)
	s.cSym = eui.NewComboBox(symTitle[1], symData, 1, func(cb *eui.ComboBox) {
		if cb.Value() == symData[1] {
			s.selSym = true
			s.selAri = false
			cb.SetText(symTitle[1])
		} else if cb.Value() == symData[2] {
			s.selAri = true
			s.selSym = false
			cb.SetText(symTitle[2])
		} else if cb.Value() == symData[0] {
			s.selSym = false
			s.selAri = false
			cb.SetText(symTitle[0])
		}
	})
	s.Add(s.cSym)

	gridSz := func() (arr []interface{}) {
		for i := 2; i <= 9; i += 1 {
			arr = append(arr, i)
		}
		return arr
	}()
	s.gridSz = gridSz[1].(int)
	s.cGridSz = eui.NewComboBox("Размер сетки", gridSz, 1, func(cb *eui.ComboBox) {
		s.gridSz = cb.Value().(int)
	})
	s.Add(s.cGridSz)

	s.optGrid = eui.NewCheckbox("Показать сетку", func(c *eui.Checkbox) {
		s.showGrid = c.IsChecked()
	})
	s.Add(s.optGrid)

	s.optUseCenter = eui.NewCheckbox("Использовать центральную ячейку", func(c *eui.Checkbox) {
		s.useCenterCell = c.IsChecked()
	})
	s.Add(s.optUseCenter)

	s.optCrossHair = eui.NewCheckbox("Показать прицел", func(c *eui.Checkbox) {
		s.showCrossHair = c.IsChecked()
	})
	s.Add(s.optCrossHair)

	s.optReset = eui.NewCheckbox("До первой ошибки", func(c *eui.Checkbox) {
		s.resetOnWrong = c.IsChecked()
	})
	s.Add(s.optReset)

	s.optAddSub = eui.NewCheckbox("Арифметика сложение и вычитание", func(c *eui.Checkbox) {
		s.useAddSub = c.IsChecked()
	})
	s.Add(s.optAddSub)

	s.optMulDiv = eui.NewCheckbox("Арифметика умножение и деление", func(c *eui.Checkbox) {
		s.useMulDiv = c.IsChecked()
	})
	s.Add(s.optMulDiv)

	s.optShowGameLbl = eui.NewCheckbox("Показ метки и отклик метки, кнопок в игре", func(c *eui.Checkbox) {
		s.showGameLabel = c.IsChecked()
	})
	s.Add(s.optShowGameLbl)

	s.lblSelectMoves = eui.NewText("Выбор и настройка ходов")
	s.Add(s.lblSelectMoves)

	dataMoves := []interface{}{0, 1, 2, 3}
	s.cMoves = eui.NewComboBox("Ходов "+dtTitleMoves[0], dataMoves, 0, func(cb *eui.ComboBox) {
		s.movesConfIndex = cb.Value().(int)
		str := fmt.Sprintf("Ходов %v на 7м уровне %v", dtTitleMoves[s.movesConfIndex], s.totalMoves(7))
		s.cMoves.SetText(str)
	})
	s.Add(s.cMoves)

	s.lblSelectThreshold = eui.NewText("Выбор и настройка порогов перехода")
	s.Add(s.lblSelectThreshold)

	thresholds := func() (arr []interface{}) {
		for i := 0; i <= 100; i += 5 {
			arr = append(arr, i)
		}
		return arr
	}()
	s.thresholdUp = thresholds[18].(int)
	s.cThUp = eui.NewComboBox("Процент перехода вверх", thresholds, 18, func(cb *eui.ComboBox) {
		s.thresholdUp = cb.Value().(int)
	})
	s.Add(s.cThUp)

	s.thresholdDown = thresholds[15].(int)
	s.cThDown = eui.NewComboBox("Процент перехода вниз", thresholds, 15, func(cb *eui.ComboBox) {
		s.thresholdDown = cb.Value().(int)
	})
	s.Add(s.cThDown)

	lives := func() (arr []interface{}) {
		for i := 1; i <= 10; i++ {
			arr = append(arr, i)
		}
		return arr
	}()

	s.thresholdAdv = lives[0].(int)
	s.cThAdv = eui.NewComboBox("Сколько игр для перехода вверх", lives, 0, func(cb *eui.ComboBox) {
		s.thresholdAdv = cb.Value().(int)
	})
	s.Add(s.cThAdv)

	s.thresholdFall = lives[0].(int)
	s.cThFall = eui.NewComboBox("Сколько доп. попыток до перехода вниз", lives, 0, func(cb *eui.ComboBox) {
		s.thresholdFall = cb.Value().(int)
	})
	s.Add(s.cThFall)

	rrs := func() (arr []interface{}) {
		for i := 10; i <= 50; i += 5 {
			arr = append(arr, i)
		}
		return arr
	}()

	s.rr = rrs[1].(int)
	s.cRR = eui.NewComboBox("Процент обязательных повторов", rrs, 1, func(cb *eui.ComboBox) {
		s.rr = cb.Value().(int)
	})
	s.Add(s.cRR)

	maxNums := []interface{}{10, 20, 50, 100, 200, 500, 1000}

	s.maxNum = maxNums[0].(int)
	s.cMaxNum = eui.NewComboBox("Максимальное число арифметика, цифры", maxNums, 0, func(cb *eui.ComboBox) {
		s.maxNum = cb.Value().(int)
	})
	s.Add(s.cMaxNum)

	s.btnApply = eui.NewButton(bNew, s.checkOptions)
	s.Add(s.btnApply)
	s.btnReset = eui.NewButton(bReset, s.checkOptions)
	s.Add(s.btnReset)
	s.btnTest = eui.NewButton(bTest, s.checkOptions)
	s.Add(s.btnTest)

	s.iColors = eui.NewIcon(nil)
	s.Add(s.iColors)
	s.examples = NewExamplesFrame(s.exLogic)
	s.Add(s.examples)
	return s
}

func (s *SceneCreateGame) exLogic(b *eui.Button) {
	switch b.GetText() {
	case btnJ:
		gc := game.NewGameConf()
		gc.Set(game.Modals, game.Pos+game.Col) // по умолчанию модальность цифры
		gc.Set(game.DefaultLevel, 1)
		gc.Set(game.MoveTime, 1.5)
		gc.Set(game.ShowCellPercent, 0.65)
		gc.Set(game.RandomRepition, 30)
		gc.Set(game.GridSize, 3)
		gc.Set(game.ShowGrid, false)
		gc.Set(game.UseCenterCell, false)
		gc.Set(game.ShowCrossHair, true)
		gc.Set(game.ResetOnFirstWrong, false)
		gc.Set(game.ThresholdAdvance, 90)
		gc.Set(game.ThresholdFallback, 75)
		gc.Set(game.ThresholdAdvanceSessions, 1)
		gc.Set(game.ThresholdFallbackSessions, 1)
		gc.Set(game.Trials, 20)
		gc.Set(game.TrialsFactor, 1)
		gc.Set(game.TrialsExponent, 1)
		gc.Set(game.MaxNumber, 10)
		gc.Set(game.UseAddSub, true)
		gc.Set(game.UseMulDiv, false)
		gc.Set(game.ShowGameLabel, true)
		s.resetOpt(&gc)

	case btnB:
		gc := game.NewGameConf()
		gc.Set(game.Modals, game.Pos+game.Sym) // по умолчанию модальность цифры
		gc.Set(game.DefaultLevel, 1)
		gc.Set(game.MoveTime, 3.0)
		gc.Set(game.ShowCellPercent, 0.5)
		gc.Set(game.RandomRepition, 30)
		gc.Set(game.GridSize, 3)
		gc.Set(game.ShowGrid, true)
		gc.Set(game.UseCenterCell, false)
		gc.Set(game.ShowCrossHair, true)
		gc.Set(game.ResetOnFirstWrong, false)
		gc.Set(game.ThresholdAdvance, 80)
		gc.Set(game.ThresholdFallback, 50)
		gc.Set(game.ThresholdAdvanceSessions, 1)
		gc.Set(game.ThresholdFallbackSessions, 3)
		gc.Set(game.Trials, 20)
		gc.Set(game.TrialsFactor, 1)
		gc.Set(game.TrialsExponent, 2)
		gc.Set(game.MaxNumber, 10)
		gc.Set(game.UseAddSub, true)
		gc.Set(game.UseMulDiv, false)
		gc.Set(game.ShowGameLabel, true)
		s.resetOpt(&gc)

	case bntQ:
		gc := game.NewGameConf()
		gc.Set(game.Modals, game.Sym+game.Col) // по умолчанию модальность цифры
		gc.Set(game.DefaultLevel, 1)
		gc.Set(game.MoveTime, 2.5)
		gc.Set(game.ShowCellPercent, 0.65)
		gc.Set(game.RandomRepition, 30)
		gc.Set(game.GridSize, 3)
		gc.Set(game.ShowGrid, true)
		gc.Set(game.UseCenterCell, false)
		gc.Set(game.ShowCrossHair, true)
		gc.Set(game.ResetOnFirstWrong, false)
		gc.Set(game.ThresholdAdvance, 90)
		gc.Set(game.ThresholdFallback, 0)
		gc.Set(game.ThresholdAdvanceSessions, 1)
		gc.Set(game.ThresholdFallbackSessions, 1)
		gc.Set(game.Trials, 20)
		gc.Set(game.TrialsFactor, 5)
		gc.Set(game.TrialsExponent, 1)
		gc.Set(game.MaxNumber, 10)
		gc.Set(game.UseAddSub, true)
		gc.Set(game.UseMulDiv, false)
		gc.Set(game.ShowGameLabel, true)
		s.resetOpt(&gc)

	case btnP:
		gc := game.NewGameConf()
		gc.Set(game.Modals, game.Ari) // по умолчанию модальность цифры
		gc.Set(game.DefaultLevel, 1)
		gc.Set(game.MoveTime, 2.0)
		gc.Set(game.ShowCellPercent, 0.8)
		gc.Set(game.RandomRepition, 30)
		gc.Set(game.GridSize, 3)
		gc.Set(game.ShowGrid, false)
		gc.Set(game.UseCenterCell, false)
		gc.Set(game.ShowCrossHair, true)
		gc.Set(game.ResetOnFirstWrong, false)
		gc.Set(game.ThresholdAdvance, 90)
		gc.Set(game.ThresholdFallback, 0)
		gc.Set(game.ThresholdAdvanceSessions, 3)
		gc.Set(game.ThresholdFallbackSessions, 1)
		gc.Set(game.Trials, 10)
		gc.Set(game.TrialsFactor, 1)
		gc.Set(game.TrialsExponent, 1)
		gc.Set(game.MaxNumber, 10)
		gc.Set(game.UseAddSub, true)
		gc.Set(game.UseMulDiv, false)
		gc.Set(game.ShowGameLabel, true)
		s.resetOpt(&gc)
	}
}

func (s *SceneCreateGame) Entered() {
	s.Resize()
	theme := eui.GetUi().GetTheme()
	bg := theme.Get(app.GameColorBg)
	fg := theme.Get(app.GameColorFg)
	s.examples.setup(bg, fg)
	icon := NewColorsBar(bg, fg)
	icon.Resize(s.iColors.GetRect().GetArr())
	icon.Setup()
	s.iColors.SetIcon(icon.Image())
	if !s.inTesting {
		s.resetOpt(nil)
	}
	s.inTesting = false
}

func (s *SceneCreateGame) setModals(modals game.ModalType) {
	s.selPos = false
	s.selCol = false
	s.selSym = false
	s.selAri = false
	s.cSym.SetValue(" ")
	s.cSym.SetText(symTitle[0])
	s.optPos.SetChecked(false)
	s.optCol.SetChecked(false)
	for _, mod := range modals {
		switch game.ModalType(mod) {
		case game.Pos:
			s.selPos = true
			s.optPos.SetChecked(s.selPos)
		case game.Col:
			s.selCol = true
			s.optCol.SetChecked(s.selCol)
		case game.Sym:
			s.selSym = true
			s.cSym.SetValue(game.Sym.String())
			s.cSym.SetText(symTitle[1])
		case game.Ari:
			s.selAri = true
			s.cSym.SetValue(game.Ari.String())
			s.cSym.SetText(symTitle[2])
		}
	}
}

func (s *SceneCreateGame) setMoves(trials, factor, exp int) {
	// "Новичёк", "Начинающий", "Профессионал", "Мастер"
	if trials == 10 && factor == 1 && exp == 1 {
		s.cMoves.SetValue(0)
		s.movesConfIndex = 0
	}
	if trials == 20 && factor == 1 && exp == 1 {
		s.cMoves.SetValue(1)
		s.movesConfIndex = 1
	}
	if trials == 20 && factor == 5 && exp == 1 {
		s.cMoves.SetValue(2)
		s.movesConfIndex = 2
	}
	if trials == 20 && exp == 2 {
		s.cMoves.SetValue(3)
		s.movesConfIndex = 3
	}
	str := fmt.Sprintf("Ходов %v на 4м уровне %v", dtTitleMoves[s.movesConfIndex], s.totalMoves(4))
	s.cMoves.SetText(str)
}

func (s *SceneCreateGame) resetOpt(conf *game.GameConf) {
	if conf == nil {
		conf = game.DefaultSettings()
	}
	mod := conf.Get(game.Modals).(game.ModalType)
	s.setModals(mod)
	s.defLevel = conf.Get(game.DefaultLevel).(int)
	s.cDefLev.SetValue(s.defLevel)
	s.moveTime = conf.Get(game.MoveTime).(float64)
	s.cMoveTime.SetValue(s.moveTime)
	s.showCellTime = conf.Get(game.ShowCellPercent).(float64)
	s.cShowCellTm.SetValue(s.showCellTime)
	s.rr = conf.Get(game.RandomRepition).(int)
	s.cRR.SetValue(s.rr)
	s.gridSz = conf.Get(game.GridSize).(int)
	s.cGridSz.SetValue(s.gridSz)
	s.showGrid = conf.Get(game.ShowGrid).(bool)
	s.optGrid.SetChecked(s.showGrid)
	s.showCrossHair = conf.Get(game.ShowCrossHair).(bool)
	s.optCrossHair.SetChecked(s.showCrossHair)
	s.useCenterCell = conf.Get(game.UseCenterCell).(bool)
	s.optUseCenter.SetChecked(s.useCenterCell)
	s.thresholdUp = conf.Get(game.ThresholdAdvance).(int)
	s.cThUp.SetValue(s.thresholdUp)
	s.thresholdDown = conf.Get(game.ThresholdFallback).(int)
	s.cThDown.SetValue(s.thresholdDown)
	s.thresholdAdv = conf.Get(game.ThresholdAdvanceSessions).(int)
	s.cThAdv.SetValue(s.thresholdAdv)
	s.thresholdFall = conf.Get(game.ThresholdFallbackSessions).(int)
	s.cThFall.SetValue(s.thresholdFall)
	s.resetOnWrong = conf.Get(game.ResetOnFirstWrong).(bool)
	s.optReset.SetChecked(s.resetOnWrong)
	s.useAddSub = conf.Get(game.UseAddSub).(bool)
	s.optAddSub.SetChecked(s.useAddSub)
	s.useMulDiv = conf.Get(game.UseMulDiv).(bool)
	s.optMulDiv.SetChecked(s.useMulDiv)
	s.maxNum = conf.Get(game.MaxNumber).(int)
	s.cMaxNum.SetValue(s.maxNum)
	trials := conf.Get(game.Trials).(int)
	factor := conf.Get(game.TrialsFactor).(int)
	exp := conf.Get(game.TrialsExponent).(int)
	s.setMoves(trials, factor, exp)
	s.showGameLabel = conf.Get(game.ShowGameLabel).(bool)
	s.optShowGameLbl.SetChecked(s.showGameLabel)
	s.inpName.SetText(s.genName())
}

func (s *SceneCreateGame) checkOptions(b *eui.Button) {
	profileName := s.genName()
	switch b.GetText() {
	case bNew:
		s.inpName.SetText(profileName)
		s.profile.AddGameProfile(profileName, s.LoadConf())
		eui.GetUi().Pop()
	case bReset:
		s.resetOpt(nil)
	case bTest:
		s.inTesting = true
		s.inpName.SetText(profileName)
		s.profile.AddGameProfile(profileName, s.LoadConf())
		sc := scene_game.New()
		GamesData := s.profile.GetGamesData(profileName)
		sc.Setup(*GamesData.Conf, GamesData.Last())
		eui.GetUi().Push(sc)
		delete(*s.profile, profileName)
	}
}

func (s *SceneCreateGame) LoadConf() *game.GameConf {
	gc := game.NewGameConf()
	_, m := s.getModals()
	gc.Set(game.Modals, m)
	gc.Set(game.DefaultLevel, s.defLevel)
	gc.Set(game.MoveTime, s.moveTime)
	gc.Set(game.ShowCellPercent, s.showCellTime)
	gc.Set(game.RandomRepition, s.rr)
	gc.Set(game.GridSize, s.gridSz)
	gc.Set(game.ShowGrid, s.showGrid)
	gc.Set(game.UseCenterCell, s.useCenterCell)
	gc.Set(game.ShowCrossHair, s.showCrossHair)
	gc.Set(game.ResetOnFirstWrong, s.resetOnWrong)
	gc.Set(game.ThresholdAdvance, s.thresholdUp)
	gc.Set(game.ThresholdFallback, s.thresholdDown)
	gc.Set(game.ThresholdAdvanceSessions, s.thresholdAdv)
	gc.Set(game.ThresholdFallbackSessions, s.thresholdFall)
	gc.Set(game.Trials, movesArr[s.movesConfIndex][0])
	gc.Set(game.TrialsFactor, movesArr[s.movesConfIndex][1])
	gc.Set(game.TrialsExponent, movesArr[s.movesConfIndex][2])
	gc.Set(game.MaxNumber, s.maxNum)
	gc.Set(game.UseAddSub, s.useAddSub)
	gc.Set(game.UseMulDiv, s.useMulDiv)
	gc.Set(game.ShowGameLabel, s.showGameLabel)
	return &gc
}

func (s *SceneCreateGame) genName() (result string) {
	s1, _ := s.getModals()
	result = fmt.Sprintf("%v Ходов %v (%v/%v) время хода(%vсек)", s1, dtTitleMoves[s.movesConfIndex], s.thresholdUp, s.thresholdDown, s.moveTime)
	if s.thresholdAdv > 1 {
		result += fmt.Sprintf(" попыток вверх(%v)", s.thresholdAdv)
	}
	if s.thresholdDown > 0 {
		result += fmt.Sprintf(" доп.попыток(%v)", s.thresholdFall)
	}
	if s.resetOnWrong {
		result += " до первой ошибки"
	}
	for _, v := range s.profile.GetProfilesName() {
		if v == result {
			result += "_ещё"
		}
	}
	return result
}

func (s *SceneCreateGame) getModals() (string, game.ModalType) {
	s1 := ""
	modals := 0
	var mt game.ModalType
	if s.selPos {
		s1 += fmt.Sprintf("Позиции[%v(%vx%v)]", game.Pos, s.gridSz, s.gridSz)
		modals++
		mt += game.Pos
	}
	if s.selSym {
		s1 += fmt.Sprintf("Цифры[%v]", game.Sym)
		modals++
		mt += game.Sym
	}
	if s.selAri {
		s1 += fmt.Sprintf("Арифметика[%v]", game.Ari)
		modals++
		mt += game.Ari
	}
	if s.selCol {
		s1 += fmt.Sprintf("Цвета[%v]", game.Col)
		modals++
		mt += game.Col
	}
	if modals == 0 {
		s1 += "Выбрать модальность"
	}
	s2 := ""
	switch modals {
	case 1:
		s2 += "Single"
	case 2:
		s2 += "Dual"
	case 3:
		s2 += "Triple"
	case 4:
		s2 += "Quad"
	}
	return fmt.Sprintf("%v(%v)", s2, s1), mt
}

func (s *SceneCreateGame) Resize() {
	w0, h0 := eui.GetUi().Size()
	rect := eui.NewRect([]int{0, 0, w0, h0})
	hTop := int(float64(rect.GetLowestSize()) * 0.05)
	margin := int(float64(rect.GetLowestSize()) * 0.008)
	x, y := 0, 0
	s.topBar.Resize([]int{x, y, w0, hTop})
	h := (rect.H - hTop) / 17 // на сколько строк деление

	w1 := w0 - w0/5
	w2 := w1 / 2 // в 2 столбика
	w3 := w1 / 3 // в 3 столбика
	w4 := w1 / 4 // в 4 столбика

	y += hTop + margin
	s.inpName.Resize([]int{x, y, w1 - margin, h - margin})
	s.examples.Resize([]int{x + w1, y, w0/5 - margin, hTop*5 - margin})

	y += h
	s.cDefLev.Resize([]int{x, y, w2 - margin, h - margin})
	s.cRR.Resize([]int{x + w2, y, w2 - margin, h - margin})

	y += h
	s.cMoveTime.Resize([]int{x, y, w2 - margin, h - margin})
	s.cShowCellTm.Resize([]int{x + w2, y, w2 - margin, h - margin})

	y += h
	s.optReset.Resize([]int{x, y, w3 - margin, h - margin})
	s.optCrossHair.Resize([]int{x + w3, y, w3 - margin, h - margin})
	s.optShowGameLbl.Resize([]int{x + w3*2, y, w3 - margin, h - margin})

	y += h
	s.lblSelectModal.Resize([]int{x, y, w1 - margin, h - margin})
	y += h
	s.optPos.Resize([]int{x, y, w1 - margin, h - margin})
	y += h
	s.optGrid.Resize([]int{x, y, w3 - margin, h - margin})
	s.cGridSz.Resize([]int{x + w3, y, w3 - margin, h - margin})
	s.optUseCenter.Resize([]int{x + w3*2, y, w3 - margin, h - margin})

	y += h
	s.optCol.Resize([]int{x, y, w1 - margin, h - margin})
	y += h
	s.iColors.Resize([]int{x, y, w1 - margin, h - margin})
	y += h
	s.cSym.Resize([]int{x, y, w1 - margin, h - margin})

	y += h
	s.cMaxNum.Resize([]int{x, y, w2 - margin, h - margin})
	s.optAddSub.Resize([]int{x + w2, y, w4 - margin, h - margin})
	s.optMulDiv.Resize([]int{x + w2 + w4, y, w4 - margin, h - margin})

	y += h
	s.lblSelectMoves.Resize([]int{x, y, w1 - margin, h - margin})
	y += h
	s.cMoves.Resize([]int{x, y, w1 - margin, h - margin})
	y += h
	s.lblSelectThreshold.Resize([]int{x, y, w1 - margin, h - margin})
	y += h
	s.cThUp.Resize([]int{x, y, w2 - margin, h - margin})
	s.cThDown.Resize([]int{x + w2, y, w2 - margin, h - margin})
	y += h

	s.cThAdv.Resize([]int{x, y, w2 - margin, h - margin})
	s.cThFall.Resize([]int{x + w2, y, w2 - margin, h - margin})

	y += h
	s.btnApply.Resize([]int{x, y, w3 - margin, h - margin})
	s.btnReset.Resize([]int{x + w3, y, w3 - margin, h - margin})
	s.btnTest.Resize([]int{x + w3*2, y, w3 - margin, h - margin})
}

func (s *SceneCreateGame) totalMoves(level int) int {
	trials := movesArr[s.movesConfIndex][0]
	factor := movesArr[s.movesConfIndex][1]
	exponent := movesArr[s.movesConfIndex][2]
	return trials + factor*int(math.Pow(float64(level), float64(exponent)))
}
