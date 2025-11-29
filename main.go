package main

import (
	"anchor/cmd"
	"github.com/lmittmann/tint"
	"log/slog"
	"os"
	"time"
)

func main() {

	w := os.Stderr
	slog.SetDefault(
		slog.New(
			tint.NewHandler(w, &tint.Options{
				AddSource:   false,
				ReplaceAttr: nil,
				TimeFormat:  time.Kitchen,
				NoColor:     false,
			})))
	cmd.Execute()
}
