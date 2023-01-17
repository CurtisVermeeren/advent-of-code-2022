package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

type dir int

// Define each direction a player can be facing. Values also represent their score according to the problem.
const right dir = 0
const down dir = 1
const left dir = 2
const up dir = 3
const invalidDir dir = -1

// Convert a direction integer back to the string it represents
func (d dir) String() string {
	switch d {
	case 0:
		return "right"
	case 1:
		return "down"
	case 2:
		return "left"
	case 3:
		return "up"
	default:
		return "invalid"
	}
}

// Turn clockwise
func (d dir) next() dir { return (d + 1) % 4 }

// Turn counter-clockwise
func (d dir) prev() dir { return (d + 3) % 4 }

// Turn 180 degrees
func (d dir) flipped() dir { return (d + 2) % 4 }

// turn returns the direction the player faces based on turnDir input
func (d dir) turn(turnDir byte) dir {
	switch turnDir {
	case 'R':
		return d.next()
	case 'L':
		return d.prev()
	default:
		panic("Invalid turn")
	}
}

// player represents the players current position on the board in cell and the direction they are facing.
type player struct {
	cell      *cell
	direction dir
}

// move is used to move a player along the board.
// Returns false if the destination cell is a wall
// Returns true if the move was successful
func (p *player) move() bool {
	next := p.cell.neighbors[p.direction]
	// Destination is a wall
	if next.val == '#' {
		return false
	}
	// Destination is on another face so the direction must be changed
	if p.cell.face.id != next.face.id {
		p.direction = (next.face.entrySide(p.cell.face) + 2) % 4
	}
	// Move the player
	p.cell = next
	return true
}

// password calculates the password for the problem as defined in the problem description
func (p *player) password() int {
	return 1000*(p.cell.row+1) + 4*(p.cell.col+1) + int(p.direction)
}

// As defined in the problem the board in broken up into several faces for part 2.
// id is used as a unique identifier for each face
// originRow and originCol are the x and y coordinates of the top-left corner of the face
// cells is a 2d array of all cells on the face
// neighbors is a slice of all faces neighbouring this face on a cube
type face struct {
	id, originRow, originCol int
	cells                    [][]*cell
	neighbors                [4]*face
}

// NewFace is a constructor for a face of the board
func NewFace(id, originRow, originCol, sideLength int) *face {
	f := face{
		id:        id,
		originRow: originRow,
		originCol: originCol,
		cells:     make([][]*cell, sideLength),
	}
	for i := range f.cells {
		f.cells[i] = make([]*cell, sideLength)
	}
	return &f
}

// sideClockwise is used to obtain the cells of the face f in a specific order based on the side passed as input.
// This can be used when the player is moving along a face this method can retrieve the cells in the order the player will encounter them.
func (f face) sideClockwise(side dir) []*cell {
	cells := make([]*cell, len(f.cells))
	n := len(cells) - 1
	var row, col, rowInc, colInc int
	switch side {
	case right:
		row, col, rowInc, colInc = 0, n, 1, 0
	case down:
		row, col, rowInc, colInc = n, n, 0, -1
	case left:
		row, col, rowInc, colInc = n, 0, -1, 0
	case up:
		row, col, rowInc, colInc = 0, 0, 0, 1
	}
	for i := 0; i <= n; i++ {
		cells[i] = f.cells[row][col]
		row += rowInc
		col += colInc
	}

	return cells
}

// entrySide returns a direction representing the side of the current face f that the player entered from the source face.
// This can be used when the player is moves from one face to another. It will determine the side of the new face that the player entered from.
func (f face) entrySide(source *face) dir {
	for entry := right; entry <= up; entry++ {
		if f.neighbors[entry] == source {
			return entry
		}
	}
	return invalidDir
}

// cell is a singular cell on the board
// face is a pointer to the face object that this cell belongs to
// row and col are the x y coordinates of the cell on the board.
// val is the value of the call. It could be a wall or the player
// neighbors are the 4 cells up, down, left, and right of the cell
type cell struct {
	face      *face
	row, col  int
	val       byte
	neighbors [4]*cell
}

// board represents the entire game board
// faces holds the 6 faces of the cube board as defined in part 2
// p points to the player state
// moves is a slice of integers that represent the number of cells the player must move on each step to solve the problem
// tursn is a slioe of bytes each of which indicate a direction L or R that the player must turn each step to solve the problem
type board struct {
	faces [6]*face
	p     *player
	moves []int
	turns []byte
}

// run is the main loop of the player
// It will iterate over all moves and turns to move the player around the board
func (b *board) run() {
	for i, m := range b.moves {
		for j := 0; j < m; j++ {
			if !b.p.move() {
				break
			}
		}
		if i < len(b.turns) {
			b.p.direction = b.p.direction.turn(b.turns[i])
		} else {
			break
		}
	}
}

// connectSides is used to connect faces f1 and f2 along edges s1 and s2
// cells in side s1 and s2 are pointed to each other
func connectSides(f1, f2 *face, s1, s2 dir) {
	side1 := f1.sideClockwise(s1)
	side2 := f2.sideClockwise(s2)
	n := len(side1) - 1
	for i := range side1 {
		cell1, cell2 := side1[i], side2[n-i]
		cell1.neighbors[s1], cell2.neighbors[s2] = cell2, cell1
	}
	f1.neighbors[s1], f2.neighbors[s2] = f2, f1
}

