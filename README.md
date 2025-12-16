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

- Use Go's built-in installer (requires `git`)
  ```sh
  go install github.com/Supraboy981322/d/src/dServer@latest
  ```

- Check the path the server binary was installed to (usually `~/go/bin/dServer`)

- Create a systemd service (replace `/home/username/go/bin/dServer` with your install path, and `/your/library/dir` with the path to your library)
  ```service
  [Unit]
  Description=post request diary server
  After=network.target
  
  [Service]
  ExecStart=/home/username/go/bin/dServer
  WorkingDirectory=/your/library/dir
  Restart=always
  
  [Install]
  WantedBy=multi-user.target
  ```

### Client

- Using Go's built-in installer
  ```sh
  go install github.com/Supraboy981322/d/src/d@latest
  ```

- Put your `d` server address in your config:
  ```toml
  [server]
  address = "https://your.server.address/"
  ```

- Ensure that GOBIN (the directory used by `go install`) is in your `$PATH`

  eg: Bash with `.bashrc`
  ```sh
  printf "PATH=$PATH:$HOME/go/bin\n" | tee -a ~/.bashrc
  ```
