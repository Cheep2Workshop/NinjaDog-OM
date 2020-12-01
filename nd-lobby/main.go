package main

import (
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("Start lobby server")
	http.HandleFunc("/startmatchmake", startMatchMake)
	// http.HandleFunc("/getmatchmakeprocess", getMatchMakeProcess)
	// http.HandleFunc("/cancelmatchmake", cancelMatchMake)

	// go refreshTickets()

	// http.HandleFunc("/hello1", hello1)

	http.ListenAndServe(":8080", nil)
}
