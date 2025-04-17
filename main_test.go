package main

import (
	"strings"
	"testing"
)

func TestGridInitialization(t *testing.T) {
	size := 5
	mines := 3
	game, err := NewGameState(size, mines)
	if err != nil {
		t.Fatalf("Error creating game state: %v", err)
	}
	mineCount := 0
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			if game.board[i][j] == mine {
				mineCount++
			}
		}
	}
	if mineCount != mines {
		t.Errorf("Expected %d mines, got %d", mines, mineCount)
	}

	if game.gridSize != size {
		t.Errorf("Expected size %d, got %d", size, game.gridSize)
	}

	// test invalid input for NewGameState
	_, err = NewGameState(0, 0)
	if err == nil {
		t.Errorf("Expected error for invalid size, got nil")
	}
	_, err = NewGameState(size, size*size+1)
	if err == nil {
		t.Errorf("Expected error for too many mines, got nil")
	}
}

func TestInputValidation(t *testing.T) {
	size := 3
	mines := 1
	game, err := NewGameState(size, mines)
	if err != nil {
		t.Fatalf("NewGameState failed: %v", err)
	}

	// simulate out-of-bounds input
	_, _ = game.handleInput("3,0")
	if !strings.Contains(game.message, "out of bounds") {
		t.Errorf("Expected out of bounds message, got: %s", game.message)
	}

	// simulate invalid row/col value
	_, _ = game.handleInput("a,b")
	if !strings.Contains(game.message, "Invalid row or column value") {
		t.Errorf("Expected invalid row/column message, got: %s", game.message)
	}

	// simulate a duplicate move
	game.revealCell(0, 0)
	_, _ = game.handleInput("0,0")
	if !strings.Contains(game.message, "Cell already revealed") {
		t.Errorf("Expected already revealed message, got: %s", game.message)
	}

	// simulate valid input
	_, _ = game.handleInput("1,1")
	if game.message == "Invalid input format. Use row,col (e.g., 0,1)" {
		t.Errorf("Expected message to change after valid input, got %s", game.message)
	}
}

func TestGameOverMine(t *testing.T) {
	size := 2
	mines := 1
	game, err := NewGameState(size, mines)
	if err != nil {
		t.Fatalf("NewGameState failed: %v", err)
	}
	// find the mine.
	mineRow, mineCol := -1, -1
	for r := 0; r < size; r++ {
		for c := 0; c < size; c++ {
			if game.board[r][c] == mine {
				mineRow, mineCol = r, c
				break
			}
		}
		if mineRow != -1 {
			break
		}
	}
	if mineRow == -1 || mineCol == -1 {
		t.Fatalf("Could not find the mine in the board")
	}

	game.revealCell(mineRow, mineCol) // Reveal the mine.
	if !game.gameOver || !strings.Contains(game.message, "BOOM") {
		t.Errorf("Expected game over due to mine, got: gameOver=%v, message='%s'", game.gameOver, game.message)
	}
}

func TestGameOverWin(t *testing.T) {
	size := 2
	mines := 0
	game, err := NewGameState(size, mines)
	if err != nil {
		t.Fatalf("NewGameState failed: %v", err)
	}
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			game.revealCell(i, j)
		}
	}
	if !game.gameOver || !game.win || !strings.Contains(game.message, "Congratulations") {
		t.Errorf("Expected win, got: gameOver=%v, win=%v, message='%s'", game.gameOver, game.win, game.message)
	}
}

func TestCountAdjacentMines(t *testing.T) {
	// create a simple board for testing.
	board := [][]string{
		{"*", "1", " "},
		{"2", " ", "1"},
		{" ", "1", "*"},
	}

	// test the center cell.
	count := countAdjacentMines(board, 1, 1)
	if count != 2 {
		t.Errorf("Expected 2 adjacent mines, got %d", count)
	}

	// test the top-left corner.
	count = countAdjacentMines(board, 1, 0)
	if count != 1 {
		t.Errorf("Expected 1 adjacent mine, got %d", count)
	}

	// test the bottom-right corner.
	count = countAdjacentMines(board, 0, 2)
	if count != 0 {
		t.Errorf("Expected 0 adjacent mine, got %d", count)
	}

	// test an edge cell.
	count = countAdjacentMines(board, 0, 1)
	if count != 1 {
		t.Errorf("Expected 1 adjacent mine, got %d", count)
	}
}
