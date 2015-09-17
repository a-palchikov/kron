package proc

import (
	"io"
	"os/exec"

	pb "github.com/a-palchikov/kron/proto/servicepb"
)

type Cmd struct {
	*exec.Cmd
	quota *pb.Quota
}

func (cmd *Cmd) Run() error {
	return nil
}

func (cmd *Cmd) Start() error {
	return nil
}

func (cmd *Cmd) StdoutPipe() (io.ReadCloser, error) {
	return nil, nil
}

func (cmd *Cmd) StderrPipe() (io.ReadCloser, error) {
	return nil, nil
}

func (cmd *Cmd) Wait() error {
	return nil
}

func Command(name string, quota *pb.Quota, args ...string) *Cmd {
	return &Cmd{
		quota: quota,
	}
}
