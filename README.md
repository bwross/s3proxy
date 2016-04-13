`s3proxy` is a simple request signing proxy for S3. It allows applications to
access S3 without having to know how to sign S3 requests.

    Usage: s3proxy [options] <url>
      -addr string
        	addr to bind to (default "127.0.0.1:8080")
      -id string
        	access key
      -key string
        	secret key

#### Examples
Run the proxy locally, proxying to `http://s3.amazonaws.com`:

    s3proxy -id <id> -key <secret> http://s3.amazonaws.com

Run the proxy publicly accessible over port 80:

    s3proxy -id <id> -key <secret> -addr 0.0.0.0:80 http://s3.amazonaws.com
