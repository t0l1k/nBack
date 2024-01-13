package data

import (
	"fmt"

	"github.com/t0l1k/eui"
)

const (
	Pos string = "p" // Позиции
	Sym string = "s" // Символы
	Col string = "c" // Цвета
	Ari string = "a" // Ариифметика
	Fig string = "f" // Фигуры
	Aud string = "w" // Звуки букв цифр
)

const (
	Regular    = "regular"
	AddCorrect = "correct added"
	AddWrong   = "wrong added"
	AddMissed  = "missed added"
)

type Modality struct {
	eui.SubjectBase
	sym                    string
	correct, wrong, missed int
	field                  []int
	moveStatus             []string
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
	m.moveStatus = nil
}

func (m *Modality) SetCorrect() {
	m.correct++
	m.SetValue([]string{m.sym, AddCorrect})
	m.moveStatus = append(m.moveStatus, AddCorrect)
}

func (m *Modality) SetWrong() {
	m.wrong++
	m.SetValue([]string{m.sym, AddWrong})
	m.moveStatus = append(m.moveStatus, AddWrong)
}

func (m *Modality) SetMissed() {
	m.missed++
	m.SetValue([]string{m.sym, AddMissed})
	m.moveStatus = append(m.moveStatus, AddMissed)
}

func (m *Modality) SetRegular() {
	m.moveStatus = append(m.moveStatus, Regular)
}

func (m *Modality) CheckMove(userMove bool, last, test int) (str string) {
	lastValue, testValue := m.GetField()[last], m.GetField()[test]
	str = fmt.Sprintf("progress for modal[%v] moves[%v-%v] values:[%v-%v]", m.GetSym(), last, test, testValue, lastValue)
	if userMove {
		if lastValue == testValue {
			m.SetCorrect()
			str += "correct answer!"
		} else {
			m.SetWrong()
			str += "wrong answer!"
		}
	} else if lastValue == testValue {
		m.SetMissed()
		str += "missed answer!"
	} else {
		m.SetRegular()
		str += "regular move!"
	}
	return str
}

func (m Modality) GetSym() string {
	return m.sym
}

func (m Modality) GetMovesStatus() []string {
	return m.moveStatus
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
