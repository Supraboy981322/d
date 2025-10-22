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
	http.HandleFunc("/", post)

	log.Fatal(http.ListenAndServe(":8008", nil))
}

func post(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusTeapot)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Err reading request body", http.StatusInternalServerError)
		return
	}

	curTime := time.Now()
	year := curTime.Year()
	month := curTime.Month()
	day := curTime.Day()
	hour := curTime.Hour()
	minute := curTime.Minute()
	fileDir := fmt.Sprintf("%d/%s", year, month)
	filePath := fmt.Sprintf("%d/%s/%d.md", year, month, day)
	line := fmt.Sprintf("`%d:%d` - %s\n", hour, minute, body)
	fmt.Printf("%s:\n", filePath)
	fmt.Printf("  %s", line)
	_, err = os.Stat(fileDir)
	if os.IsNotExist(err) {
		err = os.MkdirAll(fileDir, 0777) 
		if err != nil {
			log.Fatalf("err creating directory:  %v\n", err)
			return
		}
	}
	if err != nil {
		log.Fatalf("err checking if directory even exists:  %v\n", err)
		return
	}
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0777)
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
