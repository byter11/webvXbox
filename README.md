# webvXbox
Use your phone browser to use a virtual xbox controller

The vXboxInterface DLL is generated in Windows Temp directory and cleaned up when process ends.

Static html/css/js files are embedded into the binary using `//go:embed`

## How to Build and Run
- `go generate`
- `go build`
- `webvxbox.exe`

Up to 4 devices can connect, for every device a virtual xbox controller is plugged in.

## Note
- The frontend is very basic at the moment
- I might create a GUI interface for the server in the future
