package mcp

import (
	"os/exec"
)

type Transport interface {
	Connect() (*Client, error)
	Close() error
}

type STDIOTransport struct {
	cmd *exec.Cmd
}

func NewSTDIOTransport(command string, args []string, env map[string]string) *STDIOTransport {
	cmd := exec.Command(command, args...)
	if env != nil {
		for k, v := range env {
			cmd.Env = append(cmd.Env, k+"="+v)
		}
	}
	return &STDIOTransport{cmd: cmd}
}

func (t *STDIOTransport) Connect() (*Client, error) {
	stdin, err := t.cmd.StdinPipe()
	if err != nil {
		return nil, err
	}
	stdout, err := t.cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	if err := t.cmd.Start(); err != nil {
		return nil, err
	}
	return NewClient(stdout, stdin), nil
}

func (t *STDIOTransport) Close() error {
	if t.cmd.Process != nil {
		return t.cmd.Process.Kill()
	}
	return nil
}
