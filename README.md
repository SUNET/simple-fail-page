[![Build Status](https://travis-ci.org/SUNET/simple-fail-page.svg?branch=master)](https://travis-ci.org/SUNET/simple-fail-page)
# simple-fail-page
A very simple web server for when something goes wrong


## Docker
To build a docker image use these steps

```bash
go get github.com/mitchellh/gox
$GOPATH/bin/gox -os="linux" -arch="amd64" -output "simple-fail-page.{{.OS}}.{{.Arch}}" -ldflags "-X main.Rev=`git rev-parse --short HEAD`" -verbose
docker build -t local/simple-fail-page .
```
