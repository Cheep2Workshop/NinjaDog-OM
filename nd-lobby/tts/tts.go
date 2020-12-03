package tts

import "time"

const idleTicketTolerance float64 = 10

var ticketTimestamps map[string]time.Time

func initTicketTimestampes() {
	ticketTimestamps = make(map[string]time.Time)
}

func addTicketTimestamp(ticketID string) {

	ticketTimestamps[ticketID] = time.Now()
}

func refreshTicketTimestamps() []string {
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

func deleteTicketTimestamp(ticketID string) {
	delete(ticketTimestamps, ticketID)
}
