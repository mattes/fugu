build:	
	$(eval VERSION := $(shell go run main.go --version))

	GOOS=linux GOARCH=amd64 go build -o build/fugu.$(VERSION).linux.x86_64
	(cd build && tar -cvzf fugu.$(VERSION).linux.x86_64.tar.gz fugu.$(VERSION).linux.x86_64)
	rm build/fugu.$(VERSION).linux.x86_64

	GOOS=darwin GOARCH=amd64 go build -o build/fugu.$(VERSION).darwin.x86_64
	(cd build && tar -cvzf fugu.$(VERSION).darwin.x86_64.tar.gz fugu.$(VERSION).darwin.x86_64)
	rm build/fugu.$(VERSION).darwin.x86_64

	# TODO returns error: docker/docker/pkg/term/term.go:16: undefined: Termios
	# GOOS=windows GOARCH=amd64 go build -o build/fugu.$(VERSION).windows.x86_64
	# cd build && tar -cvzf fugu.$(VERSION).windows.x86_64.tar.gz fugu.$(VERSION).windows.x86_64)
	# rm build/fugu.$(VERSION).windows.x86_64

clean:
	rm -r build/*

test:
	go test ./...

.PHONY: build clean test