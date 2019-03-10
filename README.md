Sudoku solver in GO


How it works:

First step is creating a NewSudoku, this gets you a grid of 9 by 9 cells. 
Each cell reference 3 sets (column, row, square) of possible numbers.
Initially all sets contain the numbers 1..9 (the search space)

Next a puzzle is loaded. 
This assigns numbers to the cells and removes those numbers from the search space.

Now to solve the puzzle:
1. go over all the cells, if a cell has a single possible solution assign it.
2. keep repeating 1. until no assignments are made anymore.
3. take a cell with multi possible solutions, make a copy of the sudoku,
assign the first possibility and start solving this sodoku (recursion)
4. repeat 3 for the next possible solution.

On my laptop Puzzle 0 is solved in around 3us and Puzzle 6 in 2ms.
