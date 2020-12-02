package main

import (
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("Start lobby server")
	http.HandleFunc("/startmatchmake", startMatchMake)
	http.HandleFunc("/getmatchmake", getMatchMakeProcess)
	http.HandleFunc("/cancelmatchmake", cancelMatchMake)

	// refresh and delete tickets periodly
	go refreshTickets()

	http.HandleFunc("/hello1", hello1)

	http.ListenAndServe(":8080", nil)
}
