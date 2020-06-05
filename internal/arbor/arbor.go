package arbor

import (
	"fmt"
	"strings"
)

type status int

const (
	pending = status(iota)
	failed
	succeeded
	skipped
)

type test struct {
	name string
	deps []*test
	run  func() error

	status     status
	failReason string
}

func New(name string, f func() error, deps ...*test) *test {
	return &test{
		name:   name,
		run:    f,
		status: pending,
		deps:   deps,
	}
}

func (ts *test) Run() {
	for _, dep := range ts.deps {
		switch dep.status {
		case succeeded:
			continue
		case failed:
			return
		case pending:
			fallthrough
		default:
			if dep.status == pending {
				dep.Run()
				if dep.status != succeeded {
					return
				}
			}
		}
	}

	err := ts.run()

	ts.status = succeeded

	if err != nil {
		ts.status = failed
		ts.failReason = err.Error()
	}
}

func (ts *test) String() string {
	var buffer strings.Builder

	switch ts.status {
	case succeeded:
		buffer.WriteString(fmt.Sprintf("[%v] ok\n", ts.name))
		break
	case failed:
		buffer.WriteString(fmt.Sprintf("[%v] failed: %v\n", ts.name, ts.failReason))
		break
	case pending:
		fallthrough
	default:
		buffer.WriteString(fmt.Sprintf("[%v] not ran\n", ts.name))
	}

	for _, t := range ts.deps {
		buffer.WriteString(fmt.Sprintf("\u2BA1 %s", t.String()))
	}

	return buffer.String()
}
