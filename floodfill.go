package main

import (
	"fmt"
	"runtime"
)

type Canvas struct {
	contents [][]byte
	visited  [][]bool
}

type Node struct {
	X     int
	Y     int
	Color byte
}

func (c *Canvas) Init(width int, height int, blankChar byte) {
	c.contents = make([][]byte, width)
	for i := 0; i < width; i++ {
		c.contents[i] = make([]byte, height)
		for j := 0; j < height; j++ {
			c.contents[i][j] = blankChar
		}
	}
}

func (c *Canvas) Print() {
	for _, row := range c.contents {
		fmt.Println(string(row))
	}
}

func (c *Canvas) setVisitedMatrixToFalse() {
	width := len(c.contents)
	height := len(c.contents[0])
	c.visited = make([][]bool, width)
	for i := 0; i < width; i++ {
		c.visited[i] = make([]bool, height)
		for j := 0; j < height; j++ {
			c.visited[i][j] = false
		}
	}
}

func (c *Canvas) getNeighbors(x int, y int) []Node {
	var (
		neighbors []Node
		color     byte
	)
	if x+1 < len(c.contents) {
		color = c.contents[x+1][y]
		neighbors = append(neighbors, Node{x + 1, y, color})
	}
	if x-1 >= 0 {
		color = c.contents[x-1][y]
		neighbors = append(neighbors, Node{x - 1, y, color})
	}
	if y+1 < len(c.contents[0]) {
		color = c.contents[x][y+1]
		neighbors = append(neighbors, Node{x, y + 1, color})
	}
	if y-1 >= 0 {
		color = c.contents[x][y-1]
		neighbors = append(neighbors, Node{x, y - 1, color})
	}
	return neighbors
}

func (c *Canvas) floodFill(x int, y int, color byte, originalColor byte, toVisit chan Node, visitDone chan bool) {
	c.contents[x][y] = color
	neighbors := c.getNeighbors(x, y)
	for _, neighbor := range neighbors {
		if neighbor.Color == originalColor {
			toVisit <- neighbor
		}
	}
	visitDone <- true
}

func (c *Canvas) FloodFill(x int, y int, color byte) {
	// If unbuffered, this channel will block when we go to send the
	// initial nodes to visit (at most 4).  Not cool man.
	toVisit := make(chan Node, 4)
	visitDone := make(chan bool)
	originalColor := c.contents[x][y]
	c.setVisitedMatrixToFalse()
	go c.floodFill(x, y, color, originalColor, toVisit, visitDone)
	remainingVisits := 1
	for {
		select {
		case nextVisit := <-toVisit:
			if !c.visited[nextVisit.X][nextVisit.Y] {
				c.visited[nextVisit.X][nextVisit.Y] = true
				remainingVisits++
				go c.floodFill(nextVisit.X, nextVisit.Y, color, originalColor, toVisit, visitDone)
			}
		case <-visitDone:
			remainingVisits--
		default:
			if remainingVisits == 0 {
				return
			}
		}
	}
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	canvas := Canvas{}
	canvas.Init(100, 120, '_')
	for i := 1; i < 7; i++ {
		for j := 8; j < 50; j++ {
			canvas.contents[i][j] = '/'
		}
	}
	for i := 80; i < 90; i++ {
		for j := 80; j < 99; j++ {
			canvas.contents[i][j] = '-'
		}
	}
	for i := 30; i < 70; i++ {
		for j := 25; j < 65; j++ {
			canvas.contents[i][j] = '\\'
		}
	}
	canvas.FloodFill(2, 3, 'G')
	canvas.Print()
}
