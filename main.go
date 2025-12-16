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
	//chk the mtd type
	//  if GET, it's valid
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusTeapot)
		return
	}

	//open the client binary
	file, err := os.Open("dClient")
	if err != nil {
		http.Error(w,
			"ERR! Cannot find client binary",
			http.StatusNotFound)
		merr("err openning binary:  ", err)
		return
	}

	//close file when fn ends
	defer file.Close()

	//let client know it's a binary 
	w.Header().Set("Content-Type", 
		"application/octet-stream")

	//send the binary to the client
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
	//chk the mtd type
	//  if POST, it's valid 
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusTeapot)
		return
	}

	//read req body 
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w,
			"Err reading request body",
			http.StatusInternalServerError)
		return
	}

	curTimeR := time.Now()                  //get current time then
	curTime := curTimeR.Format("15:04:05")  //  set the format
	month := curTimeR.Month()
	day := curTimeR.Day()
	year := curTimeR.Year()
	fileDir := fmt.Sprintf("%d/%s",         //construct file dir from 
		year, month)                          //  eg: 2025/November
	filePath := fmt.Sprintf("%s/%d.md",     //construct filepath
		fileDir, day)                         //  eg: 2025/November/2.md
	line := fmt.Sprintf("`%s` - %s\n",//construct the line
		curTime, string(body))                        //  eg: `9:27:34` - foo
	wrl(filePath + ":")                     //log the filepath
	wrl("  " + line)                        //log the line

	//chk if dir exists
	_, err = os.Stat(fileDir)
	if os.IsNotExist(err) {
		//create dir if not
		err = os.MkdirAll(fileDir, 0777) 
		hanFrr(err)
	}
	hanFrr(err)
	
	//open file in append mode,
	//  create if not exist
	file, err := os.OpenFile(filePath,
		os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0777)
	hanFrr(err)
	//close file when fn ends 
	defer file.Close()

	//write the line to the file
	_, err = file.WriteString(line)
	hanFrr(err)

	//finally, respond to client
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("recieved"))
}
