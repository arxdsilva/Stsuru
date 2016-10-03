package shortener

import (
	"fmt"
	"net/url"
	"reflect"
	"testing"
)

var custom = "localhost:8080"

func TestShorten(t *testing.T) {
	type args struct {
		u          *url.URL
		customHost string
	}
	tests := []struct {
		name    string
		args    args
		want    *url.URL
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "test01",
			args: args{
				u: &url.URL{
					Scheme: "https",
					Host:   "github.com",
					Path:   "arxdsilva/stsuru",
				},
				customHost: custom,
			},
			want: &url.URL{
				Scheme: "https",
				Host:   "localhost:8080",
				Path:   "e6944260aeb6d6bc819abcc9c2e2ab61",
			},
			wantErr: false,
		},
		{
			name: "test02",
			args: args{
				u: &url.URL{
					Scheme: "https",
					Host:   "github.com",
					Path:   "tsuru/tsuru",
				},
				customHost: "",
			},
			want: &url.URL{
				Scheme: "https",
				Host:   "github.com",
				Path:   "8214e81107a57c97d827272f8ef77c04",
			},
			wantErr: false,
		},
		// {
		// 	name: "test03",
		// 	args: args{
		// 		u: &url.URL{
		// 			Scheme: "https",
		// 			Host:   "notvalid",
		// 			Path:   "",
		// 		},
		// 		customHost: custom,
		// 	},
		// 	want:    &url.URL{},
		// 	wantErr: true,
		// },
	}
	for _, tt := range tests {
		got, err := Shorten(tt.args.u, tt.args.customHost)
		if (err != nil) != tt.wantErr {
			fmt.Println(tt.args.u)
			fmt.Println(got)
			t.Errorf("1 - %q. Shorten() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("2 - %q. Shorten() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func Test_validateURL(t *testing.T) {
	type args struct {
		u *url.URL
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "01",
			args: args{
				u: &url.URL{
					Scheme: "https",
					Host:   "notvalid",
					Path:   "",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		err := validateURL(tt.args.u)
		fmt.Println(err)
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. validateURL() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
	}
}
