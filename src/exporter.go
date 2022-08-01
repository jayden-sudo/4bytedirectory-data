package src

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"sync"
	"time"
)

const (
	SignatureExportFile = "exports/signatures.json"
	CompletedPagesFile  = "exports/complete.json"
	FailedPagesFile     = "exports/failed.json"
	SaveDelay           = 60 // 1 minute
)

type (
	Exporter struct {
		Signatures      map[string][]string
		SignaturesMutex sync.Mutex
		FailedPages     []int
		CompletedPages  []int

		SignaturesChannel     chan []FourByteSignature
		FailedPagesChannel    chan int
		CompletedPagesChannel chan int
	}
)

func NewExporter() *Exporter {
	e := Exporter{
		SignaturesChannel:     make(chan []FourByteSignature),
		FailedPagesChannel:    make(chan int),
		CompletedPagesChannel: make(chan int),
	}

	e.loadFiles()

	return &e
}

func (e *Exporter) Start() {
	go e.saveWorker()
	go e.signaturesListener()
	go e.failedListener()
	go e.completedListener()
}

func (e *Exporter) loadFiles() {
	log.Println("Loading files...")

	var signatures = make(map[string][]string)
	var failed []int
	var completed []int

	signatureData, err := ioutil.ReadFile(SignatureExportFile)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(signatureData, &signatures)
	if err != nil {
		log.Fatal(err)
	}

	failedData, err := ioutil.ReadFile(FailedPagesFile)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(failedData, &failed)
	if err != nil {
		log.Fatal(err)
	}

	completedData, err := ioutil.ReadFile(CompletedPagesFile)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(completedData, &completed)
	if err != nil {
		log.Fatal(err)
	}

	e.Signatures = signatures
	e.FailedPages = failed
	e.CompletedPages = completed

	log.Println("Files loaded!")
}

func (e *Exporter) forceSaveAndExit() {
	log.Println("Force saving...")
	e.saveSignatures()
	e.saveFailed()
	e.saveCompleted()
	log.Println("Files saved, exiting!")
	os.Exit(0)
}

// saves signatures to file
func (e *Exporter) saveWorker() {
	for {
		time.Sleep(SaveDelay * time.Second)

		e.saveSignatures()
		e.saveFailed()
		e.saveCompleted()
	}
}

func (e *Exporter) findMissing() []int {
	var expected []int
	for i := 1; i <= e.CompletedPages[len(e.CompletedPages)-1]; i++ {
		expected = append(expected, i)
	}

	sort.Ints(e.CompletedPages)
	return difference(expected, e.CompletedPages)
}

func difference(a, b []int) []int {
	mb := make(map[int]struct{}, len(b))
	for _, x := range b {
		mb[x] = struct{}{}
	}
	var diff []int
	for _, x := range a {
		if _, found := mb[x]; !found {
			diff = append(diff, x)
		}
	}
	return diff
}

func (e *Exporter) saveSignatures() {
	log.Println("Saving signatures...")
	e.SignaturesMutex.Lock()
	defer e.SignaturesMutex.Unlock()

	data, err := json.MarshalIndent(e.Signatures, "", "\t")
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile(SignatureExportFile, data, 0644)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(fmt.Sprintf("Saved %d signatures to %s", len(e.Signatures), SignatureExportFile))
}

func (e *Exporter) saveFailed() {
	data, err := json.MarshalIndent(e.FailedPages, "", "\t")
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile(FailedPagesFile, data, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func (e *Exporter) saveCompleted() {
	data, err := json.MarshalIndent(e.CompletedPages, "", "\t")
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile(CompletedPagesFile, data, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func (e *Exporter) signaturesListener() {
	for signatures := range e.SignaturesChannel {
		e.SignaturesMutex.Lock()
		for _, sig := range signatures {
			// Check if this hex signature is already stored (could be a collision)
			if values, found := e.Signatures[sig.Hex]; found {

				// Check if this text signature is already stored
				if !e.contains(values, sig.Text) {
					e.Signatures[sig.Hex] = append(values, sig.Text)
				}
			} else {
				e.Signatures[sig.Hex] = []string{
					sig.Text,
				}
			}
		}
		e.SignaturesMutex.Unlock()
	}
}

func (e *Exporter) failedListener() {
	for page := range e.FailedPagesChannel {
		if !e.containsInt(e.FailedPages, page) {
			e.FailedPages = append(e.FailedPages, page)
		}
	}
}

func (e *Exporter) completedListener() {
	for page := range e.CompletedPagesChannel {
		if !e.containsInt(e.CompletedPages, page) {
			e.CompletedPages = append(e.CompletedPages, page)
		}
	}
}

func (e *Exporter) contains(slice []string, value string) bool {
	for _, i := range slice {
		if i == value {
			return true
		}
	}

	return false
}

func (e *Exporter) containsInt(slice []int, value int) bool {
	for _, i := range slice {
		if i == value {
			return true
		}
	}

	return false
}
