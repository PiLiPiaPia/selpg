## description
This is an Unix style cli-tool that helps you get a specified range of pages in a document, written in golang.

## how to install
$ git clone $GOPATH/src
$ go get github.com/spf13/pflag
$ cd $GOPATH/src
$ go build selpg.go

## usage

USAGE: selpg -s startpage -e endpage [OPTION...] [file]
  -s, --startpage int      start page (default 1)
  -e, --endpage int        end page (default 1)
  -l, --pagelen int        page length (default 72), only valid if '-f' not set
  -f, --pagetype           specify the page type that pages are separated by '\f' 
  -d, --printdest string   printer destination, the specified range of pages will be printed to the destination
  
