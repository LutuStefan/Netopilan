package cli

import (
	"fmt"
	"strings"

	"netopiland/domain"
)

// DisplayWelcome prints the full welcome screen: token ID, how to play, gate entry, and HUD.
func DisplayWelcome(token *domain.Token) {
	fmt.Println()
	fmt.Printf("Token ID: %s\n", token.ID)
	fmt.Println()
	displayHowToPlay()
	fmt.Println()
	displayGateEntry(token)
	fmt.Println()
	DisplayHUD(token)
}

func displayHowToPlay() {
	fmt.Println("HOW TO PLAY")
	fmt.Println("You are a payment token traveling through 5 zones.")
	fmt.Println("Reach the Issuer Throne with low risk to get APPROVED.")
	fmt.Println()
	for _, a := range domain.AllActions() {
		fmt.Printf("  %-12s %s\n", a, a.Description())
	}
}

func displayGateEntry(token *domain.Token) {
	gate := token.Position
	fmt.Printf("You enter %s.\n", gate)
	fmt.Println(gate.Description())
}

// DisplayHUD prints the heads-up display: path progress, health, energy, and risk.
func DisplayHUD(token *domain.Token) {
	displayPath(token.Position)
	fmt.Println()
	displayBar("HP", token.Health, token.HealthLabel())
	displayBar("EN", token.Energy, token.EnergyLabel())
	risk := token.RiskLevel()
	fmt.Printf("Risk %d  /%d  %s: %s\n", token.RiskScore.Value, token.RiskScore.Max, risk.Label, risk.Message)
	fmt.Println()
	fmt.Println("> 'move' to advance, 'scan' to scout ahead, 'help' for all commands")
}

func displayPath(current domain.Gate) {
	fmt.Println(strings.Repeat("_", 70))

	gates := []domain.Gate{
		domain.GateMerchant, domain.GateGateway, domain.GateRiskEngine,
		domain.GateAcquirer, domain.GateIssuer,
	}

	var parts []string
	for _, g := range gates {
		name := g.String()
		if g == current {
			parts = append(parts, fmt.Sprintf("[%s]", name))
		} else if g < current {
			parts = append(parts, fmt.Sprintf(" %s>", name))
		} else {
			parts = append(parts, fmt.Sprintf(" %s...", name))
		}
	}
	fmt.Println(strings.Join(parts, "  "))
	fmt.Println(strings.Repeat("_", 70))
}

// DisplayGameOver prints the final verdict banner.
func DisplayGameOver(approved bool) {
	fmt.Println()
	fmt.Println(strings.Repeat("=", 40))
	if approved {
		fmt.Println("  TRANSACTION APPROVED")
	} else {
		fmt.Println("  TRANSACTION DECLINED")
	}
	fmt.Println(strings.Repeat("=", 40))
	fmt.Println()
}

func displayBar(label string, cv domain.ClampedValue, status string) {
	barWidth := 20
	filled := barWidth * cv.Value / cv.Max
	if filled < 0 {
		filled = 0
	}
	if filled > barWidth {
		filled = barWidth
	}
	bar := strings.Repeat("#", filled) + strings.Repeat(".", barWidth-filled)
	fmt.Printf("%s [%s] %d  %s\n", label, bar, cv.Value, status)
}
