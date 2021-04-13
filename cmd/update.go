package main

import (
	"fmt"
	"github.com/zofan/go-country"
	"github.com/zofan/go-fwrite"
	"github.com/zofan/go-req"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"
)

func main() {
	fmt.Println(Update())
}

// https://www.globalfirepower.com/countries-listing.asp
// https://www.numbeo.com/quality-of-life/rankings_current.jsp
// https://en.wikipedia.org/wiki/List_of_next_general_elections

func Update() error {
	var (
		httpClient = req.New(req.DefaultConfig)
		list       = make(map[string]*country.Country)
	)

	nameRe := regexp.MustCompile(` \([^)]+\)`)

	var tmp = []struct {
		Alpha2Code     string    `json:"alpha2Code"`
		Alpha3Code     string    `json:"alpha3Code"`
		Area           float64   `json:"area"`
		Borders        []string  `json:"borders"`
		CallingCodes   []string  `json:"callingCodes"`
		Capital        string    `json:"capital"`
		Flag           string    `json:"flag"`
		Latlng         []float64 `json:"latlng"`
		Name           string    `json:"name"`
		NativeName     string    `json:"nativeName"`
		NumericCode    string    `json:"numericCode"`
		Population     float64   `json:"population"`
		Region         string    `json:"region"`
		Subregion      string    `json:"subregion"`
		Timezones      []string  `json:"timezones"`
		TopLevelDomain []string  `json:"topLevelDomain"`

		Languages []struct {
			Iso639_1   string `json:"iso639_1"`
			Iso639_2   string `json:"iso639_2"`
			Name       string `json:"name"`
			NativeName string `json:"nativeName"`
		} `json:"languages"`

		Currencies []struct {
			Code   string `json:"code"`
			Name   string `json:"name"`
			Symbol string `json:"symbol"`
		} `json:"currencies"`
	}{}

	resp := httpClient.Get(`https://restcountries.eu/rest/v2/all`)
	if resp.Error() != nil {
		return resp.Error()
	}

	err := resp.ReadJSON(&tmp)
	if err != nil {
		return err
	}

	for _, tc := range tmp {
		if len(tc.Latlng) == 0 {
			continue
		}

		c := &country.Country{
			Alpha2:  strings.TrimSpace(tc.Alpha2Code),
			Alpha3:  strings.TrimSpace(tc.Alpha3Code),
			Numeric: strings.TrimSpace(tc.NumericCode),

			Name:       strings.TrimSpace(nameRe.ReplaceAllString(tc.Name, ` `)),
			NativeName: tc.NativeName,
			FlagURL:    `https://www.countryflags.io/` + tc.Alpha2Code + `/shiny/64.png`,

			Area:       tc.Area,
			Population: tc.Population,
			Latitude:   tc.Latlng[0],
			Longitude:  tc.Latlng[1],

			Region:    tc.Region,
			SubRegion: tc.Subregion,
			Capital:   tc.Capital,

			Callings: tc.CallingCodes,
			Borders:  tc.Borders,
			TLDs:     tc.TopLevelDomain,
		}

		for _, v := range tc.Languages {
			c.Languages = append(c.Languages, strings.ToUpper(v.Iso639_2))
		}

		for _, v := range tc.Currencies {
			if strings.Contains(v.Code, `(`) {
				continue
			}
			c.Currencies = append(c.Currencies, strings.ToUpper(v.Code))
		}

		for _, v := range tc.Timezones {
			c.TimeZones = append(c.TimeZones, v)
		}

		list[c.Alpha3] = c
	}

	// ---

	updateTags(list)

	var tpl []string

	tpl = append(tpl, `package country`)
	tpl = append(tpl, ``)
	tpl = append(tpl, `// Updated at: `+time.Now().String())
	tpl = append(tpl, `var List = []Country{`)

	for _, c := range list {
		s := fmt.Sprintf(`%#v`, *c) + `,`
		s = strings.ReplaceAll(s, `country.Country`, ``)
		tpl = append(tpl, s)
	}

	tpl = append(tpl, `}`)
	tpl = append(tpl, ``)

	_, file, _, _ := runtime.Caller(0)
	dir := filepath.Dir(file)

	return fwrite.WriteRaw(dir+`/../db.go`, []byte(strings.Join(tpl, "\n")))
}

func updateTags(list map[string]*country.Country) {
	wordSplitRe := regexp.MustCompile(`[^\p{L}\p{N}]+`)
	wordMap := map[string][]*country.Country{}

	for _, c := range list {
		name := strings.ToLower(c.Name + ` ` + c.NativeName + ` ` + strings.Join(c.AltNames, ` `))
		words := wordSplitRe.Split(name, -1)
		for _, w := range words {
			if len(w) > 0 {
				wordMap[w] = append(wordMap[w], c)
			}
		}
		c.Tags = []string{}
	}

	for w, cs := range wordMap {
		if len(cs) == 1 {
			cs[0].Tags = append(cs[0].Tags, w)
		}
	}
}
