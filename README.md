# d

>[!WARNING]
>This was made specifically for my usecase. I am only open-sourcing this in case it benefits anyone else

Short story long: This is a simple HTTP server and client written in Go that only excepts `POST` requests. It then takes the body of the post request and appends it to the end of a markdown file with a path matching the current date (ex: `2025/10/21.md` for the 21st of October in 2025). This is purely made to serve my usecase.

---

## Installation:

### server 
- download the repo
- compile the server
  (make sure you're in the root directory of the repo):
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

### client
- In order to download the client binary from your server (optional, but recommended):
  - compile the client binary
    ```sh
    cd client && go build .
    ```
  - move the client binary to your library dir (and rename it to `dClient` 
    ```sh
    mv d /your/library/dir/dClient
    ```
  - download binary from server
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
