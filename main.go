package main

import (
	"fmt"
	"math/bits"
	"strconv"
	"strings"
	"time"
)

//  0  1  2   3  4  5   6  7  8
//  9 10 11  12 13 14  15 16 17
// 18 19 20  21 22 23  24 25 26
// 27
// 36
// 45
// 54
// 63
// 72 73 74  75 76 77  78 79 80
type sudoku [9*9]cell

// NewSudoku returns a blanco sudoku data structure.
func NewSudoku() *sudoku {
	s := &sudoku{}

	cs := new9set()
	for i,_ := range cs {
		for r:= 0; r < 9; r++ {
			s[r*9+i].col = &cs[i]
		}
	}

	rs := new9set()
	for i,_ := range rs {
		for r:= 0; r < 9; r++ {
			s[i*9+r].row = &rs[i]
		}
	}

	sqs := new9set()
	for i,sn := range []int{0,3,6,27,30,33,54,57,60} {
		for _, n := range []int{0,1,2,9,10,11,18,19,20} {
			s[sn+n].sq = &sqs[i]
		}
	}

	return s
}

// Copy returns a copy of the receiver.
func (s *sudoku) copy() *sudoku {
	c := NewSudoku()
	for i,_ := range s {
		if n := s[i].num; n > 0 {
			c[i].assign(n)
		}
	}
	return c
}

// Print.
func (s *sudoku) print() {
	for r := 0; r < 9; r++ {
		for c := 0; c < 9; c++ {
			if n := s[r*9+c].num; n > 0 {
				fmt.Print(n, " ")
			} else {
				fmt.Print(". ")
			}
		}
		fmt.Println()
	}
	fmt.Println()
}

// Read reads a puzzle from a string.
func (s *sudoku) read(in string) {
	eol := strings.NewReplacer("\n", "", "\r\n", "")
	for i,v := range strings.Split(eol.Replace(in), " ") {
		if v == "." {
			continue
		}
		n,_ := strconv.Atoi(eol.Replace(v))
		s[i].assign(uint(n))
	}
}

// Solve sudoku.
func (s *sudoku) solve(l *log) {
	l.solveEnter()

	var tryCandidates set
	var tryCell int

	for {
		l.iterate()
		tryCandidates = 0
		progress := false
		allCellsHaveNumbers := true
		for i := 0; i < 9*9; i++ {
			if s[i].hasNum() {
				continue
			}

			allCellsHaveNumbers = false

			if s[i].fails() {
				l.noSolution(i, s)
				return
			}
			all, first, n := s[i].possibilities()
			if n == 1 {
				s[i].assign(first)
				progress = true
			} else {
				tryCandidates = all
				tryCell = i
			}
		}
		if allCellsHaveNumbers {
			l.solution(s)
			break
		}
		if !progress {
			// no more unambiguous cells
			break
		}
	}

	// try possibilities of last ambiguous cell
	for _,i := range tryCandidates.asSlice() {
		s1 := s.copy()
		s1[tryCell].assign(i)
		s1.solve(l)
		if l.isDone() {
			break
		}
	}

	l.solveExit()
}

// Cell of the sudoku grid.
type cell struct {
	num uint	// solution (0 for no solution)
	col *set
	row *set
	sq *set
}

// HasNum returns true if the cell has a solution.
func (c *cell) hasNum() bool {
	return c.num > 0
}

// Assign assigns num to the receiver cell and removes it from the search space.
func (c *cell) assign(num uint) {
	c.num = num
	c.col.remove(num)
	c.row.remove(num)
	c.sq.remove(num)
}

// Fails returns true if the receiver cell has no possible solutions.
func (c *cell) fails() bool {
	return c.num == 0 && (*c.col & *c.row & *c.sq) == 0
}

// Possibilities returns the all possibilities, the number of possibilities
// and if num == 1, the single possibility.
// For a cel with a number assigned no possibilities are returned.
func (c *cell) possibilities() (all set, single uint, num int) {
	if c.num != 0 {
		return 0, 0, 0
	}

	all = *c.col & *c.row & *c.sq
	num = bits.OnesCount(uint(all))
	if num == 1 {
		single = bitAsNum(uint(all))
	}

	return
}

// BitAsNum returns the first set bit in bitmap
func bitAsNum(i uint) uint {
	for e := uint(0); e < 9; e++ {
		if i & (1 << e) > 0 {
			return e+1
		}
	}
	return 0
}

// Set of 1..9 numbers represented as a bitmap.
type set uint

// Remove num from set.
func (s *set) remove(num uint) {
	*s &= ^(1 << (num-1))
}

// ForEach calls fn for each number in set.
func (s *set) forEach(fn func(uint)) {
	for e := uint(0); e < 9; e++ {
		if *s & (1 << e) > 0 {
			fn(e+1)
		}
	}
}

// AsSlice return the receiver as a slice of uint.
func (s *set) asSlice() (r []uint) {
	for e := uint(0); e < 9; e++ {
		if *s & (1 << e) > 0 {
			r = append(r, e+1)
		}
	}
	return
}

