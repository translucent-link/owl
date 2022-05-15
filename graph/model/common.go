package model

func Stores() (*ChainStore, *ProtocolStore, *ProtocolInstanceStore, error) {
	protocolStore, err := NewProtocolStore()
	if err != nil {
		return &ChainStore{}, &ProtocolStore{}, &ProtocolInstanceStore{}, err
	}
	chainStore, err := NewChainStore()
	if err != nil {
		return &ChainStore{}, &ProtocolStore{}, &ProtocolInstanceStore{}, err
	}
	protocolInstanceStore, err := NewProtocolInstanceStore()
	if err != nil {
		return &ChainStore{}, &ProtocolStore{}, &ProtocolInstanceStore{}, err
	}

	return chainStore, protocolStore, protocolInstanceStore, nil
}
