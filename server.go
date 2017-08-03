package lanes

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/apcera/termtables"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/briandowns/spinner"
)

type Server struct {
	ID   string
	Name string
	Lane string
	IP   string
}

func (this *Server) SortKey() string {
	return fmt.Sprintf("%s %s %s", this.Lane, this.Name, this.ID)
}

func (this *Server) String() string {
	return fmt.Sprintf("%s (%s)", this.Name, this.ID)
}

func DisplayServers(servers []*Server) {
	termtables.EnableUTF8()
	table := termtables.CreateTable()
	table.AddTitle("AWS Servers")
	table.AddHeaders("IDX", "LANE", "SERVER", "IP ADDRESS", "ID")

	for idx, svr := range servers {
		table.AddRow(idx+1, svr.Lane, svr.Name, svr.IP, svr.ID)
	}

	fmt.Printf(table.Render())
}

func FetchServers(svc *ec2.EC2) ([]*Server, error) {
	return FetchServersBy(svc, nil)
}

func FetchServersInLane(svc *ec2.EC2, lane string) ([]*Server, error) {
	var input *ec2.DescribeInstancesInput

	if lane != "" {
		input = &ec2.DescribeInstancesInput{
			Filters: []*ec2.Filter{{
				Name:   aws.String("tag-key"),
				Values: []*string{aws.String("Lane")},
			}, {
				Name:   aws.String("tag-value"),
				Values: []*string{aws.String(lane)},
			}},
		}
	}

	return FetchServersBy(svc, input)
}

func FetchServersBy(svc *ec2.EC2, input *ec2.DescribeInstancesInput) (servers []*Server, err error) {
	var out *ec2.DescribeInstancesOutput

	fmt.Printf("Fetching servers... ")
	defer fmt.Println("done")
	spin := spinner.New(spinner.CharSets[21], 100*time.Millisecond)
	spin.Start()
	defer spin.Stop()

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

				switch strings.ToLower(*tag.Key) {
				case "name":
					svr.Name = *tag.Value
				case "lane":
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
