# go-dcgi

DCGI is a technique for serving web pages dynamically with Docker. As you may know, World Wide Web servers can only serve static files off disk. A DCGI server allows you to *execute code in real-time*, so the Web page can contain *dynamic* information.

For each HTTP request that a DGCI server receives, a Docker container is spun up to serve the HTTP request. Inside the Docker container is a [CGI](https://en.wikipedia.org/wiki/Common_Gateway_Interface) executable which handles the request. That executable could do anything – and could be written in any language or framework.

Wow! No longer do we have to build Web sites which just serve static content. For example, you could "hook up" your Unix database to the World Wide Web so people all over the world could query it. Or, you could create HTML forms to allows people to transmit information into your database engine. The possibilities are limitless.

So what's this library for? go-dcgi is a library for writing DGCI servers. It includes a Go handler, `dcgi.Handler`, which serves that HTTP request by running a Docker container.

## Usage

Say you've got a really simple CGI script, `script.pl`:

```perl
print "Content-Type: text/html\n\n";
print "<h1>Hello World!</h1>\n";
```

And a `Dockerfile` to put it inside a container:

```
FROM perl
ADD script.pl /code/script.pl
ENTRYPOINT ["perl", "/code/script.pl"]
```

Build this into a container:

```bash
$ docker build -t bfirsh/example-dcgi-app
```

You can serve this container over HTTP with go-dcgi:

```go
package main

import (
	"net/http"
	dcgi "github.com/bfirsh/go-dcgi"
	"github.com/docker/engine-api/client"
)

func main() {
	cli, err := client.NewClient("unix:///var/run/docker.sock", "v1.23", nil, nil)
	if err != nil {
		panic(err)
	}

	http.Handle("/", &dcgi.Handler{
		Image:      "bfirsh/example-dcgi-app",
		Client:     cli,
	})
	http.ListenAndServe(":80", nil)
}
```

