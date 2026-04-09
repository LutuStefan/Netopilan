package domain

const (
	DefaultHealth     = 100
	DefaultEnergy     = 100
	DefaultResistance = 30
	MaxStat           = 100
	MinStat           = 0

	RiskThresholdLow    = 20
	RiskThresholdMedium = 50
	RiskThresholdHigh   = 80

	HealthThresholdHealthy = 70
	HealthThresholdWounded = 30

	EnergyThresholdOK  = 50
	EnergyThresholdLow = 20
)

// Token is the core entity — the player's payment token traveling through gates.
// It is identified by ID (UUID string) and has mutable state.
type Token struct {
	ID        string
	Position  Gate
	Health     ClampedValue
	Energy     ClampedValue
	Resistance ClampedValue
	RiskScore  ClampedValue
	ShieldTTL int  // turns remaining on active shield; 0 = no shield
	Blocked   bool // true when a Duplicate Demon blocks the token; cleared by identify
}

// NewToken creates a Token at the starting gate with default stats.
func NewToken(id string) *Token {
	return &Token{
		ID:        id,
		Position:  GateMerchant,
		Health:     NewClampedValue(DefaultHealth, MinStat, MaxStat),
		Energy:     NewClampedValue(DefaultEnergy, MinStat, MaxStat),
		Resistance: NewClampedValue(DefaultResistance, MinStat, MaxStat),
		RiskScore:  NewClampedValue(MinStat, MinStat, MaxStat),
	}
}

// IsShielded returns true if the token has an active shield.
func (t *Token) IsShielded() bool {
	return t.ShieldTTL > 0
}

// TickShield decrements the shield timer by one turn.
func (t *Token) TickShield() {
	if t.ShieldTTL > 0 {
		t.ShieldTTL--
	}
}

// RiskLevel represents a qualitative risk assessment.
type RiskLevel struct {
	Label   string
	Message string
}

// RiskLevel returns the current risk assessment based on the risk score.
func (t *Token) RiskLevel() RiskLevel {
	switch {
	case t.RiskScore.Value <= RiskThresholdLow:
		return RiskLevel{Label: "LOW", Message: "Safe passage expected"}
	case t.RiskScore.Value <= RiskThresholdMedium:
		return RiskLevel{Label: "MEDIUM", Message: "Caution advised"}
	case t.RiskScore.Value <= RiskThresholdHigh:
		return RiskLevel{Label: "HIGH", Message: "Danger ahead"}
	default:
		return RiskLevel{Label: "CRITICAL", Message: "Rejection imminent"}
	}
}

// HealthLabel returns a human-readable label for the current health.
func (t *Token) HealthLabel() string {
	switch {
	case t.Health.Value > HealthThresholdHealthy:
		return "HEALTHY"
	case t.Health.Value > HealthThresholdWounded:
		return "WOUNDED"
	default:
		return "CRITICAL"
	}
}

// EnergyLabel returns a human-readable label for the current energy.
func (t *Token) EnergyLabel() string {
	switch {
	case t.Energy.Value > EnergyThresholdOK:
		return "OK"
	case t.Energy.Value > EnergyThresholdLow:
		return "LOW"
	default:
		return "DEPLETED"
	}
}
