# vim:ft=make

BINARY=stream

all:
	go build -o ${BINARY} .

test:
	go test ./...

clean:
	go clean
	rm -f ${BINARY}
	rm -f *.db *.dat