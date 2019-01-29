GOCMD=go

install: clean
	go-bindata \
		-pkg template \
		-o web/template/index.go \
		web/template/...
	$(GOCMD) install

clean:
	$(GOCMD) clean
	rm -f web/template/*.go
	rm -f bin/satisfied-vegan

test:
	$(GOCMD) test

deps:
	GOGET=$(GOCMD) get
	$(GOGET) github.com/mattn/go-sqlite3    // Sqlite driver
	$(GOGET) github.com/jteeuwen/go-bindata // Static asset wrapper-upper
	$(GOGET) gopkg.in/yaml.v2               // YAML parser
