package models

import (
	"encoding/json"
	"log"
)

//
// PlayerUpdate is used by both client and server to notify about player change
//
type PlayerUpdate struct {

	// The position of the player in the y axis
	PlayerPositionY float64 `json:"posY"`

	// True if the player fired a bullet at this position
	Fire bool `json:"fire"`

	IsOpponent bool `json:"isOpp"`
}

// ScoreUpdate ..
type ScoreUpdate struct {
	// filled only by server
	MyLives  int `json:"myLives"`
	OppLives int `json:"oppLives"`
}

//
// CollisionRequest is sent when a collision is detected
//
type CollisionRequest struct {

	// The position of the opponent
	OpponentLocation float64 `json:"posY"`

	// 1: player, 2: opponent
	Character int `json:"ch"`
}

// GameEnd is used to notify both clients that game has ended
type GameEnd struct {

	// The winner of the game
	Winner string `json:"winner"`
}

// GameStart is sent on successful match making
type GameStart struct {
	// Name of the opponent
	Opponent string `json:"opp"`

	// MaxLives
	MaxLives int `json:"lives"`
}

// SocketMessage ..
type SocketMessage struct {
	Type    int             `json:"type"`
	Message json.RawMessage `json:"msg"`
}

// ParseSocketMessage ..
func ParseSocketMessage(msgBytes []byte) *SocketMessage {
	msg := SocketMessage{}
	err := json.Unmarshal(msgBytes, &msg)
	if err != nil {
		log.Println("ERROR", "Error in recieving message", err)
	}
	return &msg
}

// ToBytes returns the socket message in bytes
func (msg *SocketMessage) ToBytes() (returnMsg []byte) {
	returnMsg, err := json.Marshal(msg)
	if err != nil {
		returnMsg = []byte(err.Error())
	}
	return
}

// ToSocketBytes returns the game start message to socket message in bytes
func (msg *GameStart) ToSocketBytes() []byte {
	gameStartBytes, err := json.Marshal(msg)
	if err != nil {
		gameStartBytes = []byte(err.Error())
	}
	sm := SocketMessage{
		Type:    GameStartMsg,
		Message: gameStartBytes,
	}
	return sm.ToBytes()
}

// GetScoreUpdateSocketBytes ..
func GetScoreUpdateSocketBytes(myLives int, oppLives int) []byte {
	su := ScoreUpdate{
		MyLives:  myLives,
		OppLives: oppLives,
	}
	suBytes, err := json.Marshal(su)
	if err != nil {
		suBytes = []byte(err.Error())
	}
	sm := SocketMessage{
		Type:    ScoreUpdateMsg,
		Message: suBytes,
	}
	return sm.ToBytes()
}
