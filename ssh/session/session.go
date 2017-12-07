package session

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/codekoala/go-aws-lanes"
)

type Session struct {
	*bytes.Buffer

	svr *lanes.Server
	cmd *exec.Cmd
	ctx context.Context
}

func New(svr *lanes.Server) *Session {
	return &Session{
		Buffer: bytes.NewBuffer(nil),
		svr:    svr,
	}
}

func (this *Session) GetName() string {
	return this.svr.Name
}

func (this *Session) Run(ctx context.Context, args ...string) (err error) {
	this.Buffer.Reset()

	this.cmd = this.svr.GetSSHCommand(ctx, args)
	this.cmd.Stdin = nil
	this.cmd.Stderr = this.Buffer
	this.cmd.Stdout = this.Buffer

	this.WriteString(fmt.Sprintf("Executing on %s (%s): %s\n", this.svr.Name, this.svr.IP, args))

	if err = this.cmd.Start(); err == nil {
		go this.cmd.Wait()
	}

	return err
}

func (this *Session) GetLastLine() string {
	line := strings.TrimSpace(this.String())
	if pos := strings.LastIndex(line, "\n"); pos != -1 {
		line = line[pos:]
	}

	return strings.TrimSpace(line)
}

func (this *Session) IsClosed() bool {
	if this.cmd == nil || this.cmd.ProcessState == nil {
		return false
	}

	return this.cmd.ProcessState.Exited()
}

func (this *Session) Close() (err error) {
	if this.cmd.Process != nil {
		this.cmd.Process.Signal(os.Interrupt)
		err = this.cmd.Wait()
	}

	return err
}
