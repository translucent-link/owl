# owl
An EVM blockchain indexer


## Examples

To find a block from a year ago

    go run main.go blk -u wss://mainnet.infura.io/ws/v3/ea9cb700bd834903ab2e0567faa630a9 --days 365  # returns 12309583

To find the current block

    go run main.go blk -u wss://mainnet.infura.io/ws/v3/ea9cb700bd834903ab2e0567faa630a9  # returns 14660614

To scan from a specific block to the current block

    go run main.go scan -u wss://mainnet.infura.io/ws/v3/ea9cb700bd834903ab2e0567faa630a9 -a abis/ 12309583 14660614
