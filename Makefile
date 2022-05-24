run:
	docker-compose up --build
proto:
	buf generate
test:
	go test -v ./... -cover
