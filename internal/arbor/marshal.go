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

func Marshal(ts ...*test) string {
	var nodes = make([]node, len(ts))

	statuses := []string{"pending", "fail", "pass"}
	for i, t := range ts {
		n := node{
			ID:     t.name,
			Group:  int(t.status),
			Status: statuses[t.status],
		}

		nodes[i] = n
	}

	var links = make([]link, 0)

	for _, t := range ts {
		for _, d := range t.deps {
			l := link{
				Source: t.name,
				Target: d.name,
				Value:  3,
			}

			links = append(links, l)
		}
	}

	out := output{
		Nodes: nodes,
		Links: links,
	}

	bytes, err := json.Marshal(out)

	if err != nil {
		panic(err)
	}

	return string(bytes)
}
