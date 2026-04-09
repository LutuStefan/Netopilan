package effects

import (
	"fmt"

	"netopiland/domain"
)

const (
	MerchantRisk   = 2
	GatewayRisk    = 5
	RiskEngineRisk = 8
	AcquirerRisk   = 6
	IssuerRisk     = 3
)

// MerchantEffect is the gate effect for Merchant Gate.
type MerchantEffect struct{}

func (e MerchantEffect) Description() string {
	return "Echoes of a previous request reverberate through the gate."
}

func (e MerchantEffect) Apply(t *domain.Token) string {
	t.RiskScore.Add(MerchantRisk)
	return fmt.Sprintf("The system detects a repeated pattern — suspicion rises. Risk +%d.", MerchantRisk)
}

// GatewayBridgeEffect is the gate effect for Gateway Bridge.
type GatewayBridgeEffect struct{}

func (e GatewayBridgeEffect) Description() string {
	return "The bridge shudders as your message crosses. Was it sent once, or twice?"
}

func (e GatewayBridgeEffect) Apply(t *domain.Token) string {
	t.RiskScore.Add(GatewayRisk)
	return fmt.Sprintf("A duplicate signal echoes across the bridge. Risk +%d.", GatewayRisk)
}

// RiskEngineEffect is the gate effect for Risk Engine Woods.
type RiskEngineEffect struct{}

func (e RiskEngineEffect) Description() string {
	return "The forest eyes scan your history and score your every step."
}

func (e RiskEngineEffect) Apply(t *domain.Token) string {
	t.RiskScore.Add(RiskEngineRisk)
	return fmt.Sprintf("The forest judges you harshly — your risk profile deepens. Risk +%d.", RiskEngineRisk)
}

// AcquirerPassEffect is the gate effect for Acquirer Pass.
type AcquirerPassEffect struct{}

func (e AcquirerPassEffect) Description() string {
	return "The mountain winds carry your token forward — but at what cost?"
}

func (e AcquirerPassEffect) Apply(t *domain.Token) string {
	t.RiskScore.Add(AcquirerRisk)
	return fmt.Sprintf("The pass extracts its toll as you push through. Risk +%d.", AcquirerRisk)
}

// IssuerThroneEffect is the gate effect for Issuer Throne.
type IssuerThroneEffect struct{}

func (e IssuerThroneEffect) Description() string {
	return "The throne awaits your arrival. Final judgment is near."
}

func (e IssuerThroneEffect) Apply(t *domain.Token) string {
	t.RiskScore.Add(IssuerRisk)
	return fmt.Sprintf("The throne's gaze weighs upon you one last time. Risk +%d.", IssuerRisk)
}
