package remote_utils

import (
	"fmt"
	log "github.com/cihub/seelog"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"
)

func Upload(client *ssh.Client, source io.Reader, remote_path string, remote_name string, mode os.FileMode) error {
	sftp, err := sftp.NewClient(client)
	if err != nil {
		return err
	}
	defer sftp.Close()
	tmp_filename := fmt.Sprint(remote_name, rand.Int())

	if s, err := client.NewSession(); err == nil {
		output, err := s.CombinedOutput("mkdir -p " + remote_path)
		if err != nil {
			return log.Error(err.Error() + (string)(output))
		}
	} else {
		return err
	}

	f, err := sftp.Create(remote_path + "/" + tmp_filename)
	if err != nil {
		return log.Error("Create ", remote_path+"/"+tmp_filename, " failed:", err)
	}
	_, err = io.Copy(f, source)
	if err != nil {
		return log.Error(err)
	}
	err = f.Chmod(mode)
	if err != nil {
		return log.Error("Chmod failed:", err)
	}
	err = f.Close()
	if err != nil {
		return log.Error("Close failed:", err)
	}

	if s, err := client.NewSession(); err == nil {
		output, err := s.CombinedOutput(fmt.Sprintf("rm -f \"%s\"", strings.Replace(remote_path+"/"+remote_name, "\"", "\\\"", -1)))
		if err != nil {
			return log.Error(err.Error() + (string)(output))
		}
	} else {
		return err
	}
	if s, err := client.NewSession(); err == nil {
		output, err := s.CombinedOutput(fmt.Sprintf("mv \"%s\" \"%s\"", strings.Replace(remote_path+"/"+tmp_filename, "\"", "\\\"", -1), strings.Replace(remote_path+"/"+remote_name, "\"", "\\\"", -1)))
		if err != nil {
			return log.Error(err.Error() + (string)(output))
		}
	} else {
		return err
	}
	return nil
}

func ReadFiles(client *ssh.Client, filenames []string) []interface{} {
	res := make([]interface{}, len(filenames))

	sftp, err := sftp.NewClient(client)
	if err != nil {
		for i := 0; i < len(res); i++ {
			res[i] = err
		}
	}
	defer sftp.Close()

	for i, filename := range filenames {
		f, err := sftp.Open(filename)
		if err != nil {
			res[i] = err
		} else {
			bytes, err := ioutil.ReadAll(f)
			if err != nil {
				res[i] = err
			} else {
				res[i] = bytes
			}
		}
	}
	return res
}
