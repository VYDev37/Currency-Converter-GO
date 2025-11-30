## Currency Converter in Golang
Description: Only a simple currency converter made in Golang for learning purpose.

## Which part does this cover?
- [x] API Fetching
- [x] More about Slices
- [x] More about parsing and writing JSON
- [x] Frontend Development (React + axios + TailwindCSS) (Currently working on)
- [x] Backend Development (gRPC + REST Server)
- [x] Expiration and caching

## How to run
- Run the app with command: `go run .`
- Note: If you would like to modify the protobuf file, make sure to update it after the modification with command below:
  `protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/converter.proto`

Note: There'll be three versions of this project.
## Versions:
- Command Line Interface (CLI)
- Terminal User Interface (TUI) 
- Backend of Website / API (Current branch)