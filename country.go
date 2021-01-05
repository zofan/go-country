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
