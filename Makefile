up:
	docker-compose -f ./microservices/docker-compose.yml up -d

down:
	docker-compose -f ./microservices/docker-compose.yml down --remove-orphans

build:
	make service
	make docker

service:
	cd ./src/saiEthIndexer/cmd/app && go build -o ../../../../microservices/saiEthIndexer/build/sai-eth-indexer

docker:
	docker-compose -f ./microservices/docker-compose.yml up -d --build

logs:
	docker-compose -f ./microservices/docker-compose.yml logs -f

logn:
	docker-compose -f ./microservices/docker-compose.yml logs -f sai-gn-monitor