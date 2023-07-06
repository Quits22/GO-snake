package game

/*
this is a simple snake game using termbox
*/
import (
	"fmt"
	"math/rand"
	"time"

	"github.com/nsf/termbox-go"
)

// Constants defining the game dimensions
// Change at your own risk lol jk have fun
const (
	Width              = 40
	Height             = 20
	SnakeInitialLength = 3
	FoodChar           = '@'
	SnakeChar          = 'O'
	Speed              = time.Millisecond * 10
)

// Direction represents the current movement direction of the snake
type Direction int

// Possible movement directions

const (
	Up Direction = iota
	Down
	Left
	Right
)

// Point represents a point in the game grid
type Point struct {
	X, Y int
}

// Game represents the state of the game
type Game struct {
	Snake     []Point
	Food      Point
	Direction Direction
	IsOver    bool
	Score     int
	Name      string
}

// Initialize the game states
func (g *Game) Init(name string) {
	g.Snake = make([]Point, SnakeInitialLength)
	g.Snake[0] = Point{Width / 2, Height / 2}
	g.Food = g.generateFood()
	g.Direction = Right
	g.IsOver = false
	g.Score = 0
	g.Name = name
}

// Generate a new random point for the food
func (g *Game) generateFood() Point {

	var food Point
	for {
		food = Point{rand.Intn(Width), rand.Intn(Height)}
		if !g.isPointOnSnake(food) {
			break
		}
	}

	return food
}

// Check if the point is on the snake's body
func (g *Game) isPointOnSnake(point Point) bool {
	for _, p := range g.Snake {
		if p == point {
			return true
		}
	}
	return false
}

// Update the game state on each frame
func (g *Game) Update() {
	head := g.Snake[0]

	// Update the head position based on the current direction
	switch g.Direction {
	case Up:
		head.Y--
	case Down:
		head.Y++
	case Left:
		head.X--
	case Right:
		head.X++
	}

	// Check if the snake hit itself
	if g.isPointOnSnake(head) {
		g.IsOver = true
		return
	}
	// If the snake reached a border bring it to the other side
	if head.X < 0 {
		head.X = Width
	}

	if head.X > Width {
		head.X = 0
	}

	if head.Y < 0 {
		head.Y = Height
	}

	if head.Y > Height {
		head.Y = 0
	}
	//create the snake body i don't even know how it works anymore :/
	g.Snake = append([]Point{head}, g.Snake[:len(g.Snake)-1]...)

	// Check if the snake ate food
	if head == g.Food {
		g.Score++
		g.Snake = append(g.Snake, g.Snake[len(g.Snake)-1])
		g.Food = g.generateFood()
	}
}

// Render the game state on the terminal with termbox
func (g *Game) Render() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	// drawing the borders
	// Draw the top border
	for i := 0; i < Width+2; i++ {
		termbox.SetCell(i, 0, '─', termbox.ColorDefault, termbox.ColorDefault)
	}

	// Draw the bottom border
	for i := 0; i < Width+2; i++ {
		termbox.SetCell(i, Height+1, '─', termbox.ColorDefault, termbox.ColorDefault)
	}

	// Draw the side borders
	for i := 1; i < Height+1; i++ {
		termbox.SetCell(0, i, '│', termbox.ColorDefault, termbox.ColorDefault)
		termbox.SetCell(Width+1, i, '│', termbox.ColorDefault, termbox.ColorDefault)
	}

	// Draw the top-left corner
	termbox.SetCell(0, 0, '┌', termbox.ColorDefault, termbox.ColorDefault)

	// Draw the top-right corner
	termbox.SetCell(Width+1, 0, '┐', termbox.ColorDefault, termbox.ColorDefault)

	// Draw the bottom-left corner
	termbox.SetCell(0, Height+1, '└', termbox.ColorDefault, termbox.ColorDefault)

	// Draw the bottom-right corner
	termbox.SetCell(Width+1, Height+1, '┘', termbox.ColorDefault, termbox.ColorDefault)

	// Render the snake
	for i, p := range g.Snake {
		if i == 0 {
			termbox.SetCell(p.X, p.Y, SnakeChar, termbox.ColorDefault, termbox.ColorYellow)
		} else {
			termbox.SetCell(p.X, p.Y, SnakeChar, termbox.ColorDefault, termbox.ColorGreen)
		}
	}

	// Render the food
	termbox.SetCell(g.Food.X, g.Food.Y, FoodChar, termbox.ColorDefault, termbox.ColorRed)

	// Render the score
	scoreText := fmt.Sprintf("Score: %d", g.Score)
	for i, ch := range scoreText {
		termbox.SetCell(i, Height+2, ch, termbox.ColorDefault, termbox.ColorMagenta)
	}
	// Rener the name
	nameText := fmt.Sprintf("Name: %s", g.Name)
	for i, ch := range nameText {
		termbox.SetCell(i, Height+3, ch, termbox.ColorDefault, termbox.ColorMagenta)
	}

	termbox.Flush()
}

