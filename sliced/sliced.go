//go:generate protoc -I. --go_out=.  --go_opt=paths=source_relative --go-grpc_opt=require_unimplemented_servers=false --go-grpc_opt=paths=source_relative --go-grpc_out=. *.proto

package sliced
