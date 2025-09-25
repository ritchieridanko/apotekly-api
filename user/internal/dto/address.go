package dto

type RespAddress struct {
	ID          int64   `json:"id"`
	Receiver    string  `json:"receiver"`
	Phone       string  `json:"phone"`
	Label       string  `json:"label"`
	Notes       *string `json:"notes"`
	IsPrimary   bool    `json:"is_primary"`
	Country     string  `json:"country"`
	AdminLevel1 *string `json:"admin_level_1"`
	AdminLevel2 *string `json:"admin_level_2"`
	AdminLevel3 *string `json:"admin_level_3"`
	AdminLevel4 *string `json:"admin_level_4"`
	Street      string  `json:"street"`
	PostalCode  string  `json:"postal_code"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
}

type ReqNewAddress struct {
	Receiver    string  `json:"receiver" binding:"required"`
	Phone       string  `json:"phone" binding:"required"`
	Label       string  `json:"label" binding:"required"`
	Notes       *string `json:"notes"`
	IsPrimary   bool    `json:"is_primary"`
	Country     string  `json:"country" binding:"required"`
	AdminLevel1 *string `json:"admin_level_1"`
	AdminLevel2 *string `json:"admin_level_2"`
	AdminLevel3 *string `json:"admin_level_3"`
	AdminLevel4 *string `json:"admin_level_4"`
	Street      string  `json:"street" binding:"required"`
	PostalCode  string  `json:"postal_code" binding:"required"`
	Latitude    float64 `json:"latitude" binding:"required"`
	Longitude   float64 `json:"longitude" binding:"required"`
}

type RespNewAddress struct {
	Created        RespAddress `json:"created"`
	UnsetPrimaryID *int64      `json:"unset_primary_id,omitempty"`
}

type ReqUpdateAddress struct {
	Receiver    *string  `json:"receiver"`
	Phone       *string  `json:"phone"`
	Label       *string  `json:"label"`
	Notes       *string  `json:"notes"`
	Country     *string  `json:"country"`
	AdminLevel1 *string  `json:"admin_level_1"`
	AdminLevel2 *string  `json:"admin_level_2"`
	AdminLevel3 *string  `json:"admin_level_3"`
	AdminLevel4 *string  `json:"admin_level_4"`
	Street      *string  `json:"street"`
	PostalCode  *string  `json:"postal_code"`
	Latitude    *float64 `json:"latitude"`
	Longitude   *float64 `json:"longitude"`
}

type RespUpdateAddress struct {
	Updated RespAddress `json:"updated"`
}

type RespChangePrimaryAddress struct {
	NewPrimaryID   int64 `json:"new_primary_id"`
	UnsetPrimaryID int64 `json:"unset_primary_id"`
}

type RespDeleteAddress struct {
	DeletedID    int64  `json:"deleted_id"`
	NewPrimaryID *int64 `json:"new_primary_id,omitempty"`
}
