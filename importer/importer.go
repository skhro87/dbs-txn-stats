package importer

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/skhro87/dbs-txn-stats/common"
	"gopkg.in/go-playground/validator.v9"
	)

type Config struct {
	FilePath string `validate:"required"`
	DbPath   string `validate:"required"`
}

func ImportFromFile(config Config) error {
	// validate config
	err := validator.New().Struct(config)
	if err != nil {
		return fmt.Errorf("err validating config : %v", err.Error())
	}

	// connect db
	db, err := sqlx.Connect("sqlite3", config.DbPath)
	if err != nil {
		return fmt.Errorf("err connecting to db : %v", err.Error())
	}

	// setup table
	err = common.SetupTable(db)
	if err != nil {
		return fmt.Errorf("err setting up table : %v", err.Error())
	}

	// read txns from file
	txns, err := readTxns(config.FilePath)
	if err != nil {
		return fmt.Errorf("err reading txns from file %v : %v", config.FilePath, err.Error())
	}

	// create sql txn and write txns to db
	sqlTxn, err := db.Beginx()
	if err != nil {
		return fmt.Errorf("err beginning sql txn : %v", err.Error())
	}

	for _, txn := range txns {
		err = common.SaveTxn(sqlTxn, txn)
		if err != nil {
			return fmt.Errorf("err saving txn to db : %v", err.Error())
		}
	}

	// commit sql txn
	err = sqlTxn.Commit()
	if err != nil {
		return fmt.Errorf("err committing sql txn : %v", err.Error())
	}

	fmt.Printf("\n\nimported %v transactions!\n\ndone.", len(txns))

	return nil
}
