package main

import (
	"github.com/skhro87/dbs-txn-stats/importer"
	"github.com/skhro87/dbs-txn-stats/tagger"
	"github.com/urfave/cli"
	"log"
	"os"
	"time"
)

func main() {
	app := cli.NewApp()
	app.Name = "DBS Transaction Stats Calculator"
	app.Version = "0.1.0"
	app.Usage = "import and tag dbs transactions from csv file"

	flagFile := cli.StringFlag{
		Name:   "file",
		Value:  "./data.csv",
		Usage:  "path to the DBS csv file with the transactions",
		EnvVar: "file",
	}

	flagDb := cli.StringFlag{
		Name:   "db",
		Value:  "txns.db",
		Usage:  "path to the sqlite3 db file (will be automatically created if not exist)",
		EnvVar: "DB",
	}

	app.Commands = []cli.Command{
		{
			Name:    "import",
			Aliases: []string{""},
			Usage:   "import txns from file",
			Action:  importMode,
			Flags:   []cli.Flag{flagFile, flagDb},
		},
		{
			Name:    "tag",
			Aliases: []string{""},
			Usage:   "tag txns in db",
			Action:  tagMode,
			Flags:   []cli.Flag{flagDb},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatalf("err running : %v", err.Error())
	}

	time.Sleep(1 * time.Second)
}

func importMode(c *cli.Context) error {
	config := importer.Config{
		FilePath: c.String("file"),
		DbPath:   c.String("db"),
	}

	return importer.ImportFromFile(config)
}

func tagMode(c *cli.Context) error {
	config := tagger.Config{
		DbPath: c.String("db"),
	}

	return tagger.TagTransactions(config)
}
