package model

import "github.com/ethereum/go-ethereum/ethclient"

func (c Chain) EthClient() (*ethclient.Client, error) {
	return ethclient.Dial(c.RPCURL)
}
