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

In order to build and use s3proxy, you'll need to first install go (version 1.5 or later), as well as set the `GOPATH` environment variable.

#### Installing go

1.  Start here: [https://golang.org/doc/install](https://golang.org/doc/install)
2.  Once installed, you'll want to set `GOPATH` to a directory structure where you want to download and build code.
	3.  Read about GOPATH here: [https://github.com/golang/go/wiki/GOPATH](https://github.com/golang/go/wiki/GOPATH)
	4.  Basically, you'll want to set `GOPATH` at the top level of where you'll be pulling down and buliding gocode packages. For our purposes ,:
		
			mkdir -p /Users/andypern/gocode
			export GOPATH=/Users/andypern/gocode
		Obviously you can put this into your profile/.bashrc, etc.
	5. Verify its set via:
			
			echo $GOPATH

####Install/build s3proxy

Now that you have GO all set and ready to 'go', its as simple as:

1.  Change your working directory to $GOPATH:

		cd $GOPATH

2.  Run the following to grab and build the s3proxy binary:

		go get github.com/bwross/s3proxy
		
If all went well, you should see some stuff:

	ls -R

Which should show you something like:

	bin src
	
	./bin:
	s3proxy
	
	./src:
	github.com
	
	./src/github.com:
	bwross
	
	./src/github.com/bwross:
	s3proxy
	
	./src/github.com/bwross/s3proxy:
	LICENSE    README.md  s3proxy.go
	
The file in ./bin should be an executable suitable for running on your current system.  On an OSX host:
	
	file bin/s3proxy
	bin/s3proxy: Mach-O 64-bit executable x86_64



#### Examples

Run the proxy locally, proxying to `http://s3.amazonaws.com`:

    s3proxy -id <id> -key <secret> http://s3.amazonaws.com

Run the proxy publicly accessible over port 80 using `-addr`:

    s3proxy -id <id> -key <secret> -addr 0.0.0.0:80 http://s3.amazonaws.com

Run the proxy in read-only mode using `-ro`:

    s3proxy -id <id> -key <secret> -ro http://s3.amazonaws.com

Serve from a bucket or subdirectory by appending a path to the URL:

    s3proxy -id <id> -key <secret> http://s3.amazonaws.com/bucket
