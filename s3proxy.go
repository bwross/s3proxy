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

func main() {
	proxy := new(Proxy)
	proxy.Director = proxy.Direct
	server := http.Server{Handler: proxy}

	// Parse command line flags.
	flag.StringVar(&proxy.ID, "id", "", "access key")
	flag.StringVar(&proxy.Key, "key", "", "secret key")
	flag.StringVar(&server.Addr, "addr", "127.0.0.1:8080", "addr to bind to")
	flag.BoolVar(&proxy.ReadOnly, "ro", false, "only allow GETs (read-only)")
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
	proxy.SetURL(flag.Arg(0))

	log.Fatalln(server.ListenAndServe())
}

// A Proxy is an http.Handler that proxies to a URL and signs requests.
type Proxy struct {
	httputil.ReverseProxy
	*url.URL
	ID, Key  string
	ReadOnly bool
}

// SetURL sets the base URL, exiting if it is invalid.
func (p *Proxy) SetURL(u string) {
	if p.URL, _ = url.Parse(flag.Arg(0)); p.URL == nil || p.URL.Scheme == "" {
		log.Fatalf("bad URL: %q", flag.Arg(0))
	}
	p.URL.Path = strings.TrimSuffix(p.URL.Path, "/")
}

// ServeHTTP implements http.Handler.
func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if p.ReadOnly && r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	p.ReverseProxy.ServeHTTP(w, r)
}

// Direct the incoming request to the proxy target.
func (p *Proxy) Direct(req *http.Request) {
	log.Println("Request:", req.Method, req.URL)

	// Re-route the request.
	req.URL.Scheme = p.Scheme
	req.URL.Host = p.Host
	req.URL.Path = p.Path + req.URL.Path

	if req.Header["Date"] == nil {
		req.Header.Set("Date", time.Now().Format(time.RFC1123Z))
	}

	// Extract the date if X-Amz-Date is unset.
	date := ""
	if req.Header["X-Amz-Date"] == nil {
		date = req.Header.Get("Date")
	}

	// Sign the request.
	hmac := hmac.New(sha1.New, []byte(p.Key))
	io.WriteString(hmac, req.Method+"\n"+
		req.Header.Get("Content-MD5")+"\n"+
		req.Header.Get("Content-Type")+"\n"+
		date+"\n"+
		canonicalizedAmzHeaders(req.Header)+
		canonicalizedResource(req.URL))
	sig := base64.StdEncoding.EncodeToString(hmac.Sum(nil))
	req.Header.Set("Authorization", "AWS "+p.ID+":"+sig)
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
