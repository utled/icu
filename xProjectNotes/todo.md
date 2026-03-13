# TODO
## Initial
- [x] set up traversal and metadata collection
- [x] set up file content reading
- [x] set up db initialization (sqlite)
- [] set up service initialization
- [x] write collected data to db
- [] strategy for more efficient DB write on initial scan
  - [x] persistent DB Con per scan run
- [x] set up cli to manage program
- [x] set up orchestration
- [x] set up file system change monitoring
  - [x] decide sync and monitoring strategy
  - [] set up prioritization
  - [x] implement workflow
  - [] implement change logging
- [x] implement basic search
- [x] implement TUI for search
- [x] implement additional tagging
- [] figure out what else it's supposed to do
- [] figure out what to do in life

## General
- [] add setup of .conf instead of hardcoded paths (and decide lang for .conf)
- [] (if other than .json) write parser for .conf
- [] configure DB for WAL
- [] set up separate DB for contents fit FTS5 for text search
- [] how to solve sudo permissions (if needed)
- [] set up autostart
- [] how to manage index status (log file(?) with 'last fresh index', 'last index sync' et.c.)
- [] set up robust error handling/logging
- [] explore options for defining file types for content reading
- [] explore options for defining excluded objects


# BUG FIXES & CHANGES
## General
- [x] change time representations from combined Sec+Nsec to time.Time objects
- [] fix the multiplied creation of new directories
- [x] store content snippets without regex. only regex full content
- [] update db writes to separate metadata and contents
