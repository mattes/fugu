build:	
	$(eval VERSION := $(shell go run main.go --version))

	GOOS=linux GOARCH=amd64 go build -o build/fugu.v$(VERSION).linux.x86_64
	(cd build && tar -cvzf fugu.v$(VERSION).linux.x86_64.tar.gz fugu.v$(VERSION).linux.x86_64)
	rm build/fugu.v$(VERSION).linux.x86_64

	GOOS=darwin GOARCH=amd64 go build -o build/fugu.v$(VERSION).darwin.x86_64
	(cd build && tar -cvzf fugu.v$(VERSION).darwin.x86_64.tar.gz fugu.v$(VERSION).darwin.x86_64)
	rm build/fugu.v$(VERSION).darwin.x86_64

	# TODO returns error: docker/docker/pkg/term/term.go:16: undefined: Termios
	# GOOS=windows GOARCH=amd64 go build -o build/fugu.v$(VERSION).windows.x86_64
	# cd build && tar -cvzf fugu.v$(VERSION).windows.x86_64.tar.gz fugu.v$(VERSION).windows.x86_64)
	# rm build/fugu.v$(VERSION).windows.x86_64

clean:
	rm -r build/*

test:
	go test ./...

.PHONY: build clean test