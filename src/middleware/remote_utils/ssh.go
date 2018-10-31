package remote_utils

import (
	"fmt"
	"golang.org/x/crypto/ssh"
)

func Connect(addr string, user string, password string) (*ssh.Client, error) {
	authMethods := []ssh.AuthMethod{}
	keyboardInteractiveChallenge := func(
		u,
		instruction string,
		questions []string,
		echos []bool,
	) (answers []string, err error) {
		if len(questions) == 0 {
			return []string{}, nil
		}
		return []string{password}, nil
	}
	authMethods = append(authMethods, ssh.KeyboardInteractive(keyboardInteractiveChallenge))
	authMethods = append(authMethods, ssh.Password(password))
	sshConfig := &ssh.ClientConfig{
		User: user,
		Auth: authMethods,
	}
	sshConfig.Config.SetDefaults()
	sshConfig.Config.Ciphers = append(sshConfig.Config.Ciphers, "arcfour")
	client, err := ssh.Dial("tcp", fmt.Sprintf("%s", addr), sshConfig)
	if err != nil {
		return nil, err
	}
	return client, nil
}
