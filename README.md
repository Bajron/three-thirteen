# three-thirteen

A simple game card to try out golang.

There is nothing usable here yet.

To setup a server:
# is there a better way? somthing using go?
wget https://golang.org/src/crypto/tls/generate_cert.go?m=text -O generate_cert.go
go build generate_cert.go
./generate_cert --host localhost
go build
./three-thirteen

