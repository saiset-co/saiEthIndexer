up:
	docker-compose -f ./microservices/docker-compose.yml up -d

down:
	docker-compose -f ./microservices/docker-compose.yml down --remove-orphans

build:
	make service
	make docker

service:
		cd ./src/saiEthIndexer/app && go mod tidy && go build -o ./../../../microservices/saiEthIndexer/build/sai-eth-indexer
		cp ./src/saiEthIndexer/config/config.json ./microservices/saiEthIndexer/build/config/config.json
		cp ./src/saiEthIndexer/config/contracts.json ./microservices/saiEthIndexer/build/config/contracts.json	
docker:
	docker-compose -f ./microservices/docker-compose.yml up -d --build

log:
	docker-compose -f ./microservices/docker-compose.yml logs -f

logidx:
	docker-compose -f ./microservices/docker-compose.yml logs -f sai-eth-indexer


