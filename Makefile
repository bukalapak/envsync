test:
	go test -v -race ./...

dep:
	dep ensure

cover:
	go test -coverprofile=coverage.out && go tool cover -html=coverage.out