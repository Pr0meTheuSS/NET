// package main

// import (
// 	"math/rand"
// 	"time"

// 	"github.com/faiface/pixel"
// 	"github.com/faiface/pixel/pixelgl"
// 	"golang.org/x/image/colornames"
// )

// const (
// 	Width     = 800
// 	Height    = 600
// 	BlockSize = 20
// )

// type Snake struct {
// 	Body []pixel.Vec
// 	Vel  pixel.Vec
// }

// func (s *Snake) Move() {
// 	newHead := s.Body[0].Add(s.Vel)
// 	s.Body = append([]pixel.Vec{newHead}, s.Body[:len(s.Body)-1]...)
// }

// func main() {
// 	pixelgl.Run(run)
// }

// func drawCenteredGreenSquare(win *pixelgl.Window, x, y float64) {
// 	// centerX := x - 100/2
// 	// centerY := y - 100/2
// 	c := pixel.Circle{
// 		Center: pixel.Vec{X: 50, Y: 50},
// 		Radius: 10,
// 	}
// 	mat := pixel.IM
// 	pixel.NewSprite(c, pixel.Rect{Min: pixel.Vec{X: 0, Y: 0}, Max: pixel.Vec{X: 100, Y: 100}})
// 	win.Canvas().DrawColorMask(win, mat, colornames.Green)
// }

// func run() {
// 	cfg := pixelgl.WindowConfig{
// 		Title:  "Snake Game",
// 		Bounds: pixel.R(0, 0, Width, Height),
// 		VSync:  true,
// 	}
// 	win, err := pixelgl.NewWindow(cfg)
// 	if err != nil {
// 		panic(err)
// 	}

// 	snake := Snake{Body: []pixel.Vec{{X: Width / 2, Y: Height / 2}}, Vel: pixel.V(0, 0)}

// 	rand.NewSource(time.Now().UnixNano())
// 	food := pixel.V(rand.Float64()*Width, rand.Float64()*Height)

// 	for !win.Closed() {

// 		snake.Vel = pixel.ZV
// 		if win.Pressed(pixelgl.KeyLeft) {
// 			snake.Vel.X -= BlockSize
// 		}
// 		if win.Pressed(pixelgl.KeyRight) {
// 			snake.Vel.X += BlockSize
// 		}
// 		if win.Pressed(pixelgl.KeyUp) {
// 			snake.Vel.Y += BlockSize
// 		}
// 		if win.Pressed(pixelgl.KeyDown) {
// 			snake.Vel.Y -= BlockSize
// 		}

// 		snake.Move()

// 		if snake.Body[0].To(food).Len() < BlockSize {
// 			food = pixel.V(rand.Float64()*Width, rand.Float64()*Height)
// 			snake.Body = append(snake.Body, snake.Body[len(snake.Body)-1])
// 		}

// 		win.Clear(colornames.White)
// 		drawCenteredGreenSquare(win, 10, 10)
// 		// mat := pixel.IM.Scaled(pixel.ZV, BlockSize).Moved(pixel.V(100, 100))
// 		// win.Canvas().DrawColorMask(win, mat, colornames.Green)
// 		win.Update()
// 		time.Sleep(time.Second)
// 	}
// }

package main

import (
	"fmt"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Pixel Rocks!",
		Bounds: pixel.R(0, 0, 1024, 768),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	imd := imdraw.New(nil)

	i := 0
	prevTime := time.Now()
	for !win.Closed() {
		i++
		FPS := 1 / (time.Now().Sub(prevTime).Seconds())
		fmt.Println(FPS)
		prevTime = time.Now()
		win.Clear(colornames.Aliceblue)
		imd.Draw(win)
		win.Update()

		imd.Clear()
		step := float64(i)/100
		imd.Color = pixel.RGB(1-step, step, 0)
		imd.Push(pixel.V(200, 100))
		imd.Color = pixel.RGB(0, 1-step, step)
		imd.Push(pixel.V(800, 100))
		imd.Color = pixel.RGB(step, 0, 1-step)
		imd.Push(pixel.V(500, 700))
		imd.Polygon(0)
	}
}

func main() {
	pixelgl.Run(run)
}
