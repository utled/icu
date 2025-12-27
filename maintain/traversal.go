package maintain

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"slices"
	"snafu/data"
	"snafu/utils"
	"sync"
	"syscall"
)

func traverseNewDir(readJobs chan<- data.SyncJob, startPath string, dbPath string) error {
	/*
		allInodeMappedEntries, err := data.GetInodeMappedEntries("", dbPath)
		if err != nil {
			return err
		}
	*/
	inodeMappedEntries, err := data.GetInodeMappedEntries(dbPath)
	if err != nil {
		return err
	}

	/*
		allInodesSorted := make([]uint64, 0, len(allInodeMappedEntries))
		for inode := range allInodeMappedEntries {
			allInodesSorted = append(allInodesSorted, inode)
		}
		sort.Slice(allInodesSorted, func(i, j int) bool {
			return allInodesSorted[i] < allInodesSorted[j]
		})
	*/
	err = filepath.WalkDir(startPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		var syncJob data.SyncJob
		entryStat, err := os.Stat(path)
		if err != nil {
			return err
		}
		entryStatT := entryStat.Sys().(*syscall.Stat_t)
		if inode, ok := inodeMappedEntries[entryStatT.Ino]; ok {
			entryMtim := entryStatT.Mtim.Sec + entryStatT.Mtim.Nsec
			indexedMtim := inode.ModificationTime
			if entryStat.IsDir() || entryMtim == indexedMtim {
				syncJob = data.SyncJob{Path: path, IsIndexed: true, IsContentChange: false}
			} else {
				syncJob = data.SyncJob{Path: path, IsIndexed: true, IsContentChange: true}
			}
		} else {
			if entryStat.IsDir() {
				syncJob = data.SyncJob{Path: path, IsIndexed: false, IsContentChange: false}
			} else {
				syncJob = data.SyncJob{Path: path, IsIndexed: false, IsContentChange: true}
			}
		}

		/*

			_, exists := slices.BinarySearch(allInodesSorted, entryStatT.Ino)
			if exists {
				entryMtim := entryStatT.Mtim.Sec + entryStatT.Mtim.Nsec
				indexedMtim := allInodeMappedEntries[entryStatT.Ino].ModificationTime
				if entryStat.IsDir() || entryMtim == indexedMtim {
					syncJob = data.SyncJob{Path: path, IsIndexed: true, IsContentChange: false}
				} else {
					syncJob = data.SyncJob{Path: path, IsIndexed: true, IsContentChange: true}
				}
			} else {
				if entryStat.IsDir() {
					syncJob = data.SyncJob{Path: path, IsIndexed: false, IsContentChange: false}
				} else {
					syncJob = data.SyncJob{Path: path, IsIndexed: false, IsContentChange: true}
				}
			}

		*/
		readJobs <- syncJob
		return nil
	})

	return nil
}

func traverseDirectories(
	startPath string,
	scanJobs chan<- data.InodeHeader,
	newDirJobs chan<- string,
	readJobs chan<- data.SyncJob,
	wg *sync.WaitGroup,
	inodeMappedEntries map[uint64]data.InodeHeader,
) {
	defer wg.Done()

	err := filepath.WalkDir(startPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		entryStat, err := os.Stat(path)
		if err != nil {
			return err
		}

		if d.IsDir() && slices.Contains(utils.ExcludedEntries, filepath.Base(path)) {
			return filepath.SkipDir
		}

		statT := entryStat.Sys().(*syscall.Stat_t)

		if d.IsDir() {
			if _, ok := inodeMappedEntries[statT.Ino]; !ok {
				newDirJobs <- path
			} else {
				for inode, values := range inodeMappedEntries {
					if inode != statT.Ino {
						continue
					}
					mtim := statT.Mtim.Sec + statT.Mtim.Nsec
					ctim := statT.Ctim.Sec + statT.Ctim.Nsec
					if values.ModificationTime != mtim || values.MetaDataChangeTime != ctim {
						readJobs <- data.SyncJob{Path: path, IsIndexed: true, IsContentChange: false}
						scanJobs <- values
						continue
					}
				}
			}
		}

		return nil
	})

	if err != nil {
		log.Printf("Fatal error during directory traversal: %v", err)
	}
}
