package ssh

import (
	"fmt"
	"io/ioutil"

	"golang.org/x/crypto/ssh"

	"github.com/abrander/alerto/logger"
)

type (
	Ssh struct {
		Host     string `json:"host"`
		Port     uint16 `json:"port" default:22`
		Username string `json:"username"`
	}
)

func (s *Ssh) Connect() (*ssh.Client, error) {
	// FIXME: Support default
	if s.Port == 0 {
		s.Port = 22
	}

	dialString := fmt.Sprintf("%s:%d", s.Host, s.Port)
	logger.Yellow("ssh", "Connecting to %s as %s", dialString, s.Username)

	pemBytes, err := ioutil.ReadFile("/home/abrander/src/github.com/abrander/alerto/private.pem")
	if err != nil {
		return nil, err
	}

	signer, err := ssh.ParsePrivateKey(pemBytes)
	if err != nil {
		return nil, err
	}
	config := &ssh.ClientConfig{
		User: s.Username,
		Auth: []ssh.AuthMethod{ssh.PublicKeys(signer)},
	}
	client, err := ssh.Dial("tcp", dialString, config)
	if err != nil {
		return nil, err
	}

	return client, nil
}
