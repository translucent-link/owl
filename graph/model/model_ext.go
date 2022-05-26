package model

func (pi *ProtocolInstance) ScanStartBlock() int {
	if pi.FirstBlockToRead > pi.LastBlockRead {
		return pi.FirstBlockToRead
	}
	return pi.LastBlockRead
}
