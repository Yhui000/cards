package main

import (
	"cards/cards"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Player struct {
	id              bool
	Conn            *websocket.Conn
	UpdateBoardChan chan bool
	ResultsChan     chan bool
}

type Match struct {
	players map[bool]*Player
	board   *cards.Board
}

func (m *Match) NewPlayer(conn *websocket.Conn) *Player {
	for i, p := range m.players {
		if p == nil {
			p = &Player{
				id:              i,
				Conn:            conn,
				UpdateBoardChan: make(chan bool, 1),
				ResultsChan:     make(chan bool, 1),
			}
			m.players[i] = p
			return p
		}
	}
	return nil
}

func (m *Match) DisconnectPlayer(id bool) {
	if m.players[id] != nil && m.players[id].Conn != nil {
		m.players[id].Conn.Close()
	}
	m.players[id] = nil
}

func (m *Match) BroadcastBoard() {
	for _, p := range m.players {
		if p != nil {
			p.Conn.WriteJSON(map[string]any{
				"type": "board",
				"data": TranslateBoard(m.board, p.id),
			})
		}
	}
}

func (m *Match) SendBoard(id bool) {
	p := m.players[id]
	if p != nil {
		p.Conn.WriteJSON(map[string]any{
			"type": "board",
			"data": TranslateBoard(m.board, p.id),
		})
	}
}

func (m *Match) SendError(id bool, err error) {
	p := m.players[id]
	log.Println("sendind error")
	if p != nil {
		p.Conn.WriteJSON(map[string]any{
			"type": "error",
			"data": err.Error(),
		})
	}
}

func (m *Match) UpdatePlayersBoard() {
	for _, p := range m.players {
		if p != nil && len(p.UpdateBoardChan) == 0 {
			p.UpdateBoardChan <- true
		}
	}
}

var match Match = Match{board: b, players: map[bool]*Player{true: nil, false: nil}}

func Connect(ctx *gin.Context) {
	// TODO: rework this to find a match in a list of matches (only when the project be about to be finished)
	if match.board == nil {
		match.board = b
	}

	if match.players[true] != nil && match.players[false] != nil {
		ctx.Status(400)
		return
	}
	upgrader := websocket.Upgrader{CheckOrigin: allowAnyOrigin}
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Println(err)
		ctx.Status(500)
		return
	}
	defer conn.Close()

	p := match.NewPlayer(conn)
	if p == nil {
		ctx.Status(400)
		return
	}

	p.UpdateBoardChan <- true
	log.Println("Player", p.id, "connected")

	func() {
		for {
			select {
			case r := <-p.ResultsChan:
				if r == p.id {
					p.Conn.WriteJSON(map[string]any{
						"type": "result",
						"data": "You win!",
					})
				} else {
					p.Conn.WriteJSON(map[string]any{
						"type": "result",
						"data": "You lose",
					})
				}
				match.DisconnectPlayer(p.id)
				return
			case <-p.UpdateBoardChan:
				match.SendBoard(p.id)
			}

			for {
				if match.board.PlayerTurn == p.id {
					if w := <-b.WaitingActionChan; w != nil {
						p.ResultsChan <- *w
						match.players[!*w].ResultsChan <- !*w
						break
					}
					mt, msg, err := conn.ReadMessage()
					if err != nil || mt != websocket.TextMessage {
						return
					}
					act := &cards.Action{}
					json.Unmarshal(msg, act)
					match.board.ActionChan <- act
					err = <-b.ActionEndChan
					if err != nil {
						match.SendError(p.id, err)
					}
					match.UpdatePlayersBoard()
				}
				break
			}
		}
	}()

	match.DisconnectPlayer(p.id)
	b.WaitingActionChan <- nil
	match.UpdatePlayersBoard()
	log.Println("Player", p.id, "disconnected")
}

func allowAnyOrigin(r *http.Request) bool {
	return true
}
