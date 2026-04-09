package application_test

import (
	"netopiland/application"
	"netopiland/domain"
	"netopiland/domain/creatures"
	"netopiland/domain/effects"
	"strings"
	"testing"
)

// --- helpers ---

// stubChallenge is a deterministic creature for controlled testing.
type stubChallenge struct {
	name   string
	onExec func(t *domain.Token) domain.ChallengeResult
}

func (s stubChallenge) Name() string        { return s.name }
func (s stubChallenge) Description() string { return "stub" }
func (s stubChallenge) Execute(t *domain.Token) domain.ChallengeResult {
	if s.onExec != nil {
		return s.onExec(t)
	}
	return domain.ChallengeResult{Passed: true, Message: "stub passed"}
}

// alwaysSpawner returns a spawner that always spawns the given creature.
func alwaysSpawner(c domain.Challenge) *application.Spawner {
	return application.NewSpawner(1.0, []application.CreatureWeight{
		{Creature: c, Weight: 1.0},
	})
}

// neverSpawner returns a spawner that never spawns creatures.
func neverSpawner() *application.Spawner {
	return application.NewSpawner(0.0, []application.CreatureWeight{
		{Creature: stubChallenge{name: "none"}, Weight: 1.0},
	})
}

// defaultGateEffects returns the standard gate effects map.
func defaultGateEffects() map[domain.Gate]domain.GateEffect {
	return map[domain.Gate]domain.GateEffect{
		domain.GateMerchant:   effects.MerchantEffect{},
		domain.GateGateway:    effects.GatewayBridgeEffect{},
		domain.GateRiskEngine: effects.RiskEngineEffect{},
		domain.GateAcquirer:   effects.AcquirerPassEffect{},
		domain.GateIssuer:     effects.IssuerThroneEffect{},
	}
}

// newTestEngine creates an engine with no zone events and a given spawner.
func newTestEngine(spawner *application.Spawner) *application.Engine {
	return application.NewEngine(
		domain.NewToken("test-id"),
		domain.NewJournal(),
		defaultGateEffects(),
		map[domain.Gate]domain.ZoneEvent{}, // no zone events for deterministic tests
		spawner,
	)
}

// --- Move Action ---

func TestMove_HappyPath_AdvancesGate(t *testing.T) {
	eng := newTestEngine(neverSpawner())

	msg, err := eng.Move()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if eng.Token.Position != domain.GateGateway {
		t.Errorf("expected GateGateway, got %v", eng.Token.Position)
	}
	if !strings.Contains(msg, "Gateway Bridge") {
		t.Errorf("move message should mention destination gate, got: %s", msg)
	}
}

func TestMove_CostsEnergy(t *testing.T) {
	eng := newTestEngine(neverSpawner())

	eng.Move()

	expected := domain.DefaultEnergy - application.MoveCost
	if eng.Token.Energy.Value != expected {
		t.Errorf("expected energy %d after move, got %d", expected, eng.Token.Energy.Value)
	}
}

func TestMove_NotEnoughEnergy_ReturnsError(t *testing.T) {
	eng := newTestEngine(neverSpawner())
	eng.Token.Energy = domain.NewClampedValue(application.MoveCost-1, domain.MinStat, domain.MaxStat)

	_, err := eng.Move()
	if err != domain.ErrNotEnoughEnergy {
		t.Errorf("expected ErrNotEnoughEnergy, got %v", err)
	}
	// Position should not change
	if eng.Token.Position != domain.GateMerchant {
		t.Error("position should not change when move fails")
	}
}

func TestMove_AtFinish_ReturnsError(t *testing.T) {
	eng := newTestEngine(neverSpawner())
	eng.Token.Position = domain.GateIssuer

	_, err := eng.Move()
	if err != domain.ErrAlreadyAtFinish {
		t.Errorf("expected ErrAlreadyAtFinish, got %v", err)
	}
}

func TestMove_AppliesGateEffect(t *testing.T) {
	eng := newTestEngine(neverSpawner())

	eng.Move()

	if eng.Token.RiskScore.Value != effects.GatewayRisk {
		t.Errorf("expected risk %d after gateway effect, got %d", effects.GatewayRisk, eng.Token.RiskScore.Value)
	}
}

