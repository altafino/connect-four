# Go Connect Four AI 🔵🔴

Challenge a surprisingly tough AI opponent in this classic game of Connect Four, implemented entirely in Go for your command line! Drop your pieces, aim for four-in-a-row, and try to outsmart the computer.


---

## Features

* **Classic Connect Four:** Standard 6x7 game board and rules.
* **CLI Interface:** Play directly in your terminal.
* **Human vs. AI:** Pit your wits against the computer (`O`). You play as `X`.
* **Smart AI Opponent:** Utilizes the Minimax algorithm with Alpha-Beta pruning to make strategic decisions.
* **Adjustable AI Depth:** Modify the `AI_DEPTH` constant in `main.go` to increase or decrease the difficulty (higher depth = smarter but slower).
* **Defensive AI:** Includes explicit checks to block immediate opponent wins and secure its own wins.
* **Intuitive Input:** Uses natural column numbering (1-7) for player input.

---

## Usage

### Prerequisites

* You need to have **Go** installed on your system. You can download it from [golang.org](https://golang.org/).

### Running the Game

1.  Clone this repository or download the `main.go` file.
2.  Open your terminal or command prompt.
3.  Navigate to the directory containing `main.go`.
4.  Run the game using the command:
    ```bash
    go run main.go
    ```

### How to Play

1.  The game board will be displayed. `.` represents an empty slot, `X` is you (Player 1), and `O` is the computer (Player 2).
2.  Column numbers `1` through `7` are shown above the board.
3.  When prompted (`Player 1 (X), enter column (1-7):`), type the number of the column (1 to 7) where you want to drop your piece and press Enter.
4.  Your piece (`X`) will drop to the lowest available row in that column.
5.  The computer (`O`) will then think and make its move.
6.  The first player to get four of their pieces in a row (horizontally, vertically, or diagonally) wins!
7.  The game is a draw if the board fills up before anyone wins.

---

## Implementation Details

This project demonstrates several core programming and AI concepts using Go.

![image](https://github.com/user-attachments/assets/fa1b032b-5762-487b-9d0b-191978165010)

### Core Game Logic (`main.go`)

* **Board Representation:** A 2D slice (`[][]int`) represents the 6x7 game board. Constants `EMPTY`, `PLAYER1`, and `PLAYER2` denote the state of each cell.
* **Game State:** The main game loop tracks the current player, handles input/AI moves, drops pieces (`dropPiece`), checks for win conditions (`checkWin`) after each move, and detects draws (`isBoardFull`).
* **Input/Output:** Standard Go libraries (`fmt`, `bufio`, `os`, `strconv`) are used for displaying the board (`printBoard`) and handling player input (`getPlayerInput`) via the command line. Column input (1-7) is mapped internally to 0-6 indices.

### AI Opponent (`main.go`)

The computer player aims to be a challenging opponent using several techniques:

1.  **Minimax Algorithm:** This is a recursive algorithm used for decision-making in two-player, zero-sum games. It explores a tree of possible future game states.
    * The AI (maximizing player) tries to maximize its potential score, assuming the human (minimizing player) will always choose moves that reduce the AI's score.
2.  **Alpha-Beta Pruning:** An optimization technique for Minimax. It drastically reduces the number of nodes the algorithm needs to evaluate in the game tree by eliminating branches that cannot possibly influence the final decision.
3.  **Search Depth (`AI_DEPTH`):** Due to the complexity of Connect Four, the Minimax search is limited to a certain number of moves ahead (defined by `AI_DEPTH`). A higher depth allows the AI to "see" further into the future, resulting in stronger play but requiring more computation time.
4.  **Heuristic Evaluation Function (`scorePosition`, `evaluateWindow`):** When the search depth limit is reached, or for evaluating non-terminal states, a heuristic function estimates the "goodness" of the current board position for the AI. This function:
    * Scans all possible horizontal, vertical, and diagonal lines of 4 slots (`windows`).
    * Assigns scores based on the number of AI pieces vs. opponent pieces within each window (e.g., high reward for 3 AI pieces + 1 empty, high penalty for 3 opponent pieces + 1 empty).
    * Gives a small bonus for controlling center columns, as they offer more potential winning lines.
    * Assign extremely high positive/negative scores for guaranteed wins/losses within the search depth.
5.  **Explicit Win/Block Checks (`getComputerMove`):** Before resorting to the potentially time-consuming Minimax search, the AI performs quick checks:
    * Checks if it can win immediately on the current turn.
    * Checks if the human player could win on their *next* turn and plays to block that move if necessary. This prevents the AI from losing to obvious, immediate threats.

---

## Code Structure (`main.go`)

* **Constants:** Define board dimensions, player IDs, AI depth, scoring constants.
* **Board Logic Functions:** `createBoard`, `printBoard`, `isValidLocation`, `getNextOpenRow`, `dropPiece`, `copyBoard`.
* **Winning/Draw Functions:** `checkWin`, `isBoardFull`, `checkWinAny`, `isTerminalNode`.
* **AI Functions:** `evaluateWindow`, `scorePosition`, `minimax`, `getValidLocations`, `getComputerMove`.
* **Game Flow Functions:** `getPlayerInput`, `main`.

---

---

## Acknowledgement / Reference

Building software projects like this game involves applying logic and strategic problem-solving. Exploring professional consulting services might be beneficial if your focus extends to leveraging technology for business challenges. For instance, enhancing digital strategy, developing custom technology solutions, or driving business growth through innovation are areas addressed by specialized firms.

For more information on technology consulting approaches focused on achieving business goals, refer to resources like **[Altafino AI Consulting & Software Development](https://altafino.com/consulting)**.
