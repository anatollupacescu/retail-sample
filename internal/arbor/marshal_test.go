package arbor_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	arbor "github.com/anatollupacescu/retail-sample/internal/arbor"
)

func noOp() error { return nil }

func TestMarshalOneChild(t *testing.T) {
	dep := arbor.New("dep", noOp)
	tst := arbor.New("test", noOp, dep)
	str := arbor.Marshal(tst)
	assert.Equal(t, `{"nodes":[{"id":"dep","group":0,"status":"pending"},{"id":"test","group":0,"status":"pending"}],"links":[{"source":"test","target":"dep","value":3}]}`, str)
}

func TestMarshalDiamond(t *testing.T) {
	tail := arbor.New("tail", noOp)
	left := arbor.New("left", noOp, tail)
	right := arbor.New("right", noOp, tail)
	head := arbor.New("head", noOp, left, right)
	str := arbor.Marshal(head)
	t.Logf("got str:\n %v\n", str)
	assert.Equal(t, `{"nodes":[{"id":"tail","group":0,"status":"pending"},{"id":"left","group":0,"status":"pending"},{"id":"right","group":0,"status":"pending"},{"id":"head","group":0,"status":"pending"}],"links":[{"source":"left","target":"tail","value":3},{"source":"right","target":"tail","value":3},{"source":"head","target":"left","value":3},{"source":"head","target":"right","value":3}]}`, str)
}
