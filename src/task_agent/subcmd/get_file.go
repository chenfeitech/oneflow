package subcmd

import (
	"errors"
	"fmt"
	"github.com/urfave/cli"
	"utils/helper"
	"io"
	"math/rand"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
)

func init() {
	Register(cli.Command{
		Name:    "getfile",
		Aliases: []string{"g"},
		Usage:   "Get file from flow file server",
		Action:  getfile,
		Flags: []cli.Flag{
			cli.StringSliceFlag{
				Name:  "file, f",
				Usage: "file path",
			},
		},
	})
}

func getfile(c *cli.Context) {
	files := c.StringSlice("file")

	succeed := true
	for _, require_file := range files {
		err := RequireFile(require_file)
		if err != nil {
			fmt.Fprintln(os.Stderr, "getfile: File ", require_file, ":", err, ".")
			succeed = false
		}
	}

	if succeed {
		return
	} else {
		os.Exit(1)
	}
}

func RequireFile(require_file string) error {
	if strings.Trim(require_file, " ") == "" {
		return nil
	}
	fmt.Fprintln(os.Stdout, "run: Process require file "+require_file+".")
	file, err := os.Open(require_file)
	if err == nil {
		// 检查本地文件
		fstat, err := file.Stat()
		if err == nil {
			if fstat.IsDir() {
				return errors.New("File is exists and is a directory.")
			}
			resp, err := http.Head(GetServerUrl("/file" + require_file))
			if err != nil {
				return err
			}
			defer resp.Body.Close()
			if resp.StatusCode == http.StatusNotFound {
				return nil
			}

			local_md5 := helper.GetFileMd5(file)
			if resp.Header.Get("MD5") == local_md5 {
				mode, err := strconv.ParseUint(resp.Header.Get("Mode"), 8, 32)
				if err == nil {
					if (uint32)(mode) != (uint32)(fstat.Mode()) {
						err := os.Chmod(require_file, (os.FileMode)((uint32)(mode)))
						if err != nil {
							fmt.Fprintln(os.Stdout, "run: Chmod ", resp.Header.Get("Mode"), " failed: ", err, ".")
						}
					}
				} else {
					fmt.Fprintln(os.Stdout, "run: Parse mode ", resp.Header.Get("Mode"), " failed: ", err, ".")
				}
				return nil
			}
		}
		file.Close()
	}
	// Download file
	resp, err := http.Get(GetServerUrl("/file" + require_file))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return errors.New("Get file return:" + resp.Status)
	}

	mode, err := strconv.ParseUint(resp.Header.Get("Mode"), 8, 32)
	if err != nil {
		return errors.New("Parse mode " + resp.Header.Get("Mode") + " failed: " + err.Error())
	}
	dir, _ := path.Split(require_file)
	if err := os.MkdirAll(dir, os.ModeDir|os.ModePerm); err != nil {
		return err
	}

	tmp_file := fmt.Sprint(require_file, "_tmp_", rand.Int())
	file, err = os.OpenFile(tmp_file, os.O_RDWR|os.O_CREATE|os.O_TRUNC, (os.FileMode)((uint32)(mode)))
	if err != nil {
		return err
	}
	written, err := io.Copy(file, resp.Body)
	if err != nil {
		return err
	}

	file.Close()
	file_size, err := strconv.ParseInt(resp.Header.Get("Size"), 10, 64)
	if err != nil {
		if file_size != written {
			return fmt.Errorf("Download file size error expect(%d) got(%d)", file_size, written)
		}
	}

	os.Remove(require_file)

	if err = os.Rename(tmp_file, require_file); err != nil {
		return fmt.Errorf("Rename require file failed:", err)
	}
	return nil
}