func TestMove_CreatureEncounter_AppliesDamage(t *testing.T) {
	// Use DuplicateDemon for deterministic behavior (no randomness)
	eng := newTestEngine(alwaysSpawner(creatures.DuplicateDemon{}))

	eng.Move()

	if !eng.Token.Blocked {
		t.Error("DuplicateDemon should block token during move")
	}
}

func TestMove_HealthReachesZero_GameOver(t *testing.T) {
	killer := stubChallenge{
		name: "killer",
		onExec: func(t *domain.Token) domain.ChallengeResult {
			t.Health.Add(-200) // overkill
			return domain.ChallengeResult{Passed: false, Message: "killed"}
		},
	}
	eng := newTestEngine(alwaysSpawner(killer))

	msg, err := eng.Move()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !eng.GameOver {
		t.Error("game should be over when health reaches 0")
	}
	if eng.Token.Health.Value != 0 {
		t.Errorf("health should clamp to 0, got %d", eng.Token.Health.Value)
	}
	if !strings.Contains(msg, "DECLINED") {
		t.Error("death message should contain DECLINED")
	}
}

func TestMove_CreatureGameOver_SetsFlag(t *testing.T) {
	declineCreature := stubChallenge{
		name: "decline",
		onExec: func(t *domain.Token) domain.ChallengeResult {
			return domain.ChallengeResult{Passed: false, GameOver: true, Message: "declined"}
		},
	}
	eng := newTestEngine(alwaysSpawner(declineCreature))

	eng.Move()

	if !eng.GameOver {
		t.Error("creature GameOver result should set engine GameOver")
	}
}

func TestMove_ReachFinish_LowRisk_Approved(t *testing.T) {
	eng := newTestEngine(neverSpawner())
	// Put token one gate before finish
	eng.Token.Position = domain.GateAcquirer
	// Keep risk low enough that even IssuerThroneEffect (+3) won't push over 30
	eng.Token.RiskScore = domain.NewClampedValue(10, 0, 100)

	msg, _ := eng.Move()

	if !eng.GameOver {
		t.Error("reaching finish should set GameOver")
	}
	if !eng.Approved {
		t.Error("low risk at finish should be APPROVED")
	}
	if !strings.Contains(msg, "APPROVED") {
		t.Error("message should contain APPROVED")
	}
}

func TestMove_ReachFinish_HighRisk_Declined(t *testing.T) {
	eng := newTestEngine(neverSpawner())
	eng.Token.Position = domain.GateAcquirer
	eng.Token.RiskScore = domain.NewClampedValue(50, 0, 100) // +3 from gate = 53

	msg, _ := eng.Move()

	if !eng.GameOver {
		t.Error("reaching finish should set GameOver")
	}
	if eng.Approved {
		t.Error("high risk at finish should NOT be approved")
	}
	if !strings.Contains(msg, "DECLINED") {
		t.Error("message should contain DECLINED")
	}
}

func TestMove_ReachFinish_RiskExactly31_Declined(t *testing.T) {
	eng := newTestEngine(neverSpawner())
	eng.Token.Position = domain.GateAcquirer
	eng.Token.RiskScore = domain.NewClampedValue(28, 0, 100) // +3 from IssuerThrone = 31

	eng.Move()

	if eng.Approved {
		t.Error("risk 31 at finish should be DECLINED (threshold is >30)")
	}
}

func TestMove_ReachFinish_RiskExactly30_Approved(t *testing.T) {
	eng := newTestEngine(neverSpawner())
	eng.Token.Position = domain.GateAcquirer
	eng.Token.RiskScore = domain.NewClampedValue(27, 0, 100) // +3 from IssuerThrone = 30

	eng.Move()

	if !eng.Approved {
		t.Error("risk 30 at finish should be APPROVED (threshold is >30)")
	}
}

func TestMove_ShieldTicksDownAfterEncounter(t *testing.T) {
	eng := newTestEngine(neverSpawner())
	eng.Token.ShieldTTL = application.ShieldDuration

	eng.Move()

	if eng.Token.ShieldTTL != application.ShieldDuration-1 {
		t.Errorf("shield should tick from %d to %d after move, got %d", application.ShieldDuration, application.ShieldDuration-1, eng.Token.ShieldTTL)
	}
}

func TestMove_ShieldDoesNotGoNegative(t *testing.T) {
	eng := newTestEngine(neverSpawner())
	eng.Token.ShieldTTL = 0

	eng.Move()

	if eng.Token.ShieldTTL != 0 {
		t.Errorf("shield at 0 should stay 0, got %d", eng.Token.ShieldTTL)
	}
}

