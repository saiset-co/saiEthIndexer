**saiEthIndexer** allows to index all transactions associated with specific smart contracts and send corresponding requests to other services to store and handl these transactions.


# SYNOPSIS

## Run in Docker
`make up`

## Run as standalone application
`./sai-eth-indexer` 

# CONFIGURATION

## config/config.json (application configuration)
- common(http_server,socket_server, web_socket) - common server options for http,socket and web socket servers
- geth-server - geth-server address
- storage - options for saiStorage
- start_block - number of block to start parsing 
- operations - commands under special control
- sleep - duration after which we get next block from geth server

## config/contracts.json (stored on control contracts)
- address - address of contract
- abi - Application Binary Interface (ABI) of a smart contract 
- start_block - number of block, from which contract is valid

#API

## Add contract 
curl -X POST <host:port>/v1/add_contract  -H "Content-Type: application/json" -d '{"contracts": [{"address": "0x9fe3Ace9629468AB8858660f765d329273D94D6D","abi": "324234","start_block":123},{"address": "0x9fe3Ace9629468AB8858660f765d329273D94D6E","abi":"test","start_block":34}]}'

## Delete contract 
curl -X POST <host:port>/v1/delete_contract  -H "Content-Type: application/json" -d '{"addresses": ["0x9fe3Ace9629468AB8858660f765d329273D94D6E","0x9fe3Ace9629468AB8858660f765d329273D94D6W"]}'
