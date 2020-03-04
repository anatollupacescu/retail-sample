package retailsampleapp1

import "time"

type ( //log
	ProvisionEntry struct {
		Time time.Time
		ID   int
		Qty  int
	}

	ProvisionLog interface {
		Add(ProvisionEntry)
		List() []ProvisionEntry
	}
)
