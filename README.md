# webvXbox
Use your phone browser to use a virtual xbox controller

## How it works
- The server written in Go serves the controller layout on port 8080 and establishes a Websocket connection with each client.
- Up to 4 devices can connect, for every device a virtual xbox controller is plugged in.
- The [vXboxInterface](https://github.com/shauleiz/vXboxInterface) DLL is generated in Windows Temp directory and cleaned up when process ends.
- Static html/css/js files are embedded into the binary using `//go:embed`

## How to Build and Run
### Windows
#### Install file2byteslice and 2goarray
- `go install github.com/hajimehoshi/file2byteslice`
- `go install github.com/cratonica/2goarray`
#### Build webvXbox
- `go generate`
- `go build -ldflags -H=windowsgui`
- `webvxbox.exe`

### Linux
- `go build main.go server.go`
