package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	unrevealed = "#"
	mine       = "*"
)

type GameState struct {
	gridSize  int
	numMines  int
	board     [][]string
	revealed  [][]bool
	gameOver  bool
	win       bool
	input     textinput.Model
	message   string
	moveCount int
}

func NewGameState(size, mines int) (*GameState, error) {
	if size < 1 || mines >= size*size || mines < 0 {
		return nil, fmt.Errorf("invalid grid size or number of mines")
	}

	rand.Seed(time.Now().UnixNano())

	board := make([][]string, size)
	revealed := make([][]bool, size)
	// generate the board
	for i := 0; i < size; i++ {
		board[i] = make([]string, size)
		revealed[i] = make([]bool, size)
		for j := 0; j < size; j++ {
			board[i][j] = unrevealed
			revealed[i][j] = false
		}
	}

	// place the mines randomly
	minesPlaced := 0
	for minesPlaced < mines {
		row := rand.Intn(size)
		col := rand.Intn(size)
		if board[row][col] != mine {
			board[row][col] = mine
			minesPlaced++
		}
	}

	// calculate adjacent mine counts
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			if board[i][j] != mine {
				count := countAdjacentMines(board, i, j)
				if count > 0 {
					board[i][j] = strconv.Itoa(count)
				} else {
					// empty cell without neighboring mines
					board[i][j] = " "
				}
			}
		}
	}

	ti := textinput.New()
	ti.Placeholder = "row,col (e.g., 0,1)"
	ti.Focus()

	return &GameState{
		gridSize:  size,
		numMines:  mines,
		board:     board,
		revealed:  revealed,
		gameOver:  false,
		win:       false,
		input:     ti,
		message:   "Enter your move:",
		moveCount: 0,
	}, nil
}

func countAdjacentMines(board [][]string, row, col int) int {
	count := 0
	rows := len(board)
	cols := len(board[0])

	// iterate through the 3x3 neighborhood around the cell (row, col).
	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			// skip the cell itself.
			if i == 0 && j == 0 {
				continue
			}

			// get the real coordinate of a neighbor cell
			r, c := row+i, col+j
			// check if the neighbor cell is within the board bounds and contains a mine.
			if r >= 0 && r < rows && c >= 0 && c < cols && board[r][c] == mine {
				count++
			}
		}
	}
	return count
}

// revealCell reveals a cell and handles game over conditions.
func (g *GameState) revealCell(row, col int) {
	if row < 0 || row >= g.gridSize || col < 0 || col >= g.gridSize || g.revealed[row][col] {
		return
	}

	g.revealed[row][col] = true
	g.moveCount++

	if g.board[row][col] == mine {
		g.gameOver = true
		g.message = "BOOM! Game Over."
		return
	}

	if g.board[row][col] == " " {
		// automatically reveal adjacent empty cells
		g.revealAdjacentEmpty(row, col)
	}

	if g.checkWin() {
		g.gameOver = true
		g.win = true
		g.message = fmt.Sprintf("Congratulations! You won in %d moves.", g.moveCount)
	}
}

// revealAdjacentEmpty recursively reveals adjacent empty cells.
func (g *GameState) revealAdjacentEmpty(row, col int) {
	rows := g.gridSize
	cols := g.gridSize

	// Iterate through the 3x3 neighborhood around the current cell.
	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			// get the real coordinate of a neighbor cell
			r, c := row+i, col+j
			// check bounds, if the cell is not revealed and is not a mine
			if r >= 0 && r < rows && c >= 0 && c < cols && !g.revealed[r][c] && g.board[r][c] != mine {
				g.revealed[r][c] = true
				if g.board[r][c] == " " {
					g.revealAdjacentEmpty(r, c)
				}
			}
		}
	}
}

