package create

import (
	"fmt"

	"github.com/t0l1k/eui"
	"github.com/t0l1k/nBack/app/data"
	"github.com/t0l1k/nBack/app/game"
)

var (
	dtTitleMoves = []string{"Новичёк", "Начинающий", "Профессионал", "Мастер"}
	movesArr     = [][]int{{10, 1, 1}, {20, 1, 1}, {20, 5, 1}, {20, 1, 2}}
)

type SceneCreateGame struct {
	eui.SceneBase
	topBar                                                  *eui.TopBar
	profile                                                 *data.GameProfiles
	inpName                                                 *eui.InputBox
	lblSelectModal, lblSelectMoves, lblSelectThreshold      *eui.Text
	optPos, optCol, optAddSub, optMulDiv                    *eui.Checkbox
	optCrossHair, optGrid, optUseCenter, optReset           *eui.Checkbox
	cDefLev, cSym, cGridSz, cMoves, cMoveTime, cShowCellTm  *eui.ComboBox
	cThUp, cThDown, cThAdv, cThFall, cRR, cMaxNum           *eui.ComboBox
	btnApply, btnReset                                      *eui.Button
	selPos, selCol, selNum, selAri                          bool
	defLevel, movesConfIndex, rr, gridSz, maxNum            int
	moveTime, showCellTime                                  float64
	thresholdUp, thresholdDown, thresholdAdv, thresholdFall int
	showGrid, showCrossHair, useCenterCell, resetOnWrong    bool
	useMulDiv, useAddSub                                    bool
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
	s.optPos = eui.NewCheckbox("Позиции", func(c *eui.Checkbox) {
		if c.IsChecked() {
			s.selPos = true
		} else {
			s.selPos = false
		}
	})
	s.Add(s.optPos)
	s.optCol = eui.NewCheckbox("Цвета", func(c *eui.Checkbox) {
		if c.IsChecked() {
			s.selCol = true
		} else {
			s.selCol = false
		}
	})
	s.Add(s.optCol)
	var dtTitle = []string{"Не использовать", "Цифры", "Арифметика"}
	var dataSym = []interface{}{"", game.Sym, game.Ari}
	s.cSym = eui.NewComboBox(dtTitle[1], dataSym, 1, func(cb *eui.ComboBox) {
		if cb.Value() == dataSym[1] {
			s.selNum = true
			s.selAri = false
			cb.SetText(dtTitle[1])
			fmt.Println("nums on", cb.Value())
		} else if cb.Value() == dataSym[2] {
			s.selAri = true
			s.selNum = false
			cb.SetText(dtTitle[2])
			fmt.Println("ari on", cb.Value())
		} else {
			s.selNum = false
			s.selAri = false
			cb.SetText(dtTitle[0])
			fmt.Println("reset", cb.Value())
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

	s.lblSelectMoves = eui.NewText("Выбор и настройка ходов")
	s.Add(s.lblSelectMoves)

	dataMoves := []interface{}{0, 1, 2, 3}
	s.cMoves = eui.NewComboBox(dtTitleMoves[0], dataMoves, 0, func(cb *eui.ComboBox) {
		s.movesConfIndex = cb.Value().(int)
		s.cMoves.SetText(dtTitleMoves[s.movesConfIndex])
		fmt.Println("ходов:", dtTitleMoves[s.movesConfIndex])
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

	s.btnApply = eui.NewButton("Применить", s.checkOptions)
	s.Add(s.btnApply)
	s.btnReset = eui.NewButton("Обнулить", s.checkOptions)
	s.Add(s.btnReset)
	return s
}

func (s *SceneCreateGame) Entered() {
	s.Resize()
	s.resetOpt()
}

func (s *SceneCreateGame) resetOpt() {
	conf := game.DefaultSettings()
	s.selNum = true
	s.cSym.SetValue(game.Sym)
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
	s.cMoves.SetValue(0)
	s.resetOnWrong = conf.Get(game.ResetOnFirstWrong).(bool)
	s.optReset.SetChecked(s.resetOnWrong)
	s.useAddSub = conf.Get(game.UseAddSub).(bool)
	s.optAddSub.SetChecked(s.useAddSub)
	s.useMulDiv = conf.Get(game.UseMulDiv).(bool)
	s.optMulDiv.SetChecked(s.useMulDiv)
	s.maxNum = conf.Get(game.MaxNumber).(int)
	s.cMaxNum.SetValue(s.maxNum)
	s.inpName.SetText(s.genName())
}

func (s *SceneCreateGame) checkOptions(b *eui.Button) {
	if b.GetText() == "Обнулить" {
		s.resetOpt()
		return
	}
	profileName := s.genName()
	s.inpName.SetText(profileName)
	s.profile.AddGameProfile(profileName, s.LoadConf())
	fmt.Println(s.genName())
	eui.GetUi().Pop()
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
	return &gc
}

func (s *SceneCreateGame) genName() (result string) {
	s1, _ := s.getModals()
	result = fmt.Sprintf("%v %v (%v/%v) ход(%vсек)", s1, dtTitleMoves[s.movesConfIndex], s.thresholdUp, s.thresholdDown, s.moveTime)
	if s.thresholdAdv > 1 {
		result += fmt.Sprintf(" попыток вверх(%v) доп.попыток(%v)", s.thresholdAdv, s.thresholdFall)
	} else {
		result += fmt.Sprintf(" доп.попыток(%v)", s.thresholdFall)
	}
	if s.resetOnWrong {
		result += " до первой ошибки"
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
	if s.selNum {
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
	margin := int(float64(rect.GetLowestSize()) * 0.003)
	x, y := 0, 0
	s.topBar.Resize([]int{x, y, w0, hTop})
	y += hTop
	h := rect.H / 30
	w1 := w0 - hTop
	x += (w0 - w1) / 2
	y += hTop / 2
	s.inpName.Resize([]int{x, y, w1, h})
	y += h + margin
	s.cDefLev.Resize([]int{x, y, w1, h})
	y += h + margin
	s.cMoveTime.Resize([]int{x, y, w1, h})
	y += h + margin
	s.cShowCellTm.Resize([]int{x, y, w1, h})
	y += h + margin
	s.cRR.Resize([]int{x, y, w1, h})

	y += h + margin
	s.lblSelectModal.Resize([]int{x, y, w1, h})
	y += h + margin
	s.optPos.Resize([]int{x, y, w1, h})
	y += h + margin
	s.optCol.Resize([]int{x, y, w1, h})
	y += h + margin
	s.cSym.Resize([]int{x, y, w1, h})
	y += h + margin
	s.optGrid.Resize([]int{x, y, w1, h})

	y += h + margin
	s.cGridSz.Resize([]int{x, y, w1, h})
	y += h + margin
	s.optUseCenter.Resize([]int{x, y, w1, h})
	y += h + margin
	s.optCrossHair.Resize([]int{x, y, w1, h})
	y += h + margin
	s.optAddSub.Resize([]int{x, y, w1, h})
	y += h + margin
	s.optMulDiv.Resize([]int{x, y, w1, h})
	y += h + margin
	s.cMaxNum.Resize([]int{x, y, w1, h})

	y += h + margin
	s.lblSelectMoves.Resize([]int{x, y, w1, h})
	y += h + margin
	s.cMoves.Resize([]int{x, y, w1, h})
	y += h + margin
	s.lblSelectThreshold.Resize([]int{x, y, w1, h})
	y += h + margin
	s.cThUp.Resize([]int{x, y, w1, h})
	y += h + margin
	s.cThDown.Resize([]int{x, y, w1, h})
	y += h + margin
	s.cThAdv.Resize([]int{x, y, w1, h})
	y += h + margin
	s.cThFall.Resize([]int{x, y, w1, h})
	y += h + margin
	s.optReset.Resize([]int{x, y, w1, h})
	y += h + margin
	w2 := w1 / 2
	y = rect.H - hTop - hTop/2
	s.btnApply.Resize([]int{x, y, w2, h})
	s.btnReset.Resize([]int{x + w2, y, w2, h})
}
