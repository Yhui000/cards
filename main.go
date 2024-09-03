package main

import (
	"cards/cards"

	"github.com/gin-gonic/gin"
)

var b *cards.Board

func main() {
	// ? Game Setup
	b = &cards.Board{
		PlayerTurn:        true,
		TurnCount:         0,
		Players:           map[bool]*cards.Player{true: cards.PaladinControlZoth(), false: cards.WarlockMurloc()},
		ActionChan:        make(chan *cards.Action),
		WaitingActionChan: make(chan *bool, 1),
		ActionEndChan:     make(chan error, 1),
	}

	go b.Start()

	// ? Server Setup
	eng := gin.Default()
	eng.Any("/ws/connect", Connect)
	eng.StaticFile("/", "./frontend/dist/index.html")
	eng.Static("/assets", "./frontend/dist/assets/")
	eng.Static("/imgs", "./imgs/")
	eng.StaticFile("/favicon.ico", "./frontend/dist/favicon.ico")
	eng.Run(":8123")
}
