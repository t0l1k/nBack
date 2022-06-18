package ui

type Scene interface {
	Entered()
	Quit()
	Container
}
