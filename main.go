package main

import (
	"fmt"
	"fourbytedirectory-data/src"
	"log"
	"os"
)

func main() {
	cfg := src.InitConfig()

	exporter := src.NewExporter()
	exporter.Start()

	processFlags(cfg, exporter)

	pool := src.NewWorkerPool(cfg, exporter)
	pool.Start()
}

func processFlags(cfg src.Config, exporter *src.Exporter) {
	if cfg.CountsOnly {
		log.Println(
			fmt.Sprintf(
				"Completed pages: %d, Failed pages: %d, Total Signatures: %d",
				len(exporter.CompletedPages),
				len(exporter.FailedPages),
				len(exporter.Signatures),
			))
		os.Exit(0)
	}
}
