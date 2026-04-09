package domain

import (
	"fmt"
	"time"
)

// JournalEntry records a single event in the token's journey.
type JournalEntry struct {
	Timestamp time.Time
	Gate      Gate
	Action    string
	Message   string
}

// String formats the entry for display.
func (e JournalEntry) String() string {
	return fmt.Sprintf("[%s] %s - %s: %s",
		e.Timestamp.Format("15:04:05"),
		e.Gate,
		e.Action,
		e.Message,
	)
}

// Journal is an ordered, append-only collection of journey events.
type Journal struct {
	entries []JournalEntry
}

// NewJournal creates an empty journal.
func NewJournal() *Journal {
	return &Journal{}
}

// Record appends a new entry to the journal.
func (j *Journal) Record(gate Gate, action, message string) {
	j.entries = append(j.entries, JournalEntry{
		Timestamp: time.Now(),
		Gate:      gate,
		Action:    action,
		Message:   message,
	})
}

// Entries returns a copy of all journal entries.
func (j *Journal) Entries() []JournalEntry {
	out := make([]JournalEntry, len(j.entries))
	copy(out, j.entries)
	return out
}

// Len returns the number of entries.
func (j *Journal) Len() int {
	return len(j.entries)
}
