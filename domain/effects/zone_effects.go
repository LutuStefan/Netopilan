package effects

import (
	"fmt"
	"math/rand"

	"netopiland/domain"
)

const (
	GatewayWindMinDrain = 5
	GatewayWindMaxDrain = 10
	BlessingResistance  = 10
)

// GatewayWindEffect is a zone event for Gateway Bridge.
// A powerful wind drains energy randomly between -5 and -10.
type GatewayWindEffect struct{}

func (e GatewayWindEffect) Apply(t *domain.Token) string {
	drain := -(rand.Intn(6) + GatewayWindMinDrain)
	t.Energy.Add(drain)
	return fmt.Sprintf("A powerful wind batters your token! Energy %d.", drain)
}

// RiskEngineBlessingEffect is a zone event for Risk Engine Woods.
// The forest blesses the token, adding +10 resistance.
type RiskEngineBlessingEffect struct{}

func (e RiskEngineBlessingEffect) Apply(t *domain.Token) string {
	t.Resistance.Add(BlessingResistance)
	return fmt.Sprintf("The forest spirits bless your token — your defenses grow stronger. Resistance +%d.", BlessingResistance)
}
