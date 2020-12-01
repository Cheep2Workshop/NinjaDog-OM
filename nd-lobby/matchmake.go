package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"google.golang.org/grpc"
	"open-match.dev/open-match/pkg/pb"
)

const (
	omFrontendEndpoint = "om-frontend.open-match.svc.cluster.local:50504"
)

// start match making
func startMatchMake(res http.ResponseWriter, req *http.Request) {
	conn, err := grpc.Dial(omFrontendEndpoint, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to Open Match, got %v", err)
	}
	defer conn.Close()

	fe := pb.NewFrontendServiceClient(conn)
	tReq := &pb.CreateTicketRequest{
		Ticket: generateTicket(),
	}

	resp, err := fe.CreateTicket(context.Background(), tReq)
	if err != nil {
		log.Fatalf("Failed to create ticket, got %v", err)
	}
	r := StartMatchMakeRes{
		TicketID: resp.Id,
	}
	code, err := json.Marshal(r)
	if err != nil {
		log.Fatalf("Failed encode json, got %s", err.Error())
	}
	res.Write(code)
	log.Println("Ticket created successfully, id:", resp.Id)
}

// check if ticket is still existed
func getMatchMakeProcess(res http.ResponseWriter, req *http.Request) {
	// get the ticket id of player
	var ticketID string
	if err := json.NewDecoder(req.Body).Decode(&ticketID); err != nil {
		req.Body.Close()
	}

	// connect to open-match frontend
	conn, err := grpc.Dial(omFrontendEndpoint, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to Open Match, got %v", err)
	}
	defer conn.Close()

	// create client of open-match frontend
	fe := pb.NewFrontendServiceClient(conn)
	tReq := &pb.GetTicketRequest{
		TicketId: ticketID,
	}

	// create request for getting ticket from open-match
	resp, err := fe.GetTicket(context.Background(), tReq)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(res, "Faild to get ticket in pool.")
		log.Printf("Failed to get ticket (%s), got %v\n", ticketID, err.Error())
	} else if resp.Assignment == nil {
		fmt.Fprintf(res, "Keep match making")
	} else if resp.Assignment != nil {
		fmt.Fprintf(res, fmt.Sprintf("Find match:%s", resp.Assignment.Connection))
	}
}

// cancel match making
func cancelMatchMake(res http.ResponseWriter, req *http.Request) {
	// get the ticket id of player
	var ticketID string
	if err := json.NewDecoder(req.Body).Decode(&ticketID); err != nil {
		req.Body.Close()
	}

	msg, err := deleteTicket(ticketID)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(res, err.Error())
	} else {
		fmt.Fprintf(res, msg)
	}

	// conn, err := grpc.Dial(omFrontendEndpoint, grpc.WithInsecure())
	// if err != nil {
	// 	log.Fatalf("Failed to connect to Open Match, got %v", err)
	// }
	// defer conn.Close()

	// // create client of open-match frontend
	// fe := pb.NewFrontendServiceClient(conn)
	// tReq := &pb.DeleteTicketRequest{
	// 	TicketId: ticketID,
	// }

	// // create request for getting ticket from open-match
	// _, err = fe.DeleteTicket(context.Background(), tReq)
	// if err != nil {
	// 	res.WriteHeader(http.StatusInternalServerError)
	// 	fmt.Fprintf(res, "Falied to delete ticket")
	// 	log.Fatalf("Failed to delete ticket, got %v", err)
	// } else {
	// 	fmt.Fprintf(res, "Ticket delete successfully")
	// 	log.Println("Ticket deleted successfully, id:", ticketID)
	// }
}

func deleteTicket(ticketID string) (string, error) {
	conn, err := grpc.Dial(omFrontendEndpoint, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to Open Match, got %v", err)
	}
	defer conn.Close()

	// create client of open-match frontend
	fe := pb.NewFrontendServiceClient(conn)
	tReq := &pb.DeleteTicketRequest{
		TicketId: ticketID,
	}

	// create request for getting ticket from open-match
	_, err = fe.DeleteTicket(context.Background(), tReq)
	if err != nil {
		log.Fatalf("Failed to delete ticket, got %v", err)
		return "", fmt.Errorf("Falied to delete ticket")
	} else {
		log.Println("Ticket deleted successfully, id:", ticketID)
		return "Ticket delete successfully", nil
	}
}

func refreshTickets() {
	for range time.Tick(time.Second * 1) {
		for _, ticket := range refreshTicketTimestamps() {
			deleteTicket(ticket)
		}
	}
}

type StartMatchMakeRes struct {
	TicketID string
	ErrMsg   string
}
