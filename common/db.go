package common

import "github.com/jmoiron/sqlx"

func SetupTable(db *sqlx.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS txn
		  (
		  	as_of DATETIME NOT NULL,
			txn_type TEXT NOT NULL,
			amount NUMERIC NOT NULL,
			ref1 TEXT NOT NULL,
			ref2 TEXT,
			ref3 TEXT,
			hash TEXT PRIMARY KEY ON CONFLICT IGNORE,
			category TEXT
		  )
	`)
	return err
}

func SaveTxn(sqlTxn *sqlx.Tx, txn Txn) error {
	_, err := sqlTxn.NamedExec(`
		INSERT INTO txn
		(as_of, txn_type, amount, ref1, ref2, ref3, hash)
		VALUES (:as_of, :txn_type, :amount, :ref1, :ref2, :ref3, :hash);
	`, &txn)
	return err
}

func LoadUntaggedTxns(db *sqlx.DB) ([]Txn, error) {
	var txns []Txn
	err := db.Select(&txns, `
		SELECT * FROM txn
		WHERE category IS NULL
		ORDER BY as_of DESC
	`)
	return txns, err
}

func UpdateTxnCategory(db *sqlx.DB, hash string, category string) error {
	_, err := db.NamedExec(`
		UPDATE txn
		SET category = :category
		WHERE hash = :hash
	`, map[string]interface{} {
		"category": category,
		"hash": hash,
	})
	return err
}
