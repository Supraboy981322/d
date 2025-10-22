# d

>[!WARNING]
>This is not intended for other people to use, but it is available anyways incase it is of help

Short story long: This is a simple HTTP server and client written in Go that only excepts `POST` requests. It then takes the body of the post request and appends it to the end of a markdown file with a path matching the current date (ex: `2025/10/21.md` for the 21st of October in 2025). This will likely never be updated, and is purely made to serve my usecase.

---

### in case you, for whatever reason, want to run it:
- download the repo
- edit the `main.go` file the `client` directory to replace the value of the `url` var with your own (please note that, by default, the port is `8008`):
```go
    url string = "http://your.server.address:8008"
```
- compile the server (make sure you're in the root directory of the repo):
```go
go build .
```
- compile the client (make sure you're in the directory for the client):
```go
go build .
```
- create a systemd service
```service
[Unit]
Description=post request diary server
After=network.target

[Service]
ExecStart=/your/server/dir/d
WorkingDirectory=/your/library/dir
Restart=always

[Install]
WantedBy=multi-user.target
```
- copy the client binary to your client computer in any way you like
