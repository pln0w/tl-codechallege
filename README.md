# Files watcher [GoLang + WebSockets]

The scope of this project is to build and learn GoLang and others stuff like microservices, Docker and WebSockets as well.

There is one HTTP server able to recive commands. We spawn multiple workers listening on some directories. 
* We call server for list of files withing those directories  
* Server communicates with watchers and concurrently collects data to send back
* There is a dashboard (fron-end app) that periodically ask server about files tree

## How to build
`make build` - builds images and runs and attaches server container

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
