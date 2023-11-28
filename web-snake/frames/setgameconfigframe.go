package frames

import (
	"errors"
	"log"
	"main/game"
	"main/webnodes"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func createAndShowSetGameConfigFrame(app fyne.App, username string) {
	configWindow := app.NewWindow("Конфигурация новой игры")
	configWindow.Resize(fyne.NewSize(400, 400))
	configWindow.CenterOnScreen()
	configWindow.SetContent(InitSetGameConfigWindowContent(app, username))

	configWindow.Show()
}

func InitSetGameConfigWindowContent(app fyne.App, username string) *fyne.Container {
	// // Fields for input
	widthEntry := widget.NewEntry()
	widthEntry.Text = "20"
	heightEntry := widget.NewEntry()
	heightEntry.Text = "20"
	foodEntry := widget.NewEntry()
	foodEntry.Text = "2"
	nameEntry := widget.NewEntry()
	nameEntry.Text = "Game"
	speedEntry := widget.NewEntry()
	speedEntry.Text = "300"

	// Set validation functions for each entry
	widthEntry.Validator = sizeParamsValidator
	heightEntry.Validator = sizeParamsValidator

	// Create form container
	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Ширина", Widget: widthEntry},
			{Text: "Высота", Widget: heightEntry},
			{Text: "Еда", Widget: foodEntry},
			{Text: "Название игры", Widget: nameEntry},
			{Text: "Скорость игры (ms)", Widget: speedEntry},
		},
		OnSubmit: func() {
			width, err := strconv.ParseInt(widthEntry.Text, 10, 32)
			if err != nil {
				log.Fatal(err)
			}
			height, err := strconv.ParseInt(widthEntry.Text, 10, 32)
			if err != nil {
				log.Fatal(err)
			}
			food, err := strconv.ParseInt(foodEntry.Text, 10, 32)
			if err != nil {
				log.Fatal(err)
			}
			delay, err := strconv.ParseInt(speedEntry.Text, 10, 32)
			if err != nil {
				log.Fatal(err)
			}

			newGame := game.CreateGame(app, username, nameEntry.Text, int32(width), int32(height), int32(food), int32(delay))
			webNode := webnodes.NewWebSnakeMasterNode(newGame)
			defer webNode.DestroyNode()

			go webNode.SendMultiAnnouncment()
			go webNode.ListenAndServe()
			newGame.MainLoop()
		},
	}

	// Combine form container
	return container.NewVBox(form)
}

func sizeParamsValidator(s string) error {
	if s == "" {
		return nil // Empty input is considered valid
	}

	// Try to parse the input as an integer
	val, err := strconv.Atoi(s)
	if err != nil {
		return errors.New("width must be an integer")
	}

	if val > 100 || val < 0 {
		return errors.New("width must be a positive integer between 10 and 100")
	}
	return nil
}
