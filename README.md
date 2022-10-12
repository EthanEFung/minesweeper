# vim minesweeper

a simple project to learn the [bubbletea](https://github.com/charmbracelet/bubbletea) tui framework. This is currently in MVP, with future plans to improve.

![](./minesweeper.gif)

# to run
- make sure to have a recent version
of go installed
- from your terminal navigate to this directory  
```bash
go install
go run .
```

# to play
- using h, j, k, l navigate the cursor
- press x to select the cell
- press d on a revealed number to select all adjacent cells that have not been flagged
- press f to flag the cell
- press r to reset the board
- press q to quit

# todos
- [x] create classic games "l+r" click functionality (clears all cells around a cell without flags)
- [x] create menu to configure the game
- [x] add timer
- [x] add view of the number of "potential" mines (# of mines - # of flags placed)
- [ ] ~~add mouse control~~
- [x] create scoreboard
- [ ] make this into vim go! (Where all the operations are exclusively std vim operations)
- [ ] create light and dark mode
- [x] add how to play menu
- [ ] allow users to jump multiple rows or columns
- [ ] rank the scoreboard and have seperate views for the rankings of each mode
- [x] after a score has been saved, place the scores table cursor on the most recent score
      (the game that was just played).
- [x] bug: timer not reset after each won game