package domain

// GateEffect defines the environmental effect that a gate applies to a token upon entry.
// It extends Effect with a description shown when entering the gate.
type GateEffect interface {
	Effect

	// Description returns what the player sees when entering the gate.
	Description() string
}
