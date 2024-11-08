
clean:
	rm -Rfv bin
	mkdir bin

build: clean
	# go build -o bin/wpclone main.go
	go install

docker-build:
	docker build -t rg.nl-ams.scw.cloud/wpclone/wpclone-cli docker_files/cli
	docker build -t rg.nl-ams.scw.cloud/wpclone/wpclone-db docker_files/db
	docker build -t rg.nl-ams.scw.cloud/wpclone/wpclone-wp docker_files/wp
	docker build -t rg.nl-ams.scw.cloud/wpclone/wpclone-dnsmasq docker_files/dnsmasq

docker-push:
	docker push rg.nl-ams.scw.cloud/wpclone/wpclone-cli
	docker push rg.nl-ams.scw.cloud/wpclone/wpclone-db
	docker push rg.nl-ams.scw.cloud/wpclone/wpclone-wp
	docker push rg.nl-ams.scw.cloud/wpclone/wpclone-dnsmasq

docker: docker-build docker-push

build-all: clean
	GOOS="linux"   GOARCH="amd64"       go build -o bin/wpclone__linux-amd64 main.go
	GOOS="linux"   GOARCH="arm64"       go build -o bin/wpclone__linux-arm64   main.go
	GOOS="freebsd" GOARCH="amd64"       go build -o bin/wpclone__freebsd-amd64 main.go
	GOOS="freebsd" GOARCH="arm64"       go build -o bin/wpclone__freebsd-arm64 main.go
	GOOS="darwin"  GOARCH="amd64"       go build -o bin/wpclone__macos-amd64 main.go
	GOOS="darwin"  GOARCH="arm64"       go build -o bin/wpclone__macos-arm64 main.go

release: build-all
	chmod a+x bin/*
	rsync -az -e 'ssh -p 2807' --partial --stats --progress bin/* tux@skrova.noltech.net:/opt/file.noltech.net/files/wpclone/
	rsync -az -e 'ssh -p 2807' --partial --stats --progress install.sh tux@skrova.noltech.net:/opt/file.noltech.net/files/wpclone/
	
