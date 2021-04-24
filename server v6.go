package main

import (
	"./leap"
	"fmt"
	"log"
	"encoding/json"
	//"github.com/yyasha/LeapMotionGo"
	"net/http"
	"runtime"
)

var Result []uint8

func main() {
	numcpu := runtime.NumCPU()
	runtime.GOMAXPROCS(numcpu)

	c, err := leap.Connect("localhost")
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	go http.HandleFunc("/api", jsonServer)
	go http.ListenAndServe(":80", nil)

	for {
		f, err := c.Frame()
		if err != nil {
			log.Fatal(err)
		}
		Result, _ = json.Marshal(f)
		fmt.Println(string(Result))
	}
}

func jsonServer(w http.ResponseWriter, r *http.Request) {

    switch r.Method {
    case "GET":
            w.Write(Result)
    default:
            w.WriteHeader(http.StatusMethodNotAllowed)
            fmt.Fprintf(w, "I can't do that.")
    }
}
