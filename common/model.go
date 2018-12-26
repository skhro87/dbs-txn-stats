package common

import "time"

type Txn struct {
	AsOf     time.Time `db:"as_of"`
	TxnType  string    `db:"txn_type"`
	Amount   float64   `db:"amount"`
	Ref1     string    `db:"ref1"`
	Ref2     *string   `db:"ref2"`
	Ref3     *string   `db:"ref3"`
	Hash     string    `db:"hash"`
	Category *string   `db:"category"`
}