// New9set returns 9 sets all containing numbers 1..9.
func new9set() *[9]set {
	r := new([9]set)
	for i,_ := range r {
		r[i] = 0x1ff
	}
	return r
}

// Log the stats and result when solving a puzzle.
type log struct{
	depth int
	iter int
	solve int
	sol *sudoku
	solDepth int
	startTime time.Time
	solveTime time.Time
}

func NewLog() *log {
	return &log{
		startTime: time.Now(),
	}
}

func (l *log) solveEnter() {
	l.solve++
	l.depth++
}

func (l *log) solveExit() {
	l.depth--
}

func (l *log) iterate() {
	l.iter++
}

func (l *log) print() {
	fmt.Printf("total recursions: %d iterations: %d\n", l.solve, l.iter)
	if l.sol != nil {
		fmt.Printf("solution recursion: %d time: %s\n",
			l.solDepth,
			l.solveTime.Sub(l.startTime).String())
		l.sol.print()
	}
}

func (l *log)  solution(sudoku *sudoku) {
	l.solveTime = time.Now()
	l.sol = sudoku
	l.solDepth = l.depth
}

func (l *log)  noSolution(ci int, sudoku *sudoku) {
	//ind := strings.Repeat( " ", l.depth)
	//c := (ci % 9)+1
	//r := (ci / 9)+1
	//fmt.Printf("%sNo solution (c%d,r%d)\n", ind, c, r)
	//sudoku.print()
}

func (l *log) isDone() bool  {
	return l.sol != nil
}


func main() {
	for i,t := range puzzle {
		fmt.Printf("Puzzle %d\n", i)
		s := NewSudoku()
		s.read(t)
		m := NewLog()
		s.solve(m)
		m.print()
	}
}


var puzzle = []string{
`
. 9 1 . 5 . 2 7 . 
. . . 2 . 4 6 . 1 
2 7 . 6 9 1 . 5 8 
. . 3 1 . 5 7 4 . 
. 8 5 . . . 9 6 . 
. . 2 9 6 . 1 8 5 
. . 6 . 3 . 4 . . 
. 3 9 . 2 . . 1 . 
4 2 . . . . 8 3 9
`,/* 0 ** */


`
3 . . . . . . . . 
. . . 2 6 1 . 9 . 
2 . 1 . . 4 . 7 . 
. 3 . . . 7 . 1 . 
. 8 . 6 . . 9 . . 
. . 2 . 3 5 4 . . 
7 4 3 8 . 9 . . 6 
. 5 9 . 2 6 7 . . 
. . . . . 3 . . 9
`,/* 1 *** */


`
8 3 . . 4 . . . . 
. 4 2 . . 7 . . . 
. . 7 . . . 6 . . 
. . 6 . . . 9 . 1 
5 . . . . . . 8 6 
. 7 . 3 . . 5 2 . 
. . . 2 . . . 1 . 
2 6 . 7 . 9 . . . 
. . 9 . . 8 . . .
`,/* 2 **** */


`
1 . 2 . . . . . 7 
3 . 4 . 9 5 . 8 . 
. . . 1 . 7 9 3 . 
9 . 8 . . . 1 6 . 
4 . 6 . 5 3 2 . . 
. . . 9 1 . . 5 . 
. . . . . 2 . 9 6 
. 2 3 4 . 9 7 . . 
5 . 9 6 . . 3 . 4
`,/* 3 * */


`
. . . . . . . . .
. . . . . 3 . 8 5
. . 1 . 2 . . . .
. . . 5 . 7 . . .
. . 4 . . . 1 . .
. 9 . . . . . . .
5 . . . . . . 7 3
. . 2 . 1 . . . .
. . . . 4 . . . 9
`,/* 4 */


`
2 . . . 9 . . . 1 
3 . 9 . . 7 . . . 
. . 1 . 4 . . 7 . 
. 6 . . . . . . . 
. . . . . 3 . . . 
. . 8 6 . . 7 9 . 
6 . . 7 . . 8 . . 
1 2 3 . . 8 . . . 
. 8 7 . . 4 3 . .
`,/* 5 */


`
. 9 . . . 3 2 8 . 
5 . . . . . 3 9 . 
. . 6 . . . . . . 
6 . . 8 . . 4 . 5 
. 5 . 7 . . . 3 . 
. . 9 . 4 6 . . . 
. . . . 2 . 5 . . 
. 8 . 6 . . . . . 
. . . . . . . 2 8
`,/* 6 hard */


`
. . 6 . . . 9 . . 
. . . 6 . . . 7 5 
. 5 8 . . 7 . 1 4 
. 6 . . 1 . . 5 7 
. . . 7 5 . . 6 . 
. . . . . . 3 . . 
. . . 1 . . 7 8 3 
. 1 . . . 3 . 4 . 
. . . . 6 . 5 . .
`, /* 7 */


/*`
. . . . . . . . .
. . . . . . . . .
. . . . . . . . .
. . . . . . . . .
. . . . . . . . .
. . . . . . . . .
. . . . . . . . .
. . . . . . . . .
. . . . . . . . .
`,*/
}