# Files watcher [GoLang + WebSockets]

The scope of this project is to build complete working solution, just to learn GoLang and other stuff like WebSockets.

There is single HTTP server, able to answer with full list of files within _/data_ directory. We spawn multiple workers listening on particular directories. 
* We call server via HTTP for connected workers list
* HTTP server periodically (10s) communicates with watchers in separate goroutines to refresh list of files per each directory, by calling worker via WebSocket
* If there is a change on disk, worker calls server to update it's list immediately (not implemented yet)

## How to build
`make build` - builds images and runs and attaches server container  
`make run` - only runs and attaches server container


## How to test entire flow
`make full-test` - automatically presents entire flow

## Makefile commands to play with

* Add watcher  
`make add-watcher dir=/data/anydir` - spawn another watcher node

* Run watcher  
`make run-watcher dir=/data/anydir` - spawn another watcher node and attach container to see whats in

* Dump all watchers  
`make dump-watchers` - call server to response with all connected watchers

* Call watchers end point  
`make simple-test` - call server to response with all files within all watchers directories
