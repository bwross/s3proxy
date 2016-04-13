`s3proxy` is a simple request signing proxy for S3. It allows applications to
access S3 without having to know how to sign S3 requests.

    Usage: s3proxy [options] <url>
      -addr string
        	addr to bind to (default "127.0.0.1:8080")
      -id string
        	access key
      -key string
        	secret key

#### Installing

Go 1.5 or later is required, and `GOPATH` must be configured.

    go get github.com/bwross/s3proxy

#### Examples

Run the proxy locally, proxying to `http://s3.amazonaws.com`:

    s3proxy -id <id> -key <secret> http://s3.amazonaws.com

Run the proxy publicly accessible over port 80 using `-addr`:

    s3proxy -id <id> -key <secret> -addr 0.0.0.0:80 http://s3.amazonaws.com

Run the proxy in read-only mode using `-ro`:

    s3proxy -id <id> -key <secret> -ro http://s3.amazonaws.com
