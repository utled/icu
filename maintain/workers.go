package maintain

import (
	"log"
	"snafu/data"
	"sync"
)

func scanWorker(scanJobs <-chan data.InodeHeader, readJobs chan<- data.SyncJob, inodeMappedEntries map[uint64]data.InodeHeader, wg *sync.WaitGroup) {
	defer wg.Done()
	for job := range scanJobs {
		err := scanUpdatedDir(readJobs, job.Path, inodeMappedEntries)
		if err != nil {
			log.Fatal(err)
		}
	}
}
func readWorker(readJobs <-chan data.SyncJob, dbPath string, wg *sync.WaitGroup) {
	defer wg.Done()
	for job := range readJobs {
		readEntry(job, dbPath)
	}
}
func newDirWorker(newDirJobs <-chan string, readJobs chan<- data.SyncJob, dbPath string, wg *sync.WaitGroup) {
	defer wg.Done()
	for path := range newDirJobs {
		err := traverseNewDir(readJobs, path, dbPath)
		if err != nil {
			log.Fatal(err)
		}
	}
}
func deletionWorker(delJobs <-chan string, dbPath string, wg *sync.WaitGroup) {
	defer wg.Done()
	for path := range delJobs {
		err := checkDelete(path, dbPath)
		if err != nil {
			log.Fatal(err)
		}
	}
}
