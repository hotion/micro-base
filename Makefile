default:
	@echo 'Usage of make: [ build | linux_build | windows_build | clean ]'

build: 
	@go build -ldflags "-X main.VERSION=1.0.0 -X main.GIT_HASH=`git rev-parse HEAD` -s" -o ./bin/micro-base ./

linux_build: 
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-X main.VERSION=1.0.0 -X main.GIT_HASH=`git rev-parse HEAD` -s" -o ./bin/micro-base ./

windows_build: 
	@CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags "-X main.VERSION=1.0.0 -X main.GIT_HASH=`git rev-parse HEAD` -s" -o ./bin/micro-base.exe ./

clean: 
	@rm -f ./bin/micro-base*

.PHONY: default build linux_build windows_build clean