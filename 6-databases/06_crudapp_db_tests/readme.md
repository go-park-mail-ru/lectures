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

----

docker-compose up

go test -v -coverprofile=cover.out -coverpkg=./... ./...
go tool cover -html=cover.out -o cover.html

-----

находять в папке pkg/handlers
mockgen -source=items.go -destination=items_mock.go -package=handlers ItemRepositoryInterface

go test -v -run Handler -coverprofile=handler.out && go tool cover -html=handler.out -o handler.html && rm handler.out
