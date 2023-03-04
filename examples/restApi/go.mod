module examples/restApi

go 1.18

replace alox.sh => ../..

require alox.sh v0.0.0-00010101000000-000000000000

require (
	github.com/mattn/go-sqlite3 v1.14.16 // indirect
	go.mongodb.org/mongo-driver v1.11.2 // indirect
	golang.org/x/net v0.7.0 // indirect
)
