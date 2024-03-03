package options

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/t0l1k/eui"
)

type InputKey struct {
	eui.DrawableBase
	lbl    *eui.Text
	btn    *eui.Button
	active bool
	value  ebiten.Key
}

func NewInputKey(title string) *InputKey {
	i := &InputKey{}
	i.lbl = eui.NewText(title)
	i.Add(i.lbl)
	i.btn = eui.NewButton("(?)", func(b *eui.Button) {
		if b.IsPressed() {
			i.active = true
			i.btn.Bg(eui.Yellow)
		}
	})
	i.Add(i.btn)
	eui.GetUi().GetInputKeyboard().Attach(i)
	return i
}

func (i *InputKey) Value() ebiten.Key {
	return i.value
}

func (i *InputKey) SetValue(value ebiten.Key) {
	i.value = value
	i.btn.SetText(i.value.String())
}

func (i *InputKey) UpdateInput(value interface{}) {
	switch v := value.(type) {
	case eui.KeyboardData:
		if i.active {
			i.btn.SetText(v.GetKeys()[0].String())
			i.value = v.GetKeys()[0]
		}
	}
}

func (i *InputKey) Update(dt int) {
	i.DrawableBase.Update(dt)
	if i.btn.GetState() == eui.ViewStateNormal {
		i.active = false
		i.btn.Bg(eui.Silver)
	}
}

func (i *InputKey) Resize(rect []int) {
	i.Rect(eui.NewRect(rect))
	w0, h0 := i.GetRect().Size()
	x0, y0 := i.GetRect().Pos()
	i.btn.Resize([]int{x0, y0, h0 * 2, h0})
	i.lbl.Resize([]int{x0 + h0*2, y0, w0 - h0*2, h0})
	i.ImageReset()
}

func (i *InputKey) String() string {
	return fmt.Sprintf("%v: %v", i.lbl.GetText(), i.value)
}

func (i *InputKey) Close() { eui.GetUi().GetInputKeyboard().Detach(i) }
