module github.com/josealecrim/audiototext/web

go 1.21

require (
	github.com/josealecrim/audiototext v0.0.0
	google.golang.org/grpc v1.61.0
)

replace github.com/josealecrim/audiototext => ../../

require (
	github.com/golang/protobuf v1.5.3 // indirect
	golang.org/x/net v0.20.0 // indirect
	golang.org/x/sys v0.16.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240125205218-1f4bbc51befe // indirect
	google.golang.org/protobuf v1.32.0 // indirect
)
