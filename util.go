package sofafeed

import (
	"io"
	"log/slog"
)

func mustClose(c io.Closer) {
	err := c.Close()
	if err != nil {
		slog.Error("Failed to close", "error", err)
	}
}
