package main

import (
	"fmt"
	"net/http"
	"os"
	"time"
	"io/ioutil"
	"log"
)

func main() {
	fmt.Println("starting `d`")
	http.HandleFunc("/d", serveClient)
	http.HandleFunc("/", post)
	
	log.Fatal(http.ListenAndServe(":8008", nil))
}

func serveClient(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusTeapot)
		return
	}

	file, err := os.Open("dClient")
	if err != nil {
		http.Error(w,
			"ERR! Cannot find client binary",
			http.StatusNotFound)
		log.Printf("err opening file:  %v", err)
		return
	}

	defer file.Close()

	w.Header().Set("Content-Type", 
		"application/octet-stream")

	_, err = io.Copy(w, file)
	if err != nil {
		http.Error(w,
			fmt.Sprintf("err serving binary:  %v", err),
			http.StatusInternalServerError)
		log.Printf("err serving binary:  %v", err)
		return
	}
}

func post(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusTeapot)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w,
			"Err reading request body",
			http.StatusInternalServerError)
		return
	}

	curTime := time.Now()
	year := curTime.Year()
	month := curTime.Month()
	day := curTime.Day()
	hour := curTime.Hour()
	minute := curTime.Minute()
	second := curTime.Second()
	fileDir := fmt.Sprintf("%d/%s",
		year, month)
	filePath := fmt.Sprintf("%d/%s/%d.md",
		year, month, day)
	line := fmt.Sprintf("`%d:%d:%d` - %s\n",
		hour, minute, second, body)
	fmt.Printf("%s:\n", filePath)
	fmt.Printf("  %s", line)
	_, err = os.Stat(fileDir)
	if os.IsNotExist(err) {
		err = os.MkdirAll(fileDir, 0777) 
		if err != nil {
			log.Fatalf(
				"err creating directory:  %v\n", err)
			return
		}
	}
	if err != nil {
		log.Fatalf(
			"err checking if directory even exists:  %v\n",
			err)
		return
	}
	file, err := os.OpenFile(filePath,
		os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		log.Fatalf("err openning file:  %v", err)
		return
	}
	defer file.Close()

	_, err = file.WriteString(line)
	if err != nil {
		log.Fatalf("err writting to file:  %v\n", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("recieved"))
}
