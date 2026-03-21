package dto

type City struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Region struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Cities []City `json:"cities"`
}

type Country struct {
	ID      string   `json:"id"`
	Name    string   `json:"name"`
	Regions []Region `json:"regions"`
}

type CountryDTO struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type RegionDTO struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

var geoData = []Country{
	{
		ID: "US", Name: "United States",
		Regions: []Region{
			{ID: "CA", Name: "California", Cities: []City{
				{ID: "LA", Name: "Los Angeles"},
				{ID: "SF", Name: "San Francisco"},
			}},
			{ID: "NY", Name: "New York", Cities: []City{
				{ID: "NYC", Name: "New York City"},
				{ID: "BUF", Name: "Buffalo"},
			}},
		},
	},
	{
		ID: "BR", Name: "Brazil",
		Regions: []Region{
			{ID: "SP", Name: "São Paulo", Cities: []City{
				{ID: "SAO", Name: "São Paulo"},
				{ID: "CAMP", Name: "Campinas"},
			}},
			{ID: "RJ", Name: "Rio de Janeiro", Cities: []City{
				{ID: "RIO", Name: "Rio de Janeiro"},
				{ID: "NIT", Name: "Niterói"},
			}},
		},
	},
	{
		ID: "RU", Name: "Russia",
		Regions: []Region{
			{ID: "MOW", Name: "Moscow", Cities: []City{
				{ID: "MSK", Name: "Moscow"},
			}},
			{ID: "SPE", Name: "Saint Petersburg", Cities: []City{
				{ID: "SPB", Name: "Saint Petersburg"},
			}},
			{ID: "NSK", Name: "Novosibirsk", Cities: []City{
				{ID: "NSK", Name: "Novosibirsk"},
				{ID: "BD", Name: "Berdsk"},
			}},
		},
	},
	{
		ID: "DE", Name: "Germany",
		Regions: []Region{
			{ID: "BE", Name: "Berlin", Cities: []City{
				{ID: "BER", Name: "Berlin"},
			}},
			{ID: "BY", Name: "Bavaria", Cities: []City{
				{ID: "MUC", Name: "Munich"},
			}},
		},
	},
	{
		ID: "FR", Name: "France",
		Regions: []Region{
			{ID: "IDF", Name: "Île-de-France", Cities: []City{
				{ID: "PAR", Name: "Paris"},
			}},
			{ID: "PAC", Name: "Provence-Alpes-Côte d'Azur", Cities: []City{
				{ID: "MRS", Name: "Marseille"},
			}},
		},
	},
	{
		ID: "IT", Name: "Italy",
		Regions: []Region{
			{ID: "LAZ", Name: "Lazio", Cities: []City{
				{ID: "ROM", Name: "Rome"},
			}},
			{ID: "LOM", Name: "Lombardy", Cities: []City{
				{ID: "MIL", Name: "Milan"},
			}},
		},
	},
	{
		ID: "ES", Name: "Spain",
		Regions: []Region{
			{ID: "MD", Name: "Madrid", Cities: []City{
				{ID: "MAD", Name: "Madrid"},
			}},
			{ID: "CAT", Name: "Catalonia", Cities: []City{
				{ID: "BCN", Name: "Barcelona"},
			}},
		},
	},
	{
		ID: "CN", Name: "China",
		Regions: []Region{
			{ID: "BJ", Name: "Beijing", Cities: []City{
				{ID: "BEI", Name: "Beijing"},
			}},
			{ID: "SH", Name: "Shanghai", Cities: []City{
				{ID: "SHA", Name: "Shanghai"},
			}},
		},
	},
	{
		ID: "JP", Name: "Japan",
		Regions: []Region{
			{ID: "TK", Name: "Tokyo", Cities: []City{
				{ID: "TOK", Name: "Tokyo"},
			}},
			{ID: "OS", Name: "Osaka", Cities: []City{
				{ID: "OSA", Name: "Osaka"},
			}},
		},
	},
	{
		ID: "IN", Name: "India",
		Regions: []Region{
			{ID: "DL", Name: "Delhi", Cities: []City{
				{ID: "DEL", Name: "New Delhi"},
			}},
			{ID: "MH", Name: "Maharashtra", Cities: []City{
				{ID: "MUM", Name: "Mumbai"},
			}},
		},
	},
	{
		ID: "GB", Name: "United Kingdom",
		Regions: []Region{
			{ID: "ENG", Name: "England", Cities: []City{
				{ID: "LON", Name: "London"},
			}},
			{ID: "SCT", Name: "Scotland", Cities: []City{
				{ID: "EDI", Name: "Edinburgh"},
			}},
		},
	},
	{
		ID: "CA", Name: "Canada",
		Regions: []Region{
			{ID: "ON", Name: "Ontario", Cities: []City{
				{ID: "TOR", Name: "Toronto"},
			}},
			{ID: "BC", Name: "British Columbia", Cities: []City{
				{ID: "VAN", Name: "Vancouver"},
			}},
		},
	},
	{
		ID: "AU", Name: "Australia",
		Regions: []Region{
			{ID: "NSW", Name: "New South Wales", Cities: []City{
				{ID: "SYD", Name: "Sydney"},
			}},
			{ID: "VIC", Name: "Victoria", Cities: []City{
				{ID: "MEL", Name: "Melbourne"},
			}},
		},
	},
	{
		ID: "MX", Name: "Mexico",
		Regions: []Region{
			{ID: "CMX", Name: "Mexico City", Cities: []City{
				{ID: "MEX", Name: "Mexico City"},
			}},
		},
	},
	{
		ID: "AR", Name: "Argentina",
		Regions: []Region{
			{ID: "BA", Name: "Buenos Aires", Cities: []City{
				{ID: "BUE", Name: "Buenos Aires"},
			}},
		},
	},
	{
		ID: "ZA", Name: "South Africa",
		Regions: []Region{
			{ID: "GP", Name: "Gauteng", Cities: []City{
				{ID: "JHB", Name: "Johannesburg"},
			}},
		},
	},
	{
		ID: "EG", Name: "Egypt",
		Regions: []Region{
			{ID: "C", Name: "Cairo", Cities: []City{
				{ID: "CAI", Name: "Cairo"},
			}},
		},
	},
	{
		ID: "TR", Name: "Turkey",
		Regions: []Region{
			{ID: "34", Name: "Istanbul", Cities: []City{
				{ID: "IST", Name: "Istanbul"},
			}},
		},
	},
	{
		ID: "KR", Name: "South Korea",
		Regions: []Region{
			{ID: "SEO", Name: "Seoul", Cities: []City{
				{ID: "SEL", Name: "Seoul"},
			}},
		},
	},
}

func GetCountries() []Country {
	return geoData
}

func GetRegions(countryID string) []Region {
	for _, region := range geoData {
		if region.ID == countryID {
			return region.Regions
		}
	}
	return nil
}

func GetCities(countryID, regionID string) []City {
	for _, c := range geoData {
		if c.ID == countryID {
			for _, r := range c.Regions {
				if r.ID == regionID {
					return r.Cities
				}
			}
		}
	}
	return nil
}
