package helper

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"os"
)

func GetFileMd5(file *os.File) string {
	md5_sum := md5.New()
	io.Copy(md5_sum, file)
	md5_str := hex.EncodeToString(md5_sum.Sum(nil))
	return md5_str
}
