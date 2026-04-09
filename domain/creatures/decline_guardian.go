package creatures

import (
	"math/rand"

	"netopiland/domain"
)

const DeclineChance = 0.3

// DeclineGuardian has a 30% chance to end the game by declining the token.
type DeclineGuardian struct{}

func (dg DeclineGuardian) Name() string { return "Decline Guardian" }

func (dg DeclineGuardian) Description() string {
	return "A solemn spirit, sometimes present just before the destination. It may decide to end your journey prematurely, depending on your path and the events you've been through."
}

func (dg DeclineGuardian) Execute(t *domain.Token) domain.ChallengeResult {
	if rand.Float64() < DeclineChance {
		return domain.ChallengeResult{
			Message:  "The Decline Guardian raises its hand — your token is DECLINED. Journey over.",
			Passed:   false,
			GameOver: true,
		}
	}
	return domain.ChallengeResult{
		Message: "The Decline Guardian watches you pass... this time.",
		Passed:  true,
	}
}
