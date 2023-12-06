package webnodes

import (
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
