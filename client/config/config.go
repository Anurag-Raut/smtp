package config

import (
	"os"
)

var ClientDomain = os.Getenv("CLIENT_DOMAIN")
