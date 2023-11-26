# owl
An EVM blockchain indexer


## Local Setup

Run the migrations

```bash
migrate -source file://./db/migrations -database "postgres://postgres:hoothoo@localhost:5432/owl_dev?sslmode=disable" 
up
migrate -source file://./db/migrations -database "postgres://postgres:hoothoo@localhost:5432/owl_test?sslmode=disable" up
```


## Examples

To find a block from a year ago

    owl blk --chain ethereum --days 365

To find the current block

    owl blk  --chain ethereum

To scan from a specific block to the current block

    owl scan -a abis/ 12309583 14660614


To convert an ABI event into a topic hash

    owl abi topicHash "Transfer(address,address,uint256,bytes)"

To register a new chain

    owl chain register --name polygon --url wss://polygon-mainnet.g.alchemy.com/v2/<APIKEY> --chainID 137 --nativeToken 0x0d500b1d8e8ef31e21c99d1db9a6444d3adf1270 --nativeTokenDecimals 18 --nativeTokenSymbol MATIC

To delete a chain

    owl chain delete --name polygon


POLYGON

To find a block from a year ago

    owl blk -u wss://polygon-mainnet.g.alchemy.com/v2/<APIKEY> --days 365  # returns 12309583

To find the current block

    owl blk -u wss://polygon-mainnet.g.alchemy.com/v2/<APIKEY>  # returns 14660614


To scan from a specific block to the current block on Polygon

    owl scan -u wss://polygon-mainnet.g.alchemy.com/v2/<APIKEY> -a abis/ 1203314 14014378
         
         


             owl scan -u wss://polygon-mainnet.g.alchemy.com/v2/<APIKEY> -a abis/ 20000000 26648474


AVALANCHE             


To find a block from a year ago

    owl blk -u https://api.avax.network/ext/bc/C/rpc --days 365  

To find the current block

    owl blk -u https://api.avax.network/ext/bc/C/rpc  


To scan from a specific block to the current block on Polygon

    owl scan -u https://api.avax.network/ext/bc/C/rpc -a abis/ 14000000 14169608
         
         


             owl scan -u https://api.avax.network/ext/bc/C/rpc -a abis/ 20000000 26648474

