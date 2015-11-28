package ssh

import (
	"bytes"
	"strings"

	"github.com/abrander/alerto/agent"
	"github.com/abrander/alerto/logger"
)

type (
	SshCommand struct {
		Ssh
		Command string `json:"command"`
	}
)

func init() {
	agent.Register("ssh-command", NewSshCommand)
}

func NewSshCommand() agent.Agent {
	return new(SshCommand)
}

func (s SshCommand) Execute(request agent.Request) agent.Result {
	logger.Yellow("ssh", "Executing command '%s' on %s:%d as %s", s.Command, s.Ssh.Host, s.Ssh.Port, s.Username)
	conn, err := pool.Get(s.Ssh)
	if err != nil {
		return agent.NewResult(agent.Failed, nil, err.Error())
	}
	defer pool.Done(s.Ssh)

	session, err := conn.NewSession()
	if err != nil {
		return agent.NewResult(agent.Failed, nil, err.Error())
	}
	defer session.Close()
	var stdoutBuf bytes.Buffer
	session.Stdout = &stdoutBuf

	err = session.Run(s.Command)
	if err != nil {
		return agent.NewResult(agent.Failed, nil, err.Error())
	}

	return agent.NewResult(agent.Ok, nil, strings.TrimSpace(stdoutBuf.String()))
}
