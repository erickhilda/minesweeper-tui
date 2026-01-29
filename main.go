package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	unrevealed = "#"
	mine       = "*"
	flag       = "F"
)

var (
	mineColor       = lipgloss.Color("#FF0000") // Red
	oddColor        = lipgloss.Color("#0000FF") // Blue
	evenColor       = lipgloss.Color("#008000") // Green
	unrevealedColor = lipgloss.Color("#FFFFFF") // White
	emptyColor      = lipgloss.Color("#FFFFFF")
	flagColor       = lipgloss.Color("#FFFF00") // Yellow
	cursorBgColor   = lipgloss.Color("#333333") // Gray
)

// Define styles using the colors
var (
	mineStyle = lipgloss.NewStyle().
			Foreground(mineColor).
			Bold(true)
	oddStyle        = lipgloss.NewStyle().Foreground(oddColor)
	evenStyle       = lipgloss.NewStyle().Foreground(evenColor)
	unrevealedStyle = lipgloss.NewStyle().Foreground(unrevealedColor)
	emptyStyle      = lipgloss.NewStyle().Foreground(emptyColor)
	flagStyle       = lipgloss.NewStyle().Foreground(flagColor).Bold(true)
)

type GameState struct {
	gridSize  int
	numMines  int
	board     [][]string
	revealed  [][]bool
	flagged   [][]bool
	gameOver  bool
	win       bool
	cursorRow int
	cursorCol int
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
	flagged := make([][]bool, size)
	// generate the board
	for i := 0; i < size; i++ {
		board[i] = make([]string, size)
		revealed[i] = make([]bool, size)
		flagged[i] = make([]bool, size)
		for j := 0; j < size; j++ {
			board[i][j] = unrevealed
			revealed[i][j] = false
			flagged[i][j] = false
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

	return &GameState{
		gridSize:  size,
		numMines:  mines,
		board:     board,
		revealed:  revealed,
		flagged:   flagged,
		gameOver:  false,
		win:       false,
		cursorRow: 0,
		cursorCol: 0,
		message:   "Use arrow keys to move, Space to reveal, F to flag, ? for help",
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
			// check bounds, if the cell is not revealed, not flagged, and is not a mine
			if r >= 0 && r < rows && c >= 0 && c < cols && !g.revealed[r][c] && !g.flagged[r][c] && g.board[r][c] != mine {
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

// moveCursor moves the cursor by the given delta, with bounds checking.
func (g *GameState) moveCursor(deltaRow, deltaCol int) {
	newRow := g.cursorRow + deltaRow
	newCol := g.cursorCol + deltaCol

	// Clamp to grid bounds
	if newRow >= 0 && newRow < g.gridSize {
		g.cursorRow = newRow
	}
	if newCol >= 0 && newCol < g.gridSize {
		g.cursorCol = newCol
	}
}

// revealAtCursor reveals the cell at the current cursor position.
func (g *GameState) revealAtCursor() {
	row, col := g.cursorRow, g.cursorCol

	if g.revealed[row][col] {
		g.message = "Cell already revealed"
		return
	}

	if g.flagged[row][col] {
		g.message = "Cannot reveal flagged cell. Press F to unflag first."
		return
	}

	g.revealCell(row, col)

	if !g.gameOver {
		g.message = "Use arrow keys to move, Space to reveal, F to flag, ? for help"
	}
}

// toggleFlag toggles the flag state at the current cursor position.
func (g *GameState) toggleFlag() {
	row, col := g.cursorRow, g.cursorCol

	if g.revealed[row][col] {
		g.message = "Cannot flag revealed cell"
		return
	}

	g.flagged[row][col] = !g.flagged[row][col]

	if g.flagged[row][col] {
		g.message = "Cell flagged"
	} else {
		g.message = "Flag removed"
	}
}

// showHelp displays keyboard shortcuts.
func (g *GameState) showHelp() {
	g.message = "Controls: ↑↓←→/hjkl=move | Space/Enter=reveal | F=flag | ?=help | Q/Ctrl+C=quit"
}

// Init initializes the Bubble Tea model.
func (g *GameState) Init() tea.Cmd {
	return nil
}

// Update handles UI updates based on events.
func (g *GameState) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if g.gameOver {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "ctrl+c", "q":
				return g, tea.Quit
			}
		}
		return g, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc", "q":
			return g, tea.Quit

		// Cursor movement (arrow keys)
		case "up", "k":
			g.moveCursor(-1, 0)
		case "down", "j":
			g.moveCursor(1, 0)
		case "left", "h":
			g.moveCursor(0, -1)
		case "right", "l":
			g.moveCursor(0, 1)

		// Cell reveal
		case " ", "enter":
			g.revealAtCursor()

		// Toggle flag
		case "f":
			g.toggleFlag()

		// Help
		case "?":
			g.showHelp()
		}
	}

	return g, nil
}

// View renders the current game state.
func (g *GameState) View() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("Minesweeper %dx%d (%d mines) | Moves: %d\n", g.gridSize, g.gridSize, g.numMines, g.moveCount))

	// Top border
	b.WriteString("+")
	for i := 0; i < g.gridSize*2; i++ {
		b.WriteString("-")
	}
	b.WriteString("+\n")

	// Board rows with side borders
	for i := 0; i < g.gridSize; i++ {
		b.WriteString("|")
		for j := 0; j < g.gridSize; j++ {
			cellContent := ""
			style := emptyStyle
			isCursor := (i == g.cursorRow && j == g.cursorCol)

			// Determine cell content and style
			if g.gameOver {
				// Show entire board on game over
				cellContent = g.board[i][j]
				switch cellContent {
				case mine:
					style = mineStyle
				case "1", "3", "5", "7":
					style = oddStyle
				case "2", "4", "6", "8":
					style = evenStyle
				default:
					style = emptyStyle
				}
			} else {
				// During gameplay
				if g.flagged[i][j] {
					cellContent = flag
					style = flagStyle
				} else if g.revealed[i][j] {
					cellContent = g.board[i][j]
					switch cellContent {
					case " ":
						style = emptyStyle
					case "1", "3", "5", "7":
						style = oddStyle
					case "2", "4", "6", "8":
						style = evenStyle
					default:
						style = emptyStyle
					}
				} else {
					cellContent = unrevealed
					style = unrevealedStyle
				}
			}

			// Apply cursor highlighting
			if isCursor && !g.gameOver {
				// Combine cursor style with existing style
				style = style.Copy().Background(cursorBgColor).Reverse(true)
			}

			b.WriteString(style.Render(cellContent) + " ")
		}
		b.WriteString("|\n")
	}

	// Bottom border
	b.WriteString("+")
	for i := 0; i < g.gridSize*2; i++ {
		b.WriteString("-")
	}
	b.WriteString("+\n")

	b.WriteString("\n")
	b.WriteString(g.message + "\n")
	if g.gameOver {
		b.WriteString("Press Q/Ctrl+C to quit.\n")
	} else {
		b.WriteString("Press ? for help\n")
	}
	return b.String()
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	var size, mines int
	var choice int

	fmt.Println("Select difficulty:")
	fmt.Println("1. Easy   (8x8, 10 mines)")
	fmt.Println("2. Medium (12x12, 20 mines)")
	fmt.Println("3. Hard   (16x16, 40 mines)")
	fmt.Println("4. Custom (define your own)")
	fmt.Print("Enter choice (1-4): ")

	_, err := fmt.Fscanln(reader, &choice)
	if err != nil || choice < 1 || choice > 4 {
		fmt.Println("Invalid choice. Please select 1, 2, 3, or 4.")
		return
	}

	switch choice {
	case 1:
		size = 8
		mines = 10
	case 2:
		size = 12
		mines = 20
	case 3:
		size = 16
		mines = 40
	case 4:
		fmt.Print("Enter grid size (n): ")
		_, err = fmt.Fscanln(reader, &size)
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
