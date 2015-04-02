build:	
	$(eval VERSION := $(shell godep go run fugu/main.go fugu/usage.go fugu/version.go --version))

	(cd fugu && GOOS=linux GOARCH=amd64 godep go build -o ../build/fugu.v$(VERSION).linux.x86_64)
	(cd build && tar -cvzf fugu.v$(VERSION).linux.x86_64.tar.gz fugu.v$(VERSION).linux.x86_64)
	rm build/fugu.v$(VERSION).linux.x86_64

	(cd fugu && GOOS=darwin GOARCH=amd64 godep go build -o ../build/fugu.v$(VERSION).darwin.x86_64)
	(cd build && tar -cvzf fugu.v$(VERSION).darwin.x86_64.tar.gz fugu.v$(VERSION).darwin.x86_64)
	rm build/fugu.v$(VERSION).darwin.x86_64

	# TODO returns error: docker/docker/pkg/term/term.go:16: undefined: Termios
	# GOOS=windows GOARCH=amd64 godep go build -o build/fugu.v$(VERSION).windows.x86_64
	# cd build && tar -cvzf fugu.v$(VERSION).windows.x86_64.tar.gz fugu.v$(VERSION).windows.x86_64)
	# rm build/fugu.v$(VERSION).windows.x86_64

install:
	(cd fugu && godep go install)

clean:
	rm -r build/*
	rm fugu/fugu

test:
	GOTEST=1 godep go test -v ./...

usage-file:
	(cd fugu && godep go build)

	(cd fugu && ./fugu help > usage.txt 2>&1)
	echo "\n\n------------------------------------------\n\n" >> fugu/usage.txt
	(cd fugu && ./fugu help build >> usage.txt 2>&1)
	echo "\n\n------------------------------------------\n\n" >> fugu/usage.txt
	(cd fugu && ./fugu help run >> usage.txt 2>&1)
	echo "\n\n------------------------------------------\n\n" >> fugu/usage.txt
	(cd fugu && ./fugu help exec >> usage.txt 2>&1)
	echo "\n\n------------------------------------------\n\n" >> fugu/usage.txt
	(cd fugu && ./fugu help shell >> usage.txt 2>&1)
	echo "\n\n------------------------------------------\n\n" >> fugu/usage.txt
	(cd fugu && ./fugu help destroy >> usage.txt 2>&1)
	echo "\n\n------------------------------------------\n\n" >> fugu/usage.txt
	(cd fugu && ./fugu help push >> usage.txt 2>&1)
	echo "\n\n------------------------------------------\n\n" >> fugu/usage.txt
	(cd fugu && ./fugu help pull >> usage.txt 2>&1)
	echo "\n\n------------------------------------------\n\n" >> fugu/usage.txt
	(cd fugu && ./fugu help images >> usage.txt 2>&1)
	echo "\n\n------------------------------------------\n\n" >> fugu/usage.txt
	(cd fugu && ./fugu help show-data >> usage.txt 2>&1)
	echo "\n\n------------------------------------------\n\n" >> fugu/usage.txt
	(cd fugu && ./fugu help show-labels >> usage.txt 2>&1)

release: build usage-file

.PHONY: build clean test usage-file release install