# go-fiber-sse

## Install

```sh
go get github.com/cloudmatelabs/go-fiber-sse
```

## Usage

```go
package main

import (
  "bufio"
  "fmt"
  "time"

  "github.com/gofiber/fiber/v2"
  "github.com/cloudmatelabs/go-fiber-sse"
)

func main() {
  app := fiber.New()

  app.Get("/sse", func(c *fiber.Ctx) error {
		i := 0

		return sse.Handler(c, sse.Options{
			Tick: func(c *fiber.Ctx, w *bufio.Writer, close func()) {
				i++

				if i > 3 {
					close()
					return
				}

				fmt.Fprintf(w, "data: %d - %s\n\n", i, time.Now().Format(time.RFC3339Nano))
			},
			OnClosed: func(c *fiber.Ctx, err error) {
				fmt.Println("Connection is closed:", err)
			},
			Sleep: time.Second,
		})
  })

  app.Get("/", func(c *fiber.Ctx) error {
		c.Response().Header.Set("Content-Type", "text/html")
		return c.SendString(`<!DOCTYPE html>
<html>
<body>

<h1>SSE Messages</h1>
<div id="result"></div>

<script>
	if (typeof(EventSource) === 'undefined') {
		document.getElementById("result").innerHTML = "Sorry, your browser does not support server-sent events...";
	}

	const source = new EventSource("http://127.0.0.1:3000/sse");

	source.addEventListener('open', (event) => {
		document.getElementById("result").innerHTML += 'connected<br>';
	})

	source.addEventListener('message', (event) => {
		document.getElementById("result").innerHTML += event.data + "<br>";
	})

	source.addEventListener('error', (event) => {
		document.getElementById("result").innerHTML += 'connection lost<br>';

		source.close()

		document.getElementById("result").innerHTML += 'closed<br>';
	})
</script>

</body>
</html>`)
	})

  app.Listen("localhost:3000")
}
```