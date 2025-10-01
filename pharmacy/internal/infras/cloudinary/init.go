package cloudinary

import (
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/ritchieridanko/apotekly-api/pharmacy/config"
)

func Initialize() (instance *cloudinary.Cloudinary, err error) {
	cloud := config.StorageGetName()
	key := config.StorageGetAPIKey()
	secret := config.StorageGetAPISecret()

	return cloudinary.NewFromParams(cloud, key, secret)
}
