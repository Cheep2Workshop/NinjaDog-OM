package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"nd-lobby/tts"
	"net/http"
	"time"

	lobbyres "github.com/cheep2workshop/ninjadog-om/nd-lobby-res"

	"google.golang.org/grpc"
	"open-match.dev/open-match/pkg/pb"
)

const (
	omFrontendEndpoint = "om-frontend.open-match.svc.cluster.local:50504"
)

// start match making
func startMatchMake(res http.ResponseWriter, req *http.Request) {
	var mode string
	mode = req.FormValue("mode")
	// assign default mode
	if mode == "" {
		mode = "private"
	}

	fe, conn, err := newClient(omFrontendEndpoint)
	if err != nil {
		log.Fatalf("Falied create frontend service client, got %s", err.Error())
	}
	defer conn.Close()

	tReq := &pb.CreateTicketRequest{
		Ticket: generateTicket(mode),
	}

	resp, err := fe.CreateTicket(context.Background(), tReq)
	if err != nil {
		log.Fatalf("Failed to create ticket, got %v", err)
	}
	r := lobbyres.StartMatchMakeRes{
		TicketID: resp.Id,
	}

	err = json.NewEncoder(res).Encode(r)
	if err != nil {
		log.Fatalf("Failed encode json, got %s", err.Error())
	}
	log.Println("Ticket created successfully, id:", resp.Id)
	tts.UpdateTicketTimestamp(resp.Id)
}

// check if ticket is still existed
func getMatchMakeProcess(res http.ResponseWriter, req *http.Request) {
	// get the ticket id of player
	var ticketID string
	ticketID = req.FormValue("id")

	fe, conn, err := newClient(omFrontendEndpoint)
	if err != nil {
		log.Fatalf("Falied create frontend service client, got %s", err.Error())
	}
	defer conn.Close()

	tReq := &pb.GetTicketRequest{
		TicketId: ticketID,
	}

	// create request for getting ticket from open-match
	resp, err := fe.GetTicket(context.Background(), tReq)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		log.Printf("Failed to get ticket (%s), got %v\n", ticketID, err.Error())
	} else {
		// generate response
		var conn string
		if resp.Assignment == nil {
			conn = ""
		} else {
			conn = resp.Assignment.Connection
		}
		r := lobbyres.GetMatchMakeProcessRes{
			Status:     0,
			Assignment: conn,
			ErrMsg:     "Success",
		}
		err = json.NewEncoder(res).Encode(r)
		if err != nil {
			log.Fatalf("Failed encode json, got %s", err.Error())
		}
		tts.UpdateTicketTimestamp(ticketID)
	}
}

// cancel match making
func cancelMatchMake(res http.ResponseWriter, req *http.Request) {
	// get the ticket id of player
	var ticketID string
	ticketID = req.FormValue("id")

	_, err := deleteTicket(ticketID)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(res, err.Error())
	} else {
		r := lobbyres.CancelMatchMakeRes{
			Status: 0,
			ErrMsg: "Success",
		}
		code, err := json.Marshal(r)
		if err != nil {
			log.Fatalf("Falied to encode json, got %s", err.Error())
		}
		res.Write(code)
		log.Println("Cancel match making successfully, id:", ticketID)
	}
}

func deleteTicket(ticketID string) (string, error) {
	fe, conn, err := newClient(omFrontendEndpoint)
	if err != nil {
		log.Fatalf("Falied create frontend service client, got %s", err.Error())
	}
	defer conn.Close()

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
		tts.DeleteTicketTimestamp(ticketID)
		return "Ticket delete successfully", nil
	}
}

func refreshTickets() {
	for range time.Tick(time.Second * 1) {
		for _, ticket := range tts.RefreshTicketTimestamps() {
			deleteTicket(ticket)
		}
	}
}

func newClient(endPoint string) (pb.FrontendServiceClient, *grpc.ClientConn, error) {
	conn, err := grpc.Dial(endPoint, grpc.WithInsecure())
	if err != nil {
		// log.Fatalf("Failed to connect to Open Match, got %v", err)
		conn.Close()
		return nil, nil, fmt.Errorf("Failed to connect to Open Match, got %v", err)
	}

	fe := pb.NewFrontendServiceClient(conn)
	return fe, conn, nil
}
