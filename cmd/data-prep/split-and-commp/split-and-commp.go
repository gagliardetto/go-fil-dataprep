package split_and_commp

import (
	"encoding/csv"
	"github.com/anjor/go-fil-dataprep/cmd/data-prep/utils"
	"github.com/urfave/cli/v2"
	"os"
	"strconv"
	"time"
)

var Cmd = &cli.Command{

	Name:    "split-and-commp",
	Usage:   "Split CAR and calculate commp",
	Aliases: []string{"sac"},
	Action:  splitAndCommpAction,
	Flags:   splitAndCommpFlags,
}

var splitAndCommpFlags = []cli.Flag{
	&cli.IntFlag{
		Name:     "size",
		Aliases:  []string{"s"},
		Required: true,
		Usage:    "Target size in bytes to chunk CARs to.",
	},
	&cli.StringFlag{
		Name:     "output",
		Aliases:  []string{"o"},
		Required: true,
		Usage:    "optional output name for car files. Defaults to filename (stdin if streamed in from stdin).",
	},
	&cli.StringFlag{
		Name:     "metadata",
		Aliases:  []string{"m"},
		Required: false,
		Usage:    "optional metadata file name. Defaults to __metadata.csv",
		Value:    "__metadata.csv",
	},
}

func splitAndCommpAction(c *cli.Context) error {

	fi, err := utils.GetReader(c)
	if err != nil {
		return err
	}

	size := c.Int("size")
	output := c.String("output")
	meta := c.String("metadata")

	carFiles, err := utils.SplitAndCommp(fi, size, output)
	if err != nil {
		return err
	}

	f, err := os.Create(meta)
	defer f.Close()
	if err != nil {
		return err
	}

	w := csv.NewWriter(f)
	err = w.Write([]string{"timestamp", "original data", "car file", "piece cid", "padded piece size", "unpadded piece size"})
	if err != nil {
		return err
	}
	defer w.Flush()
	for _, c := range carFiles {
		err = w.Write([]string{
			time.Now().String(),
			output,
			c.CarName,
			c.CommP.String(),
			strconv.FormatUint(c.PaddedSize, 10),
			strconv.FormatUint(c.PaddedSize/128*127, 10),
		})
	}
	return nil
}