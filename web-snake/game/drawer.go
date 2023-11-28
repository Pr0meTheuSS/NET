package game

import (
	"fmt"
	"image/color"
	"main/geometry"
	"main/snake"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func GetGridToDraw(gridSize, winSize Size) []fyne.CanvasObject {
	cellSize := fyne.NewSize(float32(winSize.Width/gridSize.Width), float32(winSize.Height/gridSize.Height))

	lines := make([]fyne.CanvasObject, 0, int(gridSize.Width)*2)
	// Draw vertical lines.
	for x := int32(1); x < int32(gridSize.Width); x++ {
		line := &canvas.Line{
			Position1:   fyne.NewPos(float32(x)*float32(cellSize.Width), 0.0),
			Position2:   fyne.NewPos(float32(x)*float32(cellSize.Height), float32(winHeight)),
			Hidden:      false,
			StrokeColor: color.Color(color.Gray{}),
			StrokeWidth: 2,
		}

		lines = append(lines, line)
	}

	// Draw horizontal lines.
	for y := int32(1); y < int32(gridSize.Height); y++ {
		line := &canvas.Line{
			Position1:   fyne.NewPos(0.0, float32(y)*float32(cellSize.Height)),
			Position2:   fyne.NewPos(float32(winWidth), float32(y)*float32(cellSize.Height)),
			Hidden:      false,
			StrokeColor: color.Color(color.Gray{}),
			StrokeWidth: 2,
		}

		lines = append(lines, line)
	}

	return lines
}

func GetSnakeToDraw(s snake.Snake, food []geometry.Position, winSize, gridSize Size, clr color.Color) []fyne.CanvasObject {
	cellSize := fyne.NewSize(float32(winSize.Width/gridSize.Width), float32(winSize.Height/gridSize.Height))

	snakeCells := []fyne.CanvasObject{}
	for i, cell := range s.Body {
		cellToDraw := canvas.NewRectangle(clr)
		if i == 0 {
			cellToDraw.FillColor = color.Gray{}
		}

		cellToDraw.Resize(cellSize)
		cellToDraw.Move(fyne.NewPos(float32(cell.X)*cellSize.Width, float32(cell.Y)*cellSize.Height))
		snakeCells = append(snakeCells, cellToDraw)
	}

	for _, foodCell := range food {
		cellToDraw := canvas.NewRectangle(color.RGBA{
			R: 255,
			G: 0,
			B: 0,
			A: 255,
		})
		cellToDraw.Resize(fyne.NewSize(cellSize.Width, cellSize.Height))
		cellToDraw.Move(fyne.NewPos(float32(foodCell.X)*cellSize.Width, float32(foodCell.Y)*cellSize.Height))
		snakeCells = append(snakeCells, cellToDraw)
	}

	return snakeCells
}

func Draw(g *Game) {

	objsToDraw := []fyne.CanvasObject{}
	objsToDraw = append(objsToDraw, GetGridToDraw(g.GridSize, g.WinSize)...)

	for i, p := range g.Players {
		snakeColor := color.RGBA{
			R: 128,
			G: 128,
			B: 128,
			A: 255,
		}

		if g.mainPlayerId == p.Id {
			snakeColor = color.RGBA{
				R: 0,
				G: 255,
				B: 0,
				A: 255,
			}
		}
		objsToDraw = append(objsToDraw, GetSnakeToDraw(*p.Snake, g.Food, g.WinSize, g.GridSize, snakeColor)...)

		scoreLabel := widget.NewLabel("Score for " + strconv.Itoa(int(p.Id)) + " " + p.Name + ":" + fmt.Sprint(p.Score))
		scoreLabel.Move(fyne.NewPos(600.0, float32(i)*150.0))
		scoreLabel.Resize(fyne.NewSize(200.0, 150.0))
		objsToDraw = append(objsToDraw, scoreLabel)
	}

	g.Window.SetContent(container.NewWithoutLayout(objsToDraw...))
}