// this algorithm is super simple and will easily hit itself
func (g *Game) snakeAI() {
	snakeHead := g.Snake[0]
	food := g.Food

	// Determine the horizontal and vertical distances to the food
	dx := food.X - snakeHead.X
	dy := food.Y - snakeHead.Y

	// Adjust the direction based on the distance to the food, while avoiding the opposite direction and collisions with the snake's body
	if dx > 0 && g.Direction != Left && !g.collidesWithSelf(snakeHead.X+1, snakeHead.Y) {
		g.Direction = Right
	} else if dx < 0 && g.Direction != Right && !g.collidesWithSelf(snakeHead.X-1, snakeHead.Y) {
		g.Direction = Left
	} else if dy > 0 && g.Direction != Up && !g.collidesWithSelf(snakeHead.X, snakeHead.Y+1) {
		g.Direction = Down
	} else if dy < 0 && g.Direction != Down && !g.collidesWithSelf(snakeHead.X, snakeHead.Y-1) {
		g.Direction = Up
	}
}

// dosen't really do much
func (g *Game) collidesWithSelf(x, y int) bool {
	for _, segment := range g.Snake[1:] {
		if segment.X == x && segment.Y == y {
			return true
		}
	}
	return false
}

// Handle keyboard events for player control
func (g *Game) HandleEvents() {
	switch ev := termbox.PollEvent(); ev.Type {
	case termbox.EventKey:
		switch ev.Key {
		case termbox.KeyArrowUp:
			if g.Direction != Down {
				g.Direction = Up
			}
		case termbox.KeyArrowDown:
			if g.Direction != Up {
				g.Direction = Down
			}
		case termbox.KeyArrowLeft:
			if g.Direction != Right {
				g.Direction = Left
			}
		case termbox.KeyArrowRight:
			if g.Direction != Left {
				g.Direction = Right
			}
		case termbox.KeyEsc:
			g.IsOver = true
		}
	case termbox.EventError:
		panic(ev.Err)
	}
}

// start the game and return palyer score data after gameover
func Start(name string) (int, string, string) {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	game := Game{}
	game.Init(name)
	fmt.Print("Would you like to play or let the AI do it for you? Enter 1 for player control, or any other value for AI control: ")
	var response int
	fmt.Scan(&response)

	// Define the control function based on user response
	var controlFunc func()
	if response == 1 {
		controlFunc = game.HandleEvents
	} else {
		controlFunc = game.snakeAI
	}

	// Game loop
	for !game.IsOver {
		game.Render()
		game.Update()
		controlFunc()
		time.Sleep(Speed)
	}
	//get the current time not needed the database does it too
	currentTime := time.Now()
	mysqlDateTime := currentTime.Format("2006-01-02 15:04:05")

	return game.Score, game.Name, mysqlDateTime
}
