package ssh

import (
	"bytes"
	"io"

	"github.com/abrander/alerto/logger"
	"github.com/abrander/alerto/plugins"
)

type (
	SshCommand struct {
		Ssh
	}
)

func init() {
	plugins.Register("ssh-command", NewSshCommand)
}

func NewSshCommand() plugins.Plugin {
	return new(SshCommand)
}

func (s *SshCommand) Exec(cmd string, arguments ...string) (io.Reader, io.Reader, error) {
	for _, arg := range arguments {
		cmd += " " + arg
	}

	logger.Yellow("ssh", "Executing command '%s' on %s:%d as %s", cmd, s.Ssh.Host, s.Ssh.Port, s.Username)
	conn, err := pool.Get(s.Ssh)
	if err != nil {
		return nil, nil, err
	}
	defer pool.Done(s.Ssh)

	session, err := conn.NewSession()
	if err != nil {
		return nil, nil, err
	}
	defer session.Close()

	var stdoutBuf, stderrBuf bytes.Buffer
	session.Stdout = &stdoutBuf
	session.Stderr = &stderrBuf

	err = session.Run(cmd)
	if err != nil {
		return &stdoutBuf, &stderrBuf, err
	}

	return &stdoutBuf, &stderrBuf, nil
}

// Ensure compliance
var _ plugins.Transport = (*SshCommand)(nil)
