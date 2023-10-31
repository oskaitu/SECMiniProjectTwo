# SECMiniProjectTwo

# protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/proto.proto

openssl req -nodes -x509 -sha256 -newkey rsa:4096 -keyout priv.key -out server.crt -days 356 -subj "/C=DK/ST=Copenhagen/L=Copenhagen/O=Me/OU=mpc/CN=localhost" -addext "subjectAltName = DNS:localhost,IP:0.0.0.0"

go run main.go -port 0 -name "Hospital"
go run main.go -port 1 -name "Alice"
go run main.go -port 2 -name "Bob"
go run main.go -port 3 -name "Charlie"
