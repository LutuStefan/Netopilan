package creatures

import (
	"fmt"
	"math/rand"

	"netopiland/domain"
)

const (
	TimeoutSpiritMinRisk   = 5
	TimeoutSpiritMaxRisk   = 10
	TimeoutSpiritMinDamage = 5
	TimeoutSpiritMaxDamage = 15
)

// TimeoutSpirit distorts time around your token.
// Effect: Risk +random(5,10), Health -random(5,15).
type TimeoutSpirit struct{}

func (ts TimeoutSpirit) Name() string { return "Timeout Spirit" }

func (ts TimeoutSpirit) Description() string {
	return "Time flows strangely near this spirit. It can slow reality, skip an action, or force you to repeat what you already tried."
}

func (ts TimeoutSpirit) Execute(t *domain.Token) domain.ChallengeResult {
	riskAdd := rand.Intn(6) + TimeoutSpiritMinRisk
	dmg := -(rand.Intn(11) + TimeoutSpiritMinDamage)

	t.RiskScore.Add(riskAdd)
	t.Health.Add(dmg)

	msg := fmt.Sprintf("The Timeout Spirit warps time around you! Risk +%d, Health %d.", riskAdd, dmg)
	return domain.ChallengeResult{Message: msg, Passed: true}
}
