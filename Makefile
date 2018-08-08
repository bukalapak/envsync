test:
	go test -v -race ./...

dep:
	dep ensure

cover:
	go test -coverprofile=coverage.out && go tool cover -html=coverage.out

build:
	CGO_ENABLED=0 gox \
	-os="linux darwin windows" \
	-arch="amd64" \
	-output="build/envsync_{{.OS}}_{{.Arch}}" \
	./app/cli/