package terms_test

import (
	"fmt"
	"github.com/datapace/datapace/terms"
	"github.com/datapace/datapace/terms/mocks"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	termsText1 = "this is sample terms text1"
)

// newServer creates a http server serving mock terms
func newServer() *httptest.Server {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/nonexisting" {
			w.WriteHeader(http.StatusNotFound)
		} else {
			fmt.Fprintln(w, termsText1)
		}
	}))
	return ts
}

func newService() terms.Service {
	tl := mocks.NewTermsLedger()
	tr := mocks.NewTermsRepository()
	return terms.New(tr, tl)
}

func TestCreateTerms(t *testing.T) {
	svc := newService()
	srv := newServer()

	type args struct {
		t terms.Terms
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "create valid terms",
			args: args{t: terms.Terms{
				StreamID: "123",
				TermsURL: srv.URL,
			}},
			want:    "",
			wantErr: false,
		},
		{
			name: "try to get non existing terms",
			args: args{t: terms.Terms{
				StreamID: "456",
				TermsURL: srv.URL + "/nonexisting",
			}},
			want:    "",
			wantErr: true,
		},
		{
			name: "try to get non reachable terms",
			args: args{t: terms.Terms{
				StreamID: "456",
				TermsURL: "http://localhost:66000",
			}},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := svc.CreateTerms(tt.args.t)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateTerms() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CreateTerms() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidateTerms(t *testing.T) {
	svc := newService()
	srv := newServer()
	type args struct {
		t terms.Terms
	}
	tests := []struct {
		name      string
		wrongHash string
		args      args
		want      bool
		wantErr   bool
	}{
		{
			name: "create and validate terms",
			args: args{t: terms.Terms{
				StreamID: "789",
				TermsURL: srv.URL,
			}},
			want:    true,
			wantErr: false,
		},
		{
			name:      "validate against wrong hash",
			wrongHash: "thisiswronghash",
			args: args{t: terms.Terms{
				StreamID: "091",
				TermsURL: srv.URL,
			}},
			want:    false,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, err := svc.CreateTerms(tt.args.t)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateTerms() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			createdTerms := tt.args.t
			createdTerms.TermsHash = hash
			if tt.wrongHash != "" {
				createdTerms.TermsHash = tt.wrongHash
			}
			got, err := svc.ValidateTerms(createdTerms)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateTerms() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ValidateTerms() got = %v, want %v", got, tt.want)
			}
		})
	}
}
