package fileconfig

// ServiceBinding type is the data structure for bound config file
type ServiceBinding struct {
	Name        string
	Provider    string
	Properties  map[string]string
	BindingType string
}
