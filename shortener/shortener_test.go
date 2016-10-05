package shortener

import (
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
		{
			name: "test03",
			args: args{
				u: &url.URL{
					Scheme: "https",
					Host:   "notvalid",
					Path:   "",
				},
				customHost: custom,
			},
			want:    &url.URL{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		a := NewShorten{
			U:          tt.args.u,
			CustomHost: tt.args.customHost,
			Token:      "a",
		}
		got, err := a.Shorten()
		if (err != nil) != tt.wantErr {
			t.Errorf("1 - %q. Shorten() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if tt.wantErr == true {
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
			wantErr: true,
		},
		{
			name: "02",
			args: args{
				u: &url.URL{
					Scheme: "https",
					Host:   "",
					Path:   "notvalid/path",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		err := validateURL(tt.args.u)
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. validateURL() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
	}
}

func TestTokenGenerator(t *testing.T) {
	tests := []struct {
		Num int
		exp int
	}{
		{
			Num: 0,
			exp: 4,
		},
		{
			Num: 6,
			exp: 12,
		},
		{
			Num: 8,
			exp: 16,
		},
	}
	for _, tt := range tests {
		a := tokenGenerator(tt.Num)
		if len(a) != tt.exp {
			t.Errorf("Token of %v bytes expected but %v bytes found", tt.exp, len(a))
		}
	}
}
