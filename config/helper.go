package config

func Get(config []Value, name string) Value {
	for _, c := range config {
		for _, n := range c.Names() {
			if n == name {
				return c
			}
		}
	}
	return nil
}

func Set(config *[]Value, name string, value interface{}) bool {
	for _, c := range *config {
		for _, n := range c.Names() {
			if n == name {
				c.Set(value)
				return true
			}
		}
	}
	return false
}
