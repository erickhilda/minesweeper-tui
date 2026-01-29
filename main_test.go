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

	// Test revealing already revealed cell
	game.revealCell(0, 0)
	game.cursorRow, game.cursorCol = 0, 0
	game.revealAtCursor()
	if !strings.Contains(game.message, "Cell already revealed") {
		t.Errorf("Expected already revealed message, got: %s", game.message)
	}

	// Test valid reveal
	game.cursorRow, game.cursorCol = 1, 1
	game.revealAtCursor()
	if game.revealed[1][1] != true {
		t.Errorf("Expected cell to be revealed")
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

// TestCursorMovement tests cursor movement and bounds checking.
func TestCursorMovement(t *testing.T) {
	size := 3
	mines := 1
	game, err := NewGameState(size, mines)
	if err != nil {
		t.Fatalf("NewGameState failed: %v", err)
	}

	// Test initial position
	if game.cursorRow != 0 || game.cursorCol != 0 {
		t.Errorf("Expected cursor at (0,0), got (%d,%d)", game.cursorRow, game.cursorCol)
	}

	// Test movement within bounds
	game.moveCursor(1, 1)
	if game.cursorRow != 1 || game.cursorCol != 1 {
		t.Errorf("Expected cursor at (1,1), got (%d,%d)", game.cursorRow, game.cursorCol)
	}

	// Test bounds clamping (top-left)
	game.cursorRow, game.cursorCol = 0, 0
	game.moveCursor(-1, -1)
	if game.cursorRow != 0 || game.cursorCol != 0 {
		t.Errorf("Cursor went out of bounds: (%d,%d)", game.cursorRow, game.cursorCol)
	}

	// Test bounds clamping (bottom-right)
	game.cursorRow, game.cursorCol = size-1, size-1
	game.moveCursor(1, 1)
	if game.cursorRow != size-1 || game.cursorCol != size-1 {
		t.Errorf("Cursor went out of bounds: (%d,%d)", game.cursorRow, game.cursorCol)
	}
}

// TestFlagToggle tests flagging and unflagging cells.
func TestFlagToggle(t *testing.T) {
	size := 3
	mines := 1
	game, err := NewGameState(size, mines)
	if err != nil {
		t.Fatalf("NewGameState failed: %v", err)
	}

	// Flag a cell
	game.toggleFlag()
	if !game.flagged[0][0] {
		t.Errorf("Expected cell (0,0) to be flagged")
	}

	// Unflag the cell
	game.toggleFlag()
	if game.flagged[0][0] {
		t.Errorf("Expected cell (0,0) to be unflagged")
	}

	// Cannot flag revealed cell
	game.revealed[0][0] = true
	game.toggleFlag()
	if game.flagged[0][0] {
		t.Errorf("Should not be able to flag revealed cell")
	}
	if !strings.Contains(game.message, "Cannot flag revealed cell") {
		t.Errorf("Expected message about flagged cell, got: %s", game.message)
	}
}

// TestRevealFlaggedCell tests that flagged cells cannot be revealed.
func TestRevealFlaggedCell(t *testing.T) {
	size := 3
	mines := 1
	game, err := NewGameState(size, mines)
	if err != nil {
		t.Fatalf("NewGameState failed: %v", err)
	}

	// Flag cell at cursor
	game.flagged[0][0] = true

	// Try to reveal flagged cell
	game.revealAtCursor()
	if game.revealed[0][0] {
		t.Errorf("Flagged cell should not be revealed")
	}
	if !strings.Contains(game.message, "Cannot reveal flagged cell") {
		t.Errorf("Expected message about flagged cell, got: %s", game.message)
	}
}

// TestFlaggedCellsNotAutoRevealed tests that auto-reveal respects flags.
func TestFlaggedCellsNotAutoRevealed(t *testing.T) {
	size := 3
	mines := 0
	game, err := NewGameState(size, mines)
	if err != nil {
		t.Fatalf("NewGameState failed: %v", err)
	}

	// Flag a cell
	game.flagged[1][1] = true

	// Reveal adjacent empty cell (should auto-reveal neighbors except flagged)
	game.revealCell(0, 0)

	if game.revealed[1][1] {
		t.Errorf("Flagged cell should not be auto-revealed")
	}
}
