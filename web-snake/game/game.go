package game

import (
	"errors"
	"log"
	"main/geometry"
	"main/pubsub"
	"main/snake"
	"main/websnake"
	"time"

	"fyne.io/fyne/v2"
)

const (
	winWidth  = 600.0
	winHeight = 600.0
)

// var uname = ""

type Size struct {
	Width, Height int32
}

type Player struct {
	Name        string
	Id          int32
	IpAddress   string
	Port        int32
	Role        websnake.NodeRole
	Type        websnake.PlayerType
	Score       int32
	Snake       *snake.Snake
	LastTimeout time.Time
}

func NewGame(gamename string, win fyne.Window, plrs []Player, winsize Size, gridsize Size, delay int32, staticfood int32, food []geometry.Position) *Game {
	players := map[int32]*Player{}
	for _, p := range plrs {
		players[p.Id] = &p
	}

	newGame := &Game{
		Name:              gamename,
		Players:           players,
		WinSize:           winsize,
		GridSize:          gridsize,
		Delay:             delay,
		StaticFood:        staticfood,
		Food:              food,
		Window:            win,
		IsRun:             true,
		MainPlayerID:      new(int32),
		ConnectionChannel: make(chan *websnake.GameAnnouncement),
	}

	for i := int32(0); i < newGame.StaticFood; i++ {
		newGame.AddFood()
	}

	return newGame
}

func (g *Game) ConnectToTheGame(v *websnake.GameAnnouncement) {
	log.Println("send game announce to game channel")
	g.ConnectionChannel <- v
}

func (g *Game) SteerPlayerSnake(playerId int32, dir snake.Direction) {
	if player, ok := g.Players[playerId]; ok {
		player.Snake.SetDirection(dir)
	}
}

func (g *Game) UpdateUserTimeout(userId int32, lastTime time.Time) {
	if player, ok := g.Players[userId]; ok {
		player.LastTimeout = lastTime
	}
}

func (g *Game) UpdateGameState(message pubsub.Message) {
	updatedGame := g.GameStateToGame(message.Msg.GetState().GetState())
	g.Players = updatedGame.Players
	g.Food = updatedGame.Food
	Draw(g)
}

func GamePlayerToPlayer(gamePlayer *websnake.GamePlayer) *Player {
	return &Player{
		Name:      *gamePlayer.Name,
		Id:        *gamePlayer.Id,
		IpAddress: *gamePlayer.IpAddress,
		Port:      *gamePlayer.Port,
		Role:      *gamePlayer.Role,
		Type:      *gamePlayer.Type,
		Score:     *gamePlayer.Score,
	}
}

var NetDirToModel = map[websnake.Direction]snake.Direction{
	websnake.Direction_DOWN:  snake.DOWN,
	websnake.Direction_UP:    snake.UP,
	websnake.Direction_LEFT:  snake.LEFT,
	websnake.Direction_RIGHT: snake.RIGHT,
}

// GameState_Snake to Snake
func (g *Game) GameStateSnakeToSnake(gsSnake *websnake.GameState_Snake) *snake.Snake {
	body := []geometry.Position{}
	if len(gsSnake.Points) == 0 {
		return &snake.Snake{
			Body:     body,
			Dir:      NetDirToModel[*gsSnake.HeadDirection],
			IsZombie: gsSnake.GetState() == websnake.GameState_Snake_ZOMBIE,
			IsAlive:  true,
		}
	}

	head := geometry.Position{
		X: gsSnake.Points[0].GetX(),
		Y: gsSnake.Points[0].GetY(),
	}

	body = append(body, head)

	for _, curDirection := range gsSnake.GetPoints()[1:] {
		prev := body[len(body)-1]
		body = append(body, geometry.Position{
			X: (prev.X + curDirection.GetX() + g.GridSize.Width) % g.GridSize.Width,
			Y: (prev.Y + curDirection.GetY() + g.GridSize.Height) % g.GridSize.Height,
		})
	}

	return &snake.Snake{
		Body:     body,
		Dir:      NetDirToModel[*gsSnake.HeadDirection],
		IsZombie: gsSnake.GetState() == websnake.GameState_Snake_ZOMBIE,
		IsAlive:  true,
	}
}

func (g *Game) GameStateToGame(gs *websnake.GameState) Game {
	gamePlayers := gs.Players.Players
	players := map[int32]*Player{}
	gameSnakes := gs.Snakes

	for _, gp := range gamePlayers {
		currPlayer := GamePlayerToPlayer(gp)
		players[currPlayer.Id] = currPlayer
	}

	for _, gameSnake := range gameSnakes {
		newSnake := g.GameStateSnakeToSnake(gameSnake)
		players[*gameSnake.PlayerId].Snake = newSnake
	}

	foods := []geometry.Position{}
	for _, coord := range gs.Foods {
		foods = append(foods, geometry.Position{
			X: *coord.X,
			Y: *coord.Y,
		})
	}

	return Game{
		Players: players,
		Food:    foods,
	}
}

