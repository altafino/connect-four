// https://altafino.com AI driven software development and consulting
package main

import (
	"bufio"
	"fmt"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

// --- Constants ---

const (
	ROWS          = 6
	COLS          = 7
	EMPTY         = 0
	PLAYER1       = 1 // Human
	PLAYER2       = 2 // Computer
	WINDOW_LENGTH = 4
	AI_DEPTH      = 6 // Increased Depth for better foresight

	// Define large scores for win/loss to dominate heuristic scores
	WIN_SCORE  = 10000000
	LOSE_SCORE = -10000000
)

// Player symbols for display
var playerSymbols = map[int]string{
	EMPTY:   ".",
	PLAYER1: "X",
	PLAYER2: "O",
}

// --- Board Logic (Keep as before) ---

func createBoard() [][]int {
	board := make([][]int, ROWS)
	for r := 0; r < ROWS; r++ {
		board[r] = make([]int, COLS)
	}
	return board
}

// printBoard displays the current state of the board (with 1-7 headers)
func printBoard(board [][]int) {
	fmt.Println()
	// Print column headers (1 to COLS)
	for c := 0; c < COLS; c++ {
		fmt.Printf(" %d", c+1) // Display 1-based index
	}
	fmt.Println(" ")
	fmt.Println(strings.Repeat("--", COLS) + "-")

	// Print board rows (top-down)
	for r := 0; r < ROWS; r++ {
		fmt.Printf("|")
		for c := 0; c < COLS; c++ {
			fmt.Printf("%s|", playerSymbols[board[r][c]])
		}
		fmt.Println()
	}
	fmt.Println(strings.Repeat("==", COLS) + "=")
	fmt.Println()
}

func isValidLocation(board [][]int, col int) bool {
	if col < 0 || col >= COLS {
		return false
	}
	return board[0][col] == EMPTY
}

func getNextOpenRow(board [][]int, col int) (int, bool) {
	for r := ROWS - 1; r >= 0; r-- {
		if board[r][col] == EMPTY {
			return r, true
		}
	}
	return -1, false
}

func dropPiece(board [][]int, col int, player int) (int, int, bool) {
	row, ok := getNextOpenRow(board, col)
	if ok {
		board[row][col] = player
		return row, col, true
	}
	return -1, -1, false
}

func copyBoard(board [][]int) [][]int {
	newBoard := make([][]int, ROWS)
	for r := 0; r < ROWS; r++ {
		newBoard[r] = make([]int, COLS)
		copy(newBoard[r], board[r])
	}
	return newBoard
}

// --- Winning/Draw Logic (Keep as before, ensure checkWinAny is robust) ---

func checkWin(board [][]int, player int, lastRow, lastCol int) bool {
	// Check Horizontal
	for c := 0; c <= COLS-WINDOW_LENGTH; c++ {
		if board[lastRow][c] == player &&
			board[lastRow][c+1] == player &&
			board[lastRow][c+2] == player &&
			board[lastRow][c+3] == player {
			return true
		}
	}
	// Check Vertical
	for r := 0; r <= ROWS-WINDOW_LENGTH; r++ {
		if board[r][lastCol] == player &&
			board[r+1][lastCol] == player &&
			board[r+2][lastCol] == player &&
			board[r+3][lastCol] == player {
			return true
		}
	}
	// Check Positive Diagonal (/)
	for r := 0; r <= ROWS-WINDOW_LENGTH; r++ {
		for c := 0; c <= COLS-WINDOW_LENGTH; c++ {
			if board[r][c] == player &&
				board[r+1][c+1] == player &&
				board[r+2][c+2] == player &&
				board[r+3][c+3] == player {
				return true
			}
		}
	}
	// Check Negative Diagonal (\) - Corrected Loop Bounds if necessary
	// Check from top-left downwards to bottom-right for '\'
	for r := 0; r <= ROWS-WINDOW_LENGTH; r++ {
		 for c := 0; c <= COLS-WINDOW_LENGTH; c++ { // Original positive diag check was fine
			 if board[r][c] == player && board[r+1][c+1] == player && board[r+2][c+2] == player && board[r+3][c+3] == player {
				 return true
			 }
		 }
	 }
	 // Check from top-right downwards to bottom-left for '/'
	 for r := 0; r <= ROWS-WINDOW_LENGTH; r++ {
		 for c := WINDOW_LENGTH - 1; c < COLS; c++ {
			 if board[r][c] == player && board[r+1][c-1] == player && board[r+2][c-2] == player && board[r+3][c-3] == player {
				 return true
			 }
		 }
	 }

	return false
}

func isBoardFull(board [][]int) bool {
	for c := 0; c < COLS; c++ {
		if board[0][c] == EMPTY {
			return false
		}
	}
	return true
}

// checkWinAny is a helper for isTerminalNode, checking win across the whole board
// Ensure this covers all cases correctly for terminal checks.
func checkWinAny(board [][]int, player int) bool {
    // Check Horizontal
    for r := 0; r < ROWS; r++ {
        for c := 0; c <= COLS-WINDOW_LENGTH; c++ {
            if board[r][c] == player && board[r][c+1] == player && board[r][c+2] == player && board[r][c+3] == player {
                return true
            }
        }
    }
    // Check Vertical
    for c := 0; c < COLS; c++ {
        for r := 0; r <= ROWS-WINDOW_LENGTH; r++ {
            if board[r][c] == player && board[r+1][c] == player && board[r+2][c] == player && board[r+3][c] == player {
                return true
            }
        }
    }
     // Check Positive Diagonal (/)
    for r := 0; r <= ROWS-WINDOW_LENGTH; r++ {
        for c := 0; c <= COLS-WINDOW_LENGTH; c++ {
            if board[r][c] == player && board[r+1][c+1] == player && board[r+2][c+2] == player && board[r+3][c+3] == player {
                return true
            }
        }
    }
    // Check Negative Diagonal (\)
    for r := 0; r <= ROWS-WINDOW_LENGTH; r++ { // Check from top-right downwards
        for c := WINDOW_LENGTH - 1; c < COLS; c++ {
            if board[r][c] == player && board[r+1][c-1] == player && board[r+2][c-2] == player && board[r+3][c-3] == player {
                return true
            }
        }
    }
    return false
}


// --- AI Logic (Minimax with Alpha-Beta) ---

// *** IMPROVED evaluateWindow ***
func evaluateWindow(window []int, player int) int {
	score := 0
	opponent := PLAYER1
	if player == PLAYER1 {
		opponent = PLAYER2
	}

	playerPieces := 0
	opponentPieces := 0
	emptySlots := 0
	for _, piece := range window {
		if piece == player {
			playerPieces++
		} else if piece == opponent {
			opponentPieces++
		} else {
			emptySlots++
		}
	}

	// Assign scores based on pieces in the window
	if playerPieces == 4 {
		score += 10000 // Very high score for winning
	} else if playerPieces == 3 && emptySlots == 1 {
		score += 100 // Higher score for setting up a win
	} else if playerPieces == 2 && emptySlots == 2 {
		score += 5 // Small bonus for two pieces
	}

	// Penalize for opponent's pieces
	// Note: opponent having 4 pieces should be caught by terminal node check earlier in minimax
	if opponentPieces == 3 && emptySlots == 1 {
		score -= 500 // *** CRITICAL: Much higher penalty for opponent's 3-in-a-row threat ***
	} else if opponentPieces == 2 && emptySlots == 2 {
		score -= 3 // Small penalty
	}

	return score
}


// scorePosition calculates the heuristic score for the current board state
// (This function remains largely the same, relying on the improved evaluateWindow)
func scorePosition(board [][]int, player int) int {
	score := 0

	// Center column preference (slightly increased bonus)
	centerArray := make([]int, ROWS)
	for r := 0; r < ROWS; r++ {
		centerArray[r] = board[r][COLS/2]
	}
	centerCount := 0
	for _, piece := range centerArray {
		if piece == player {
			centerCount++
		}
	}
	score += centerCount * 6 // Slightly higher bonus for center control

	// Score Horizontal
	for r := 0; r < ROWS; r++ {
		for c := 0; c <= COLS-WINDOW_LENGTH; c++ {
			window := board[r][c : c+WINDOW_LENGTH]
			score += evaluateWindow(window, player) // Uses improved evaluation
		}
	}

	// Score Vertical
	for c := 0; c < COLS; c++ {
		for r := 0; r <= ROWS-WINDOW_LENGTH; r++ {
			window := make([]int, WINDOW_LENGTH)
			for i := 0; i < WINDOW_LENGTH; i++ {
				window[i] = board[r+i][c]
			}
			score += evaluateWindow(window, player) // Uses improved evaluation
		}
	}

	// Score Positive Diagonal (/)
	for r := 0; r <= ROWS-WINDOW_LENGTH; r++ {
		for c := 0; c <= COLS-WINDOW_LENGTH; c++ {
			window := make([]int, WINDOW_LENGTH)
			for i := 0; i < WINDOW_LENGTH; i++ {
				window[i] = board[r+i][c+i]
			}
			score += evaluateWindow(window, player) // Uses improved evaluation
		}
	}

	// Score Negative Diagonal (\)
	for r := 0; r <= ROWS-WINDOW_LENGTH; r++ {
		for c := WINDOW_LENGTH - 1; c < COLS; c++ {
			window := make([]int, WINDOW_LENGTH)
			for i := 0; i < WINDOW_LENGTH; i++ {
				window[i] = board[r+i][c-i]
			}
			score += evaluateWindow(window, player) // Uses improved evaluation
		}
	}

	return score
}

// isTerminalNode checks if the current state is a win, loss, or draw
func isTerminalNode(board [][]int) (bool, int) {
	if checkWinAny(board, PLAYER1) { // Human wins
		return true, PLAYER1
	}
	if checkWinAny(board, PLAYER2) { // AI wins
		return true, PLAYER2
	}
	if isBoardFull(board) { // Draw
		return true, 0 // 0 indicates a draw
	}
	return false, -1 // Not a terminal node
}


// minimax implements the algorithm with alpha-beta pruning
// Returns the best column and its score
func minimax(board [][]int, depth int, alpha float64, beta float64, maximizingPlayer bool) (int, int) {
	validLocations := getValidLocations(board)
	isTerminal, winner := isTerminalNode(board)

	if depth == 0 || isTerminal {
		if isTerminal {
			if winner == PLAYER2 { // AI wins
				return -1, WIN_SCORE // Use defined large score
			} else if winner == PLAYER1 { // Human wins
				return -1, LOSE_SCORE // Use defined large negative score
			} else { // Draw
				return -1, 0
			}
		} else { // Depth is zero
			// Return heuristic score for AI (PLAYER2)
			return -1, scorePosition(board, PLAYER2)
		}
	}

	if maximizingPlayer { // AI's turn (PLAYER2)
		value := math.Inf(-1)
        var bestCol int
        if len(validLocations) > 0 {
             bestCol = validLocations[rand.Intn(len(validLocations))] // Start with a random valid move
        } else {
             return -1, 0 // Should not happen if not terminal, but safety check
        }


		for _, col := range validLocations {
			tempBoard := copyBoard(board)
			// Simulate dropping the piece for the AI
			// *** FIX HERE: Use blank identifiers _ for unused row/col ***
			_, _, dropped := dropPiece(tempBoard, col, PLAYER2)
			if !dropped { continue } // Should not happen if validLocations is correct

			_, newScore := minimax(tempBoard, depth-1, alpha, beta, false)

			if newScore > int(value) {
				value = float64(newScore)
				bestCol = col
			}
			alpha = math.Max(alpha, value)
			if alpha >= beta {
				break // Beta cut-off
			}
		}
		return bestCol, int(value)

	} else { // Minimizing player's turn (Human - PLAYER1 simulation)
		value := math.Inf(1)
        var bestCol int
         if len(validLocations) > 0 {
             bestCol = validLocations[rand.Intn(len(validLocations))]
         } else {
             return -1, 0 // Safety check
         }

		for _, col := range validLocations {
			tempBoard := copyBoard(board)
			// Simulate dropping the piece for the Human
			// *** FIX HERE: Use blank identifiers _ for unused row/col ***
            _, _, dropped := dropPiece(tempBoard, col, PLAYER1)
            if !dropped { continue } // Should not happen if validLocations is correct

			_, newScore := minimax(tempBoard, depth-1, alpha, beta, true)

			if newScore < int(value) {
				value = float64(newScore)
				bestCol = col // Track best col for consistency
			}
			beta = math.Min(beta, value)
			if alpha >= beta {
				break // Alpha cut-off
			}
		}
		return bestCol, int(value)
	}
}

func getValidLocations(board [][]int) []int {
	validCols := []int{}
	for col := 0; col < COLS; col++ {
		if isValidLocation(board, col) {
			validCols = append(validCols, col)
		}
	}
	return validCols
}


// getComputerMove determines the AI's move using minimax, WITH explicit checks
func getComputerMove(board [][]int) int {
	validLocations := getValidLocations(board)
	opponent := PLAYER1

	// --- Start: Explicit Checks ---
	// 1. Check if AI can win in the next move
	for _, col := range validLocations {
		tempBoardWin := copyBoard(board)
		row, _, dropped := dropPiece(tempBoardWin, col, PLAYER2)
		if dropped && checkWin(tempBoardWin, PLAYER2, row, col) {
			// Output 1-based index
			fmt.Printf("Computer choosing immediate win in column %d\n", col+1)
			return col // Return 0-based index
		}
	}

	// 2. Check if Human could win on their next move, and block it
	for _, col := range validLocations {
		tempBoardBlock := copyBoard(board)
		row, _, dropped := dropPiece(tempBoardBlock, col, opponent)
		if dropped && checkWin(tempBoardBlock, opponent, row, col) {
			// Output 1-based index
			fmt.Printf("Computer blocking opponent's immediate win in column %d\n", col+1)
			return col // Return 0-based index
		}
	}
	// --- End: Explicit Checks ---


	// 3. If no immediate win or block needed, use Minimax for best strategic move
	fmt.Printf("Computer ('%s') thinking (Minimax Depth: %d)...\n", playerSymbols[PLAYER2], AI_DEPTH)
	startTime := time.Now()

	colInternal, score := minimax(board, AI_DEPTH, float64(LOSE_SCORE-1), float64(WIN_SCORE+1), true)

	duration := time.Since(startTime)

    // Internal check still uses 0-based index
	if !isValidLocation(board, colInternal) {
        fmt.Println("!!! AI MINIMAX returned invalid column, choosing random valid move as fallback !!!")
        validCols := getValidLocations(board)
        if len(validCols) > 0 {
            colInternal = validCols[rand.Intn(len(validCols))]
        } else {
             fmt.Println("!!! AI Error: No valid moves found on non-terminal board !!!")
             return 0 // Default fallback (0-based)
        }
    }

    // *** Output 1-based index for the computer's choice ***
	fmt.Printf("Computer chose column %d (Score: %d, Time: %v)\n", colInternal+1, score, duration)

    // Return the 0-based index for internal game logic
	return colInternal
}

// --- Game Flow (Keep as before) ---

// getPlayerInput prompts the human player for their move (using 1-7)
func getPlayerInput(board [][]int) int {
	reader := bufio.NewReader(os.Stdin)
	for {
		// Ask for 1-based column index
		fmt.Printf("Player %s (%s), enter column (1-%d): ", "1", playerSymbols[PLAYER1], COLS)
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			continue
		}
		input = strings.TrimSpace(input)
		colInput, err := strconv.Atoi(input)
		if err != nil {
			// Refer to 1-based index in error
			fmt.Printf("Invalid input. Please enter a number between 1 and %d.\n", COLS)
			continue
		}

		// Validate input is within 1-7 range
		if colInput < 1 || colInput > COLS {
			fmt.Printf("Column out of bounds. Please choose between 1 and %d.\n", COLS)
			continue
		}

		// *** Convert 1-based input to 0-based index for internal use ***
		colInternal := colInput - 1

		// Use internal 0-based index for validation check
		if !isValidLocation(board, colInternal) {
			// Refer to 1-based index in error message
			fmt.Printf("Column %d is full. Please choose another column.\n", colInput)
			continue
		}

		// Return the 0-based index for game logic
		return colInternal
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	board := createBoard()
	gameOver := false
	// Randomly choose who starts (optional, could make it PLAYER1 starts)
    // turn := rand.Intn(2) + 1
    turn := PLAYER1 // Let Human start for consistency

	fmt.Println("Welcome to Connect Four!")
     fmt.Printf("Human: %s, Computer: %s\n", playerSymbols[PLAYER1], playerSymbols[PLAYER2])
     fmt.Printf("Computer AI Depth: %d\n", AI_DEPTH)


	for !gameOver {
		printBoard(board)
		var col int
		var row int
		var dropped bool
        var playerSymbol string
        var playerName string

        if turn == PLAYER1 { // Human's turn
            playerName = "Human"
            playerSymbol = playerSymbols[PLAYER1]
			col = getPlayerInput(board)
			row, _, dropped = dropPiece(board, col, PLAYER1)

		} else { // Computer's turn
            playerName = "Computer"
            playerSymbol = playerSymbols[PLAYER2]
			col = getComputerMove(board)
			row, _, dropped = dropPiece(board, col, PLAYER2)
		}

        // Post-move checks
        if !dropped {
             // This indicates a potential serious bug if it happens.
             fmt.Printf("!!! Error: Failed to drop piece in column %d for Player %d. Exiting. !!!\n", col, turn)
             gameOver = true // Force exit on unexpected error
             continue
        }

        if checkWin(board, turn, row, col) {
            printBoard(board)
            fmt.Printf("%s (%s) WINS!\n", playerName, playerSymbol)
            gameOver = true
        } else if isBoardFull(board) {
            printBoard(board)
            fmt.Println("It's a DRAW!")
            gameOver = true
        }


		// Switch turns if game is not over
		if !gameOver {
			if turn == PLAYER1 {
				turn = PLAYER2
			} else {
				turn = PLAYER1
			}
		}
	}

	fmt.Println("Game Over.")
}
