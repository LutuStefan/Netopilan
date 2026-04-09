package domain

// ChallengeResult describes the outcome of a creature encounter.
type ChallengeResult struct {
	Message  string
	Passed   bool
	GameOver bool // if true, the game ends (e.g., Decline Guardian rejects the token)
}

// Challenge defines the behavior of a creature in Netopiland.
type Challenge interface {
	// Name returns a display name for this creature.
	Name() string

	// Description returns flavor text shown when the creature is scanned.
	Description() string

	// Execute runs the encounter, mutating the token directly, and returns the result.
	Execute(t *Token) ChallengeResult
}
