package handlers

import (
	"net/http"
	"testing"
)

func TestHandler_GetFullURL(t *testing.T) {
	type fields struct {
		db DatabaseInterface
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				db: tt.fields.db,
			}
			h.GetFullURL(tt.args.w, tt.args.r)
		})
	}
}
