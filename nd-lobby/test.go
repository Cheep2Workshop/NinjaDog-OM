package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// post
func hello1(resWriter http.ResponseWriter, req *http.Request) {
	if err := req.ParseForm(); err != nil {
		log.Fatalf("Failed parse form of post, got %s", err.Error())
	}
	msg := fmt.Sprintf("Hello, %s!", req.FormValue("name"))
	resWriter.Header().Set("Content-type", "application/json")
	res := hello1Res{
		Msg: msg,
	}

	fmt.Println(res)
	js, err := json.Marshal(res)
	if err != nil {
		log.Fatalf("Falied to parse response to json, got %s", err.Error())
	}
	resWriter.Write(js)
	fmt.Println(len(js))
}

type hello1Res struct {
	Msg string
}

// get
func hello2(resWriter http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	name, err := req.Form["name"]
	if !err {
		fmt.Println("Hello2 failed")
	}
	myname := "yoooo"
	fmt.Fprintf(resWriter, fmt.Sprintf("Hello, %v !", name[0]))
	fmt.Fprintln(resWriter)
	fmt.Fprintf(resWriter, fmt.Sprintf("Hello, %v !", myname))
}
