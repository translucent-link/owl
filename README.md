# owl
An EVM blockchain indexer


## Examples

To find a block from a year ago

    owl blk --days 365

To find the current block

    owl blk  

To scan from a specific block to the current block

    owl scan -a abis/ 12309583 14660614


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

