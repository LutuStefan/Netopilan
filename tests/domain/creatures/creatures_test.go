package creatures_test

import (
	"netopiland/domain"
	"netopiland/domain/creatures"
	"testing"
)

// --- Fraudster ---

func TestFraudster_ShieldBlocksAllDamage(t *testing.T) {
	token := domain.NewToken("test-id")
	token.ShieldTTL = 2

	initialHealth := token.Health.Value
	initialRisk := token.RiskScore.Value

	result := creatures.Fraudster{}.Execute(token)

	if !result.Passed {
		t.Error("expected Passed=true when shielded")
	}
	if token.Health.Value != initialHealth {
		t.Errorf("shield should prevent health loss, got %d want %d", token.Health.Value, initialHealth)
	}
	if token.RiskScore.Value != initialRisk {
		t.Errorf("shield should prevent risk increase, got %d want %d", token.RiskScore.Value, initialRisk)
	}
}

func TestFraudster_DamageScalesWithAttacks(t *testing.T) {
	// Run many times to cover the random 1-3 attack range.
	// Every execution must respect: per attack = Risk+5, Health-8.
	for i := 0; i < 100; i++ {
		token := domain.NewToken("test-id")
		before := token.Health.Value

		creatures.Fraudster{}.Execute(token)

		healthLost := before - token.Health.Value
		riskGained := token.RiskScore.Value

		if healthLost != 8 && healthLost != 16 && healthLost != 24 {
			t.Fatalf("unexpected health loss: %d (not 8/16/24)", healthLost)
		}
		if riskGained != 5 && riskGained != 10 && riskGained != 15 {
			t.Fatalf("unexpected risk gain: %d (not 5/10/15)", riskGained)
		}
		attacks := healthLost / 8
		if riskGained != attacks*5 {
			t.Fatalf("health/risk mismatch: %d health lost, %d risk gained", healthLost, riskGained)
		}
	}
}

func TestFraudster_CanKillToken(t *testing.T) {
	token := domain.NewToken("test-id")
	token.Health = domain.NewClampedValue(10, 0, 100) // low health

	for i := 0; i < 50; i++ {
		tk := domain.NewToken("test-id")
		tk.Health = domain.NewClampedValue(10, 0, 100)
		creatures.Fraudster{}.Execute(tk)
		if tk.Health.Value <= 0 {
			return
		}
	}
	t.Fatal("Fraudster never killed a token with 10 HP in 50 tries — statistically unlikely")
}

func TestFraudster_ShieldZeroTTL_DoesNotProtect(t *testing.T) {
	token := domain.NewToken("test-id")
	token.ShieldTTL = 0 // expired shield

	creatures.Fraudster{}.Execute(token)

	if token.Health.Value == domain.DefaultHealth {
		t.Error("expired shield (TTL=0) should not block damage")
	}
}

// --- DuplicateDemon ---

func TestDuplicateDemon_BlocksToken(t *testing.T) {
	token := domain.NewToken("test-id")

	result := creatures.DuplicateDemon{}.Execute(token)

	if token.Blocked != true {
		t.Error("DuplicateDemon should set Blocked=true")
	}
	if result.Passed {
		t.Error("expected Passed=false — the demon wins the encounter")
	}
	if result.GameOver {
		t.Error("DuplicateDemon should not end the game")
	}
}

func TestDuplicateDemon_DoesNotDamageStats(t *testing.T) {
	token := domain.NewToken("test-id")
	initialHealth := token.Health.Value
	initialRisk := token.RiskScore.Value
	initialEnergy := token.Energy.Value

	creatures.DuplicateDemon{}.Execute(token)

	if token.Health.Value != initialHealth {
		t.Error("DuplicateDemon should not change health")
	}
	if token.RiskScore.Value != initialRisk {
		t.Error("DuplicateDemon should not change risk")
	}
	if token.Energy.Value != initialEnergy {
		t.Error("DuplicateDemon should not change energy")
	}
}

func TestDuplicateDemon_BlocksAlreadyBlockedToken(t *testing.T) {
	token := domain.NewToken("test-id")
	token.Blocked = true // already blocked

	creatures.DuplicateDemon{}.Execute(token)

	if !token.Blocked {
		t.Error("token should remain blocked")
	}
}

// --- TimeoutSpirit ---

func TestTimeoutSpirit_AlwaysDealsDamageAndRisk(t *testing.T) {
	for i := 0; i < 100; i++ {
		token := domain.NewToken("test-id")
		creatures.TimeoutSpirit{}.Execute(token)

		healthLost := domain.DefaultHealth - token.Health.Value
		riskGained := token.RiskScore.Value

		if riskGained < creatures.TimeoutSpiritMinRisk || riskGained > creatures.TimeoutSpiritMaxRisk {
			t.Fatalf("risk out of range [%d,%d]: %d", creatures.TimeoutSpiritMinRisk, creatures.TimeoutSpiritMaxRisk, riskGained)
		}
		if healthLost < creatures.TimeoutSpiritMinDamage || healthLost > creatures.TimeoutSpiritMaxDamage {
			t.Fatalf("health loss out of range [%d,%d]: %d", creatures.TimeoutSpiritMinDamage, creatures.TimeoutSpiritMaxDamage, healthLost)
		}
	}
}

func TestTimeoutSpirit_IgnoresShield(t *testing.T) {
	// TimeoutSpirit does NOT check shield — this is important game behavior
	token := domain.NewToken("test-id")
	token.ShieldTTL = 2

	creatures.TimeoutSpirit{}.Execute(token)

	if token.Health.Value == domain.DefaultHealth {
		t.Error("TimeoutSpirit should deal damage even through shield")
	}
}

// --- DeclineGuardian ---

func TestDeclineGuardian_CanEndGame(t *testing.T) {
	gameOverSeen := false
	for i := 0; i < 200; i++ {
		token := domain.NewToken("test-id")
		result := creatures.DeclineGuardian{}.Execute(token)
		if result.GameOver {
			gameOverSeen = true
			if result.Passed {
				t.Fatal("GameOver=true should have Passed=false")
			}
			break
		}
	}
	if !gameOverSeen {
		t.Fatal("DeclineGuardian never triggered GameOver in 200 runs — statistically impossible at 30% rate")
	}
}

func TestDeclineGuardian_CanLetPass(t *testing.T) {
	passSeen := false
	for i := 0; i < 200; i++ {
		token := domain.NewToken("test-id")
		result := creatures.DeclineGuardian{}.Execute(token)
		if !result.GameOver && result.Passed {
			passSeen = true
			break
		}
	}
	if !passSeen {
		t.Fatal("DeclineGuardian never let token pass in 200 runs — statistically impossible at 70% rate")
	}
}

func TestDeclineGuardian_NeverModifiesTokenStats(t *testing.T) {
	for i := 0; i < 50; i++ {
		token := domain.NewToken("test-id")
		creatures.DeclineGuardian{}.Execute(token)

		if token.Health.Value != domain.DefaultHealth {
			t.Fatal("DeclineGuardian should not modify health")
		}
		if token.RiskScore.Value != domain.MinStat {
			t.Fatal("DeclineGuardian should not modify risk")
		}
		if token.Energy.Value != domain.DefaultEnergy {
			t.Fatal("DeclineGuardian should not modify energy")
		}
	}
}
