build:
	docker-compose build
	docker-compose up --scale watcher-node=0

run:
	docker-compose up --scale watcher-node=0

add-watcher:
	docker-compose run -e DIR=$(dir) -d watcher-node

run-watcher:
	docker-compose run -e DIR=$(dir) watcher-node

dump-watchers:
	curl -i -H "Accept: application/json" -H "Content-Type: application/json" -X GET http://filewatcher.local/watchers

simple-test:
	curl -i -H "Accept: application/json" -H "Content-Type: application/json" -X GET http://filewatcher.local

full-test:
	make add-watcher dir=/data/aaa
	make add-watcher dir=/data/bbb
	make add-watcher dir=/data/ccc
	mkdir ./data/aaa
	mkdir ./data/bbb
	mkdir ./data/ccc
	touch ./data/aaa/test_file1.data
	touch ./data/bbb/test_file1.data
	touch ./data/ccc/test_file1.data
	curl -i -H "Accept: application/json" -H "Content-Type: application/json" -X GET http://filewatcher.local

define do_cleanup = 
	rm -rf ./data/*
	docker rm $(docker ps -a -q) -f
endef
clean: ; @$(value do_cleanup)