func TestMove_ScanThenMove_UsesPendingCreature(t *testing.T) {
	blocker := stubChallenge{
		name: "blocker",
		onExec: func(t *domain.Token) domain.ChallengeResult {
			t.Blocked = true
			return domain.ChallengeResult{Passed: false, Message: "blocked"}
		},
	}

	eng := newTestEngine(neverSpawner()) // spawner won't be used for move
	eng.PendingCreature = blocker
	eng.Scanned = true

	eng.Move()

	if !eng.Token.Blocked {
		t.Error("move should use pending creature from scan")
	}
	if eng.Scanned {
		t.Error("Scanned flag should be cleared after move")
	}
	if eng.PendingCreature != nil {
		t.Error("PendingCreature should be nil after move")
	}
}

func TestMove_TraverseAllGates_NoCreatures_AccumulatesRisk(t *testing.T) {
	eng := newTestEngine(neverSpawner())

	for i := 0; i < 4; i++ {
		eng.Move()
	}

	expectedRisk := effects.GatewayRisk + effects.RiskEngineRisk + effects.AcquirerRisk + effects.IssuerRisk
	if eng.Token.RiskScore.Value != expectedRisk {
		t.Errorf("expected total risk %d after traversing all gates, got %d", expectedRisk, eng.Token.RiskScore.Value)
	}
	if !eng.Approved {
		t.Error("risk 22 should result in APPROVED")
	}
}

func TestMove_JournalRecordsMovement(t *testing.T) {
	eng := newTestEngine(neverSpawner())

	eng.Move()

	entries := eng.Journal.Entries()
	if len(entries) == 0 {
		t.Fatal("journal should have entries after move")
	}
	lastEntry := entries[len(entries)-1]
	if lastEntry.Action != "move" {
		t.Errorf("last journal entry action should be 'move', got %q", lastEntry.Action)
	}
}

// --- Scan ---

func TestScan_CostsEnergy(t *testing.T) {
	eng := newTestEngine(neverSpawner())

	eng.Scan()

	expected := domain.DefaultEnergy - application.ScanCost
	if eng.Token.Energy.Value != expected {
		t.Errorf("expected energy %d after scan, got %d", expected, eng.Token.Energy.Value)
	}
}

func TestScan_SetsPendingCreatureAndFlag(t *testing.T) {
	eng := newTestEngine(alwaysSpawner(creatures.Fraudster{}))

	eng.Scan()

	if !eng.Scanned {
		t.Error("Scanned flag should be true after scan")
	}
	if eng.PendingCreature == nil {
		t.Error("PendingCreature should be set when spawner returns a creature")
	}
}

func TestScan_DoesNotMoveToken(t *testing.T) {
	eng := newTestEngine(neverSpawner())

	eng.Scan()

	if eng.Token.Position != domain.GateMerchant {
		t.Error("scan should not change position")
	}
}

func TestScan_AtFinish_ReturnsError(t *testing.T) {
	eng := newTestEngine(neverSpawner())
	eng.Token.Position = domain.GateIssuer

	_, err := eng.Scan()
	if err != domain.ErrAlreadyAtFinish {
		t.Errorf("expected ErrAlreadyAtFinish, got %v", err)
	}
}

func TestScan_NotEnoughEnergy_ReturnsError(t *testing.T) {
	eng := newTestEngine(neverSpawner())
	eng.Token.Energy = domain.NewClampedValue(3, 0, 100)

	_, err := eng.Scan()
	if err != domain.ErrNotEnoughEnergy {
		t.Errorf("expected ErrNotEnoughEnergy, got %v", err)
	}
}

// --- Shield ---

func TestShield_CostsEnergyAndActivates(t *testing.T) {
	eng := newTestEngine(neverSpawner())

	msg, err := eng.Shield()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expectedEnergy := domain.DefaultEnergy - application.ShieldCost
	if eng.Token.Energy.Value != expectedEnergy {
		t.Errorf("expected energy %d, got %d", expectedEnergy, eng.Token.Energy.Value)
	}
	if eng.Token.ShieldTTL != application.ShieldDuration {
		t.Errorf("expected ShieldTTL %d, got %d", application.ShieldDuration, eng.Token.ShieldTTL)
	}
	if !strings.Contains(msg, "Shield activated") {
		t.Error("message should confirm shield activation")
	}
}