// parseCube takes the input strings of the board and returns a cube of 6 face objects
func parseCube(input []string) [6]*face {
	// Count total number of cells to derive cube face size and board width
	cellCount := 0
	boardWidth := 0
	for _, line := range input {
		if len(line) > boardWidth {
			boardWidth = len(line)
		}
		for _, r := range line {
			if r != ' ' {
				cellCount++
			}
		}
	}
	sideLength := int(math.Sqrt(float64(cellCount / 6)))

	// Create faces and parse cells
	var faces [6]*face
	faceId := 0
	faceGrid := make([][]*face, len(input)/sideLength)
	for faceRowNum := range faceGrid {
		faceRow := make([]*face, boardWidth/sideLength)
		faceGrid[faceRowNum] = faceRow
		boardRowNum := faceRowNum * sideLength

		for faceColNum := range faceGrid[faceRowNum] {
			boardColNum := faceColNum * sideLength
			if boardColNum >= len(input[boardRowNum]) || input[boardRowNum][boardColNum] == ' ' {
				continue
			}

			newFace := NewFace(faceId, boardRowNum, boardColNum, sideLength)
			for row := 0; row < sideLength; row++ {
				for col := 0; col < sideLength; col++ {
					cellRow, cellCol := newFace.originRow+row, newFace.originCol+col
					newCell := cell{
						face: newFace,
						row:  cellRow,
						col:  cellCol,
						val:  input[cellRow][cellCol],
					}
					newFace.cells[row][col] = &newCell
				}
			}

			faces[faceId] = newFace
			faceId++
			faceGrid[faceRowNum][faceColNum] = newFace
		}
	}

	// Set cell neighbors within faces
	for _, f := range faces {
		for _, row := range f.cells {
			for colNum, rightCell := range row[1:] {
				c := row[colNum]
				if c.neighbors[right] != nil || rightCell.neighbors[left] != nil {
					panic("Overwriting neighbor config")
				}
				c.neighbors[right], rightCell.neighbors[left] = rightCell, c
			}
		}
		for rowNum, downRow := range f.cells[1:] {
			for colNum, downCell := range downRow {
				c := f.cells[rowNum][colNum]
				if c.neighbors[down] != nil || downCell.neighbors[up] != nil {
					panic("Overwriting neighbor config")
				}
				c.neighbors[down], downCell.neighbors[up] = downCell, c
			}
		}
	}

	// Set face neighbors based on board before folding
	for _, fRow := range faceGrid {
		for colNum := range fRow[1:] {
			if fRow[colNum] != nil && fRow[colNum+1] != nil {
				connectSides(fRow[colNum], fRow[colNum+1], right, left)
			}
		}
	}
	for rowNum, fRow := range faceGrid[1:] {
		for colNum, f := range fRow {
			if faceGrid[rowNum][colNum] != nil && f != nil {
				connectSides(faceGrid[rowNum][colNum], f, down, up)
			}
		}
	}

	// Finally, fold the cube and set new neighbors
	// (6 faces)*(4 edges per face) = 24 edges forms 12 pairs
	// 5 pairs are handled in the last section, 7 remain
	// Strategy: Repeatedly check for L-shapes and fold them
	for disconnectedPairs := 7; disconnectedPairs > 0; {
		for _, f := range faces {
			for side := right; side <= up; side++ {
				if f.neighbors[side] != nil {
					continue
				}

				next, prev := side.next(), side.prev()
				if mid := f.neighbors[next]; mid != nil { // clockwise search
					next = mid.entrySide(f).next()
					if fold := mid.neighbors[next]; fold != nil {
						next = fold.entrySide(mid).next()
						if fold.neighbors[next] == nil {
							connectSides(f, fold, side, next)
							disconnectedPairs--
						}
					}
				} else if mid := f.neighbors[prev]; mid != nil { // anticlockwise search
					prev = mid.entrySide(f).prev()
					if fold := mid.neighbors[prev]; fold != nil {
						prev = fold.entrySide(mid).prev()
						if fold.neighbors[prev] == nil {
							connectSides(f, fold, side, prev)
							disconnectedPairs--
						}
					}
				}
			}
		}
	}

	return faces
}

// parseInstructions takes the movement string from the input file and creates two slices for actions
// moves is a slice of integers each of which is the number of cells to move the player that step
// turns is a slice of bytes each of which is R or L indictating a direction to turn
// Both slices are in order from the input with index 0 being the first action to perform
func parseInstruction(instruction string) ([]int, []byte) {
	moves, turns := make([]int, 0), make([]byte, 0)
	for len(instruction) > 0 {
		nextTurn := strings.IndexAny(instruction, "LR")
		if nextTurn == -1 {
			if tiles, err := strconv.Atoi(instruction); err == nil {
				moves = append(moves, tiles)
				break
			} else {
				panic(err)
			}
		} else {
			if tiles, err := strconv.Atoi(instruction[:nextTurn]); err == nil {
				moves = append(moves, tiles)
			} else {
				panic(err)
			}
			turns = append(turns, instruction[nextTurn])
			instruction = instruction[nextTurn+1:]
		}
	}
	return moves, turns
}

// parseInput is used to build the board object from the input strings
// It calls the parseCube method and parseInstruction to get the parts of the boardt.
func parseInput(input []string) *board {
	var emptyLine int
	for lineNum, line := range input {
		if line == "" {
			emptyLine = lineNum
			break
		}
	}

	b := new(board)
	b.faces = parseCube(input[:emptyLine])
	b.moves, b.turns = parseInstruction(input[emptyLine+1])

	for colNum := 0; colNum < len(b.faces[0].cells); colNum++ {
		if startingCell := b.faces[0].cells[0][colNum]; startingCell.val == '.' {
			b.p = &player{cell: startingCell, direction: right}
			break
		}
	}

	return b
}

// readInput takes a filename from stdin and reads it line by line into a slice of strings
func readInput() []string {
	lines := make([]string, 0)
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if scanner.Err() != nil {
		panic(scanner.Err())
	}
	return lines
}

func main() {
	input := readInput()
	cube := parseInput(input)
	cube.run()
	fmt.Println(cube.p.password())
}
