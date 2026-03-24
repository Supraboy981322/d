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
	"errors"
	"slices"
	"strings"
	"encoding/json"
	"compress/gzip"
	"github.com/yuin/goldmark"
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


func render_md(src []byte) ([]byte, error) {
	var b bytes.Buffer
	e := goldmark.Convert(src, &b)
	if e != nil { return nil, e }
	return b.Bytes(), nil
}

func main() {
	goto start
	
	err: {
		fmt.Fprintln(os.Stderr, "startup err")
		os.Exit(1)
		panic("failed to fail")
	}

	start: {
		todays_lines = []Msg{}
		today = time.Now().Day()

		j, e := os.ReadFile("scrollback.json")
		if e != nil {
			if errors.Is(e, os.ErrNotExist) {
				e = os.WriteFile("scrollback.json", []byte{ '[', ']' }, 0644)
				if e != nil {
					fmt.Fprintf(os.Stderr, "failed to create scrollback.json: %v\n", e)
					goto err
				}
				j = []byte{ '[', ']' }
			} else {
				fmt.Fprintf(os.Stderr, "failed to read scrollback.json: %v\n", e)
				goto err
			}
		}

		e = json.Unmarshal(j, &todays_lines)
		if e != nil {
			fmt.Fprintf(os.Stderr, "failed to unmarshall scrollback.json: %v\n", e)
			goto err
		}
	}

	go ticker()

	wrl("starting \033[32md\033[0m...")
	http.HandleFunc("/d", serveClient)
	http.HandleFunc("/post", post)
	http.HandleFunc("/today", send_today)
	http.HandleFunc("/sync", send_new)
	http.HandleFunc("/", web_ui)
	
	wrl("started.")
	ferr(http.ListenAndServe(":8008", nil))
}

func ticker() {
	previous_length := len(todays_lines)
	for {
		if today != time.Now().Day() {
			todays_lines = []Msg{}
			today = time.Now().Day()
		}
		if len(todays_lines) != previous_length {
			j, e := json.Marshal(todays_lines)
			if e != nil { panic(e.Error()) }

			e = os.WriteFile("scrollback.json", j, 0644)
			if e != nil { panic(e.Error()) }
		}
		previous_length = len(todays_lines)
		time.Sleep(1 * time.Second)
	}
}

func compress(w http.ResponseWriter, r *http.Request, og []byte) ([]byte, error) {
	picked_idx := -1 
	{
		accepted := r.Header.Get("Accept-Encoding")

		//new slice without any whitespace
		accepted_trimmed := []string{}
		for _, thing := range strings.Split(accepted, ",") {
			accepted_trimmed = append(accepted_trimmed, strings.TrimSpace(thing))
		}

		loop: for i, enc := range compression_priority {
			if slices.Contains(accepted_trimmed, enc) {
				picked_idx = i; break loop
			}
		}
	}

	//could've been a simple ternary
	var picked_name string
	if picked_idx > -1 {
		picked_name = compression_priority[picked_idx];
	}

	var page bytes.Buffer
	if picked_idx < 0 { page.Write(og) } else {
		var e error
		switch picked_name {
			case "gzip", "gz": {
				wr := gzip.NewWriter(&page)
				_, e = wr.Write(og)
				wr.Flush()
				goto err
			}
		  case "brotli", "br": {
				opts := brotli.WriterOptions {
					Quality: 11,
					LGWin: 24,
				}
				wr := brotli.NewWriter(&page, opts)
				_, e = wr.Write(og)
				wr.Flush()
				goto err
			}
			default:
				panic("unknown compression picked: " + picked_name)
		}
		err: if e != nil {
			http.Error(w, "failed to compress: " + e.Error(), 500)
			return nil, e
		}
	}

	w.Header().Set("Content-Encoding", picked_name)
	return page.Bytes(), nil
}

func web_ui(w http.ResponseWriter, r *http.Request) {
	page, e := compress(w, r, web_spa)
	if e != nil { return } //handled by compressor

	w.Header().Set("Content-Length", strconv.Itoa(len(page)))
	w.Write(page)
}

func send_today(w http.ResponseWriter, r *http.Request) {
	j, e := json.Marshal(todays_lines)
	if e != nil {
		http.Error(w, "failed to marshal into json: " + e.Error(), 500)
		return
	}
	json_compressed, e := compress(w, r, j)
	if e != nil { return }
	w.Header().Set("Content-Type", "text/json")
	w.Write(json_compressed)
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
		w.Header().Set("have", strconv.Itoa(len(todays_lines)))
		http.Error(w, msg, 403)
		return
	}
	j, e := json.Marshal(todays_lines[has:])
	if e != nil {
		http.Error(w, "failed to marshal into json: " + e.Error(), 500)
		return
	}

	json_compressed, e := compress(w, r, j)
	if e != nil { return }
	w.Header().Set("Content-Type", "text/json")
	w.Write(json_compressed)
}

func serveClient(w http.ResponseWriter, r *http.Request) {
	//chk the mtd type
	//  if GET, it's valid
	if r.Method != http.MethodGet {
		w.WriteHeader(405)
		w.Write([]byte("method not allowed"))
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
	//CORS
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Headers", "echo, Content-Type")
		return
	}
	//chk the mtd type
	//  if POST, it's valid 
	if r.Method != http.MethodPost {
		w.WriteHeader(405)
		w.Write([]byte("bad method"));
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

	rendered_raw, e := render_md(body)
	if e != nil {
		http.Error(w, "failed to render markdown", 500)
		return
	}
	rendered := strings.TrimSpace(string(rendered_raw))

	var resp []byte
	if r.Header.Get("echo") == "HTML" {
		resp = []byte(rendered)
	} else {
		resp = []byte("recieved")
	}

	//finally, respond to client
	w.WriteHeader(http.StatusOK)
	w.Write(resp)

	new_line := Msg{
		Timestamp: curTime,
		Msg: rendered,
	}

	todays_lines = append(todays_lines, new_line)
}
