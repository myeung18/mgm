package fileconfig

import "os"

var bindingRootEnvVar = os.Getenv("SERVICE_BINDING_ROOT")

// GetBindingRootDirectory returns the root directory where service binding operator will save the config files
func GetBindingRootDirectory() string {
	if bindingRootEnvVar == "" {
		return "/bindings"
	}
	return bindingRootEnvVar
}
