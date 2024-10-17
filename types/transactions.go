package types

type Status int

const (
	STATUS_OK Status = iota
	STATUS_FAIL
)

type CommittedTx struct {
	Tx TransferTx

	TxStatus Status
}

type TransferTx struct {
	Source      string
	Destination string
	Amount      uint32
}
