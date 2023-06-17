log-env:
	sudo mkdir -p /var/log/dream-mail-go/
	sudo chmod -R 777 /var/log/dream-mail-go/

build:
	CGO_ENABLED=0 go build -a -installsuffix cgo -o dream-mail-go main.go

run:
	go run main.go

compose-up:
	docker-compose up -d

compose-down:
	docker-compose down

docker-build:
	docker build -t dream-mail-go .

docker-run:
	docker run --rm --name=dream-mail-go \
		--network="host" \
		-p 8080:8080 \
		dream-mail-go