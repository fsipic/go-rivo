# go-rivo
web server u Go-u  za real time obavijesti o cijenama goriva.


upute za korištenje:

mkdir certs
openssl genrsa -out certs/server.key 2048
openssl req -new -x509 -sha256 -key certs/server.key -out certs/server.crt -days 365


go build -o myapp cmd/server/main.go
./myapp

or

go run cmd/server/main.go