type Game struct {
	Name              string
	Players           map[int32]*Player
	WinSize           Size
	GridSize          Size
	Delay             int32
	StaticFood        int32
	Food              []geometry.Position
	Window            fyne.Window
	IsRun             bool
	MainPlayerID      *int32
	ConnectionChannel chan *websnake.GameAnnouncement
}

func (g *Game) Run() {
	g.IsRun = true
}

func (g *Game) SetMainPlayer(id int32) {
	*g.MainPlayerID = id
}

func (g *Game) Close() {
	g.IsRun = false
}

func (g *Game) GetMainPlayer() *Player {
	if len(g.Players) == 1 {
		for _, p := range g.Players {
			return p
		}
	}

	if player, ok := g.Players[*g.MainPlayerID]; ok {
		return player
	}

	return nil
}

func (g *Game) AddPlayer(username, ipAddress string, port int, role websnake.NodeRole, tp websnake.PlayerType) (*Player, error) {
	newSnake, err := g.ImplaceSnake()
	if err != nil {
		return nil, err
	}

	newId := generatePlayerId()
	for p, ok := g.Players[newId]; ok && p != nil; newId = generatePlayerId() {
		p, ok = g.Players[newId]
		log.Println(p, ok)
		log.Println(len(g.Players))
		log.Println(newId)
	}

	newPlayer := &Player{
		Name:        username,
		Id:          newId,
		IpAddress:   ipAddress,
		Port:        int32(port),
		Role:        role,
		Type:        tp,
		Score:       0,
		Snake:       newSnake,
		LastTimeout: time.Now(),
	}

	g.Players[newPlayer.Id] = newPlayer
	// каждому игроку по еде.
	g.AddFood()
	return newPlayer, nil
}

func (g *Game) AddMainPlayer(username, ipAddress string, port int, role websnake.NodeRole, tp websnake.PlayerType) error {
	newSnake, err := g.ImplaceSnake()
	if err != nil {
		return err
	}

	newPlayer := &Player{
		Name:        username,
		Id:          int32(generatePlayerId()),
		IpAddress:   ipAddress,
		Port:        int32(port),
		Role:        role,
		Type:        tp,
		Score:       0,
		Snake:       newSnake,
		LastTimeout: time.Now(),
	}

	g.Players[newPlayer.Id] = newPlayer
	*g.MainPlayerID = newPlayer.Id
	return nil
}

func HandleUserInput(ke *fyne.KeyEvent, s *snake.Snake) {
	log.Println("Handle user input")
	keyToDir := map[fyne.KeyName]snake.Direction{
		fyne.KeyW: snake.UP,
		fyne.KeyD: snake.RIGHT,
		fyne.KeyS: snake.DOWN,
		fyne.KeyA: snake.LEFT,
	}

	if newdir, ok := keyToDir[ke.Name]; ok {
		s.SetDirection(newdir)
	}
}

func (g *Game) handleJoin(message pubsub.Message) {
	msg := message.Msg

	joinMsg := msg.GetJoin()
	log.Println("Catch join msg: ", joinMsg)
	log.Println("From: ", message.From)
	log.Println("To: ", message.To)

	newPlayer, err := g.AddPlayer(
		*joinMsg.PlayerName,
		message.From.IP.String(),
		message.From.Port,
		*joinMsg.RequestedRole,
		*joinMsg.PlayerType)

	if err != nil {
		// Отправляем сообщение с ошибкой с замененными отправителем и получателем.
		errorMessage := err.Error()
		pubsub.GetGlobalPubSubService().Publish("senderror", pubsub.Message{
			Msg: &websnake.GameMessage{
				MsgSeq:     new(int64),
				SenderId:   new(int32),
				ReceiverId: new(int32),
				Type: &websnake.GameMessage_Error{
					Error: &websnake.GameMessage_ErrorMsg{
						ErrorMessage: &errorMessage,
					},
				},
			},
			To:   message.From,
			From: message.To,
		})
		return
	}

	pubsub.GetGlobalPubSubService().Publish("sendack", pubsub.Message{
		Msg: &websnake.GameMessage{
			MsgSeq:     new(int64),
			SenderId:   new(int32),
			ReceiverId: &newPlayer.Id,
			Type:       &websnake.GameMessage_Ack{Ack: &websnake.GameMessage_AckMsg{}},
		},
		To:   message.From,
		From: message.To,
	})
}

