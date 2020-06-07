package arbor

import (
	"net/http"
	"strings"
)

func Upload(url, data string) {
	if url == "" {
		return
	}

	req, err := http.NewRequest("POST", url, strings.NewReader(data))

	if err != nil {
		panic(err)
	}

	hc := http.Client{}

	r, err := hc.Do(req)

	if err != nil {
		panic(err)
	}

	if r.StatusCode != 200 {
		panic("bad response status: " + r.Status)
	}
}
