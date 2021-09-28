package fileconfig

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
)

var bindingRootDirectory = GetBindingRootDirectory()

const (
	// Provider is the database provider name
	Provider = "provider"
	// BindingType is the database binding type e.g. mongodb, postgresql
	BindingType = "type"
)

// BindingFileReader read the file config
type BindingFileReader struct {
	ReadDir  func(filename string) ([]fs.FileInfo, error)
	ReadFile func(filename string) ([]byte, error)
	Stat     func(fileanme string) (fs.FileInfo, error)
}

//NewBindingReader creates and returns and NewBindingReader object
func NewBindingReader() *BindingFileReader {
	return &BindingFileReader{
		ReadDir:  ioutil.ReadDir,
		ReadFile: ioutil.ReadFile,
		Stat:     os.Stat,
	}
}

// ReadServiceBindingConfig reads binding config files and converts them into a list of ServiceBinding objects
func (bfr *BindingFileReader) ReadServiceBindingConfig() ([]ServiceBinding, error) {
	fs, err := bfr.Stat(bindingRootDirectory)
	if err != nil {
		return nil, err
	}
	if !fs.IsDir() {
		return nil, fmt.Errorf("service Binding root %s is not a directory", bindingRootDirectory)
	}

	lstFile, err := bfr.ReadDir(bindingRootDirectory)
	if err != nil {
		return nil, err
	}

	// each directory is a set of binding config
	bindings := []ServiceBinding{}
	for _, dir := range lstFile {
		if dir.IsDir() {
			provider, bindingType, content, err := bfr.readBindingContent(dir)
			if err != nil {
				return nil, err
			}
			bindings = append(bindings, ServiceBinding{
				Name:        dir.Name(),
				Provider:    provider,
				Properties:  content,
				BindingType: bindingType,
			})
		}
	}
	return bindings, err
}

func (bfr *BindingFileReader) readBindingContent(bindigDir os.FileInfo) (string, string, map[string]string, error) {
	subDir := filepath.Join(bindingRootDirectory, bindigDir.Name())
	subDirFiles, err := bfr.ReadDir(subDir)
	if err != nil {
		return "", "", nil, err
	}
	var provider, bindingType string
	content := make(map[string]string)
	for _, prop := range subDirFiles {
		//skip sub-dir and hidden files
		if prop.IsDir() || prop.Name()[0:1] == "." {
			continue
		}
		buf, err := bfr.ReadFile(filepath.Join(subDir, prop.Name()))
		if err != nil {
			return "", "", nil, err
		}
		if prop.Name() == BindingType {
			bindingType = string(buf)
		} else if prop.Name() == Provider {
			provider = string(buf)
		} else {
			content[prop.Name()] = string(buf)
		}
	}
	return provider, bindingType, content, nil
}
