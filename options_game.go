package main

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/t0l1k/nBack/ui"
)

type OptGame struct {
	rect          *ui.Rect
	Image         *ebiten.Image
	Dirty, Visibe bool
	bg, fg        color.Color
}

func NewOptGame(rect []int) *OptGame {
	return &OptGame{
		rect:   ui.NewRect(rect),
		bg:     getApp().theme.bg,
		fg:     getApp().theme.fg,
		Dirty:  true,
		Visibe: true,
	}
}
func (r *OptGame) Layout() *ebiten.Image {
	if !r.Dirty {
		return r.Image
	}
	w, h := r.rect.Size()
	cellWidth, cellHeight := w, h/9
	image := ebiten.NewImage(w, h)
	image.Fill(r.bg)
	x, y := 0, 0
	rect := []int{x, y, cellWidth / 2, cellHeight}
	ttncLbl := ui.NewLabel(fmt.Sprintf("Time to next Cell: %v", getApp().preferences.timeToNextCell), rect, app.theme.bg, app.theme.fg)
	ttncLbl.DrawRect = true
	ttncLbl.Draw(image)
	x, y = cellWidth/2, 0
	rect = []int{x, y, cellWidth / 2, cellHeight}
	ttsc := ui.NewLabel(fmt.Sprintf("Time to show Cell: %v", getApp().preferences.timeShowCell), rect, app.theme.bg, app.theme.fg)
	ttsc.DrawRect = true
	ttsc.Draw(image)

	x, y = 0, cellHeight
	rect = []int{x, y, cellWidth / 3, cellHeight}
	trialsLbl := ui.NewLabel(fmt.Sprintf("Trials: %v", getApp().preferences.trials), rect, app.theme.bg, app.theme.fg)
	trialsLbl.DrawRect = true
	trialsLbl.Draw(image)
	x, y = cellWidth/3, cellHeight
	rect = []int{x, y, cellWidth / 3, cellHeight}
	trialsFactorLbl := ui.NewLabel(fmt.Sprintf("Factor: %v", getApp().preferences.trialsFactor), rect, app.theme.bg, app.theme.fg)
	trialsFactorLbl.DrawRect = true
	trialsFactorLbl.Draw(image)
	x, y = cellWidth/3*2, cellHeight
	rect = []int{x, y, cellWidth / 3, cellHeight}
	trialsExpLbl := ui.NewLabel(fmt.Sprintf("Exponent: %v", getApp().preferences.trialsExponent), rect, app.theme.bg, app.theme.fg)
	trialsExpLbl.DrawRect = true
	trialsExpLbl.Draw(image)

	x, y = 0, cellHeight*2
	rect = []int{x, y, cellWidth / 3, cellHeight}
	advLbl := ui.NewLabel(fmt.Sprintf("Advance: %v", getApp().preferences.thresholdAdvance), rect, app.theme.bg, app.theme.fg)
	advLbl.DrawRect = true
	advLbl.Draw(image)
	x = cellWidth / 3
	rect = []int{x, y, cellWidth / 3, cellHeight}
	fallLbl := ui.NewLabel(fmt.Sprintf("Fallback: %v", getApp().preferences.thresholdFallback), rect, app.theme.bg, app.theme.fg)
	fallLbl.DrawRect = true
	fallLbl.Draw(image)
	x = cellWidth / 3 * 2
	rect = []int{x, y, cellWidth / 3, cellHeight}
	fallSessionsLbl := ui.NewLabel(fmt.Sprintf("Fallback sessions: %v", getApp().preferences.thresholdFallbackSessions), rect, app.theme.bg, app.theme.fg)
	fallSessionsLbl.DrawRect = true
	fallSessionsLbl.Draw(image)

	x, y = 0, cellHeight*3
	rect = []int{x, y, cellWidth / 3, cellHeight}
	defLevelLbl := ui.NewLabel(fmt.Sprintf("Default level: %v", getApp().preferences.defaultLevel), rect, app.theme.bg, app.theme.fg)
	defLevelLbl.DrawRect = true
	defLevelLbl.Draw(image)
	x = cellWidth / 3
	rect = []int{x, y, cellWidth / 3, cellHeight}
	manualLbl := ui.NewLabel(fmt.Sprintf("Manual: %v", getApp().preferences.manual), rect, app.theme.bg, app.theme.fg)
	manualLbl.DrawRect = true
	manualLbl.Draw(image)
	x = cellWidth / 3 * 2
	rect = []int{x, y, cellWidth / 3, cellHeight}
	manualAdvLbl := ui.NewLabel(fmt.Sprintf("Manual Advance: %v", getApp().preferences.manualAdv), rect, app.theme.bg, app.theme.fg)
	manualAdvLbl.DrawRect = true
	manualAdvLbl.Draw(image)

	x, y = 0, cellHeight*4
	rect = []int{x, y, cellWidth, cellHeight}
	resetLbl := ui.NewLabel(fmt.Sprintf("Reset on first wrong: %v", getApp().preferences.resetOnFirstWrong), rect, app.theme.bg, app.theme.fg)
	resetLbl.DrawRect = true
	resetLbl.Draw(image)

	x, y = 0, cellHeight*5
	rect = []int{x, y, cellWidth, cellHeight}
	rrLbl := ui.NewLabel(fmt.Sprintf("Random repition: %v %%", getApp().preferences.rr), rect, app.theme.bg, app.theme.fg)
	rrLbl.DrawRect = true
	rrLbl.Draw(image)

	x, y = 0, cellHeight*6
	rect = []int{x, y, cellWidth / 3, cellHeight}
	centerLbl := ui.NewLabel(fmt.Sprintf("Use center cell: %v", getApp().preferences.usecentercell), rect, app.theme.bg, app.theme.fg)
	centerLbl.DrawRect = true
	centerLbl.Draw(image)
	x = cellWidth / 3
	rect = []int{x, y, cellWidth / 3, cellHeight}
	feedbackLbl := ui.NewLabel(fmt.Sprintf("Feedback on move: %v", getApp().preferences.feedbackOnUserMove), rect, app.theme.bg, app.theme.fg)
	feedbackLbl.DrawRect = true
	feedbackLbl.Draw(image)
	x = cellWidth / 3 * 2
	rect = []int{x, y, cellWidth / 3, cellHeight}
	gridSzLbl := ui.NewLabel(fmt.Sprintf("Grid size: %v", getApp().preferences.gridSize), rect, app.theme.bg, app.theme.fg)
	gridSzLbl.DrawRect = true
	gridSzLbl.Draw(image)

	x, y = 0, cellHeight*7
	rect = []int{x, y, cellWidth, cellHeight}
	restLbl := ui.NewLabel(fmt.Sprintf("Pause to Rest after game: %v second", getApp().preferences.pauseRest/1000), rect, app.theme.bg, app.theme.fg)
	restLbl.DrawRect = true
	restLbl.Draw(image)

	x, y = 0, cellHeight*8
	rect = []int{x, y, cellWidth, cellHeight}
	fullScrLbl := ui.NewLabel(fmt.Sprintf("Fullscreen on app start, toggle by F11: %v", getApp().preferences.fullScreen), rect, app.theme.bg, app.theme.fg)
	fullScrLbl.DrawRect = true
	fullScrLbl.Draw(image)

	r.Dirty = false
	return image
}
func (r *OptGame) Update(dt int) {}
func (r *OptGame) Draw(surface *ebiten.Image) {
	if r.Dirty {
		r.Image = r.Layout()
	}
	if r.Visibe {
		op := &ebiten.DrawImageOptions{}
		x, y := r.rect.Pos()
		op.GeoM.Translate(float64(x), float64(y))
		surface.DrawImage(r.Image, op)
	}
}

func (r *OptGame) Resize(rect []int) {
	r.rect = ui.NewRect(rect)
	r.Dirty = true
}
