package src

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"log"
	"os"
	"sync"
	"time"
)

const (
	MaxAttempts = 50
)

type WorkerPool struct {
	api       FourByteApi
	page      int
	exporter  *Exporter
	pageMutex sync.Mutex
	logMutex  sync.Mutex
	waitGroup sync.WaitGroup
	pages     []int
	config    Config
}

func NewWorkerPool(config Config, exporter *Exporter) WorkerPool {
	if config.FailedOnly {
		return newFailedPagesWorkerPool(config, exporter)
	}

	if config.FindMissing {
		return newMissingPagesWorkerPool(config, exporter)
	}

	log.Println(fmt.Sprintf("Starting workers from page %d", config.StartPage))

	return WorkerPool{
		api:       config.Api,
		page:      config.StartPage,
		exporter:  exporter,
		pageMutex: sync.Mutex{},
		waitGroup: sync.WaitGroup{},
		config:    config,
	}
}

func newFailedPagesWorkerPool(config Config, exporter *Exporter) WorkerPool {
	log.Println("Starting workers in --failed-only mode")

	return WorkerPool{
		api:       config.Api,
		exporter:  exporter,
		pageMutex: sync.Mutex{},
		waitGroup: sync.WaitGroup{},
		pages:     exporter.FailedPages,
		config:    config,
	}
}

func newMissingPagesWorkerPool(config Config, exporter *Exporter) WorkerPool {
	missing := exporter.findMissing()
	if len(missing) <= 0 {
		log.Println("No missing pages were found, exiting!")
		os.Exit(0)
	}

	log.Println(fmt.Sprintf("Found %d missing pages, fetching...", len(missing)))
	log.Println("Starting workers in --missing mode")

	spew.Dump(missing)

	return WorkerPool{
		api:       config.Api,
		exporter:  exporter,
		pageMutex: sync.Mutex{},
		waitGroup: sync.WaitGroup{},
		pages:     missing,
		config:    config,
	}
}

func (p *WorkerPool) Start() {
	for i := 1; i <= p.config.Threads; i++ {
		p.waitGroup.Add(1)
		go p.worker(p.exporter.SignaturesChannel, p.exporter.FailedPagesChannel, p.exporter.CompletedPagesChannel)
	}

	p.waitGroup.Wait()
	p.exporter.forceSaveAndExit()
}

func (p *WorkerPool) worker(signatures chan<- []FourByteSignature, failed chan<- int, completed chan<- int) {
	for {
		attempt := 1
		found := false
		page := p.getNewPage()

		if page == 0 {
			break
		}

		// Due to the API being fairly unstable we retry any failed attempts
		for i := attempt; i <= p.config.MaxRetries; i++ {
			results, err := p.api.FetchPage(page)

			if err != nil {
				time.Sleep(3 * time.Second)
				continue
			}

			p.log(fmt.Sprintf("Successfully fetched page %d", page))
			signatures <- results
			completed <- page
			found = true
			break
		}

		if !found {
			p.log(fmt.Sprintf("Page %d failed after %d attempts", page, p.config.MaxRetries))
			failed <- page
		}
	}
}

// todo: maybe improve this so that failed/missing use the same slice
func (p *WorkerPool) getNewPage() (page int) {
	p.pageMutex.Lock()
	defer p.pageMutex.Unlock()

	if p.config.FailedOnly {
		// If there are no more pages
		if len(p.exporter.FailedPages) <= 0 {
			p.waitGroup.Done()
			return
		}

		page, p.exporter.FailedPages = p.exporter.FailedPages[0], p.exporter.FailedPages[1:]

		return page
	}

	if p.config.FindMissing {
		if len(p.pages) <= 0 {
			p.waitGroup.Done()
			return
		}

		page, p.pages = p.pages[0], p.pages[1:]
		return page
	}

	page = p.page
	p.page++

	return page
}

func (p *WorkerPool) log(message string) {
	p.logMutex.Lock()
	log.Println(message)
	p.logMutex.Unlock()
}
