package result

import "github.com/t0l1k/eui"

type SceneResults struct {
	eui.SceneBase
	topbar *eui.TopBar
}

func NewSceneResults() *SceneResults {
	s := &SceneResults{}
	s.topbar = eui.NewTopBar("Итоги в игре нназад", nil)
	s.Add(s.topbar)
	return s
}

func (s *SceneResults) Entered() {
	s.Resize()
}

func (s *SceneResults) Resize() {
	w0, h0 := eui.GetUi().Size()
	// w1 := int(float64(w0) * 0.68)
	h1 := int(float64(h0) * 0.068)
	s.topbar.Resize([]int{0, 0, w0, h1})
}
