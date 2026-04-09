package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"netopiland/application"
	"netopiland/domain"
)

// RunGameLoop reads user input and dispatches actions to the engine until the game ends.
func RunGameLoop(engine *application.Engine) {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(strings.ToLower(scanner.Text()))
		if input == "" {
			continue
		}

		action, ok := domain.ActionFromString(input)
		if !ok {
			fmt.Printf("Warning: '%s' is not a valid action. Type 'help' for available actions.\n", input)
			continue
		}

		if engine.Token.Blocked && action != domain.ActionStatus && action != domain.ActionHelp && action != domain.ActionQuit && action != domain.ActionIdentify {
			fmt.Println(domain.ErrTokenBlocked)
			continue
		}

		switch action {
		case domain.ActionHelp:
			fmt.Println("=== Available Actions ===")
			for _, a := range domain.AllActions() {
				fmt.Printf("  %-12s %s\n", a, a.Description())
			}

		case domain.ActionMove:
			msg, err := engine.Move()
			if err != nil {
				fmt.Println(err)
				continue
			}
			fmt.Println(msg)
			fmt.Println()
			DisplayHUD(engine.Token)
			if engine.GameOver {
				DisplayGameOver(engine.Approved)
				return
			}

		case domain.ActionScan:
			msg, err := engine.Scan()
			if err != nil {
				fmt.Println(err)
				continue
			}
			fmt.Println(msg)

		case domain.ActionShield:
			msg, err := engine.Shield()
			if err != nil {
				fmt.Println(err)
				continue
			}
			fmt.Println(msg)

		case domain.ActionIdentify:
			msg, err := engine.Identify()
			if err != nil {
				fmt.Println(err)
				continue
			}
			fmt.Println(msg)

		case domain.ActionStatus:
			DisplayHUD(engine.Token)

		case domain.ActionWait:
			msg := engine.Wait()
			fmt.Println(msg)

		case domain.ActionJournal:
			fmt.Println(engine.JournalView())

		case domain.ActionQuit:
			fmt.Println("Your token fades from the network. Goodbye!")
			return
		}
	}
}
