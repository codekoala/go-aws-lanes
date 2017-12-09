package session

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/codekoala/go-aws-lanes"
	"github.com/codekoala/go-aws-lanes/ssh"
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

func (this *Session) Profile() *ssh.Profile {
	return this.svr.Profile().NoTunnels()
}

func (this *Session) GetName() string {
	return this.svr.Name
}

func (this *Session) Run(ctx context.Context, profile *ssh.Profile, args ...string) error {
	this.Buffer.Reset()

	this.cmd = this.svr.GetSSHCommandWithProfile(ctx, profile, args)
	this.cmd.Stdin = nil
	this.cmd.Stderr = this.Buffer
	this.cmd.Stdout = this.Buffer

	this.WriteString(fmt.Sprintf("Executing on %s (%s): %s\n", this.svr.Name, this.svr.IP, strings.Join(args, " ")))

	return this.cmd.Start()
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
		err = this.Wait()
	}

	return err
}

func (this *Session) Wait() error {
	return this.cmd.Wait()
}
