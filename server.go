package lanes

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/briandowns/spinner"

	"github.com/codekoala/go-aws-lanes/ssh"
)

type Server struct {
	ID    string
	Name  string
	Lane  string
	IP    string
	State string

	profile *ssh.Profile
}

// Profile returns the SSH profile used to access this server.
func (this *Server) Profile() *ssh.Profile {
	return this.profile
}

// Login attempts to SSH into the server using the default profile.
func (this *Server) Login(ctx context.Context, args []string) error {
	if this.profile == nil {
		return ErrMissingSSHProfile
	}

	return this.LoginWithProfile(ctx, this.profile, args)
}

// LoginWithProfile attempts to SSH into the server using a custom profile.
func (this *Server) LoginWithProfile(ctx context.Context, profile *ssh.Profile, args []string) error {
	return this.GetSSHCommandWithProfile(ctx, profile, args).Run()
}

func (this *Server) GetSSHCommand(ctx context.Context, args []string) *exec.Cmd {
	return this.GetSSHCommandWithProfile(ctx, this.profile, args)
}

// GetSSHCommandWithProfile sets up a command to SSH into the server using a custom profile.
func (this *Server) GetSSHCommandWithProfile(ctx context.Context, profile *ssh.Profile, args []string) *exec.Cmd {
	sshArgs := append(profile.SSHArgs(this.IP), args...)

	cmd := exec.CommandContext(ctx, "ssh", sshArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return cmd
}

// Push attempts to copy files to the server using the default profile.
func (this *Server) Push(dest string, sources ...string) error {
	if this.profile == nil {
		return ErrMissingSSHProfile
	}

	return this.PushWithProfile(this.profile, dest, sources...)
}

// PushWithProfile attempts to copy files to the server using a custom profile.
func (this *Server) PushWithProfile(profile *ssh.Profile, dest string, sources ...string) error {
	scpArgs := []string{"-i", profile.Identity, "-r"}
	scpArgs = append(scpArgs, sources...)
	scpArgs = append(scpArgs, fmt.Sprintf("%s:%s", profile.UserAt(this.IP), dest))

	//fmt.Printf("scp %s\n", strings.Join(scpArgs, " "))

	cmd := exec.Command("scp", scpArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return cmd.Run()
}

func (this *Server) SortKey() string {
	return strings.ToLower(fmt.Sprintf("%s %s %s", this.Lane, this.Name, this.ID))
}

func (this *Server) String() string {
	return fmt.Sprintf("%s (%s)", this.Name, this.ID)
}

func CreateLaneFilter(lane string) (input *ec2.DescribeInstancesInput) {
	if lane != "" {
		input = &ec2.DescribeInstancesInput{
			Filters: []*ec2.Filter{{
				Name:   aws.String("tag-key"),
				Values: []*string{aws.String(config.Tags.Lane)},
			}, {
				Name:   aws.String("tag-value"),
				Values: []*string{aws.String(lane)},
			}},
		}
	}

	return
}

func FetchServers(svc *ec2.EC2) ([]*Server, error) {
	return FetchServersBy(svc, nil, "")
}

func FetchServersInLane(svc *ec2.EC2, lane string) ([]*Server, error) {
	return FetchServersInLaneByKeyword(svc, lane, "")
}

func FetchServersInLaneByKeyword(svc *ec2.EC2, lane, keyword string) ([]*Server, error) {
	return FetchServersBy(svc, CreateLaneFilter(lane), keyword)
}

func FetchServersBy(svc *ec2.EC2, input *ec2.DescribeInstancesInput, keyword string) (servers []*Server, err error) {
	var (
		out    *ec2.DescribeInstancesOutput
		buf    = &bytes.Buffer{}
		exists bool
	)

	fmt.Fprintf(os.Stderr, "Fetching servers... ")
	spin := spinner.New(spinner.CharSets[21], 50*time.Millisecond)
	spin.Writer = os.Stderr
	spin.Start()

	defer func() {
		spin.Stop()
		fmt.Fprintln(os.Stderr, "done")
		fmt.Fprintf(os.Stderr, buf.String())
	}()

	if out, err = svc.DescribeInstances(input); err != nil {
		return
	}

	for _, rez := range out.Reservations {
		for _, inst := range rez.Instances {
			if inst.PublicIpAddress == nil || *inst.PublicIpAddress == "" {
				continue
			}

			svr := &Server{
				ID:    *inst.InstanceId,
				IP:    *inst.PublicIpAddress,
				State: *inst.State.Name,
			}

			for _, tag := range inst.Tags {
				if tag == nil {
					continue
				}

				switch *tag.Key {
				case config.Tags.Name:
					svr.Name = *tag.Value
				case config.Tags.Lane:
					svr.Lane = *tag.Value
				}
			}

			// filter servers by keyword
			if !svr.Matches(keyword) {
				continue
			}

			if config.profile != nil {
				// assign appropriate profile to server
				if svr.profile, exists = config.profile.SSH.Mods[svr.Lane]; !exists {
					fmt.Fprintf(buf, "WARNING: no profile found for %s in lane %q\n", svr, svr.Lane)
					svr.profile = config.profile.SSH.Default
					if svr.profile == nil {
						svr.profile = &ssh.DefaultProfile
					}
				}
			}

			servers = append(servers, svr)
		}
	}

	sort.Slice(servers, func(i, j int) bool {
		return servers[i].SortKey() < servers[j].SortKey()
	})

	return servers, nil
}

func (this *Server) Matches(keyword string) bool {
	if keyword == "" {
		return true
	}

	keyword = strings.ToUpper(keyword)

	return (strings.Contains(strings.ToUpper(this.Name), keyword) ||
		strings.Contains(strings.ToUpper(this.ID), keyword) ||
		strings.Contains(strings.ToUpper(this.Lane), keyword) ||
		strings.Contains(this.IP, keyword))
}
