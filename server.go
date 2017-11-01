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
	ID    string
	Name  string
	Lane  string
	IP    string
	State string

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

func DisplayServersCols(servers []*Server, columns ColumnSet) error {
	return DisplayServersColsWriter(os.Stdout, servers, columns)
}

func DisplayServersWriter(writer io.Writer, servers []*Server) (err error) {
	return DisplayServersColsWriter(writer, servers, DefaultColumnSet)
}

func DisplayServersColsWriter(writer io.Writer, servers []*Server, columns ColumnSet) (err error) {
	if len(servers) == 0 {
		return fmt.Errorf("No servers found.")
	}

	if len(columns) == 0 {
		return nil
	}

	if !config.DisableUTF8 {
		termtables.EnableUTF8()
	}

	table := termtables.CreateTable()
	if config.Table.HideBorders {
		table.Style.SkipBorder = true
		table.Style.BorderY = ""
		table.Style.PaddingLeft = 0
	}

	if !config.Table.HideTitle {
		table.AddTitle("AWS Servers")
	}

	for idx, svr := range servers {
		row := table.AddRow()

		for _, col := range columns {
			switch col {
			case ColumnIndex:
				row.AddCell(idx + 1)
			case ColumnLane:
				row.AddCell(svr.Lane)
			case ColumnServer:
				row.AddCell(svr.Name)
			case ColumnIP:
				row.AddCell(svr.IP)
			case ColumnState:
				row.AddCell(svr.State)
			case ColumnID:
				row.AddCell(svr.ID)
			default:
				continue
			}
		}
	}

	if !config.Table.HideHeaders {
		// add headers after all rows because cell alignment only applies to cells that exist when SetAlign is called
		for idx, col := range columns {
			table.AddHeaders(col)

			switch col {
			case ColumnIndex:
				table.SetAlign(termtables.AlignRight, idx+1)
			}
		}
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
	return FetchServersBy(svc, nil, "")
}

func FetchServersInLane(svc *ec2.EC2, lane string) ([]*Server, error) {
	return FetchServersInLaneByKeyword(svc, lane, "")
}

func FetchServersInLaneByKeyword(svc *ec2.EC2, lane, keyword string) ([]*Server, error) {
	return FetchServersBy(svc, CreateLaneFilter(lane), keyword)
}

func FetchServersBy(svc *ec2.EC2, input *ec2.DescribeInstancesInput, keyword string) (servers []*Server, err error) {
	var out *ec2.DescribeInstancesOutput

	fmt.Fprintf(os.Stderr, "Fetching servers... ")
	spin := spinner.New(spinner.CharSets[21], 50*time.Millisecond)
	spin.Writer = os.Stderr
	spin.Start()
	defer func() {
		spin.Stop()
		fmt.Fprintln(os.Stderr, "done")
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
