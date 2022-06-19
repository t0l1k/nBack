package ui

type Scene interface {
	Entered()
	Quit()
	Resize()
	Container
}