func TestShield_NotEnoughEnergy(t *testing.T) {
	eng := newTestEngine(neverSpawner())
	eng.Token.Energy = domain.NewClampedValue(application.ShieldCost-1, domain.MinStat, domain.MaxStat)

	_, err := eng.Shield()
	if err != domain.ErrNotEnoughEnergy {
		t.Errorf("expected ErrNotEnoughEnergy, got %v", err)
	}
	if eng.Token.ShieldTTL != 0 {
		t.Error("shield should not activate when energy is insufficient")
	}
}

func TestShield_OverwritesExistingShield(t *testing.T) {
	eng := newTestEngine(neverSpawner())
	eng.Token.ShieldTTL = 1 // existing shield about to expire

	eng.Shield()

	if eng.Token.ShieldTTL != application.ShieldDuration {
		t.Errorf("shield should reset to %d, got %d", application.ShieldDuration, eng.Token.ShieldTTL)
	}
}

// --- Identify ---

func TestIdentify_ClearsBlock(t *testing.T) {
	eng := newTestEngine(neverSpawner())
	eng.Token.Blocked = true

	msg, err := eng.Identify()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if eng.Token.Blocked {
		t.Error("identify should clear blocked state")
	}
	expectedEnergy := domain.DefaultEnergy - application.IdentifyCost
	if eng.Token.Energy.Value != expectedEnergy {
		t.Errorf("expected energy %d, got %d", expectedEnergy, eng.Token.Energy.Value)
	}
	if !strings.Contains(msg, "Identity confirmed") {
		t.Error("message should confirm identity")
	}
}

func TestIdentify_NotBlocked_DoesNotCostEnergy(t *testing.T) {
	eng := newTestEngine(neverSpawner())
	// token is NOT blocked

	msg, err := eng.Identify()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if eng.Token.Energy.Value != domain.DefaultEnergy {
		t.Errorf("identify on non-blocked token should not cost energy, got %d", eng.Token.Energy.Value)
	}
	if !strings.Contains(msg, "not blocked") {
		t.Error("should tell user they're not blocked")
	}
}

func TestIdentify_NotEnoughEnergy(t *testing.T) {
	eng := newTestEngine(neverSpawner())
	eng.Token.Blocked = true
	eng.Token.Energy = domain.NewClampedValue(application.IdentifyCost-1, domain.MinStat, domain.MaxStat)

	_, err := eng.Identify()
	if err != domain.ErrNotEnoughEnergy {
		t.Errorf("expected ErrNotEnoughEnergy, got %v", err)
	}
	if !eng.Token.Blocked {
		t.Error("block should not be cleared when identify fails")
	}
}

// --- Wait ---

func TestWait_RestoresEnergy(t *testing.T) {
	eng := newTestEngine(neverSpawner())
	eng.Token.Energy = domain.NewClampedValue(50, domain.MinStat, domain.MaxStat)

	eng.Wait()

	if eng.Token.Energy.Value != 50+application.WaitRestore {
		t.Errorf("expected energy %d, got %d", 50+application.WaitRestore, eng.Token.Energy.Value)
	}
}

func TestWait_EnergyCapsAt100(t *testing.T) {
	eng := newTestEngine(neverSpawner())
	eng.Token.Energy = domain.NewClampedValue(domain.MaxStat-5, domain.MinStat, domain.MaxStat)

	eng.Wait()

	if eng.Token.Energy.Value != domain.MaxStat {
		t.Errorf("energy should cap at %d, got %d", domain.MaxStat, eng.Token.Energy.Value)
	}
}

func TestWait_TicksShieldDown(t *testing.T) {
	eng := newTestEngine(neverSpawner())
	eng.Token.ShieldTTL = application.ShieldDuration

	eng.Wait()

	if eng.Token.ShieldTTL != application.ShieldDuration-1 {
		t.Errorf("wait should tick shield from %d to %d, got %d", application.ShieldDuration, application.ShieldDuration-1, eng.Token.ShieldTTL)
	}
}

