package country

import "strings"

type Country struct {
	Alpha2  string
	Alpha3  string
	Numeric string

	Name       string
	NativeName string
	AltNames   []string
	Tags       []string
	FlagURL    string

	Area       float64
	Population float64
	Latitude   float64
	Longitude  float64

	Region    string
	SubRegion string
	Capital   string

	Callings   []string
	Borders    []string
	TimeZones  []string
	Languages  []string
	Currencies []string
	TLDs       []string
}

func Get(v string) *Country {
	for _, c := range List {
		if c.Alpha3 == v || c.Alpha2 == v || c.Numeric == v {
			return &c
		}
	}

	return nil
}

func ByName(v string) *Country {
	fn := func(v string) string {
		v = strings.ReplaceAll(strings.ToLower(v), ` `, ``)
		v = strings.Replace(v, `St.`, ``, 1)
		v = strings.Replace(v, `g`, `q`, 1)
		return v
	}

	v = fn(v)

	for _, c := range List {
		n := fn(c.Name)
		n2 := fn(c.NativeName)

		if strings.Contains(n, v) || strings.Contains(n2, v) || strings.Contains(v, n) || strings.Contains(v, n2) {
			return &c
		}

		for _, n := range c.AltNames {
			n = fn(n)

			if strings.Contains(n, v) || strings.Contains(v, n) {
				return &c
			}
		}

		for _, n := range c.Tags {
			n = fn(n)

			if strings.Contains(n, v) || strings.Contains(v, n) {
				return &c
			}
		}
	}

	return nil
}
