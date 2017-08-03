package lanes_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/codekoala/go-aws-lanes"
)

func TestServerSortKey(t *testing.T) {
	s := &lanes.Server{
		ID:   "z",
		Name: "r",
		Lane: "u",
		IP:   "5",
	}
	assert.Equal(t, s.SortKey(), "u r z")
	assert.Equal(t, s.String(), "r (z)")
}

func TestDisplayServersWriter(t *testing.T) {
	var (
		ss  []*lanes.Server
		buf = bytes.NewBuffer(nil)
		s   = &lanes.Server{
			ID:   "i-id-i",
			Name: "n-name-n",
			Lane: "l-lane-l",
			IP:   "i-5-i",
		}
	)

	assert.NotNil(t, lanes.DisplayServersWriter(buf, ss))
	assert.Equal(t, buf.String(), "")

	ss = append(ss, s)
	assert.Nil(t, lanes.DisplayServersWriter(buf, ss))

	tbl := buf.String()
	assert.True(t, strings.Contains(tbl, s.ID))
	assert.True(t, strings.Contains(tbl, s.Name))
	assert.True(t, strings.Contains(tbl, s.Lane))
	assert.True(t, strings.Contains(tbl, s.IP))

	idIdx := strings.Index(tbl, s.ID)
	nameIdx := strings.Index(tbl, s.Name)
	laneIdx := strings.Index(tbl, s.Lane)
	ipIdx := strings.Index(tbl, s.IP)

	assert.True(t, laneIdx < nameIdx)
	assert.True(t, nameIdx < ipIdx)
	assert.True(t, ipIdx < idIdx)
}
