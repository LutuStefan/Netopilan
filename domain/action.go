package domain

// Action represents a user action in the game.
type Action int

const (
	ActionHelp     Action = iota
	ActionMove
	ActionScan
	ActionShield
	ActionIdentify
	ActionStatus
	ActionWait
	ActionJournal
	ActionQuit
)

// actionInfo holds the display metadata for a single action.
type actionInfo struct {
	Name        string
	Description string
}

// actionRegistry maps each Action to its name and description.
// This is the single source of truth for action metadata.
var actionRegistry = map[Action]actionInfo{
	ActionMove:     {Name: "move", Description: "Advance to the next zone on the sacred path (-5 energy)"},
	ActionScan:     {Name: "scan", Description: "Sense what lies ahead in the next zone (-5 energy)"},
	ActionShield:   {Name: "shield", Description: "Activate a protective barrier for 2 turns (-20 energy)"},
	ActionIdentify: {Name: "identify", Description: "Prove your identity against a Duplicate Demon (-10 energy)"},
	ActionStatus:   {Name: "status", Description: "Display your current state and attributes"},
	ActionWait:     {Name: "wait", Description: "Rest and restore 15 energy"},
	ActionJournal:  {Name: "journal", Description: "Review your journey log"},
	ActionHelp:     {Name: "help", Description: "Show available commands"},
	ActionQuit:     {Name: "quit", Description: "End your journey"},
}

// actionsByName is a reverse lookup from name string to Action.
var actionsByName map[string]Action

func init() {
	actionsByName = make(map[string]Action, len(actionRegistry))
	for action, info := range actionRegistry {
		actionsByName[info.Name] = action
	}
}

// ActionFromString parses user input into an Action.
// Returns ok=false if the input is not a recognized action.
func ActionFromString(s string) (Action, bool) {
	action, ok := actionsByName[s]
	return action, ok
}

// AllActions returns all available actions in display order.
func AllActions() []Action {
	return []Action{
		ActionMove, ActionScan, ActionShield, ActionIdentify,
		ActionStatus, ActionWait, ActionJournal, ActionHelp, ActionQuit,
	}
}

// String returns the action keyword.
func (a Action) String() string {
	if info, ok := actionRegistry[a]; ok {
		return info.Name
	}
	return "unknown"
}

// Description returns a human-readable explanation of the action.
func (a Action) Description() string {
	if info, ok := actionRegistry[a]; ok {
		return info.Description
	}
	return "unknown action"
}
