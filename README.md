<h1 align="center">d</h1>

>[!WARNING]
>This was made specifically for my usecase. I am only open-sourcing this in case it benefits anyone else

This is purely made to serve my usecase.

### How it works
Short story long:
- This is a simple HTTP server and client written in Go that only excepts `POST` requests.
- It then takes the body of the post request and appends it to the end of a markdown file
- Then it writes the markdown file to a path with the current date (eg: `2025/October/21.md`).

Short story shorter:
- Recieve POST request
- Write the body of request to a file

---

## Installation:

### Server 
- Download the repo
- Compile the server
  (make sure you're in the root directory of the repo):
  ```go
  go build .
  ```
- Create a systemd service
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

### Client
- In order to download the client binary from your server (optional, but recommended):
  - Compile the client binary
    ```sh
    cd client && go build .
    ```
  - Move the client binary to your library dir (and rename it to `dClient` 
    ```sh
    mv d /your/library/dir/dClient
    ```
  - Download binary from server
    ```sh
    wget https://your.server.address:8008/d
    ```
- Move the `d` client to a location in your path:
  
    eg: `/usr/bin`
    ```sh
    mv d /usr/bin
    ```
- Make the binary executable (may need `su`):
    ```sh
    chmod a+x d
    ```
- Put your `d` server address in your config:
    ```toml
    [server]
    address = "https://your.server.address/"
    ```
