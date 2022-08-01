# 4byte.directory data

A full collection of signatures scraped from https://www.4byte.directory/

## Exports

All signatures can be found in `exports/signatures.json`

Due to signature collisions, the signatures are stored as an array.

## Usage

```
usage: export [-h|--help] [-p|--page <integer>] [-t|--threads <integer>]
              [-r|--retries <integer>] [-m|--missing] [-f|--failed-only]
              [-c|--counts]

              Exports function signature data from 4byte.directory

Arguments:

  -h  --help         Print help information
  -p  --page         The page to start scraping from. Default: 1
  -t  --threads      The number of threads. Default: 10
  -r  --retries      The number of times a page should be retried before it is
                     considered failed. Default: 25
  -m  --missing      Checks the completed pages array and looks for any pages
                     that may be missing and fetches them
  -f  --failed-only  If set only pages that have previously failed will be
                     processed
  -c  --counts       Counts the number of scraped pages, failed pages and
                     signatures
```

### Standard scrape

To scrape all the pages starting from the first page

```
export -p 1
```
---
### Failed pages only

Due to the instability of the 4byte.directory API, if a page could not be fetched after the retry value specified it is considered failed.

To retry all previously failed pages:

```
export --failed-only
```

---
### Missing pages only

Check if any pages are missing from the result set and rescrape them if needed.

```
export --missing
```
---
### View result totals

```
export --counts
```

Result totals as of 2022-08-01:

```
Completed pages: 8115, Failed pages: 0, Total Signatures: 810744
```

## TODO

- [ ] Add `--latest` arg to fetch the most recently updated signatures