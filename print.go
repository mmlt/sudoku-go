package main

import "fmt"

// PrintStructure shows what sets are shared by what cells.
// Addresses are abbreviated to make the pattern easier to spot.
func (s *sudoku) printStructure() {
	fn := func(p *set) string {
		s := fmt.Sprintf("%p", p)
		return string([]byte(s)[len(s)-3:])
	}
	fmt.Println("c   r   s     c   r   s     c   r   s     c   r   s     c   r   s     c   r   s     c   r   s     c   r   s     c   r   s")
	for r := 0; r < 9; r++ {
		for c := 0; c < 9; c++ {
			x := s[r*9+c]
			fmt.Printf("%s %s %s | ", fn(x.col), fn(x.row), fn(x.sq))
		}
		fmt.Println()
	}
	fmt.Println()
}

// PrintSets shows the columns, rows and squares sets content AKA the search space.
func (s *sudoku) printSets() {
	fmt.Println("c   r   s     c   r   s     c   r   s     c   r   s     c   r   s     c   r   s     c   r   s     c   r   s     c   r   s")
	for r := 0; r < 9; r++ {
		for c := 0; c < 9; c++ {
			x := s[r*9+c]
			fmt.Printf("%03x %03x %03x | ", *x.col, *x.row, *x.sq)
		}
		fmt.Println()
	}
	fmt.Println()
}
