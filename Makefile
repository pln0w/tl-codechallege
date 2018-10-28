build:
	docker-compose build
	docker-compose up --scale watcher-node=2

test:
	curl -i -H "Accept: application/json" -H "Content-Type: application/json" -X GET http://filewatcher.local