package main

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path"
	"sort"
	"strings"
	"time"
)

// Query parameters that should be used in signing the request.
var canonParams = map[string]bool{
	"acl":                          true,
	"delete":                       true,
	"lifecycle":                    true,
	"location":                     true,
	"logging":                      true,
	"notification":                 true,
	"partnumber":                   true,
	"policy":                       true,
	"requestpayment":               true,
	"response-cache-control":       true,
	"response-content-disposition": true,
	"response-content-encoding":    true,
	"response-content-language":    true,
	"response-content-type":        true,
	"response-expires":             true,
	"torrent":                      true,
	"uploadid":                     true,
	"uploads":                      true,
	"versionid":                    true,
	"versioning":                   true,
	"versions":                     true,
	"website":                      true,
}

var (
	id, key string
	base    *url.URL
)

func main() {
	server := http.Server{Handler: &httputil.ReverseProxy{Director: direct}}

	// Parse command line flags.
	flag.StringVar(&id, "id", "", "access key")
	flag.StringVar(&key, "key", "", "secret key")
	flag.StringVar(&server.Addr, "addr", "127.0.0.1:8080", "addr to bind to")
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: s3proxy [options] <url>")
		flag.PrintDefaults()
	}
	flag.Parse()

	// Parse the proxy target.
	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(1)
	}
	if base, _ = url.Parse(flag.Arg(0)); base == nil || base.Scheme == "" {
		log.Fatalf("bad URL: %q", flag.Arg(0))
	}

	// Extract id and key from URL if not passed explicitly.
	if base.User != nil {
		if id == "" {
			id = base.User.Username()
		}
		if key == "" {
			key, _ = base.User.Password()
		}
	}

	log.Fatalln(server.ListenAndServe())
}

// Direct the incoming request to the proxy target.
func direct(req *http.Request) {
	log.Println("Request:", req.Method, req.URL)

	// Re-route the request.
	req.URL.Scheme = base.Scheme
	req.URL.Host = base.Host
	req.URL.Path = path.Join(base.Path, req.URL.Path)

	if req.Header["Date"] == nil {
		req.Header.Set("Date", time.Now().Format(time.RFC1123Z))
	}

	// Extract the date if X-Amz-Date is unset.
	date := ""
	if req.Header["X-Amz-Date"] == nil {
		date = req.Header.Get("Date")
	}

	// Sign the request.
	hmac := hmac.New(sha1.New, []byte(key))
	io.WriteString(hmac, req.Method+"\n"+
		req.Header.Get("Content-MD5")+"\n"+
		req.Header.Get("Content-Type")+"\n"+
		date+"\n"+
		canonicalizedAmzHeaders(req.Header)+
		canonicalizedResource(req.URL))
	sig := base64.StdEncoding.EncodeToString(hmac.Sum(nil))
	req.Header.Set("Authorization", "AWS "+id+":"+sig)
}

func canonicalizedAmzHeaders(h http.Header) (s string) {
	var keys []string
	for k := range h {
		if strings.HasPrefix(k, "X-Amz") {
			keys = append(keys, strings.ToLower(k))
		}
	}
	sort.Strings(keys)
	for _, k := range keys {
		s += k + ":" + strings.Join(h[k], ",") + "\n"
	}
	return
}

func canonicalizedResource(u *url.URL) string {
	if q := canonicalizedQuery(u.Query()); q != "" {
		return u.EscapedPath() + "?" + q
	}
	return u.EscapedPath()
}

func canonicalizedQuery(query url.Values) string {
	for k := range query {
		if !canonParams[strings.ToLower(k)] {
			delete(query, k)
		}
	}
	return query.Encode()
}
