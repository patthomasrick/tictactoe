package tictactoe

import (
	"fmt"
	"strings"
)

const BOARD_EMPTY uint32 = 0
const BOARD_X uint32 = 0b01
const BOARD_O uint32 = 0b10
const BOARD_MASK uint32 = 0b11
const BOARD_MASK_SIZE = 2

type State struct {
	Board      uint32
	nextPlayer bool
}

func (s *State) At(row, col int) uint32 {
	return (s.Board >> (BOARD_MASK_SIZE * (row*3 + col))) & BOARD_MASK
}

func (s *State) CharAt(row, col int) string {
	switch s.At(row, col) {
	case BOARD_X:
		return "X"
	case BOARD_O:
		return "O"
	default:
		return " "
	}
}

func (s *State) SetAt(row, col int, val uint32) {
	offset := BOARD_MASK_SIZE * (row*3 + col)
	s.Board &= ^(BOARD_MASK << offset)
	s.Board |= (val << offset)

	// Set next player.
	switch val {
	case BOARD_X:
		s.nextPlayer = false
	case BOARD_O:
		s.nextPlayer = true
	}
}

func (s *State) GetNextPlayer() uint32 {
	if s.nextPlayer {
		return BOARD_X
	}
	return BOARD_O
}

func (s *State) SetNextPlayer(val uint32) {
	s.nextPlayer = val == BOARD_X
}

func (s *State) RotateClockwise() {
	s.rotateClockwise()
	s.rotateClockwise()
}

func (s *State) rotateClockwise() {
	tmp := s.At(0, 0)

	s.SetAt(0, 0, s.At(1, 0))
	s.SetAt(1, 0, s.At(2, 0))

	s.SetAt(2, 0, s.At(2, 1))
	s.SetAt(2, 1, s.At(2, 2))

	s.SetAt(2, 2, s.At(1, 2))
	s.SetAt(1, 2, s.At(0, 2))

	s.SetAt(0, 2, s.At(0, 1))
	s.SetAt(0, 1, tmp)
}

func (s *State) RotateCounterClockwise() {
	s.rotateCounterClockwise()
	s.rotateCounterClockwise()
}

func (s *State) rotateCounterClockwise() {
	tmp := s.At(0, 0)

	s.SetAt(0, 0, s.At(0, 1))
	s.SetAt(0, 1, s.At(0, 2))

	s.SetAt(0, 2, s.At(1, 2))
	s.SetAt(1, 2, s.At(2, 2))

	s.SetAt(2, 2, s.At(2, 1))
	s.SetAt(2, 1, s.At(2, 0))

	s.SetAt(2, 0, s.At(1, 0))
	s.SetAt(1, 0, tmp)
}

func (s *State) CheckVictory() uint32 {
	victor := BOARD_EMPTY
	for i := 0; i < 4; i++ {
		if victor == BOARD_EMPTY {
			victor = s.checkRotationVictory()
		}
		s.RotateClockwise()
	}
	return victor
}

func (s *State) checkRotationVictory() uint32 {
	// Only check first row, first diagonal, and second row.
	// First row.
	if s.At(0, 0) != BOARD_EMPTY && s.At(0, 0) == s.At(0, 1) && s.At(0, 0) == s.At(0, 2) {
		return s.At(0, 0)
	}

	// Second row.
	if s.At(1, 0) != BOARD_EMPTY && s.At(1, 0) == s.At(1, 1) && s.At(1, 0) == s.At(1, 2) {
		return s.At(1, 0)
	}

	// First diagonal.
	if s.At(0, 0) != BOARD_EMPTY && s.At(0, 0) == s.At(1, 1) && s.At(0, 0) == s.At(2, 2) {
		return s.At(0, 0)
	}
	return BOARD_EMPTY
}

func (s *State) Clone() *State {
	return &State{s.Board, s.nextPlayer}
}

func (s *State) String() string {
	output := strings.Builder{}
	for row := 0; row < 3; row++ {
		for col := 0; col < 3; col++ {
			char := s.CharAt(row, col)
			if char == " " {
				char = "-"
			}
			output.WriteString(char)
		}
		output.WriteString("\n")
	}
	return output.String()
}

func (s *State) Print() {
	// Print ASCII square around the board.
	fmt.Println("+-------+")
	for row := 0; row < 3; row++ {
		fmt.Printf("| %s %s %s |\n", s.CharAt(row, 0), s.CharAt(row, 1), s.CharAt(row, 2))
	}
	fmt.Println("+-------+")
}

func (s *State) GetValidMoves() []Point {
	moves := []Point{}
	for row := 0; row < 3; row++ {
		for col := 0; col < 3; col++ {
			if s.At(row, col) == BOARD_EMPTY {
				moves = append(moves, Point{row, col})
			}
		}
	}
	return moves
}

func (s *State) MinBoard() uint32 {
	// Rotate the board until it is in the minimum configuration.
	min := s.Board
	for i := 0; i < 4; i++ {
		if s.Board < min {
			min = s.Board
		}
		s.RotateClockwise()
	}
	return min
}

func (s *State) Minimize() {
	s.Board = s.MinBoard()
}