func TestWait_RerollsPendingCreature_WhenScanned(t *testing.T) {
	// When waiting with a scanned creature, the creature is re-rolled
	eng := newTestEngine(alwaysSpawner(creatures.Fraudster{}))
	eng.Scanned = true
	eng.PendingCreature = creatures.DuplicateDemon{}

	eng.Wait()

	// The spawner always spawns Fraudster, so pending should now be Fraudster
	if eng.PendingCreature == nil {
		t.Fatal("PendingCreature should be re-rolled on wait")
	}
	if eng.PendingCreature.Name() != "Fraudster" {
		t.Errorf("expected Fraudster from re-roll, got %s", eng.PendingCreature.Name())
	}
}

func TestWait_DoesNotSetScanned_WhenNotScanned(t *testing.T) {
	eng := newTestEngine(neverSpawner())

	eng.Wait()

	if eng.Scanned {
		t.Error("wait should not set Scanned flag")
	}
}

// --- Integration-style: full game scenarios ---

func TestScenario_ShieldThenMove_FraudsterBlocked(t *testing.T) {
	eng := newTestEngine(alwaysSpawner(creatures.Fraudster{}))

	eng.Shield()
	initialHealth := eng.Token.Health.Value

	eng.Move()

	// Fraudster should be blocked by shield
	if eng.Token.Health.Value < initialHealth {
		// Health only lost from move energy cost, not from Fraudster
		// Actually health doesn't change from energy — check that Fraudster damage was blocked
		t.Error("Fraudster should be blocked by active shield")
	}
}

func TestScenario_ShieldExpires_AfterTwoMoves(t *testing.T) {
	eng := newTestEngine(neverSpawner())

	eng.Shield()
	if eng.Token.ShieldTTL != application.ShieldDuration {
		t.Fatalf("shield should start at %d", application.ShieldDuration)
	}

	eng.Move() // tick: 2 -> 1
	if eng.Token.ShieldTTL != 1 {
		t.Errorf("expected shield 1 after first move, got %d", eng.Token.ShieldTTL)
	}

	eng.Move() // tick: 1 -> 0
	if eng.Token.ShieldTTL != 0 {
		t.Errorf("expected shield 0 after second move, got %d", eng.Token.ShieldTTL)
	}

	if eng.Token.IsShielded() {
		t.Error("shield should be expired after 2 moves")
	}
}

func TestScenario_DuplicateDemon_ThenIdentify_ThenContinue(t *testing.T) {
	demonSpawner := alwaysSpawner(creatures.DuplicateDemon{})
	eng := newTestEngine(demonSpawner)

	// Move encounters DuplicateDemon
	eng.Move()
	if !eng.Token.Blocked {
		t.Fatal("should be blocked after DuplicateDemon encounter")
	}

	// Identify clears the block
	eng.Identify()
	if eng.Token.Blocked {
		t.Fatal("identify should clear the block")
	}
}

func TestScenario_EnergyDepletion_CannotAct(t *testing.T) {
	eng := newTestEngine(neverSpawner())
	eng.Token.Energy = domain.NewClampedValue(0, domain.MinStat, domain.MaxStat)

	_, err := eng.Move()
	if err != domain.ErrNotEnoughEnergy {
		t.Errorf("move should fail: %v", err)
	}

	_, err = eng.Scan()
	if err != domain.ErrNotEnoughEnergy {
		t.Errorf("scan should fail: %v", err)
	}

	_, err = eng.Shield()
	if err != domain.ErrNotEnoughEnergy {
		t.Errorf("shield should fail: %v", err)
	}

	// Wait should still work (free action)
	msg := eng.Wait()
	if eng.Token.Energy.Value != application.WaitRestore {
		t.Errorf("wait should restore %d energy from 0, got %d", application.WaitRestore, eng.Token.Energy.Value)
	}
	if !strings.Contains(msg, "15") {
		t.Error("wait message should mention energy restored")
	}
}

func TestScenario_FullTraversal_Approved(t *testing.T) {
	eng := newTestEngine(neverSpawner())

	// Traverse all 4 gates without creatures
	for i := 0; i < 4; i++ {
		_, err := eng.Move()
		if err != nil {
			t.Fatalf("move %d failed: %v", i+1, err)
		}
	}

	if eng.Token.Position != domain.GateIssuer {
		t.Errorf("should be at Issuer, got %v", eng.Token.Position)
	}
	if !eng.GameOver {
		t.Error("game should be over at finish")
	}
	if !eng.Approved {
		t.Errorf("expected APPROVED with risk %d", eng.Token.RiskScore.Value)
	}
}
