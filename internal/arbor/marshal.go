package arbor

import "encoding/json"

/*

{
  "nodes": [
    {"id": "Myriel", "group": 1}
  ],
  "links": [
	{"source": "Napoleon", "target": "Myriel", "value": 1},
  ]
}

*/

type node struct {
	ID     string `json:"id"`
	Group  int    `json:"group"`
	Status string `json:"status"`
}

type link struct {
	Source string `json:"source"`
	Target string `json:"target"`
	Value  int    `json:"value"`
}

type output struct {
	Nodes []node `json:"nodes"`
	Links []link `json:"links"`
}

func Marshal(tests ...*test) string {
	out := output{
		Nodes: make([]node, 0),
		Links: make([]link, 0),
	}

	for _, t := range tests {
		marshal(t, &out)
	}

	bytes, err := json.Marshal(out)

	if err != nil {
		panic(err)
	}

	return string(bytes)
}

var statuses = []string{"pending", "fail", "pass"}

func marshal(ts *test, out *output) {
	for _, t := range ts.deps {
		marshal(t, out)
	}

	for _, d := range ts.deps {
		l := link{
			Source: ts.name,
			Target: d.name,
			Value:  3,
		}

		out.Links = append(out.Links, l)
	}

	if contains(out.Nodes, ts.name) {
		return
	}

	n := node{
		ID:     ts.name,
		Group:  int(ts.Status),
		Status: statuses[ts.Status],
	}

	out.Nodes = append(out.Nodes, n)
}

func contains(ns []node, name string) bool {
	for i := 0; i < len(ns); i++ {
		if ns[i].ID == name {
			return true
		}
	}

	return false
}
