package main

import (
	"fmt"

	"open-match.dev/open-match/pkg/pb"
)

func generateTicket(mode string) *pb.Ticket {
	fmt.Printf("Generating ticket (%s) ...", mode)
	fmt.Println()
	ticket := &pb.Ticket{
		SearchFields: &pb.SearchFields{
			Tags: []string{
				mode,
			},
			DoubleArgs: map[string]float64{
				"Rating": 0,
			},
		},
	}
	return ticket
}
