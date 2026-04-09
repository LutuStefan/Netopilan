package creatures

import (
	"fmt"
	"math/rand"
	"strings"

	"netopiland/domain"
)

const (
	FraudsterMaxAttacks   = 3
	FraudsterRiskPerHit   = 5
	FraudsterDamagePerHit = 8
)

// Fraudster attacks 1-3 times. Each attack: Risk +5, Health -8.
type Fraudster struct{}

func (f Fraudster) Name() string { return "Fraudster" }

func (f Fraudster) Description() string {
	return "An aggressive spirit that leaps at you trying to stop your progress. It attacks repetitively or randomly, sometimes faster than you expect."
}

func (f Fraudster) Execute(t *domain.Token) domain.ChallengeResult {
	if t.IsShielded() {
		return domain.ChallengeResult{
			Message: "The Fraudster lunges — but your shield absorbs the blow!",
			Passed:  true,
		}
	}

	attacks := rand.Intn(FraudsterMaxAttacks) + 1

	var lines []string
	for i := 1; i <= attacks; i++ {
		t.RiskScore.Add(FraudsterRiskPerHit)
		t.Health.Add(-FraudsterDamagePerHit)
		lines = append(lines, fmt.Sprintf("  Attack %d: Risk +%d, Health -%d.", i, FraudsterRiskPerHit, FraudsterDamagePerHit))
	}

	header := fmt.Sprintf("The Fraudster strikes %d time(s)!", attacks)
	msg := header + "\n" + strings.Join(lines, "\n")
	return domain.ChallengeResult{Message: msg, Passed: true}
}
