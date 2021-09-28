package convert

import (
	"github.com/kamva/mgm/v3/internal/sbo/fileconfig"
	"net/url"
	"strings"
)

// MongoDBConverter type for MongoDB
type MongoDBConverter struct{}

// Convert converts a mongodb ServiceBinding to a connection string
func (m *MongoDBConverter) Convert(binding fileconfig.ServiceBinding) string {
	prefix := "mongodb://"
	if strings.EqualFold(binding.Properties[keySrv], "true") {
		prefix = "mongodb+srv://"
	}

	database := ""
	if binding.Properties[keyOptions] != "" {
		database = "?" + binding.Properties[keyOptions]
	}
	if binding.Properties[keyDatabase] != "" {
		database = "/" + binding.Properties[keyDatabase] + database
	} else if binding.Properties[keyOptions] != "" {
		database = "/" + database
	}

	return strings.Join([]string{prefix,
		encodeIfContainsSpecialCharacters(binding.Properties[keyUsername]), ":",
		encodeIfContainsSpecialCharacters(binding.Properties[keyPassword]), "@",
		binding.Properties[keyHost],
		database}, "")
}

func encodeIfContainsSpecialCharacters(userNameOrPassword string) string {
	if strings.ContainsAny(userNameOrPassword, specialChars) {
		return url.QueryEscape(userNameOrPassword)
	}
	return userNameOrPassword
}
