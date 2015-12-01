package ssh

import (
	"bytes"
	"strings"

	"github.com/abrander/alerto/logger"
	"github.com/abrander/alerto/plugins"
)

type (
	SshCommand struct {
		Ssh
		Command string `json:"command" description:"Command to execute on remote host"`
	}
)

func init() {
	plugins.Register("ssh-command", NewSshCommand)
}

func NewSshCommand() plugins.Plugin {
	return new(SshCommand)
}

func (s SshCommand) Execute(request plugins.Request) plugins.Result {
	logger.Yellow("ssh", "Executing command '%s' on %s:%d as %s", s.Command, s.Ssh.Host, s.Ssh.Port, s.Username)
	conn, err := pool.Get(s.Ssh)
	if err != nil {
		return plugins.NewResult(plugins.Failed, nil, err.Error())
	}
	defer pool.Done(s.Ssh)

	session, err := conn.NewSession()
	if err != nil {
		return plugins.NewResult(plugins.Failed, nil, err.Error())
	}
	defer session.Close()
	var stdoutBuf bytes.Buffer
	session.Stdout = &stdoutBuf

	err = session.Run(s.Command)
	if err != nil {
		return plugins.NewResult(plugins.Failed, nil, err.Error())
	}

	return plugins.NewResult(plugins.Ok, nil, strings.TrimSpace(stdoutBuf.String()))
}
