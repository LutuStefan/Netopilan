package main

import (
	"crypto/rand"
	"fmt"

	"netopiland/application"
	"netopiland/domain"
	"netopiland/domain/creatures"
	"netopiland/domain/effects"
	"netopiland/infrastructure/cli"
)

func main() {
	token := domain.NewToken(generateUUID())

	gateEffects := map[domain.Gate]domain.GateEffect{
		domain.GateMerchant:   effects.MerchantEffect{},
		domain.GateGateway:    effects.GatewayBridgeEffect{},
		domain.GateRiskEngine: effects.RiskEngineEffect{},
		domain.GateAcquirer:   effects.AcquirerPassEffect{},
		domain.GateIssuer:     effects.IssuerThroneEffect{},
	}

	zoneEvents := map[domain.Gate]domain.ZoneEvent{
		domain.GateGateway:    {Effect: effects.GatewayWindEffect{}, Probability: 0.5},
		domain.GateRiskEngine: {Effect: effects.RiskEngineBlessingEffect{}, Probability: 0.5},
	}

	spawner := application.NewSpawner(0.5, []application.CreatureWeight{
		{Creature: creatures.Fraudster{}, Weight: 31.3},
		{Creature: creatures.DuplicateDemon{}, Weight: 31.3},
		{Creature: creatures.TimeoutSpirit{}, Weight: 31.3},
		{Creature: creatures.DeclineGuardian{}, Weight: 6.0},
	})

	engine := application.NewEngine(token, domain.NewJournal(), gateEffects, zoneEvents, spawner)

	cli.DisplayWelcome(engine.Token)
	cli.RunGameLoop(engine)
}

func generateUUID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	b[6] = (b[6] & 0x0f) | 0x40 // version 4
	b[8] = (b[8] & 0x3f) | 0x80 // variant 2
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}
