package webnodes

import (
	"log"
	"main/game"
	"main/snake"
	"main/websnake"
)

var DirToNetDir = map[snake.Direction]websnake.Direction{
	snake.UP:    websnake.Direction_UP,
	snake.DOWN:  websnake.Direction_DOWN,
	snake.LEFT:  websnake.Direction_LEFT,
	snake.RIGHT: websnake.Direction_RIGHT,
}

func playerToNetSnake(p game.Player) *websnake.GameState_Snake {
	netCoords := []*websnake.GameState_Coord{}
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

	webDir := DirToNetDir[p.Snake.Dir]

	log.Printf("Snake coords: %+v", netCoords)
	return &websnake.GameState_Snake{
		PlayerId: &p.Id,
		Points:   netCoords,
		// TODO: сейчас змея всегда живая, режим зомби не реализован
		State:         websnake.GameState_Snake_ALIVE.Enum(),
		HeadDirection: &webDir,
	}

}

func playerToNet(player game.Player) *websnake.GamePlayer {
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
