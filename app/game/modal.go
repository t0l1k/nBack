package game

import (
	"fmt"

	"github.com/t0l1k/eui"
)

type ModalType string

const (
	Pos ModalType = "p" // Позиции
	Col ModalType = "c" // Цвета
	Sym ModalType = "s" // Символы
	Ari ModalType = "a" // Ариифметика
	Fig ModalType = "f" // Фигуры
	Aud ModalType = "w" // Звуки букв цифр ещё в поле идей
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
	score                  []int
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

func (m *Modality) ResetResults() {
	m.correct = 0
	m.wrong = 0
	m.missed = 0
}

func (m *Modality) Reset() {
	m.ResetResults()
	m.field = nil
	m.moveStatus = nil
	m.score = nil
}

func (m *Modality) SetCorrect(level int) {
	m.correct++
	move := make(map[ModalType]MoveType)
	move[m.sym] = AddCorrect
	m.SetValue(move)
	m.moveStatus = append(m.moveStatus, AddCorrect)
	m.score = append(m.score, level)
}

func (m *Modality) SetWrong() {
	m.wrong++
	move := make(map[ModalType]MoveType)
	move[m.sym] = AddWrong
	m.SetValue(move)
	m.moveStatus = append(m.moveStatus, AddWrong)
	m.score = append(m.score, 0)
}

func (m *Modality) SetMissed() {
	m.missed++
	move := make(map[ModalType]MoveType)
	move[m.sym] = AddMissed
	m.SetValue(move)
	m.moveStatus = append(m.moveStatus, AddMissed)
	m.score = append(m.score, 0)
}

func (m *Modality) SetRegular(level int) {
	move := make(map[ModalType]MoveType)
	move[m.sym] = AddRegular
	m.SetValue(move)
	m.moveStatus = append(m.moveStatus, AddRegular)
	m.score = append(m.score, level)
}

func (m *Modality) CheckMove(userMove bool, last, test, level int) (str string) {
	lastValue, testValue := m.GetField()[last], m.GetField()[test]
	str = fmt.Sprintf("progress for modal[%v] moves[%v-%v] values:[%v-%v]", m.GetSym(), last, test, testValue, lastValue)
	if userMove {
		if lastValue == testValue {
			m.SetCorrect(level)
			str += "correct answer!"
		} else {
			m.SetWrong()
			str += "wrong answer!"
		}
	} else if lastValue == testValue {
		m.SetMissed()
		str += "missed answer!"
	} else {
		m.SetRegular(level)
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
	case Col:
		result = "Col"
	case Sym:
		result = "Sym"
	case Ari:
		result = "Ari"
	case Fig:
		result = "Fig"
	case Aud:
		result = "Aud"
	}
	return result
}
