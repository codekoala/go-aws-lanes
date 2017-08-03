package lanes

import (
	"fmt"
)

type Server struct {
	ID   string
	Name string
	Lane string
	IP   string
}

func (this *Server) String() string {
	return fmt.Sprintf("%s (%s)", this.Name, this.ID)
}
