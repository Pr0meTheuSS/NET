package game

import (
	"errors"
	"log"
	"net"
	"snake/geometry"
	"snake/models"
	"snake/websnake_proto_gen/main/websnake"
	"strconv"
	"time"
)

type GameControllerImpl struct {
	games map[string]*Game
}

type Size struct {
	Width  int32
	Height int32
}

type Player struct {
	Id          int32
	Name        string
	IpAddress   string
	Port        int32
	Role        websnake.NodeRole
	Type        websnake.PlayerType
	Score       int32
	Snake       *models.Snake
	LastTimeout time.Time
}

type Message struct {
	GameMessage *websnake.GameMessage
	To          *net.UDPAddr
}
type Game struct {
	Name           string
	Players        map[int32]*Player
	GridSize       Size
	Delay          int32
	StaticFood     int32
	Food           []geometry.Position
	IsRun          bool
	MainPlayerID   *int32
	MessageChannel chan Message
	Sequence       int64
	// ConnectionChannel chan *websnake.GameAnnouncement
}

type GameConfig struct {
	BoardSize  Size
	StaticFood int32
	DelayInMs  int32
}

func GenerateRandomConfig() GameConfig {
	// TODO: add random factor.
	return GameConfig{
		BoardSize: Size{
			Width:  20,
			Height: 20,
		},
		StaticFood: 0,
		DelayInMs:  300,
	}
}

var idGenerator = int32(0)

func (g *Game) generateNewPlayerId() int32 {
	idGenerator++
	return idGenerator
}

func (g *Game) newPlayer(id int32, name string, playerType websnake.PlayerType, role websnake.NodeRole, ip string, port int32) *Player {
	newSnake, err := g.implaceSnake()
	if err != nil {
		return nil
	}

	return &Player{
		Id:          id,
		Name:        name,
		IpAddress:   ip,
		Port:        port,
		Role:        role,
		Type:        playerType,
		Score:       0,
		Snake:       newSnake,
		LastTimeout: time.Now(),
	}
}

func (g *Game) JoinPlayer(playerName string, playerType websnake.PlayerType, requestedRole websnake.NodeRole, addr *net.UDPAddr) {
	log.Println("****** Join new player into the game ******")

	newPlayerId := g.generateNewPlayerId()
	newPlayer := g.newPlayer(newPlayerId, playerName, playerType, requestedRole, addr.IP.String(), int32(addr.Port))
	g.Players[newPlayerId] = newPlayer
}

func (g *Game) Run() {
	log.Println("****** Run new game ******")
	g.IsRun = true
	go g.MainLoop()
}

func (g *Game) MainLoop() {
	for g.IsRun {
		<-time.NewTimer(time.Duration(g.Delay) * time.Millisecond).C

		for _, p := range g.Players {
			log.Printf("%+v\n", p)
		}

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
		sendGameState(g)
	}
}

var globalStateOrder = int64(0)

func sendGameState(game *Game) {
	if len(game.Players) < 2 {
		return
	}

	globalStateOrder++

	log.Println("Send gamestate----------------------------")
	defer log.Println("----------------------------Gamestate sent")

	// Рассылка по всем пользователям
	for _, player := range game.Players {
		to := net.UDPAddr{
			IP:   net.ParseIP(player.IpAddress),
			Port: int(player.Port),
			Zone: "",
		}

		// исключаем мастера из рассылки
		if player.Role != websnake.NodeRole_MASTER {
			msgToSend := Message{
				GameMessage: game.buildGameStateBytes(player.Id, int32(globalStateOrder)),
				To:          &to,
			}
			game.MessageChannel <- msgToSend
		}
	}
}
func (g *Game) generateSeq() int64 {
	g.Sequence++
	return g.Sequence
}

func (g *Game) buildGameStateBytes(receiverId int32, stateOrder int32) *websnake.GameMessage {
	seq := g.generateSeq()
	return &websnake.GameMessage{
		MsgSeq:     &seq,
		SenderId:   g.MainPlayerID,
		ReceiverId: &receiverId,
		Type: &websnake.GameMessage_State{
			State: &websnake.GameMessage_StateMsg{
				State: &websnake.GameState{
					StateOrder: &stateOrder,
					Snakes:     mapModelPlayersToNetSnakes(g),
					Foods:      mapModelFoodToNetFood(g),
					Players:    mapModelPlayersToNetPlayers(g),
				},
			},
		},
	}
}

func mapModelPlayersToNetPlayers(game *Game) *websnake.GamePlayers {
	netPlayers := []*websnake.GamePlayer{}
	for _, player := range game.Players {
		netPlayers = append(netPlayers, playerToNet(*player))
	}

	return &websnake.GamePlayers{
		Players: netPlayers,
	}
}

func playerToNet(player Player) *websnake.GamePlayer {
	return &websnake.GamePlayer{
		Name:      &player.Name,
		Id:        &player.Id,
		IpAddress: &player.IpAddress,
		Port:      &player.Port,
		Role:      &player.Role,
		Type:      &player.Type,
		Score:     &player.Score,
	}
}

