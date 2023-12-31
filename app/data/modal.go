package data

import "github.com/t0l1k/eui"

const (
	Pos string = "p" // Позиции
	Sym string = "s" // Символы
	Col string = "c" // Цвета
	Ari string = "a" // Ариифметика
	Fig string = "f" // Фигуры
	Aud string = "w" // Звуки букв цифр
)

const (
	AddCorrect = "correct added"
	AddWrong   = "wrong added"
	AddMissed  = "missed added"
)

type Modality struct {
	eui.SubjectBase
	sym                    string
	correct, wrong, missed int
	field                  []int
}

func NewModality(sym string) *Modality {
	m := &Modality{sym: sym}
	return m
}

func (m *Modality) AddField(f []int) {
	m.field = f
}

func (m *Modality) GetField() []int {
	return m.field
}

func (m *Modality) Reset() {
	m.correct = 0
	m.wrong = 0
	m.missed = 0
	m.field = nil
}

func (m *Modality) SetCorrect(value int) {
	m.correct += value
	m.SetValue([]string{m.sym, AddCorrect})
}

func (m *Modality) SetWrong(value int) {
	m.wrong += value
	m.SetValue([]string{m.sym, AddWrong})
}

func (m *Modality) SetMissed(value int) {
	m.missed += value
	m.SetValue([]string{m.sym, AddMissed})
}

func (m Modality) GetSym() string {
	return m.sym
}

func (m Modality) String() (result string) {
	switch m.sym {
	case Pos:
		result = "Pos"
	case Sym:
		result = "Sym"
	case Col:
		result = "Col"
	case Ari:
		result = "Ari"
	case Fig:
		result = "Fig"
	case Aud:
		result = "Aud"
	}
	return result
}