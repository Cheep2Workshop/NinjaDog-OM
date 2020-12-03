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
	fe, err := newClient(omFrontendEndpoint)
	if err != nil {
		log.Fatalf("Falied create frontend service client, got %s", err.Error())
	}

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

	err = json.NewEncoder(res).Encode(r)
	//code, err := json.Marshal(r)
	if err != nil {
		log.Fatalf("Failed encode json, got %s", err.Error())
	}
	// res.Write(code)
	log.Println("Ticket created successfully, id:", resp.Id)
}

// check if ticket is still existed
func getMatchMakeProcess(res http.ResponseWriter, req *http.Request) {
	// get the ticket id of player
	var ticketID string
	ticketID = req.FormValue("id")

	fe, err := newClient(omFrontendEndpoint)
	if err != nil {
		log.Fatalf("Falied create frontend service client, got %s", err.Error())
	}

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
		r := GetMatchMakeProcessRes{
			Status:     0,
			Assignment: conn,
			ErrMsg:     "Success",
		}
		err = json.NewEncoder(res).Encode(r)
		if err != nil {
			log.Fatalf("Failed encode json, got %s", err.Error())
		}

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
		r := CancelMatchMakeRes{
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
	fe, err := newClient(omFrontendEndpoint)
	if err != nil {
		log.Fatalf("Falied create frontend service client, got %s", err.Error())
	}

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

func newClient(endPoint string) (pb.FrontendServiceClient, error) {
	conn, err := grpc.Dial(endPoint, grpc.WithInsecure())
	if err != nil {
		// log.Fatalf("Failed to connect to Open Match, got %v", err)
		conn.Close()
		return nil, fmt.Errorf("Failed to connect to Open Match, got %v", err)
	}
	defer conn.Close()

	fe := pb.NewFrontendServiceClient(conn)
	return fe, nil
}

// StartMatchMakeRes is response of startmatchmake
type StartMatchMakeRes struct {
	TicketID string `json:"ticketid,string,omitempty"`
	ErrMsg   string `json:"errmsg,omitempty"`
}

// CancelMatchMakeRes is response of cancelmatchmake
type CancelMatchMakeRes struct {
	Status int32  `json:"status"`
	ErrMsg string `json:"errmsg,omitempty"`
}

// GetMatchMakeProcessRes is response of getmatchmakeprocess
type GetMatchMakeProcessRes struct {
	Status     int32  `json:"status"`
	Assignment string `json:"assignment,omitempty"`
	ErrMsg     string `json:"errmsg,omitempty"`
}
