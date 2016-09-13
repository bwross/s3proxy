`s3proxy` is a simple request signing proxy for S3. It allows applications to
access S3 without having to know how to sign S3 requests.

    Usage: s3proxy [options] <url>
      -addr string
        	addr to bind to (default "127.0.0.1:8080")
      -id string
        	access key
      -key string
        	secret key
      -ro
        	only allow GETs (read-only)

### Installing

In order to build and use `s3proxy`, you'll need to first install Go (version 1.5 or later), as well as set the `GOPATH` environment variable.

#### Installing Go

1. Start here: https://golang.org/doc/install
2. Once installed, you'll want to set `GOPATH` to a directory structure where you want to download and build code.
	1. Read about `GOPATH` here: https://github.com/golang/go/wiki/GOPATH
	2. You'll want to set `GOPATH` to the top level of where you'll be pulling down and building Go packages. For our purposes, `$HOME/gocode`:

			mkdir -p $HOME/gocode
			export GOPATH=$HOME/gocode

3. Optionally, add `$GOPATH/bin` to your `PATH` so you can easily use Go binaries from the command line.

		export PATH=$PATH:$GOPATH/bin

#### Install/build s3proxy

Now that you have Go all set and ready to "go", it's as simple as:

	go get github.com/igneous-systems/s3proxy

If all went well, `$GOPATH/bin/s3proxy` should be an executable suitable for running on your current system. If you configured your `PATH`, you can now run `s3proxy` directly from the command line.

#### Examples

Run the proxy locally, proxying to `http://s3.amazonaws.com`:

    s3proxy -id <id> -key <secret> http://s3.amazonaws.com

Run the proxy publicly accessible over port 80 using `-addr`:

    s3proxy -id <id> -key <secret> -addr 0.0.0.0:80 http://s3.amazonaws.com

Run the proxy in read-only mode using `-ro`:

    s3proxy -id <id> -key <secret> -ro http://s3.amazonaws.com

Serve from a bucket or subdirectory by appending a path to the URL:

    s3proxy -id <id> -key <secret> http://s3.amazonaws.com/bucket
