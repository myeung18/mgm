package fileconfig

import (
	"io/fs"
	"reflect"
	"testing"
	"time"
)

type fakeDir struct{ fn string }

func (f fakeDir) Name() string       { return f.fn }
func (f fakeDir) Size() int64        { return 0 }
func (f fakeDir) Mode() fs.FileMode  { return 0 }
func (f fakeDir) ModTime() time.Time { return time.Unix(0, 0) }
func (f fakeDir) IsDir() bool        { return true }
func (f fakeDir) Sys() interface{}   { return nil }

type fakeFile struct{ fn string }

func (f fakeFile) Name() string       { return f.fn }
func (f fakeFile) Size() int64        { return 0 }
func (f fakeFile) Mode() fs.FileMode  { return 0 }
func (f fakeFile) ModTime() time.Time { return time.Unix(0, 0) }
func (f fakeFile) IsDir() bool        { return false }
func (f fakeFile) Sys() interface{}   { return nil }

func TestBindingFileReader_ReadServiceBindingConfig(t *testing.T) {
	type fields struct {
		ReadDir  func(filename string) ([]fs.FileInfo, error)
		ReadFile func(filename string) ([]byte, error)
		Stat     func(fileanme string) (fs.FileInfo, error)
	}
	tests := []struct {
		name    string
		fields  fields
		want    []ServiceBinding
		wantErr bool
	}{
		{
			name: "failed to read binding folder - not a directory",
			fields: fields{
				ReadDir: func(filename string) ([]fs.FileInfo, error) {
					var fi []fs.FileInfo
					if filename == "/bindings" {
						fi = []fs.FileInfo{
							fakeFile{fn: "local"},
						}
					}
					return fi, nil
				},
				ReadFile: func(filename string) ([]byte, error) {
					switch filename {
					case "/bindings/local/type":
						return []byte("mongodb"), nil
					default:
						return []byte(filename), nil
					}
				},
				Stat: func(filename string) (fs.FileInfo, error) {
					if filename == "/bindings" {
						return fakeFile{fn: filename}, nil //it is a file, error
					}
					return fakeFile{fn: filename}, nil
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "successfully reading binding data",
			fields: fields{
				ReadDir: func(filename string) ([]fs.FileInfo, error) {
					var fi []fs.FileInfo
					if filename == "/bindings" {
						fi = []fs.FileInfo{
							fakeDir{fn: "local"},
						}
					} else if filename == "/bindings/local" {
						fi = []fs.FileInfo{
							fakeFile{fn: "host"},
							fakeFile{fn: "provider"},
							fakeFile{fn: "username"},
							fakeFile{fn: "type"},
							fakeFile{fn: "password"},
						}
					}
					return fi, nil
				},
				ReadFile: func(filename string) ([]byte, error) {
					switch filename {
					case "/bindings/local/type":
						return []byte("mongodb"), nil
					default:
						return []byte(filename), nil
					}
				},
				Stat: func(filename string) (fs.FileInfo, error) {
					if filename == "/bindings" {
						return fakeDir{fn: filename}, nil
					}
					return fakeFile{fn: filename}, nil
				},
			},
			want: []ServiceBinding{
				{
					Name:        "local",
					BindingType: "mongodb",
					Provider:    "/bindings/local/provider",
					Properties: map[string]string{
						"host":     "/bindings/local/host",
						"username": "/bindings/local/username",
						"password": "/bindings/local/password",
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bfr := &BindingFileReader{
				ReadDir:  tt.fields.ReadDir,
				ReadFile: tt.fields.ReadFile,
				Stat:     tt.fields.Stat,
			}
			got, err := bfr.ReadServiceBindingConfig()
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadServiceBindingConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReadServiceBindingConfig() got = %v, want %v", got, tt.want)
			}
		})
	}
}
