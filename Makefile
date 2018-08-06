build:
	docker-compose build
	docker-compose up -d
	docker-compose up --scale node-slave=4 -d