build:
	docker-compose build
	docker-compose up -d
	docker-compose up --scale node-slave=4 -d

test:
	curl -i -H "Accept: application/json" -H "Content-Type: application/json" -X GET http://codechallenge.local