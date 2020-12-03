package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	lobbyres "github.com/cheep2workshop/ninjadog-om/nd-lobby-res"
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
	var req = flag.String("r", "startmatchmake", "Request")
	var id = flag.String("i", "", "Ticket ID")
	flag.Parse()

	if !strings.Contains(*urlflag, "http://") {
		*urlflag = fmt.Sprintf("http://%s", *urlflag)
	}

	msg := ""
	switch *req {
	case "startmatchmake":
		msg = "startmatchmake"
		startmatchmake(*urlflag)
	case "getmatchmake":
		msg = "getmatchmake"
		getmatchmakeprocess(*urlflag, *id)
	case "cancelmatchmake":
		msg = "canelmatchmake"
		cancelmatchmake(*urlflag, *id)
	}

	fmt.Printf("%s/%s", *urlflag, msg)
}

func startmatchmake(url string) {

	resp, err := http.Post(url+"/startmatchmake", "application/json", nil)
	fmt.Println("Send start match make request,", url)
	if err != nil {
		fmt.Println("Failed to start match make, got ", err.Error())
		log.Fatalf("Failed to start match make, got %s", err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		fmt.Println("Faild to connect to url: 404 not found")
		return
	}
	fmt.Println("1")

	var smmRes *lobbyres.StartMatchMakeRes = &lobbyres.StartMatchMakeRes{}
	// err = json.NewDecoder(resp.Body).Decode(smmRes)
	// err = json.Unmarshal(resp.Body, &smmRes)
	err = getJSON(*resp, smmRes)
	if err != nil {
		log.Fatalf("Falied decode json, got %s", err.Error())
	}
	fmt.Println("Got ticket:", smmRes.TicketID)
}

func cancelmatchmake(u string, ticketID string) {
	fmt.Println("Send cancel match make request")
	resp, err := http.PostForm(u+"/cancelmatchmake", url.Values{"id": {ticketID}})
	if err != nil {
		log.Fatalf("Falied to cancel match make, got %s", err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		res := &lobbyres.CancelMatchMakeRes{}
		err := json.NewDecoder(resp.Body).Decode(res)
		if err != nil {
			log.Fatalf("Falied to decode json, got %s", err.Error())
		}
		fmt.Println(res)
	}
}

func getmatchmakeprocess(u string, ticketID string) {
	resp, err := http.PostForm(u+"/getmatchmake", url.Values{"id": {ticketID}})
	if err != nil {
		log.Fatalf("Falied to get match make process, got %s", err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		res := &lobbyres.GetMatchMakeProcessRes{}
		err := json.NewDecoder(resp.Body).Decode(res)
		if err != nil {
			log.Fatalf("Falied to decode json, got %s", err.Error())
		}
		fmt.Println(res)
	} else {
		log.Fatalf("Recv status code:%d", resp.StatusCode)
	}
}

func getJSON(r http.Response, target interface{}) error {
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}
