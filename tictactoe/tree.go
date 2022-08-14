package tictactoe

import (
	"fmt"
	"os"
)

type Tree struct {
	State     *State
	Children  map[Point]*Tree
	Victories map[uint32]int
}

var known_trees map[uint32]*Tree

func FromState(s *State) (*Tree, bool) {
	// Create the known trees map if it doesn't exist.
	if known_trees == nil {
		known_trees = make(map[uint32]*Tree)
	}

	// Check if we already have a tree for this state or any rotation of the state.
	min_board := s.MinBoard()
	if tree, ok := known_trees[min_board]; ok {
		return tree, false
	}

	// Create a new tree.
	tree := &Tree{
		State:    s,
		Children: make(map[Point]*Tree),
	}
	known_trees[s.Board] = tree

	return tree, true
}

func LookupState(s *State) *Tree {
	return known_trees[s.MinBoard()]
}

func (t *Tree) AddChild(p Point, s *State) (*Tree, bool) {
	new_tree, new := FromState(s)
	t.Children[p] = new_tree
	return new_tree, new
}

func (t *Tree) GetChild(p Point) *Tree {
	return t.Children[p]
}

func (t *Tree) MiniMaxMove(as_player uint32) Point {
	moves := t.State.GetValidMoves()
	if len(moves) == 0 {
		return Point{}
	}

	// Create a map of moves to their scores.
	move_scores := make(map[Point]int)
	for _, move := range moves {
		new_state := t.State.Clone()
		new_state.SetAt(move.Row, move.Col, as_player)
		new_state.Minimize()
		// Get child tree.
		child_tree, _ := FromState(new_state)
		// Get child tree's score.
		move_scores[move] = child_tree.Victories[as_player] - child_tree.Victories[^as_player]
	}

	// Find the best move.
	best_move := Point{}
	best_score := 0
	for move, score := range move_scores {
		if score > best_score {
			best_move = move
			best_score = score
		}
	}

	return best_move
}

func (t *Tree) TallyVictories() {
	// Create victory map if needed.
	if t.Victories == nil {
		t.Victories = make(map[uint32]int)
	}

	// Assume this is run on a complete decision tree.
	// If this has no children, then we check for victories.
	if len(t.Children) == 0 {
		t.Victories[t.State.CheckVictory()]++
	}

	// Otherwise, recurse on children.
	for _, child := range t.Children {
		child.TallyVictories()
		// Add child victories to this tree's victories.
		for k, v := range child.Victories {
			t.Victories[k] += v
		}
	}
}

func (t *Tree) GraphViz(output_filename string, max_nodes int) {
	// Open/create the output file.
	f, err := os.Create(output_filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	fmt.Fprintf(f, "digraph {\n")

	// Set font to a monospaced font.
	fmt.Fprintf(f, "graph [fontname = \"Monaco\"];\n")
	fmt.Fprintf(f, "node [fontname = \"Monaco\"];\n")
	fmt.Fprintf(f, "edge [fontname = \"Monaco\"];\n")

	// Have set of known edges to not re-draw.
	known_edges := make(map[string]bool)

	node_count := 0
	queue := []*Tree{t}
	for len(queue) > 0 && node_count < max_nodes {
		node := queue[0]
		queue = queue[1:]
		node_count++

		// Generate a node label.
		label := fmt.Sprintf("state_%d", node.State.MinBoard())

		// Get node contents.
		contents := node.State.String()

		// Make victory percent string.
		child_sum := 0
		for _, v := range node.Victories {
			child_sum += v
		}
		victory_percent := fmt.Sprintf("%.2f X/%.2f O/%.2f draw",
			float64(node.Victories[BOARD_X])/float64(child_sum),
			float64(node.Victories[BOARD_O])/float64(child_sum),
			float64(node.Victories[BOARD_EMPTY])/float64(child_sum))

		// Print node.
		fmt.Fprintf(f, "%s [label=\"%s\n%s\"];\n", label, contents, victory_percent)

		// Get the best move to highlight it later.
		clone := node.State.Clone()
		best_move := node.MiniMaxMove(clone.GetNextPlayer())
		clone.SetAt(best_move.Row, best_move.Col, clone.GetNextPlayer())
		clone.Minimize()

		// Print edges.
		for _, child := range node.Children {
			child_label := fmt.Sprintf("state_%d", child.State.MinBoard())
			edge_label := ""
			// edge_label := fmt.Sprintf("%s -> %s;\n", label, child_label)
			if child.State.MinBoard() == clone.MinBoard() {
				edge_label = fmt.Sprintf("%s -> %s [color=red];\n", label, child_label)
			} else {
				edge_label = fmt.Sprintf("%s -> %s;\n", label, child_label)
			}
			if _, ok := known_edges[edge_label]; !ok {
				known_edges[edge_label] = true
				fmt.Fprint(f, edge_label)
			}
			queue = append(queue, child)
		}
	}

	fmt.Fprintf(f, "}\n")
}
