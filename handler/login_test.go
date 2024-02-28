package handler

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/tetran/go-web-app-example/testutil"
)

func TestLogin(t *testing.T) {
	type moq struct {
		token string
		err   error
	}
	type want struct {
		status  int
		fspFile string
	}
	tests := map[string]struct {
		reqFile string
		moq     moq
		want    want
	}{
		"ok": {
			reqFile: "testdata/login/ok_req.json.golden",
			moq: moq{
				token: "from_moq",
			},
			want: want{
				status:  http.StatusOK,
				fspFile: "testdata/login/ok_rsp.json.golden",
			},
		},
		"badRequest": {
			reqFile: "testdata/login/bad_request_req.json.golden",
			want: want{
				status:  http.StatusBadRequest,
				fspFile: "testdata/login/bad_request_rsp.json.golden",
			},
		},
		"internalServerError": {
			reqFile: "testdata/login/ok_req.json.golden",
			moq: moq{
				err: errors.New("error from moq"),
			},
			want: want{
				status:  http.StatusInternalServerError,
				fspFile: "testdata/login/internal_server_error_rsp.json.golden",
			},
		},
	}

	for n, tt := range tests {
		tt := tt
		t.Run(n, func(t *testing.T) {
			t.Parallel()

			w := httptest.NewRecorder()
			r := httptest.NewRequest(
				http.MethodPost,
				"/login",
				bytes.NewReader(testutil.LoadFile(t, tt.reqFile)),
			)

			moq := &LoginServiceMock{}
			moq.LoginFunc = func(ctx context.Context, userName, password string) (string, error) {
				return tt.moq.token, tt.moq.err
			}

			sut := Login{
				Service:   moq,
				Validator: validator.New(),
			}
			sut.ServeHTTP(w, r)

			reqp := w.Result()
			testutil.AssertResponse(t, reqp, tt.want.status, testutil.LoadFile(t, tt.want.fspFile))
		})
	}
}
