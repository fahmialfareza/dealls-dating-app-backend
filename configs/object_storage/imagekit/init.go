package imagekit

import (
	"github.com/fahmialfareza/deals-dating-app-backend/configs/environment"
	"github.com/imagekit-developer/imagekit-go"
)

func NewImageKit() *imagekit.ImageKit {
	ik := imagekit.NewFromParams(imagekit.NewParams{
		PublicKey:   environment.ImageKitPublicAPIKey,  // Replace with your ImageKit public key
		PrivateKey:  environment.ImageKitPrivateAPIKey, // Replace with your ImageKit private key
		UrlEndpoint: environment.ImageKitURLEndpoint,   // Replace with your ImageKit URL endpoint
	})

	return ik
}
