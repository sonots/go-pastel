DEBUG_FLAG = $(if $(DEBUG),-debug)

deps:
	go get github.com/jteeuwen/go-bindata/...
	go get -d -t ./...
	go-bindata views/... static/...

test: deps
	go test -v ./...

install: deps
	go install

pkg: deps
	go get github.com/mitchellh/gox
	mkdir -p pkg && cd pkg && gox --os=linux --os=windows ../... # dawin fails ...

run: install
	$(GOPATH)/bin/go-pastel
