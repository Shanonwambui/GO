# NFT/HTTP

Package http provides HTTP client and server implementations.

`Get`, `Head`, `Post`, and `PostForm` make HTTP (or HTTPS) requests:
```
resp, err := http.Get("http://example.com/")
...
resp, err := http.Post("http://example.com/upload", "image/jpeg", &buf)
...
resp, err := http.PostForm("http://example.com/form",
	url.Values{"key": {"Value"}, "id": {"123"}})

```

The caller must close the response body when finished with it:

```
resp, err := http.Get("http://example.com/")
if err != nil {
	// handle error
}
defer resp.Body.Close()
body, err := io.ReadAll(resp.Body)
// ...
```

## Clients and Transports
For control over HTTP client headers, redirect policy, and other settings, create a Client:
```
client := &http.Client{
	CheckRedirect: redirectPolicyFunc,
}

resp, err := client.Get("http://example.com")
// ...

req, err := http.NewRequest("GET", "http://example.com", nil)
// ...
req.Header.Add("If-None-Match", `W/"wyzzy"`)
resp, err := client.Do(req)
// ...
```
For control over proxies, TLS configuration, keep-alives, compression, and other settings, create a Transport:

```
tr := &http.Transport{
	MaxIdleConns:       10,
	IdleConnTimeout:    30 * time.Second,
	DisableCompression: true,
}
client := &http.Client{Transport: tr}
resp, err := client.Get("https://example.com")
```
Clients and Transports are safe for concurrent use by multiple goroutines and for efficiency should only be created once and re-used.

## Servers
ListenAndServe starts an HTTP server with a given address and handler. The handler is usually nil, which means to use `DefaultServeMux`. `Handle` and `HandleFunc` add handlers to `DefaultServeMux`:

```
http.Handle("/foo", fooHandler)

http.HandleFunc("/bar", func(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
})

log.Fatal(http.ListenAndServe(":8080", nil))

```
More control over the server's behavior is available by creating a custom Server:
```
s := &http.Server{
	Addr:           ":8080",
	Handler:        myHandler,
	ReadTimeout:    10 * time.Second,
	WriteTimeout:   10 * time.Second,
	MaxHeaderBytes: 1 << 20,
}
log.Fatal(s.ListenAndServe())

```

Read more: https://pkg.go.dev/net/http#Handle

### func Handle
```
func Handle(pattern string, handler Handler)
```
Handle registers the handler for the given pattern in DefaultServeMux. The documentation for ServeMux explains how patterns are matched.

### func HandleFunc
```
func HandleFunc(pattern string, handler func(ResponseWriter, *Request))
```
HandleFunc registers the handler function for the given pattern in DefaultServeMux.

### func ListenAndServe 
```
func ListenAndServe(addr string, handler Handler) error
```
ListenAndServe listens on the TCP network address addr and then calls Serve with handler to handle requests on incoming connections. Accepted connections are configured to enable TCP keep-alives.

The handler is typically nil, in which case DefaultServeMux is used.

ListenAndServe always returns a non-nil error.

### func ListenAndServeTLS
```
func ListenAndServeTLS(addr, certFile, keyFile string, handler Handler) error
```
ListenAndServeTLS acts identically to ListenAndServe, except that it expects HTTPS connections. Additionally, files containing a certificate and matching private key for the server must be provided. If the certificate is signed by a certificate authority, the certFile should be the concatenation of the server's certificate, any intermediates, and the CA's certificate.

### func ServeFile

```
func ServeFile(w ResponseWriter, r *Request, name string)
```

ServeFile replies to the request with the contents of the named file or directory.

If the provided file or directory name is a relative path, it is interpreted relative to the current directory and may ascend to parent directories. If the provided name is constructed from user input, it should be sanitized before calling ServeFile.

As a precaution, ServeFile will reject requests where r.URL.Path contains a ".." path element; this protects against callers who might unsafely use filepath.Join on r.URL.Path without sanitizing it and then use that filepath.Join result as the name argument.

As another special case, ServeFile redirects any request where r.URL.Path ends in "/index.html" to the same path, without the final "index.html". To avoid such redirects either modify the path or use ServeContent.

Outside of those two special cases, ServeFile does not use r.URL.Path for selecting the file or directory to serve; only the file or directory provided in the name argument is used.

### func FileServer
```
func FileServer(root FileSystem) Handler
```
FileServer returns a handler that serves HTTP requests with the contents of the file system rooted at root.

As a special case, the returned file server redirects any request ending in "/index.html" to the same path, without the final "index.html".

To use the operating system's file system implementation, use `http.Dir`:
```
http.Handle("/", http.FileServer(http.Dir("/tmp")))
```