package objects

import (
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// Player ..
type Player struct {
	conn          *websocket.Conn
	MsgFromClient chan []byte
	ExitChan      chan int
	posY          float64
	ready         bool
	lives         int
	name          string
	sendMutex     sync.Mutex
}

// NewPlayer creates a new instance of Player
func NewPlayer(playerName string, conn *websocket.Conn) *Player {
	return &Player{
		conn:          conn,
		MsgFromClient: make(chan []byte, 5),
		ExitChan:      make(chan int),
		name:          playerName,
	}
}

// RecieveMessages ..
func (p *Player) RecieveMessages() {
	for {
		_, message, err := p.conn.ReadMessage()
		p.MsgFromClient <- message
		if err != nil {
			log.Println(p.name, " has quit.")
			p.conn.Close()
			// notify that the player has quit
			close(p.ExitChan)
			break
		}
	}
}

// SendMessage ..
func (p *Player) SendMessage(msg []byte) error {
	p.sendMutex.Lock()
	defer p.sendMutex.Unlock()

	// log.Println("Sending message to", p.name, ":", string(msg))
	return p.conn.WriteMessage(1, msg)
}

// Close ..
func (p *Player) Close() {
	time.Sleep(1 * time.Second)
	p.conn.Close()
}
