package cloudfunctions

import (
	"encoding/base64"
	"os"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
)

const provider = "google"

func init() {
	functions.HTTP("PopulateJobs", populateJobs)
}

func jsonFromEnv(env string) ([]byte, error) {
	encoded := os.Getenv(env)
	decoded, err := base64.URLEncoding.DecodeString(encoded)

	return decoded, err
}
