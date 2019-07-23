package gums

// A Player is someone or something able to decide
// the next move.
type Player interface {
	Choose(PlayerColor, State) (canMove bool, move position)
}
