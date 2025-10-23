package main

import (
	"net/http"
	"log"
	"strings"
	"time"
	"os"
	"io"
)

var (
	//change this to your address!
	url string = "http://localhost:8008"
	line string
)

func main() {
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))
	if len(os.Args) > 1 {
		line = os.Args[1]
	} else {
		log.Printf("no input")
		return
	}

	client := &http.Client{
		Timeout: time.Second * 10,
	}

	req, err := http.NewRequest("POST", url, strings.NewReader(line))
	if err != nil {
		log.Fatalf("err sending request:  %v\n", err)
	}

	req.Header.Set("Content-Type", "text/plain")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("err sending request: %v\n", err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("failed to read body:  %v\n", err)
	}
	log.Println(string(body))

	defer resp.Body.Close()
}
