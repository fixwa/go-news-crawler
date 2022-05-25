package crawlers

import (
	"testing"
)

func TestInfobaeStoresArticle(t *testing.T) {
	// @todo Mock requests in order to test the collector
}

func TestConvertInfobaeDate(t *testing.T) {
	cases := map[string]map[string]string{
		"Jan": map[string]string{
			"givenDate": "1 de Enero de 1999",
			"expected":  "1999-01-01",
		},
		"Feb": map[string]string{
			"givenDate": "12 de Febrero de 2000",
			"expected":  "2000-02-12",
		},
		"Mar": map[string]string{
			"givenDate": "13 de Marzo de 2001",
			"expected":  "2001-03-13",
		},
		"Abr": map[string]string{
			"givenDate": "22 de Abril de 1980",
			"expected":  "1980-04-22",
		},
		"May": map[string]string{
			"givenDate": "10 de Mayo de 2009",
			"expected":  "2009-05-10",
		},
		"Jun": map[string]string{
			"givenDate": "19 de Junio de 2022",
			"expected":  "2022-06-19",
		},
		"Jul": map[string]string{
			"givenDate": "20 de Julio de 2024",
			"expected":  "2024-07-20",
		},
		"Aug": map[string]string{
			"givenDate": "29 de Agosto de 2023",
			"expected":  "2023-08-29",
		},
		"Sep": map[string]string{
			"givenDate": "18 de Septiembre de 1987",
			"expected":  "1987-09-18",
		},
		"Oct": map[string]string{
			"givenDate": "18 de Octubre de 2026",
			"expected":  "2026-10-18",
		},
		"Nov": map[string]string{
			"givenDate": "13 de Noviembre de 2050",
			"expected":  "2050-11-13",
		},
		"Dec": map[string]string{
			"givenDate": "31 de Diciembre de 3099",
			"expected":  "3099-12-31",
		},
	}

	for _, c := range cases {
		actual := ConvertInfobaeDate(c["givenDate"])
		if c["expected"] != actual {
			t.Errorf("Expected: %q, Actual %q", c["expected"], actual)
		}
	}
}
