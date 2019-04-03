package subcmd

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/urfave/cli"
	"io"
	"os"
)

func init() {
	Register(cli.Command{
		Name:    "checksum",
		Aliases: []string{"c"},
		Usage:   "Calculate task agent checksum",
		Action:  checksum,
	})
}

func checksum(c *cli.Context) {
	md5_sum := md5.New()
	bin, _ := os.Open(os.Args[0])
	io.Copy(md5_sum, bin)
	fmt.Println(hex.EncodeToString(md5_sum.Sum(nil)))
}
