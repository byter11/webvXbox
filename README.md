# webvXbox
Use your phone browser to use a virtual xbox controller

The [vXboxInterface](https://github.com/shauleiz/vXboxInterface) DLL is generated in Windows Temp directory and cleaned up when process ends.

Static html/css/js files are embedded into the binary using `//go:embed`

## How to Build and Run
### Install file2byteslice
- `git clone https://github.com/hajimehoshi/file2byteslice`
- `cd file2byteslice/cmd/file2byteslice`
- `go install`
### Build webvXbox
- `go generate`
- `go build main.go`
- `webvxbox.exe`

Up to 4 devices can connect, for every device a virtual xbox controller is plugged in.

## Notes
- Recommended to use USB tethering if only 1 controller needed, or use a 5Ghz hotspot for least latency
- The frontend is very basic at the moment
- I might create a GUI interface for the server in the future
