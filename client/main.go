package main

import (
	"net/http"
	"log"
	"strings"
	"time"
	"os"
)

var (
	args []string = os.Args[1:]
	//change this to your address!
	url string = "[address redacted]"
	line string = args[0]
)

func main() {
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

	defer resp.Body.Close()

	log.Printf("Response Status:  %s\n", resp.Status)
}
