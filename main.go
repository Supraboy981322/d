package main

import (
	"fmt"
	"net/http"
	"os"
	"time"
	"io/ioutil"
	"io"
)

func main() {
	wrl("starting \033[32md\033[0m...")
	http.HandleFunc("/d", serveClient)
	http.HandleFunc("/", post)
	
	wrl("started.")
	ferr(http.ListenAndServe(":8008", nil))
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
		merr("err opening file:  ", err)
		return
	}

	defer file.Close()

	w.Header().Set("Content-Type", 
		"application/octet-stream")

	_, err = io.Copy(w, file)
	if err != nil {
		http.Error(w,
			"err serving binary:  " + err.Error(),
			http.StatusInternalServerError)
		merr("err serving binary:  ", err)
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
		hanFrr(err)
	}
	hanFrr(err)
	file, err := os.OpenFile(filePath,
		os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0777)
	hanFrr(err)
	defer file.Close()

	_, err = file.WriteString(line)
	hanFrr(err)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("recieved"))
}
