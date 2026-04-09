package domain

// Gate represents a zone on the payment token's path.
// It is a value object — immutable, compared by value.
type Gate int

const (
	GateMerchant   Gate = iota // Start
	GateGateway                // Step 2
	GateRiskEngine             // Step 3
	GateAcquirer               // Step 4
	GateIssuer                 // Finish
)

// GateCount is the total number of gates on the path.
const GateCount = 5

// gateInfo holds the display metadata for a single gate.
type gateInfo struct {
	Name        string
	Description string
}

// gateRegistry maps each Gate to its name and description.
var gateRegistry = map[Gate]gateInfo{
	GateMerchant: {
		Name:        "Merchant Gate",
		Description: "This is where you are created. Yet sometimes echoes of previous requests linger in the air, even though they shouldn't.",
	},
	GateGateway: {
		Name:        "Gateway Bridge",
		Description: "A fragile bridge, often battered by the chaotic winds of Netopiland. Sometimes it sends your message onward, other times it loses it or resends it multiple times, just out of excessive caution.",
	},
	GateRiskEngine: {
		Name:        "Risk Engine Woods",
		Description: "A dark forest full of unseen eyes. They analyze your behavior, history, and pace. If something seems suspicious, the forest tests your very existence.",
	},
	GateAcquirer: {
		Name:        "Acquirer Pass",
		Description: "A restless mountain pass where the winds shift constantly. You may be pushed forward, turned back, or held for a moment.",
	},
	GateIssuer: {
		Name:        "Issuer Throne",
		Description: "The final destination. Here you will receive the supreme answer: Approved or Declined.",
	},
}

// String returns the human-readable name of the gate.
func (g Gate) String() string {
	if info, ok := gateRegistry[g]; ok {
		return info.Name
	}
	return "Unknown"
}

// Description returns the flavor text for this gate.
func (g Gate) Description() string {
	if info, ok := gateRegistry[g]; ok {
		return info.Description
	}
	return ""
}

// IsFinish returns true if this is the final gate.
func (g Gate) IsFinish() bool {
	return g == GateIssuer
}

// Next returns the next gate. Returns ok=false if already at the finish.
func (g Gate) Next() (Gate, bool) {
	if g.IsFinish() {
		return g, false
	}
	return g + 1, true
}
