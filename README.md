# go-pastel

Golang fork of https://github.com/bobpp/pastel, which is a copy and paste sharing web app.

# Description

go-pastel is a copy and paste sharing web application like gist written in golang. 

## Installation

Executable binaries, which contains asset files, are available at [releases](https://github.com/sonots/go-pastel/releases).

For example, for linux x86\_64, 

```bash
$ version=0.1.0
$ wget https://github.com/sonots/go-pastel/releases/download/$version/go-pastel_linux_amd64 -O go-pastel
$ chmod a+x go-pastel
```

If you have the go runtime installed, you may use go get. 

```bash
$ go get github.com/sonots/go-pastel
```

## Usage

```
$ go-pastel -h
NAME:
   go-pastel - A copy and paste sharing web application like git

USAGE:
   go-pastel [global options] command [command options] [arguments...]

COMMANDS:
   start, s     Start up
   help, h      Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --host '0.0.0.0'             Address to serve this service
   --port '5050'                Port number to serve this service
   --database_url 'pastel.db'   Path to sqlite storage file
   --help, -h                   show help
   --version, -v                print the version
```

## Build

To build, use go get and make

```
$ go get -d github.com/sonots/go-pastel
$ cd $GOPATH/src/github.com/sonots/go-pastel
$ make
$ go run main.go version.go bindata.go # run
```

To release binaries, I use [gox](https://github.com/mitchellh/gox) and [ghr](https://github.com/tcnksm/ghr)

```
go get github.com/mitchellh/gox
gox -build-toolchain # only first time
go get github.com/tcnksm/ghr

mkdir -p pkg && cd pkg && gox --os=linux --os=windows ../...
ghr <tag> .
```

## ToDo

1. write tests
2. output log to file, allow to change log_level
3. How to limit number of concurrent http requests?
4. How to limit number of concurrent db connections? => db.SetMaxIdleConns

## Contribution

1. Fork (https://github.com/sonots/go-pastel/fork)
2. Create a feature branch
3. Commit your changes
4. Rebase your local changes against the master branch
5. Run test suite with the go test ./... command and confirm that it passes
6. Run gofmt -s
7. Create new Pull Request

## Copyright

See [LICENSE](./LICENSE)
