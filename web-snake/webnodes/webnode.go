package webnodes

import "main/game"

type WebNode struct {
	game *game.Game
}

func NewWebNode(g *game.Game) *WebNode {
	return &WebNode{
		game: g,
	}
}
