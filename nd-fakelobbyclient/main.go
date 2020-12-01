package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
)

func main() {
	// test()
	lobbyClient()
}

func test() {
	var urlflag = flag.String("u", "localhost:8080/hello1", "Url of http server")
	var name = flag.String("n", "user", "Name of user")
	flag.Parse()

	if !strings.Contains(*urlflag, "http://") {
		*urlflag = fmt.Sprintf("http://%s", *urlflag)
	}
	resp, err := http.PostForm(*urlflag, url.Values{"name": {*name}})
	if err != nil {
		log.Fatalf("Failed to post, got %s", err.Error())
	}
	defer resp.Body.Close()
	// fmt.Printf("Got status code of response : %d", resp.StatusCode)
	res := &hello1Res{}
	err = json.NewDecoder(resp.Body).Decode(res)
	if err != nil {
		log.Fatalf("Falied decode json, got %s", err.Error())
	}
	fmt.Println(*res)
}

type hello1Res struct {
	Msg string
}

func lobbyClient() {
	var urlflag = flag.String("u", "localhost:8080", "Url of http server")
	var req = flag.String("r", "start", "Request")
	// var id = flag.String("i", "", "Ticket ID")
	flag.Parse()

	if !strings.Contains(*urlflag, "http://") {
		*urlflag = fmt.Sprintf("http://%s", *urlflag)
	}

	msg := ""
	switch *req {
	case "start":
		msg = "startmatchmake"
		startmatchmake(*urlflag)
	case "get":
		msg = "getmatchmakeprocess"
		// getmatchmakeprocess(*url, *id)
	case "canel":
		msg = "canelmatchmake"
	}

	fmt.Printf("%s/%s", *urlflag, msg)
}

func startmatchmake(url string) {
	resp, err := http.PostForm(url+"/startmatchmake", nil)
	if err != nil {
		log.Fatalf("Failed to start match make, got %s", err.Error())
	}
	defer resp.Body.Close()

	smmRes := &StartMatchMakeRes{}
	err = json.NewDecoder(resp.Body).Decode(smmRes)
	if err != nil {
		log.Fatalf("Falied decode json, got %s", err.Error())
	}
	fmt.Println(smmRes.TicketID)
}

type StartMatchMakeRes struct {
	TicketID string
	ErrMsg   string
}
