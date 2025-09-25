package entities

type Address struct {
	ID          int64
	Receiver    string
	Phone       string
	Label       string
	Notes       *string
	IsPrimary   bool
	Country     string
	AdminLevel1 *string
	AdminLevel2 *string
	AdminLevel3 *string
	AdminLevel4 *string
	Street      string
	PostalCode  string
	Latitude    float64
	Longitude   float64
}

type NewAddress struct {
	Receiver    string
	Phone       string
	Label       string
	Notes       *string
	IsPrimary   bool
	Country     string
	AdminLevel1 *string
	AdminLevel2 *string
	AdminLevel3 *string
	AdminLevel4 *string
	Street      string
	PostalCode  string
	Latitude    float64
	Longitude   float64
}

type AddressChange struct {
	Receiver    *string
	Phone       *string
	Label       *string
	Notes       *string
	Country     *string
	AdminLevel1 *string
	AdminLevel2 *string
	AdminLevel3 *string
	AdminLevel4 *string
	Street      *string
	PostalCode  *string
	Latitude    *float64
	Longitude   *float64
}
