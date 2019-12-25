package objects

import (
	"log"
	"sync"
	"time"
)

// gameServer ..
type gameServer struct {
	clientMap    map[string]*Player
	onlineGames  map[string]*Game
	matchChannel chan *Player
	mu           sync.Mutex
}

// GS ..
var GS *gameServer

func init() {
	StartGameServer()
}

// StartGameServer Starts a new game server
func StartGameServer() {
	GS = &gameServer{
		clientMap:    make(map[string]*Player),
		onlineGames:  make(map[string]*Game),
		matchChannel: make(chan *Player, 2),
	}
	go GS.MatchMakerRoutine()

}

func (gs *gameServer) AddClient(p *Player) {
	gs.mu.Lock()
	defer gs.mu.Unlock()

	gs.clientMap[p.name] = p
}

func (gs *gameServer) RemoveClient(playerName string, p *Player) {
	gs.mu.Lock()
	defer gs.mu.Unlock()

	delete(gs.clientMap, playerName)

}

func (gs *gameServer) AddGame(gameName string, g *Game) {
	gs.mu.Lock()
	defer gs.mu.Unlock()

	gs.onlineGames[gameName] = g
}

func (gs *gameServer) RemoveGame(gameName string, g *Game) {
	gs.mu.Lock()
	defer gs.mu.Unlock()

	delete(gs.onlineGames, gameName)
}

func (gs *gameServer) MatchMakerRoutine() {
	log.Println("Starting MatchMakerRoutine")
	for {
		// nil is recieved if the player wants to cancel
		p1 := <-gs.matchChannel
		p2 := <-gs.matchChannel

		p1Err := p1.SendMessage([]byte("Match made"))
		p2Err := p2.SendMessage([]byte("Match made"))

		if p1Err != nil && p2Err != nil {
			continue
		} else if p1Err != nil {
			gs.matchChannel <- p2
			continue
		} else if p2Err != nil {
			gs.matchChannel <- p1
			continue
		} else {
			g := NewGame(p1, p2)
			gs.AddGame("game"+time.Now().String(), g)
			continue
		}
	}
}

func (gs *gameServer) Match(p *Player) {
	gs.AddClient(p)
	gs.matchChannel <- p
}
