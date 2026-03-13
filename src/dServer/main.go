package main

import (
	"fmt"
	"net/http"
	"os"
	_"embed"
	"time"
	"io/ioutil"
	"strconv"
	"io"
	"bytes"
	"strings"
	"encoding/json"
	"compress/gzip"
	brotli "github.com/google/brotli/go/cbrotli"
)

//go:embed web_built/amalgamation.html
var web_spa []byte

type Msg struct {
	Timestamp string
	Msg string
}
var (
	todays_lines []Msg
	today int
	compression_priority = []string {
		"brotli", "br",
		"gzip", "gz",
	}
)

func main() {
	todays_lines = []Msg{}
	today = time.Now().Day()

	go time_tracker()

	wrl("starting \033[32md\033[0m...")
	http.HandleFunc("/d", serveClient)
	http.HandleFunc("/post", post)
	http.HandleFunc("/today", send_today)
	http.HandleFunc("/sync", send_new)
	http.HandleFunc("/", web_ui)
	
	wrl("started.")
	ferr(http.ListenAndServe(":8008", nil))
}

func time_tracker() {
	for {
		if today != time.Now().Day() {
			todays_lines = []Msg{}
			today = time.Now().Day()
		}
		time.Sleep(1 * time.Second)
	}
}

func web_ui(w http.ResponseWriter, r *http.Request) {
	picked_idx := -1
	for i, enc := range strings.Split(r.Header.Get("Accept-Encoding"), ",") {
		if picked_idx < idx_of_str(compression_priority, strings.TrimSpace(enc)) {
			picked_idx = i
		}
	}
	//could've been a simple ternary
	var picked_name string
	if picked_idx > -1 { 
		picked_name = compression_priority[picked_idx];
	}

	var page bytes.Buffer
	var length int
	if picked_idx < 0 { page.Write(web_spa) } else {
		var e error
		switch picked_name {
			case "gzip", "gz":
				wr := gzip.NewWriter(&page)
				length, e = wr.Write(web_spa)
				wr.Flush()
				goto err
		  case "brotli", "br":
				opts := brotli.WriterOptions {
					Quality: 11,
					LGWin: 0,
				}
				wr := brotli.NewWriter(&page, opts)
				length, e = wr.Write(web_spa)
				wr.Flush()
				goto err
			default:
				panic("unknown compression picked: " + picked_name)
		}
		err: if e != nil {
			http.Error(w, "failed to compress: " + e.Error(), 500)
			return
		}
	}

	if length == 0 { length = len(page.Bytes()) } 

	w.Header().Set("Content-Encoding", picked_name)
	w.Header().Set("Content-Length", strconv.Itoa(length))
	w.Write(page.Bytes())
}

func send_today(w http.ResponseWriter, r *http.Request) {
	j, e := json.Marshal(todays_lines)
	if e != nil {
		http.Error(w, "failed to marshal into json: " + e.Error(), 500)
		return
	}
	w.Write(j)
}

func send_new(w http.ResponseWriter, r *http.Request) {
	has_str := r.Header.Get("have")
	has, e := strconv.Atoi(has_str)
	if e != nil {
		err := strings.Split(e.Error(), ": ")
		http.Error(w, "NaN: " + err[len(err)-1], 400)
		return
	}
	if has >= len(todays_lines) {
		A_TERNARY_WOULDVE_BEEN_GREAT := "ies"
		A_TERNARY_WOULDVE_BEEN_GREAT_2 := "are"
		if len(todays_lines) == 1 {
			A_TERNARY_WOULDVE_BEEN_GREAT = "y"
			A_TERNARY_WOULDVE_BEEN_GREAT_2 = "is"
		}
		msg := fmt.Sprintf(
			"there %s only %d entr%s",
			A_TERNARY_WOULDVE_BEEN_GREAT_2, len(todays_lines),
			A_TERNARY_WOULDVE_BEEN_GREAT,
		)
		http.Error(w, msg, 403)
		return
	}
	j, e := json.Marshal(todays_lines[has:])
	if e != nil {
		http.Error(w, "failed to marshal into json: " + e.Error(), 500)
		return
	}
	w.Write(j)
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
	month := curTimeR.Month()               //get the month (file dir)
	day := curTimeR.Day()                   //get the day (filename)
	year := curTimeR.Year()                 //get the year (file dir)
	fileDir := fmt.Sprintf(                 //construct file dir from 
		"%d/%s",                              //  eg: 2025/November
		year, month,                          //
	)                                       //
	filePath := fmt.Sprintf(                //construct filepath
		"%s/%d.md",                           //  eg: 2025/November/2.md
		fileDir, day,                         //
	)                                       //
	line := fmt.Sprintf(                    //construct the line
		"`%s` - %s\n",                        //  eg: `9:27:34` - foo
		curTime, string(body),                //
	)                                       
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
	file, err := os.OpenFile(
		filePath,
		os.O_APPEND|os.O_WRONLY|os.O_CREATE,
		0777,
	)
	hanFrr(err)
	//close file when fn ends 
	defer file.Close()

	//write the line to the file
	_, err = file.WriteString(line)
	hanFrr(err)

	//finally, respond to client
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("recieved"))
	new_line := Msg{
		Timestamp: curTime,
		Msg: string(body),
	}
	todays_lines = append(todays_lines, new_line)
}
