package arbor

import (
	"fmt"
)

type status int

const (
	pending = status(iota)
	fail
	pass
)

type test struct {
	name string
	deps []*test
	run  func() error

	status     status
	failReason string

	Success bool
}

func New(name string, f func() error, deps ...*test) *test {
	return &test{
		name:   name,
		run:    f,
		status: pending,
		deps:   deps,
	}
}

func Suite(name string, deps ...*test) *test {
	return &test{
		name:   name,
		run:    func() error { return nil },
		status: pending,
		deps:   deps,
	}
}

func (ts *test) Run() {
	for _, dep := range ts.deps {
		switch dep.status {
		case pass:
			continue
		case fail:
			return
		case pending:
			fallthrough
		default:
			if dep.status == pending {
				dep.Run()
				if dep.status != pass {
					return
				}
			}
		}
	}

	err := ts.run()

	ts.status = pass
	ts.Success = true

	if err != nil {
		ts.status = fail
		ts.Success = false
		ts.failReason = err.Error()
	}
}

func (ts *test) Lines() (buffer []string) {
	var curr string

	switch ts.status {
	case pass:
		curr = fmt.Sprintf("\u2BA1[%v] passed\n", ts.name)
		break
	case fail:
		curr = fmt.Sprintf("\u2BA1[%v] failed: %v\n", ts.name, ts.failReason)
		break
	case pending:
		fallthrough
	default:
		curr = fmt.Sprintf("\u2BA1[%v] skipped\n", ts.name)
	}

	buffer = append(buffer, curr)

	for _, t := range ts.deps {
		for _, line := range t.Lines() {
			buffer = append(buffer, fmt.Sprintf("  %s", line))
		}
	}

	return buffer
}

func (ts *test) String() (out string) {
	lines := ts.Lines()

	out = lines[0][len("⮡"):len(lines[0])] //drop first ⮡

	for i := 1; i < len(lines); i++ {
		val := lines[i]
		out += val
	}

	return
}
