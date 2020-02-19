go mod init crudapp
# go mod init github.com/rvasily/crudapp
go build
go mod download
go mod verify
go mod tidy

go build  -o ./bin/crudapp ./cmd/crudapp
go test -v -coverpkg=./... ./...

go mod vendor
go build -mod=vendor -o ./bin/myapp ./cmd/myapp
go test -v -mod=vendor -coverpkg=./... ./...