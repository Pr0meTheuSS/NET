package snake

import (
	"main/geometry"
	"math/rand"
	"time"
)

type Direction int

// Константы подобраны так,
// чтобы каждое направление было уникальным и противоположные направления в сумме давали 0.
const (
	UP    = 1
	RIGHT = 2
	DOWN  = -1
	LEFT  = -2
)

type Snake struct {
	Body     []geometry.Position
	prevTail geometry.Position
	Dir      Direction
	IsAlive  bool
	score    int
	IsZombie bool
}

func getRandomDirection() int {
	rand.NewSource(time.Now().UnixNano())
	directions := []int{UP, RIGHT, DOWN, LEFT}
	randomIndex := rand.Intn(len(directions))
	return directions[randomIndex]
}

var globalBoardWidth = int32(0)
var globalBoardHeight = int32(0)

func createTail(head geometry.Position, dir Direction) geometry.Position {
	switch dir {
	case UP:
		return geometry.Position{X: head.X, Y: head.Y + 1}
	case LEFT:
		return geometry.Position{X: head.X + 1, Y: head.Y}
	case RIGHT:
		return geometry.Position{X: head.X - 1, Y: head.Y}
	case DOWN:
		return geometry.Position{X: head.X, Y: head.Y - 1}
	default:
		return head
	}
}

func NewSnake(boardWidth, boardHeight, snakePosX, snakePosY int32) *Snake {
	globalBoardWidth = boardWidth
	globalBoardHeight = boardHeight
	head := geometry.Position{X: snakePosX, Y: snakePosY}
	dir := Direction(getRandomDirection())
	tail := createTail(head, dir)

	snakeGame := &Snake{
		Body:     []geometry.Position{head, tail},
		prevTail: tail,
		Dir:      Direction(dir),
		IsAlive:  true,
		score:    0,
		IsZombie: false,
	}

	return snakeGame
}

func (s *Snake) GrowUp() {
	s.Body = append(s.Body, s.prevTail)
}

func (s *Snake) IsSnakeAlive() bool {
	return s.IsAlive
}

func (s *Snake) Head() *geometry.Position {
	if len(s.Body) == 0 {
		return nil
	}

	return &s.Body[0]
}

func (s *Snake) Move() {
	if !s.IsAlive {
		return
	}

	head := s.Body[0]
	s.prevTail = s.Body[len(s.Body)-1]

	switch s.Dir {
	case UP:
		head.Y = (head.Y - 1 + globalBoardHeight) % globalBoardHeight
	case DOWN:
		head.Y = (head.Y + 1) % globalBoardHeight
	case LEFT:
		head.X = (head.X - 1 + globalBoardWidth) % globalBoardWidth
	case RIGHT:
		head.X = (head.X + 1) % globalBoardWidth
	}

	// Если произошло столкновение с собой.
	if geometry.Find(s.Body[1:], head) != -1 {
		s.IsAlive = false
		return
	}

	s.Body = s.Body[:len(s.Body)-1]
	s.Body = append([]geometry.Position{head}, s.Body...)
}

func (s *Snake) GetScore() int {
	return s.score
}

func (s *Snake) isOppositeDirection(dir Direction) bool {
	return s.Dir+dir == 0
}

func (s *Snake) SetDirection(newdir Direction) {
	if !s.isOppositeDirection(newdir) {
		s.Dir = newdir
	}
}
