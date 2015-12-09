package common

type Collectables struct {
	Raw        map[string]uint64
	Calculated map[string]float64
}

func (c *Collectables) AddRaw(name string, value uint64) error {
	c.Raw[name] = value
	return nil
}

func (c *Collectables) AddCalculated(name string, value float64) error {
	c.Calculated[name] = value
	return nil
}
