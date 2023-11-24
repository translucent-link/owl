package model

type ProtocolInstance struct {
	ID         int `json:"id"`
	Protocol   *Protocol
	Chain      *Chain
	ProtocolID int `json:"protocolId"`
	ChainID    int `json:"chainId"`

	ContractAddress  string `json:"contractAddress"`
	FirstBlockToRead int    `json:"firstBlockToRead"`
	LastBlockRead    int    `json:"lastBlockRead"`
}
