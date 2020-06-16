package arbor

import (
	"fmt"
)

type status int

const (
	Pending = status(iota)
	Fail
	Pass
)

type test struct {
	name    string
	deps    []*test
	runFunc func() error

	Status     status
	FailReason string

	Success bool
}

func New(name string, f func() error, deps ...*test) *test {
	return &test{
		name:    name,
		runFunc: f,
		Status:  Pending,
		deps:    deps,
	}
}

func Alias(name string, aliasFor *test) *test {
	return &test{
		name:    name,
		runFunc: func() error { return nil },
		Status:  Pending,
		deps:    []*test{aliasFor},
	}
}

func Run(tt ...*test) (all []*test, success bool) {
	success = true
	for _, v := range tt {
		v.run()
		if v.Status != Pass {
			success = false
		}
	}

	return tt, success
}

func (ts *test) run() {
	for _, dep := range ts.deps {
		switch dep.Status {
		case Pass:
			continue
		case Fail:
			return
		case Pending:
			fallthrough
		default:
			if dep.Status == Pending {
				dep.run()
				if dep.Status != Pass {
					return
				}
			}
		}
	}

	err := ts.runFunc()

	ts.Status = Pass
	ts.Success = true

	if err != nil {
		ts.Status = Fail
		ts.Success = false
		ts.FailReason = err.Error()
	}
}

func (ts *test) Lines() (buffer []string) {
	var curr string

	switch ts.Status {
	case Pass:
		curr = fmt.Sprintf("\u2BA1[%v] passed\n", ts.name)
	case Fail:
		curr = fmt.Sprintf("\u2BA1[%v] failed: %v\n", ts.name, ts.FailReason)
	case Pending:
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
