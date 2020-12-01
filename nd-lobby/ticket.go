package main

import (
	"fmt"

	"open-match.dev/open-match/pkg/pb"
)

func generateTicket() *pb.Ticket {
	fmt.Println("Generating ticket ...")
	mode := "mode.dev"
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
