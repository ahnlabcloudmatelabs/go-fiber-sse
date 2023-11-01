package sse

import (
	"bufio"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
)

type Options struct {
	// Tick of each event
	Tick func(c *fiber.Ctx, w *bufio.Writer, close func())
	// Closed by unknown reason. e.g. client is closed
	OnClosed func(c *fiber.Ctx, err error)
	// Term of each tick. Default is 1 second
	Sleep time.Duration
}

/*
Example:

		i := 0

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
*/
func Handler(c *fiber.Ctx, options Options) error {
	if options.Sleep == 0 {
		options.Sleep = time.Second
	}

	if options.OnClosed == nil {
		options.OnClosed = func(c *fiber.Ctx, err error) {}
	}

	closeFlag := false
	close := func() {
		closeFlag = true
	}

	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache")
	c.Set("Connection", "keep-alive")
	c.Set("Transfer-Encoding", "chunked")

	c.Context().SetBodyStreamWriter(fasthttp.StreamWriter(func(w *bufio.Writer) {
		for {
			if closeFlag {
				return
			}

			options.Tick(c, w, close)

			if err := w.Flush(); err != nil {
				options.OnClosed(c, err)

				break
			}

			time.Sleep(options.Sleep)
		}
	}))

	return nil
}
