package create

import (
	"fmt"

	"github.com/t0l1k/eui"
	"github.com/t0l1k/nBack/app/data"
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
	var dataSym = []interface{}{"", data.Sym, data.Ari}
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
	conf := data.DefaultSettings()
	s.selNum = true
	s.cSym.SetValue(data.Sym)
	s.defLevel = conf.Get(data.DefaultLevel).(int)
	s.cDefLev.SetValue(s.defLevel)
	s.moveTime = conf.Get(data.MoveTime).(float64)
	s.cMoveTime.SetValue(s.moveTime)
	s.showCellTime = conf.Get(data.ShowCellPercent).(float64)
	s.cShowCellTm.SetValue(s.showCellTime)
	s.rr = conf.Get(data.RandomRepition).(int)
	s.cRR.SetValue(s.rr)
	s.gridSz = conf.Get(data.GridSize).(int)
	s.cGridSz.SetValue(s.gridSz)
	s.showGrid = conf.Get(data.ShowGrid).(bool)
	s.optGrid.SetChecked(s.showGrid)
	s.showCrossHair = conf.Get(data.ShowCrossHair).(bool)
	s.optCrossHair.SetChecked(s.showCrossHair)
	s.useCenterCell = conf.Get(data.UseCenterCell).(bool)
	s.optUseCenter.SetChecked(s.useCenterCell)
	s.thresholdUp = conf.Get(data.ThresholdAdvance).(int)
	s.cThUp.SetValue(s.thresholdUp)
	s.thresholdDown = conf.Get(data.ThresholdFallback).(int)
	s.cThDown.SetValue(s.thresholdDown)
	s.thresholdAdv = conf.Get(data.ThresholdAdvanceSessions).(int)
	s.cThAdv.SetValue(s.thresholdAdv)
	s.thresholdFall = conf.Get(data.ThresholdFallbackSessions).(int)
	s.cThFall.SetValue(s.thresholdFall)
	s.cMoves.SetValue(0)
	s.resetOnWrong = conf.Get(data.ResetOnFirstWrong).(bool)
	s.optReset.SetChecked(s.resetOnWrong)
	s.useAddSub = conf.Get(data.UseAddSub).(bool)
	s.optAddSub.SetChecked(s.useAddSub)
	s.useMulDiv = conf.Get(data.UseMulDiv).(bool)
	s.optMulDiv.SetChecked(s.useMulDiv)
	s.maxNum = conf.Get(data.MaxNumber).(int)
	s.cMaxNum.SetValue(s.maxNum)

	s.inpName.SetText(s.genName())
	// trials := conf.Get(data.Trials).(int)
	// fact := conf.Get(data.TrialsFactor).(int)
	// exp := conf.Get(data.TrialsExponent).(int)
	// m := trials + fact*int(math.Pow(float64(s.defLevel), float64(exp)))
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

func (s *SceneCreateGame) LoadConf() *data.GameConf {
	gc := data.NewGameConf()
	_, m := s.getModals()
	gc.Set(data.Modals, m)
	gc.Set(data.DefaultLevel, s.defLevel)
	gc.Set(data.MoveTime, s.moveTime)
	gc.Set(data.ShowCellPercent, s.showCellTime)
	gc.Set(data.RandomRepition, s.rr)
	gc.Set(data.GridSize, s.gridSz)
	gc.Set(data.ShowGrid, s.showGrid)
	gc.Set(data.UseCenterCell, s.useCenterCell)
	gc.Set(data.ShowCrossHair, s.showCrossHair)
	gc.Set(data.ResetOnFirstWrong, s.resetOnWrong)
	gc.Set(data.ThresholdAdvance, s.thresholdUp)
	gc.Set(data.ThresholdFallback, s.thresholdDown)
	gc.Set(data.ThresholdAdvanceSessions, s.thresholdAdv)
	gc.Set(data.ThresholdFallbackSessions, s.thresholdFall)
	gc.Set(data.Trials, movesArr[s.movesConfIndex][0])
	gc.Set(data.TrialsFactor, movesArr[s.movesConfIndex][1])
	gc.Set(data.TrialsExponent, movesArr[s.movesConfIndex][2])
	gc.Set(data.MaxNumber, s.maxNum)
	gc.Set(data.UseAddSub, s.useAddSub)
	gc.Set(data.UseMulDiv, s.useMulDiv)
	return &gc
}

func (s *SceneCreateGame) genName() string {
	s1, _ := s.getModals()
	return fmt.Sprintf("%v %v (%v/%v) ход(%vсек) попыток(%v)", s1, dtTitleMoves[s.movesConfIndex], s.thresholdUp, s.thresholdDown, s.moveTime, s.thresholdFall)
}

func (s *SceneCreateGame) getModals() (string, data.ModalType) {
	s1 := ""
	modals := 0
	var mt data.ModalType
	if s.selPos {
		s1 += fmt.Sprintf("Позиции[%v(%vx%v)]", data.Pos, s.gridSz, s.gridSz)
		modals++
		mt += data.Pos
	}
	if s.selNum {
		s1 += fmt.Sprintf("Цифры[%v]", data.Sym)
		modals++
		mt += data.Sym
	}
	if s.selAri {
		s1 += fmt.Sprintf("Арифметика[%v]", data.Ari)
		modals++
		mt += data.Ari
	}
	if s.selCol {
		s1 += fmt.Sprintf("Цвета[%v]", data.Col)
		modals++
		mt += data.Col
	}
	if modals == 0 {
		s1 += "Выбрать модальность"
	}
	s2 := ""
	switch modals {
	case 1:
		s2 += "Single"
	case 2:
		s2 += "Double"
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
