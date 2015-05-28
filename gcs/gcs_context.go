package gcs

import (
	"github.com/wayt/happyngine/env"
	"github.com/wayt/happyngine/log"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/cloud"
	"google.golang.org/cloud/storage"
	"io/ioutil"
)

var CloudCtx context.Context
var GoogleAccessID string
var PrivateKey []byte

func init() {
	jsonKey, err := ioutil.ReadFile(env.Get("HAPPY_GCS_KEY_JSON_FILE"))
	if err != nil {
		log.Criticalln(err)
	}
	conf, err := google.JWTConfigFromJSON(
		jsonKey,
		storage.ScopeReadWrite,
	)
	if err != nil {
		log.Criticalln(err)
	}
	CloudCtx = cloud.NewContext(env.Get("HAPPY_GCS_APP_ID"), conf.Client(oauth2.NoContext))

	GoogleAccessID = env.Get("HAPPY_GCS_ACCESS_ID")

	if gcsPrivateKeyFile := env.Get("HAPPY_GCS_PRIVATE_KEY_FILE"); len(gcsPrivateKeyFile) > 0 {
		PrivateKey, err = ioutil.ReadFile(gcsPrivateKeyFile)
	}
}
