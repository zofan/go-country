package country

type Country struct {
	Alpha2  string
	Alpha3  string
	Numeric string

	Name       string
	NativeName string
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

func ByAlpha3(v string) *Country {
	for _, c := range List {
		if c.Alpha3 == v {
			return &c
		}
	}

	return nil
}

func ByAlpha2(v string) *Country {
	for _, c := range List {
		if c.Alpha2 == v {
			return &c
		}
	}

	return nil
}

func ByNumeric(v string) *Country {
	for _, c := range List {
		if c.Numeric == v {
			return &c
		}
	}

	return nil
}
