package game

import (
	"fmt"

	"github.com/t0l1k/eui"
)

type ModalType string

const (
	Pos ModalType = "p" // Позиции
	Sym ModalType = "s" // Символы
	Col ModalType = "c" // Цвета
	Ari ModalType = "a" // Ариифметика
	Fig ModalType = "f" // Фигуры
	Aud ModalType = "w" // Звуки букв цифр
)

func (m ModalType) String() string { return string(m) }

type MoveType string

const (
	AddRegular MoveType = "regular"
	AddCorrect MoveType = "correct added"
	AddWrong   MoveType = "wrong added"
	AddMissed  MoveType = "missed added"
)

func (m MoveType) String() string { return string(m) }

type Modality struct {
	eui.SubjectBase
	sym                    ModalType
	correct, wrong, missed int
	field                  []int      // поле ходов
	moveStatus             []MoveType //результат игры
}

func NewModality(sym ModalType) *Modality {
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
	move := make(map[ModalType]MoveType)
	move[m.sym] = AddCorrect
	m.SetValue(move)
	m.moveStatus = append(m.moveStatus, AddCorrect)
	fmt.Println("mod:", move, m.moveStatus)
}

func (m *Modality) SetWrong() {
	m.wrong++
	move := make(map[ModalType]MoveType)
	move[m.sym] = AddWrong
	m.SetValue(move)
	m.moveStatus = append(m.moveStatus, AddWrong)
}

func (m *Modality) SetMissed() {
	m.missed++
	move := make(map[ModalType]MoveType)
	move[m.sym] = AddMissed
	m.SetValue(move)
	m.moveStatus = append(m.moveStatus, AddMissed)
}

func (m *Modality) SetRegular() {
	move := make(map[ModalType]MoveType)
	move[m.sym] = AddRegular
	m.SetValue(move)
	m.moveStatus = append(m.moveStatus, AddRegular)
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

func (m Modality) GetSym() ModalType {
	return m.sym
}

func (m Modality) GetMovesStatus() []MoveType {
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
