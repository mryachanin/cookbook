GOCMD=go
GOGET=${GOCMD} get

build: pack
	${GOCMD} build

install: pack
	${GOCMD} install

pack: clean
	go-bindata \
		-pkg template \
		-o web/template/assets.go \
		web/template/...

clean:
	${GOCMD} clean
	rm -f web/template/*.go
	rm -f bin/cookbook

test:
	${GOCMD} test

deps:
	${GOGET} github.com/rhinoman/couchdb-go  # CouchDB driver
	${GOGET} github.com/google/uuid          # UUID support
	${GOGET} github.com/jteeuwen/go-bindata  # Static asset wrapper-upper
	${GOGET} gopkg.in/yaml.v2                # YAML parser