func (g *Game) AddFood() {
	if len(g.Food) >= len(g.getAlivePlayers())+int(g.StaticFood) {
		return
	}

	width := g.GridSize.Width
	height := g.GridSize.Height

	new := generateFood(width, height)
	for geometry.Find(g.Food, new) != -1 {
		new = generateFood(width, height)
	}
	for _, p := range g.Players {
		snake := p.Snake
		for geometry.Find(snake.Body, new) != -1 {
			new = generateFood(width, height)
		}
	}

	g.Food = append(g.Food, new)
}

var globalPlayerCounter = int32(-1)

func generatePlayerId() int32 {
	log.Println("Generate new player id", globalPlayerCounter)
	globalPlayerCounter++
	return globalPlayerCounter
}

func (g *Game) isSquareFree(squareSize, x, y int32) bool {
	for i := x; i < x+squareSize; i = (i + 1) % g.GridSize.Width {
		for j := y; j < y+squareSize; j = (j + 1) % g.GridSize.Height {
			for _, p := range g.Players {
				if geometry.Find(p.Snake.Body, geometry.Position{X: i, Y: j}) != -1 {
					return false
				}
			}

			if geometry.Find(g.Food, geometry.Position{X: i, Y: j}) != -1 {
				return false
			}
		}
	}

	return true
}

func (g *Game) findFreePlace() (*geometry.Position, error) {
	log.Println("----------------------------------Try to find free place for snake")
	squareSize := int32(5)
	for x := int32(0); x < g.GridSize.Width; x++ {
		for y := int32(0); y < g.GridSize.Height; y++ {
			if g.isSquareFree(squareSize, x, y) {
				return &geometry.Position{
					X: x,
					Y: y,
				}, nil
			}
		}
	}
	return nil, errors.New("cannot find free place")
}

func (g *Game) ImplaceSnake() (*snake.Snake, error) {
	snakePos, err := g.findFreePlace()
	if err != nil {
		return nil, err
	}

	return snake.NewSnake(g.GridSize.Width, g.GridSize.Height, snakePos.X, snakePos.Y), nil
}

func (g Game) getAlivePlayers() []*Player {
	alivePlayers := []*Player{}
	for _, player := range g.Players {
		if player.Snake.IsAlive {
			alivePlayers = append(alivePlayers, player)
		}
	}

	return alivePlayers
}

func (g *Game) MainLoop() {
	for g.IsRun {
		log.Println("++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++")
		for _, p := range g.Players {
			log.Println(*p, "Timeout:", time.Since(p.LastTimeout).Milliseconds())
		}
		log.Println("++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++")

		<-time.NewTimer(time.Duration(g.Delay) * time.Millisecond).C
		log.Println(g.Players)
		for _, p := range g.Players {
			currentSnake := p.Snake
			currentSnake.Move()

			// Голова пересеклась с едой
			snakeHead := currentSnake.Head()
			if snakeHead == nil {
				continue
			}

			if catchedFoodPosition := geometry.Find(g.Food, *snakeHead); catchedFoodPosition != -1 {
				p.Score++
				g.Food = append(g.Food[0:catchedFoodPosition], g.Food[catchedFoodPosition+1:]...)
				g.AddFood()
				currentSnake.GrowUp()
			}
		}

		alivePlayers := g.getAlivePlayers()
		// TODO: Голова пересеклась с телом другой змеи
		for _, p := range alivePlayers {
			for _, player := range alivePlayers {
				if p != player {
					// Если в теле другой змеи найдётся голова текущей змеи
					snakeHead := p.Snake.Head()
					if snakeHead == nil {
						continue
					}

					if geometry.Find(player.Snake.Body, *snakeHead) != -1 {
						p.Snake.IsAlive = false
						// Накидываем убийце +1
						player.Score++
						g.castBodyToFood(p.Snake)
					}
				}
			}
		}

		Draw(g)
		// отправляем пустое сообщение, пока что так
		// формирование сообщения реализовано на уровне сети.
		pubsub.GetGlobalPubSubService().Publish("sendgamestate", pubsub.Message{})
	}
}

func (g *Game) castBodyToFood(s *snake.Snake) {
	for _, cell := range s.Body[1:] {
		if random(2)%2 == 0 {
			g.Food = append(g.Food, cell)
		}
	}

	s.Body = []geometry.Position{}
}

func generateFood(width, height int32) geometry.Position {
	return geometry.Position{
		X: random(width),
		Y: random(height),
	}
}

func random(max int32) int32 {
	return int32(time.Now().Nanosecond()) % max
}
