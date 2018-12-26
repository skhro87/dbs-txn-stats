package importer

import (
	"bufio"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"github.com/skhro87/dbs-txn-stats/common"
	"os"
	"strconv"
	"strings"
	"time"
)

// 18 Dec 2018
const timeLayoutAsOf = "02 Jan 2006"

func readTxns(filePath string) ([]common.Txn, error) {
	// create hasher
	hasher := sha1.New()

	// read file
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("err reading file %v : %v", filePath, err.Error())
	}
	defer file.Close()

	// iterate through lines
	var txns []common.Txn
	n := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		n++

		// ignore misc lines at start (dbs too dumb too export proper .csv file)
		if n < 21 {
			continue
		}

		line := scanner.Text()

		// ignore blank lines (see dbs rant above)
		if strings.TrimSpace(line) == "" {
			continue
		}

		// check if we have 'ok' line
		parts := strings.Split(line, ",")
		if len(parts) < 7 || len(parts) > 10 {
			return nil, fmt.Errorf("err with line %v : got %v parts, expected >= 7 and <= 20 : \nline: %v", n, len(parts), line)
		}

		// hash to avoid duplicates
		hasher.Write([]byte(line))
		hash := base64.URLEncoding.EncodeToString(hasher.Sum(nil))

		// parse date
		asOf, err := time.Parse(timeLayoutAsOf, parts[0])
		if err != nil {
			return nil, fmt.Errorf("err parsing as of %v : %v", parts[0], err.Error())
		}

		// parse amount
		amount, err := parseAmount(parts[2], parts[3])
		if err != nil {
			return nil, fmt.Errorf("err parsing amount from debit %v credit %v : %v", parts[2], parts[3], err.Error())
		}

		txns = append(txns, common.Txn{
			AsOf:    asOf,
			TxnType: parts[1],
			Amount:  amount,
			Ref1:    parts[4],
			Ref2:    &parts[5],
			Ref3:    &parts[6],
			Hash:    hash,
		})
	}

	return txns, nil
}

func parseAmount(debitRaw, creditRaw string) (float64, error) {
	debitRaw = strings.TrimSpace(debitRaw)
	creditRaw = strings.TrimSpace(creditRaw)

	var err error
	debit := 0.0
	if debitRaw != "" {
		debit, err = strconv.ParseFloat(debitRaw, 64)
		if err != nil {
			return 0.0, fmt.Errorf("err parsing debit %v : %v", debitRaw, err.Error())
		}
	}

	credit := 0.0
	if creditRaw != "" {
		credit, err = strconv.ParseFloat(creditRaw, 64)
		if err != nil {
			return 0.0, fmt.Errorf("err parsing credit %v : %v", creditRaw, err.Error())
		}
	}

	return credit - debit, nil
}