func mapModelPlayersToNetSnakes(game *Game) []*websnake.GameState_Snake {
	netSnakes := []*websnake.GameState_Snake{}
	for _, player := range game.Players {
		// if player.Snake.IsAlive {
		netSnakes = append(netSnakes, playerToNetSnake(*player))
		// }
	}

	return netSnakes
}

var DirToNetDir = map[models.Direction]websnake.Direction{
	models.UP:    websnake.Direction_UP,
	models.DOWN:  websnake.Direction_DOWN,
	models.LEFT:  websnake.Direction_LEFT,
	models.RIGHT: websnake.Direction_RIGHT,
}

func playerToNetSnake(p Player) *websnake.GameState_Snake {
	netCoords := []*websnake.GameState_Coord{}
	webDir := DirToNetDir[p.Snake.Dir]

	snakeState := websnake.GameState_Snake_ALIVE
	if p.Snake.IsZombie {
		snakeState = websnake.GameState_Snake_ZOMBIE
	}

	if len(p.Snake.Body) == 0 {
		return &websnake.GameState_Snake{
			PlayerId:      &p.Id,
			Points:        netCoords,
			State:         &snakeState,
			HeadDirection: &webDir,
		}
	}

	netCoords = append(netCoords, &websnake.GameState_Coord{
		X: &p.Snake.Body[0].X,
		Y: &p.Snake.Body[0].Y,
	})

	for i, curr := range p.Snake.Body[1:] {
		x := curr.X - p.Snake.Body[i].X
		y := curr.Y - p.Snake.Body[i].Y
		netCoords = append(netCoords, &websnake.GameState_Coord{
			X: &x,
			Y: &y,
		})
	}

	return &websnake.GameState_Snake{
		PlayerId:      &p.Id,
		Points:        netCoords,
		State:         &snakeState,
		HeadDirection: &webDir,
	}

}

func mapModelFoodToNetFood(game *Game) []*websnake.GameState_Coord {
	netFood := []*websnake.GameState_Coord{}
	for i := range game.Food {
		netFood = append(netFood, &websnake.GameState_Coord{
			X: &game.Food[i].X,
			Y: &game.Food[i].Y,
		})
	}

	return netFood
}

func (g *Game) SteerPlayerSnake(playerId int32, dir models.Direction) {
	if player, ok := g.Players[playerId]; ok {
		player.Snake.SetDirection(dir)
	}
}

func (g *Game) AddFood() {
	if len(g.Food) >= len(g.getAlivePlayers())+int(g.StaticFood) {
		return
	}

	width := g.GridSize.Width
	height := g.GridSize.Height

	new := generateFood(width, height)

	for _, p := range g.Players {
		snake := p.Snake
		if geometry.Find(snake.Body, new) != -1 {
			g.AddFood()
			return
		}
	}

	if geometry.Find(g.Food, new) != -1 {
		g.AddFood()
		return
	}

	g.Food = append(g.Food, new)
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

func (g *Game) castBodyToFood(s *models.Snake) {
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

type IGameController interface {
	FindGameByName(string) (*Game, error)
	CreateNewGame(gameName string, config GameConfig) *Game
	GetAllGames() []*Game
}

func (gc *GameControllerImpl) GetAllGames() []*Game {
	gamesSlice := make([]*Game, len(gc.games))
	for _, game := range gc.games {
		gamesSlice = append(gamesSlice, game)
	}

	return gamesSlice
}

func (gc *GameControllerImpl) FindGameByName(gamename string) (*Game, error) {
	if game, ok := gc.games[gamename]; ok {
		return game, nil
	}

	return nil, errors.New("game with name " + gamename + " does not exist")
}

func (gc *GameControllerImpl) CreateNewGame(gameName string, config GameConfig) *Game {
	newGame := &Game{
		Name:           gameName,
		Players:        map[int32]*Player{},
		GridSize:       config.BoardSize,
		Delay:          config.DelayInMs,
		StaticFood:     config.StaticFood,
		Food:           []geometry.Position{},
		IsRun:          false,
		MainPlayerID:   new(int32),
		MessageChannel: make(chan Message),
	}

	for i := 0; i < int(newGame.StaticFood); i++ {
		newGame.AddFood()
	}

	return newGame
}

func (g *Game) implaceSnake() (*models.Snake, error) {
	snakePos, err := g.findFreePlace()
	if err != nil {
		return nil, err
	}

	return models.NewSnake(g.GridSize.Width, g.GridSize.Height, snakePos.X, snakePos.Y), nil
}

func (g *Game) UpdatePlayerByAddress(address *net.UDPAddr) {
	for _, p := range g.Players {
		if net.JoinHostPort(p.IpAddress, strconv.Itoa(int(p.Port))) == net.JoinHostPort(address.IP.String(), strconv.Itoa(address.Port)) {
			p.LastTimeout = time.Now()
		}
	}
}

func (g *Game) SteerPlayerByAddress(address *net.UDPAddr, dir models.Direction) {
	for _, p := range g.Players {
		if net.JoinHostPort(p.IpAddress, strconv.Itoa(int(p.Port))) == net.JoinHostPort(address.IP.String(), strconv.Itoa(address.Port)) {
			g.SteerPlayerSnake(p.Id, dir)
		}
	}
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
	// log.Println("----------------------------------Try to find free place for snake")
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
