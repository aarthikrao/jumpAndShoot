package models

// Client to server messages
var (
	Ping = 1

	// PosUpdateMsg is recieved when the player taps to move up
	PosUpdateMsg = 2

	// Collision is recieved when the player detects a collision
	CollisionMsg = 3

	// used to notify the player that the match has started
	GameStartMsg = 4

	// GameEnd is used to notify that the game has ended
	GameEndMsg = 5

	// ScoreUpdate will be sent to both players
	ScoreUpdateMsg = 6
)
