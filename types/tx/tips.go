package tx

// TipTx defines the interface to be implemented by Txs that handle Tips.
type TipTx interface {
	GetTip() *Tip
}
