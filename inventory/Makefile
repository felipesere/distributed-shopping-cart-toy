build:
	go build cmd/server/server.go

run:
	MONGO_DB_CLUSTER="mongodb://localhost:27017" \
	MONGO_USER="root" \
	MONGO_PASSWORD="rootpassword" \
	go run cmd/server/server.go

build-seed:
	go build cmd/seed/seed.go

seed:
	go run cmd/seed/seed.go -file ./seed.json
