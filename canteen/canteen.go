package canteen

// Dish describes a single line in the CSV.
type Dish struct {
	Date          string `csv:"datum"`
	Day           string `csv:"tag"`
	Allergens     string `csv:"warengruppe"`
	Name          string `csv:"name"`
	MealType      string `csv:"kennz"`
	Price         string `csv:"preis"`
	PriceStudent  string `csv:"stud"`
	PriceEmployee string `csv:"bed"`
	PriceGuest    string `csv:"gast"`
}

// Canteens2Abbrev maps the university name to the abbreviations used in the online CSV source.
var Canteens2Abbrev = map[string]string{
	"UNI_REGENSBURG_CANTEEN":              "UNI-R",
	"UNI_REGENSBURG_GUEST_CANTEEN":        "UNI-R-Gs",
	"UNI_REGENSBURG_CAFETERIA_PT":         "Cafeteria-PT",
	"UNI_REGENSBURG_CAFETERIA_CHEMISTRY":  "Cafeteria-Chemie",
	"UNI_REGENSBURG_CAFETERIA_SPORT":      "Cafeteria-Sport",
	"OTH_REGENSBURG_CANTEEN_LUNCH":        "HS-R-tag",
	"OTH_REGENSBURG_CANTEEN_DINNER":       "HS-R-abend",
	"OTH_REGENSBURG_CAFETERIA_PRUEFENING": "Cafeteria-pruefening",
	"UNI_PASSAU_CANTEEN":                  "UNI-P",
	"UNI_PASSAU_NIKOLA_CAFETERIA":         "Cafeteria-Nikolakloster",
	"HS_DEGGENDORF_CANTEEN":               "HS-DEG",
	"HS_LANDSHUT_CANTEEN":                 "HS-LA",
	"HS_STRAUBING_CANTEEN":                "HS-SR",
	"HS_PFARRKIRCHEN_CANTEEN":             "HS-PAN",
}

var Abbrev2Canteens = reverseMap(Canteens2Abbrev)

func reverseMap(m map[string]string) map[string]string {
	n := make(map[string]string)
	for k, v := range m {
		n[v] = k
	}
	return n
}
