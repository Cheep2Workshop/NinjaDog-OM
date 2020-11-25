package main

import (
	"context"
	"errors"
	"fmt"

	agonesv1 "agones.dev/agones/pkg/apis/agones/v1"
	allocationv1 "agones.dev/agones/pkg/apis/allocation/v1"
	"agones.dev/agones/pkg/client/clientset/versioned"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"open-match.dev/open-match/pkg/pb"
)

func createOMAssignTicketRequest(match *pb.Match, gsa *allocationv1.GameServerAllocation) *pb.AssignTicketsRequest {
	tids := []string{}
	for _, t := range match.GetTickets() {
		tids = append(tids, t.GetId())
	}

	fmt.Printf("Gameserver (%s) : %s", gsa.Status.GameServerName, gsa.Status.Address)
	if len(gsa.Status.Ports) <= 0 {

		return nil
	}

	return &pb.AssignTicketsRequest{
		Assignments: []*pb.AssignmentGroup{
			{
				TicketIds: tids,
				Assignment: &pb.Assignment{
					Connection: fmt.Sprintf("%s:%d", gsa.Status.Address, gsa.Status.Ports[0].Port),
				},
			},
		},
	}
}

// Create allocation of specific fleet label and game server name
func createAgonesGameServerAllocation() *allocationv1.GameServerAllocation {
	return &allocationv1.GameServerAllocation{
		Spec: allocationv1.GameServerAllocationSpec{
			Required: metav1.LabelSelector{
				// match label = agones.GroupName + "/fleet"
				// may use other label such as : private, battle, etc ...
				MatchLabels: map[string]string{agonesv1.FleetNameLabel: "simple-udp"},
			},
		},
	}
}

// Allocate game server from agones and create assignment, then deliver the assignment to beckend.
func assign(be pb.BackendServiceClient, agonesClient *versioned.Clientset, match *pb.Match) error {
	gsa, err := agonesClient.AllocationV1().GameServerAllocations("default").Create(createAgonesGameServerAllocation())
	if err != nil {
		return err

	}
	// if gsa.Status.State != allocationv1.GameServerAllocationAllocated {
	// 	fmt.Printf("Failed to allocate game server.\n")
	// }

	assignTicketReq := createOMAssignTicketRequest(match, gsa)
	if assignTicketReq != nil {
		for _, assign := range assignTicketReq.GetAssignments() {
			connection := assign.GetAssignment().GetConnection()
			fmt.Printf("Assignment ticket req connection:%s", connection)
			fmt.Println()
		}

		fmt.Println("Assigned ticket")
		resp, err := be.AssignTickets(context.Background(), assignTicketReq)
		if err != nil {
			// Corner case where we allocated a game server for players who left the queue after some waiting time.
			// Note that we may still leak some game servers when tickets got assigned but players left the queue before game frontend announced the assignments.
			if err = agonesClient.AgonesV1().GameServers("default").Delete(gsa.Status.GameServerName, &metav1.DeleteOptions{}); err != nil {
				return err
			}

			if len(resp.GetFailures()) > 0 {
				fmt.Printf("Assignment failed : %d", len(resp.GetFailures()))
				fmt.Println()
			}
		}
		return nil
	} else {
		return errors.New("No port of allocated server")
	}
}
