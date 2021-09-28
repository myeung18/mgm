package convert

import (
	"fmt"
	"github.com/kamva/mgm/v3/internal/sbo/fileconfig"
	"testing"
)

func TestPostgreSQLConnectionStringConverter_Convert(t *testing.T) {
	type args struct {
		binding fileconfig.ServiceBinding
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Correct connection string returned",
			args: args{
				binding: fileconfig.ServiceBinding{
					Name:        "local",
					BindingType: "postgresql",
					Provider:    "Crunchy Bridges",
					Properties: map[string]string{
						"host":     "example.com:10011",
						"username": "a-db-user",
						"password": "password",
						"srv":      "true",
						"options":  "some-db-options",
						"database": "local-db",
					},
				},
			},
			want: "host=example.com:10011 user=a-db-user password=password dbname=local-db",
		},
		{
			name: "Correct connection string returned without password",
			args: args{
				binding: fileconfig.ServiceBinding{
					Name:        "local",
					BindingType: "postgresql",
					Provider:    "Crunchy Bridges",
					Properties: map[string]string{
						"host":     "example.com:10011",
						"username": "a-db-user",
						"srv":      "true",
						"options":  "some-db-options",
						"database": "local-db",
					},
				},
			},
			want: "host=example.com:10011 user=a-db-user dbname=local-db",
		},
		{
			name: "Correct connection string returned with special char escaping",
			args: args{
				binding: fileconfig.ServiceBinding{
					Name:        "local",
					BindingType: "postgresql",
					Provider:    "Crunchy Bridges",
					Properties: map[string]string{
						"host":     "example.com:10011",
						"username": "a-db-user",
						"password": "password'",
						"srv":      "true",
						"options":  "some-db-options",
						"database": "local-db",
					},
				},
			},
			want: "host=example.com:10011 user=a-db-user password=password\\' dbname=local-db",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &PostgreSQLConnectionStringConverter{}
			if got := m.Convert(tt.args.binding); got != tt.want {
				t.Errorf("Convert() = %v, want %v", got, tt.want)
			}
			fmt.Println(m.Convert(tt.args.binding))
		})
	}
}
