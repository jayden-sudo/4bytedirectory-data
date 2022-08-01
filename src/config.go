package src

import (
	"fmt"
	"github.com/akamensky/argparse"
	"os"
)

type Config struct {
	Api         FourByteApi
	StartPage   int
	Threads     int
	MaxRetries  int
	FindMissing bool
	FailedOnly  bool
	CountsOnly  bool
}

// todo: validate args, --failed cannot be run with --missing etc
func InitConfig() Config {
	parser := argparse.NewParser("export", "Exports function signature data from 4byte.directory")

	startPage := parser.Int("p", "page", &argparse.Options{
		Help: "The page to start scraping from", Default: 1,
	})

	threads := parser.Int("t", "threads", &argparse.Options{
		Help: "The number of threads", Default: 10,
	})

	retries := parser.Int("r", "retries", &argparse.Options{
		Help: "The number of times a page should be retried before it is considered failed", Default: 25,
	})

	missing := parser.Flag("m", "missing", &argparse.Options{
		Help: "Checks the completed pages array and looks for any pages that may be missing and fetches them",
	})

	failedOnly := parser.Flag("f", "failed-only", &argparse.Options{
		Help: "If set only pages that have previously failed will be processed",
	})

	countsOnly := parser.Flag("c", "counts", &argparse.Options{
		Help: "Counts the number of scraped pages, failed pages and signatures",
	})

	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Println(parser.Usage(err))
		os.Exit(1)
	}

	return Config{
		Api:         NewFourByteApi(),
		StartPage:   *startPage,
		Threads:     *threads,
		MaxRetries:  *retries,
		FindMissing: *missing,
		FailedOnly:  *failedOnly,
		CountsOnly:  *countsOnly,
	}
}
