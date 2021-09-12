module github.com/FlowKeeper/FlowAgent/v2

// +heroku goVersion go1.17
go 1.17

require (
	github.com/google/uuid v1.3.0
	github.com/gorilla/mux v1.8.0
	github.com/mattn/go-sqlite3 v1.14.8
	gitlab.cloud.spuda.net/Wieneo/golangutils/v2 v2.0.0-20210912064442-dc511fe9c28b
	go.mongodb.org/mongo-driver v1.7.2
)

require github.com/FlowKeeper/FlowUtils/v2 v2.0.0-20210911185616-289b6e6b3efd // indirect
