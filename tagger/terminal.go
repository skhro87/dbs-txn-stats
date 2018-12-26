package tagger

import (
	"fmt"
	"github.com/skhro87/dbs-txn-stats/common"
	"os"
	"os/exec"
)

func clearScreen() {
	cmd := exec.Command("clear") //Linux example, its tested
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func printCategories() {
	for i, category := range common.Categories {
		fmt.Printf("\t %v: %v\n", i+1, category)
	}
	fmt.Println()
}

func printTxn(txn common.Txn) {
	fmt.Println()
	fmt.Printf("\t%v\t%v\n", txn.AsOf.Format("02 01 2006"), txn.TxnType)
	fmt.Printf("\t%v %v %v\n", txn.Ref1, *txn.Ref2, *txn.Ref3)
	fmt.Printf("\t%v\n", txn.Amount)
	fmt.Println()
	fmt.Println()
}