package domain

// Effect is the base contract for anything that can modify a token's state.
type Effect interface {
	// Apply mutates the token and returns a message describing what happened.
	Apply(t *Token) string
}
