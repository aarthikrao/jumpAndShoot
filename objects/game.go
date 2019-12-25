package objects

import (
	"encoding/json"
	"log"

	"github.com/aarthikrao/jumpAndShoot/models"
)

const maxLives = 5

// Game ..
type Game struct {
	Player1 *Player
	Player2 *Player
}

// NewGame ..
func NewGame(p1 *Player, p2 *Player) *Game {
	log.Println("Creating a new game instance")
	p1.lives = maxLives
	p2.lives = maxLives
	game := Game{
		Player1: p1,
		Player2: p2,
	}
	// TODO : Check if the player is still available
	go game.RouteMessage()
	go game.RouteMessage()

	// Send game start message to player1
	initMsgP1 := models.GameStart{
		Opponent: p2.name,
		MaxLives: p1.lives,
	}
	p1.SendMessage(initMsgP1.ToSocketBytes())

	// Send game start message to player2
	initMsgP2 := models.GameStart{
		Opponent: p1.name,
		MaxLives: p2.lives,
	}
	p2.SendMessage(initMsgP2.ToSocketBytes())

	return &game
}

// RouteMessage ..
func (g *Game) RouteMessage() {
	log.Println("Starting RouteMessage ...")
BREAK:
	for {
		select {
		case msgBytes := <-g.Player1.MsgFromClient:
			msg := models.ParseSocketMessage(msgBytes)
			g.handleGameMove(msg, g.Player1, g.Player2)

		case msgBytes := <-g.Player2.MsgFromClient:
			msg := models.ParseSocketMessage(msgBytes)
			g.handleGameMove(msg, g.Player2, g.Player1)

		// Handle exit messages
		case <-g.Player1.ExitChan:
			log.Println(g.Player1.name, " has exited the game")
			g.Player2.SendMessage(getGameEndMessage(g.Player2.name))
			go g.Player2.Close()
			break BREAK

		case <-g.Player2.ExitChan:
			log.Println(g.Player2.name, " has exited the game")
			g.Player1.SendMessage(getGameEndMessage(g.Player1.name))
			go g.Player1.Close()
			break BREAK
		}
	}
	log.Println("Closing RouteMessage game for", g.Player1.name, " and ", g.Player2.name)
}

// handleGameMove ..
func (g *Game) handleGameMove(msg *models.SocketMessage, player *Player, opponent *Player) {
	switch msg.Type {

	case models.Ping:
		// Ping the message back to client
		player.SendMessage(msg.ToBytes())

	case models.PosUpdateMsg:
		handlePositionUpdate(msg, player, opponent)

	case models.CollisionMsg:
		handleCollisionMessage(msg, player, opponent)

	}
}

func handlePositionUpdate(msg *models.SocketMessage, player *Player, opponent *Player) {
	var pU models.PlayerUpdate
	err := json.Unmarshal(msg.Message, &pU)
	if err != nil {
		log.Println("ERROR", "Invalid message", msg.Message, err)
	}
	if pU.Fire && !player.ready {
		// The player is ready for the match
		log.Println(player.name, "is ready for the match")
		player.ready = true
	}

	// Update the player position value in server
	player.posY = pU.PlayerPositionY

	// Marshal the json and send to player
	j, _ := json.Marshal(pU)
	msg.Message = j
	player.SendMessage(msg.ToBytes())

	// Marshal the json and send to opponent
	pU.IsOpponent = true
	j, _ = json.Marshal(pU)
	msg.Message = j
	opponent.SendMessage(msg.ToBytes())
}

func handleCollisionMessage(msg *models.SocketMessage, player *Player, opponent *Player) {
	if !player.ready || !opponent.ready {
		return
	}
	// verify message
	var cM models.CollisionRequest
	err := json.Unmarshal(msg.Message, &cM)
	if err != nil {
		log.Println("ERROR", "Error in parsing collision message", err)
	}
	if cM.Character == 1 {
		// Player has collided
		log.Println("Collision detected: " + player.name)
		player.lives--
		psu := models.GetScoreUpdateSocketBytes(player.lives, opponent.lives)
		player.SendMessage(psu)

		osu := models.GetScoreUpdateSocketBytes(opponent.lives, player.lives)
		opponent.SendMessage(osu)
	}
	if player.lives <= 0 {
		msg := getGameEndMessage(opponent.name)
		player.SendMessage(msg)
		opponent.SendMessage(msg)

		go player.Close()
		go opponent.Close()
	}
}

func getGameEndMessage(winner string) []byte {
	gW := models.GameEnd{
		Winner: winner,
	}
	winMessage, _ := json.Marshal(gW)
	msg := models.SocketMessage{
		Type:    models.GameEndMsg,
		Message: winMessage,
	}
	return msg.ToBytes()
}
