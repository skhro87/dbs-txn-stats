package tagger

import (
	"bufio"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/skhro87/dbs-txn-stats/common"
	"gopkg.in/go-playground/validator.v9"
	"log"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	DbPath string `validate:"required"`
}

func TagTransactions(config Config) error {
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

	// load untagged txns
	txns, err := common.LoadUntaggedTxns(db)
	if err != nil {
		return fmt.Errorf("err loading untagged txns : %v", err.Error())
	}

	log.Printf("loaded %v txns", len(txns))

	// tag txns
	for i, txn := range txns {
		err := tagTxn(db, i, len(txns), txn)
		if err != nil {
			return fmt.Errorf("err tagging txn : %v", err.Error())
		}
	}

	return nil
}

func tagTxn(db *sqlx.DB, i, total int, txn common.Txn) error {
	fmt.Printf("%+v\n", txn)

	category, err := readCategory(i, total, txn)
	if err != nil {
		return fmt.Errorf("err reading category : %v", err.Error())
	}

	err = common.UpdateTxnCategory(db, txn.Hash, category)
	if err != nil {
		return fmt.Errorf("err updating category : %v", err.Error())
	}

	fmt.Printf("cat %v", category)

	return nil
}

func readCategory(i, total int, txn common.Txn) (string, error) {
	categoryName := ""
	var category int
	var err error
	var input string
	for categoryName == "" || err != nil {
		clearScreen()
		fmt.Printf("\n   %v / %v\n\n", i+1, total)
		printTxn(txn)
		if err != nil {
			fmt.Printf("err reading category : %v\n\n", err.Error())
		}
		printCategories()

		fmt.Print("select category: ")
		input, err = bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			continue
		}

		input = strings.Replace(input, "\n", "", -1)

		category, err = strconv.Atoi(input)
		if err != nil {
			continue
		}

		categoryName, err = common.CategoryTranslation(category)
		if err != nil {
			continue
		}
	}

	return categoryName, nil
}
