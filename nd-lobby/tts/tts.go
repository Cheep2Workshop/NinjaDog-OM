package tts

import (
	"time"
)

const idleTicketTolerance float64 = 10

var ticketTimestamps map[string]time.Time

// InitTicketTimestamps initailize ticket timestamps
func InitTicketTimestamps() {
	ticketTimestamps = make(map[string]time.Time)
}

// UpdateTicketTimestamp update timestamp of ticket
func UpdateTicketTimestamp(ticketID string) {
	ticketTimestamps[ticketID] = time.Now()
}

// RefreshTicketTimestamps refresh and check if ticket timestamp expired
func RefreshTicketTimestamps() []string {
	var deletedTickets []string
	now := time.Now()
	for ticket, timestamp := range ticketTimestamps {
		duration := now.Sub(timestamp)
		//fmt.Printf("Ticket(%s) duration = %s", ticket, duration)
		if duration.Seconds() > idleTicketTolerance {
			deletedTickets = append(deletedTickets, ticket)
		}
	}
	return deletedTickets
}

// DeleteTicketTimestamp delete ticketstamp
func DeleteTicketTimestamp(ticketID string) {
	delete(ticketTimestamps, ticketID)
}
