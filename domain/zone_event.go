package domain

import "math/rand"

// ZoneEvent wraps an Effect with a probability of triggering.
// Unlike GateEffects which always apply, ZoneEvents are optional and fire based on chance.
type ZoneEvent struct {
	Effect      Effect
	Probability float64 // 0.0 to 1.0
}

// Trigger rolls the dice and, if successful, applies the effect to the token.
// Returns the effect message and true if triggered, or empty string and false otherwise.
func (ze ZoneEvent) Trigger(t *Token) (string, bool) {
	if rand.Float64() < ze.Probability {
		return ze.Effect.Apply(t), true
	}
	return "", false
}
