package tts

import "time"

const idleTicketTolerance float64 = 10

var ticketTimestamps map[string]time.Time

// InitTicketTimestampes initailize ticket timestamps
func InitTicketTimestampes() {
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
		duration := timestamp.Sub(now)
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