// checkWin checks if all non-mine cells have been revealed.
func (g *GameState) checkWin() bool {
	nonMineCells := 0
	revealedNonMines := 0
	for i := 0; i < g.gridSize; i++ {
		for j := 0; j < g.gridSize; j++ {
			if g.board[i][j] != mine {
				nonMineCells++
				if g.revealed[i][j] {
					revealedNonMines++
				}
			}
		}
	}
	return revealedNonMines == nonMineCells
}

// Init initializes the Bubble Tea model.
func (g *GameState) Init() tea.Cmd {
	return textinput.Blink
}

// Update handles UI updates based on events.
func (g *GameState) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if g.gameOver {
		if k, ok := msg.(tea.KeyMsg); ok {
			if k.Type == tea.KeyCtrlC || k.Type == tea.KeyEnter {
				return g, tea.Quit
			}
		}
		return g, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return g, tea.Quit
		case tea.KeyEnter:
			input := g.input.Value()
			g.input.Reset()
			return g.handleInput(input)
		}
	}

	var cmd tea.Cmd
	g.input, cmd = g.input.Update(msg)
	return g, cmd
}

// handleInput processes player input.
func (g *GameState) handleInput(input string) (*GameState, tea.Cmd) {
	parts := strings.Split(input, ",")
	if len(parts) != 2 {
		g.message = "Invalid input format. Use row,col (e.g., 0,1)"
		return g, nil
	}

	rowStr := strings.TrimSpace(parts[0])
	colStr := strings.TrimSpace(parts[1])

	row, errRow := strconv.Atoi(rowStr)
	col, errCol := strconv.Atoi(colStr)

	if errRow != nil || errCol != nil {
		g.message = "Invalid row or column value. Must be numbers."
		return g, nil
	}

	if row < 0 || row >= g.gridSize || col < 0 || col >= g.gridSize {
		g.message = fmt.Sprintf("Coordinates out of bounds (0-%d).", g.gridSize-1)
		return g, nil
	}

	if g.revealed[row][col] {
		g.message = "Cell already revealed. Choose another."
		return g, nil
	}

	g.revealCell(row, col)
	g.message = "Enter your move:"
	return g, nil
}

// View renders the current game state.
func (g *GameState) View() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("Minesweeper %dx%d (%d mines)\n", g.gridSize, g.gridSize, g.numMines))
	b.WriteString("   ")
	for i := 0; i < g.gridSize; i++ {
		b.WriteString(fmt.Sprintf("%d ", i))
	}
	b.WriteString("\n")
	b.WriteString(" --")
	for i := 0; i < g.gridSize; i++ {
		b.WriteString("--")
	}
	b.WriteString("\n")
	for i := 0; i < g.gridSize; i++ {
		b.WriteString(fmt.Sprintf("%d| ", i))
		for j := 0; j < g.gridSize; j++ {
			if g.gameOver {
				b.WriteString(fmt.Sprintf("%s ", g.board[i][j]))
			} else {
				if g.revealed[i][j] {
					b.WriteString(fmt.Sprintf("%s ", g.board[i][j]))
				} else {
					b.WriteString(fmt.Sprintf("%s ", unrevealed))
				}
			}
		}
		b.WriteString("\n")
	}
	b.WriteString("\n")
	b.WriteString(g.message + "\n")
	if !g.gameOver {
		b.WriteString(g.input.View())
	} else {
		b.WriteString("Press Ctrl+C or Enter to quit.\n")
	}
	return b.String()
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	var size, mines int

	fmt.Print("Enter grid size (n): ")
	_, err := fmt.Fscanln(reader, &size)
	if err != nil || size < 1 {
		fmt.Println("Invalid grid size.")
		return
	}

	fmt.Printf("Enter number of mines (less than %d): ", size*size)
	_, err = fmt.Fscanln(reader, &mines)
	if err != nil || mines < 0 || mines >= size*size {
		fmt.Printf("Invalid number of mines. Must be less than %d.\n", size*size)
		return
	}

	game, err := NewGameState(size, mines)
	if err != nil {
		fmt.Println("Error initializing game:", err)
		return
	}

	p := tea.NewProgram(game)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
