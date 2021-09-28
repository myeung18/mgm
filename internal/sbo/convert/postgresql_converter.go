package convert

import (
	"github.com/kamva/mgm/v3/internal/sbo/fileconfig"
	"strings"
)

var escaper = strings.NewReplacer(` `, `\ `, `'`, `\'`, `\`, `\\`)

// PostgreSQLConnectionStringConverter converts a ServiceBinding object to a DSN connection string
type PostgreSQLConnectionStringConverter struct{}

// Convert converts a ServiceBinding to a DSN string
func (m *PostgreSQLConnectionStringConverter) Convert(binding fileconfig.ServiceBinding) string {
	var parts []string
	addToParts := func(k, v string) {
		parts = append(parts, k+"="+escaper.Replace(v))
	}
	addToParts("host", binding.Properties[keyHost])
	if binding.Properties[keyPort] != "" {
		addToParts("port", binding.Properties[keyPort])
	}
	if binding.Properties[keyUsername] != "" {
		addToParts("user", binding.Properties[keyUsername])
	}
	if binding.Properties[keyPassword] != "" {
		addToParts("password", binding.Properties[keyPassword])
	}
	addToParts("dbname", binding.Properties[keyDatabase])
	return strings.Join(parts, " ")
}
