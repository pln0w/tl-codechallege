# Fun project - communication between microservices, using WebSockets

## How to build
1. `make run` # build images  
2. `docker-compose up` # run stack and attach containers output

Add watchers  
`docker-compose run -d watcher-node` # run watcher   

## How to test  
`make test` # hit GET / http://filewatcher.local