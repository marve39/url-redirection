run-test:
	go test
build-linux: 
	mkdir -p dist && CGO_ENABLED=0 GOOS=linux go build -a -ldflags "-s" -o dist/url-redirection main.go 
docker-build:
	docker build --tag marve39/url-redirection:1.0.1 . 
docker-run:
	docker run -d -p 80:80 -e DOMAIN_1="test;redirect-domain;add-path;add=parameter" marve39/url-redirection:latest