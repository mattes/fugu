package data

type Data struct {
	data map[string][]string
}

func New() *Data {
	return &Data{
		data: make(map[string][]string),
	}
}

// Keys returns all data names
// Don't rely on it's order!
func (d *Data) Keys() []string {
	keys := []string{}
	for k := range d.data {
		keys = append(keys, k)
	}
	return keys
}

func (d *Data) Exists(name string) bool {
	if d.data == nil {
		return false
	}
	_, ok := d.data[name]
	return ok
}

func (d *Data) GetAll(name string) []string {
	if d.Exists(name) {
		return d.data[name]
	}
	return []string{}
}

func (d *Data) Get(name string) string {
	if d.Exists(name) {
		if len(d.data[name]) == 1 {
			return d.data[name][0]
		}
	}
	return ""
}

func (d *Data) IsTrue(name string) bool {
	return d.Get(name) == "true"
}

func (d *Data) IsFalse(name string) bool {
	return d.Get(name) == "false"
}

// PickAll gets element and deletes it afterwards
func (d *Data) PickAll(name string) []string {
	defer d.Delete(name)
	return d.GetAll(name)
}

// Pick gets element and deletes it afterwards
func (d *Data) Pick(name string) string {
	defer d.Delete(name)
	return d.Get(name)
}

// Set sets (overwrites) values for name
func (d *Data) Set(name string, value ...string) *Data {
	d.data[name] = value
	return d
}

// Add adds values for name
func (d *Data) Add(name string, value ...string) *Data {
	for _, v := range value {
		d.data[name] = append(d.data[name], v)
	}
	return d
}

func (d *Data) SetTrue(name string) *Data {
	d.Set(name, "true")
	return d
}

func (d *Data) SetFalse(name string) *Data {
	d.Set(name, "false")
	return d
}

func (d *Data) Delete(name string) {
	if d.Exists(name) {
		delete(d.data, name)
	}
}

func (d *Data) Raw() map[string][]string {
	return d.data
}

func (d *Data) RawEnhanced() map[string]interface{} {
	n := make(map[string]interface{})
	for k, v := range d.data {
		if len(v) == 1 {
			switch v[0] {
			case "true":
				n[k] = true
			case "false":
				n[k] = false
			default:
				n[k] = v[0]
			}
		} else {
			n[k] = v
		}
	}
	return n
}

// Merge merges p2 into p.data.
// Later values overwrite previous ones.
func (d *Data) Merge(p2 ...*Data) {
	for _, pp := range p2 {
		if pp != nil && pp.data != nil {
			for k, v := range pp.data {
				d.Set(k, v...)
			}
		}
	}
}

// Filter filters f from d.data
func (d *Data) Filter(f interface{}) {
	for k := range d.data {
		keep := false

		switch f.(type) {
		case string:
			if k == f.(string) {
				keep = true
			}

		case []string:
			for _, k2 := range f.([]string) {
				if k2 == k {
					keep = true
					break
				}
			}
		case *Data:
			for _, k2 := range f.(*Data).Keys() {
				if k2 == k {
					keep = true
					break
				}
			}

		default:
			panic("unsupported type")
		}

		if !keep {
			d.Delete(k)
		}
	}
}

func Filter(d *Data, f interface{}) *Data {
	dn := New()
	dn.Merge(d)
	dn.Filter(f)
	return dn
}

// Merge merges data objects.
// Later values overwrite previous ones.
func Merge(data ...*Data) *Data {
	newData := New()
	for _, p := range data {
		if p != nil && p.data != nil {
			for k, v := range p.data {
				newData.Set(k, v...)
			}
		}
	}
	return newData
}

func ToData(p map[string][]string) *Data {
	np := New()
	for k, v := range p {
		np.Set(k, v...)
	}
	return np
}
