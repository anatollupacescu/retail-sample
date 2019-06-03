package domain

type Design struct {
	designID int
	Name     string
	Price    float32
}

func (d Design) Validate() bool {
	if len(d.Name) < 5 {
		return false
	}
	if d.Price == 0 {
		return false
	}
	return true
}
