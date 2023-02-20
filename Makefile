run:
	mkdir -p docker/data && \
	docker-compose up -d --build

stop:
	docker-compose down

db-logs:
	docker-compose logs -f db

banking-app-logs:
	docker-compose logs -f banking-app

migrates-up:
	migrate \
		-source file://docker/migrations \
		-database postgres://postgres:qwerty123@localhost:5432/banking?sslmode=disable \
		up

migrates-down:
	migrate \
		-source file://docker/migrations \
		-database postgres://postgres:qwerty123@localhost:5432/banking?sslmode=disable \
		down

gofmt:
	gofmt -l -w .

goimports:
	goimports -w ./

test:
	go test -coverprofile="/tmp/go-cover.$$.tmp" ./... && \
	go tool cover -func="/tmp/go-cover.$$.tmp" && \
	unlink "/tmp/go-cover.$$.tmp"