package lanes

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"

	"github.com/apcera/termtables"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/briandowns/spinner"

	"github.com/codekoala/go-aws-lanes/ssh"
)

type Server struct {
	ID   string
	Name string
	Lane string
	IP   string

	profile *ssh.Profile
}

// Login attempts to SSH into the server using the default profile.
func (this *Server) Login(args []string) error {
	if this.profile == nil {
		return ErrMissingSSHProfile
	}

	return this.LoginWithProfile(this.profile, args)
}

// LoginWithProfile attempts to SSH into the server using a custom profile.
func (this *Server) LoginWithProfile(profile *ssh.Profile, args []string) error {
	sshArgs := append(profile.SSHArgs(this.IP), args...)

	cmd := exec.Command("ssh", sshArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return cmd.Run()
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

func DisplayServers(servers []*Server) error {
	return DisplayServersWriter(os.Stdout, servers)
}

func DisplayServersWriter(writer io.Writer, servers []*Server) (err error) {
	if len(servers) == 0 {
		return fmt.Errorf("No servers found.")
	}

	if !config.DisableUTF8 {
		termtables.EnableUTF8()
	}

	table := termtables.CreateTable()
	table.AddTitle("AWS Servers")
	table.AddHeaders("IDX", "LANE", "SERVER", "IP ADDRESS", "ID")

	for idx, svr := range servers {
		table.AddRow(idx+1, svr.Lane, svr.Name, svr.IP, svr.ID)
	}

	fmt.Fprintf(writer, table.Render())

	return nil
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
	return FetchServersBy(svc, nil)
}

func FetchServersInLane(svc *ec2.EC2, lane string) ([]*Server, error) {
	return FetchServersBy(svc, CreateLaneFilter(lane))
}

func FetchServersBy(svc *ec2.EC2, input *ec2.DescribeInstancesInput) (servers []*Server, err error) {
	var out *ec2.DescribeInstancesOutput

	fmt.Printf("Fetching servers... ")
	spin := spinner.New(spinner.CharSets[21], 50*time.Millisecond)
	spin.Start()
	defer func() {
		spin.Stop()
		fmt.Println("done")
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
				ID: *inst.InstanceId,
				IP: *inst.PublicIpAddress,
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

			servers = append(servers, svr)
		}
	}

	sort.Slice(servers, func(i, j int) bool {
		return servers[i].SortKey() < servers[j].SortKey()
	})

	return servers, nil
}
