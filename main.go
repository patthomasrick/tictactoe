package main

import (
	"fmt"

	"patthomasrick.net/TicTacToc/v2/tictactoe"
)

func main() {
	// Create a decision tree for all possible moves.
	tree_queue := make([]*tictactoe.Tree, 0)
	empty_board := tictactoe.State{}
	empty_board.SetNextPlayer(tictactoe.BOARD_X)
	root, _ := tictactoe.FromState(&empty_board)
	tree_queue = append(tree_queue, root)

	// Simple stats.
	num_nodes := 0
	num_victories := make(map[uint32]int)

	for len(tree_queue) > 0 {
		// Pop the next tree off the queue.
		tree := tree_queue[0]
		tree_queue = tree_queue[1:]
		num_nodes++

		// If the board is a winning state, stop.
		victor := tree.State.CheckVictory()
		if victor != tictactoe.BOARD_EMPTY {
			num_victories[victor]++
			continue
		}

		// Get possible moves.
		for _, p := range tree.State.GetValidMoves() {
			// Make a copy of the board.
			board := tree.State.Clone()
			// Make a move.
			board.SetAt(p.Row, p.Col, tree.State.GetNextPlayer())
			board.Minimize()
			// Add a child to the tree.
			child, new := tree.AddChild(p, board)
			if new {
				tree_queue = append(tree_queue, child)
			}
		}
	}

	fmt.Printf("%d nodes\n", num_nodes)
	fmt.Printf("%d X victories\n", num_victories[tictactoe.BOARD_X])
	fmt.Printf("%d O victories\n", num_victories[tictactoe.BOARD_O])

	// Tally the victories for each player.
	root.TallyVictories()

	// Export the decision tree.
	root.GraphViz("output.dot", 100)
}
