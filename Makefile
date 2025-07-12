GOCMD=go

build:
	${GOCMD} build

install:
	${GOCMD} install

clean:
	${GOCMD} clean
	rm -f bin/cookbook

test:
	${GOCMD} test

deps:
	${GOCMD} mod tidy
