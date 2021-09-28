package convert

import (
	"fmt"
	"github.com/kamva/mgm/v3/internal/sbo/fileconfig"
)

const (
	//specialChars special chars used in username or password
	specialChars = ":/?#[]@"
	// keyDatabase DB instance name
	keyDatabase = "database"
	// keyHost DB host
	keyHost = "host"
	// keyPort DB port
	keyPort = "port"
	// keyOptions Connection options
	keyOptions = "options"
	// keyUsername DB User name
	keyUsername = "username"
	// keyPassword  DB User Password
	keyPassword = "password"
	// keySrv true to use DNS seed list
	keySrv = "srv"
	// mongodb
	mongodb = "mongodb"
	// postgresql
	postgreSql = "postgresql"
)

// Converter converts and returns a ServiceBinding object to a DSN connection string
type Converter interface {
	// Convert converts a ServiceBinding obj to a connection string in DSN format
	Convert(sb fileconfig.ServiceBinding) string
}

// GetPostgreSQLConnectionString converts and returns the BindingService to a PostgreSQL DSN connection string
func GetPostgreSQLConnectionString() (string, error) {
	dbBinding, err := singleMatchingByType(postgreSql)
	if err != nil {
		return "", err
	}
	converter := &PostgreSQLConnectionStringConverter{}
	return converter.Convert(dbBinding), nil
}

// GetMongoDBConnectionString converts and returns a BindingService to a MongoDB DSN connection string
func GetMongoDBConnectionString() (string, error) {
	dbBinding, err := singleMatchingByType(mongodb)
	if err != nil {
		return "", err
	}
	converter := &MongoDBConverter{}
	return converter.Convert(dbBinding), nil
}

//
func matchingByType(bindingType string, serviceBindings []fileconfig.ServiceBinding) []fileconfig.ServiceBinding {
	var res []fileconfig.ServiceBinding
	if len(serviceBindings) == 0 {
		return res
	}
	for _, sb := range serviceBindings {
		if sb.BindingType == bindingType {
			res = append(res, sb)
		}
	}
	return res
}

func singleMatchingByType(bindingType string) (fileconfig.ServiceBinding, error) {
	bindingFileReader := fileconfig.NewBindingReader()
	serviceBindings, err := bindingFileReader.ReadServiceBindingConfig()
	var dbBinging fileconfig.ServiceBinding
	if err != nil {
		return dbBinging, err
	}
	matchingBindings := matchingByType(bindingType, serviceBindings)
	if len(matchingBindings) == 0 {
		return dbBinging, fmt.Errorf("no service binding config found for %s", bindingType)
	}
	return matchingBindings[0], nil
}
