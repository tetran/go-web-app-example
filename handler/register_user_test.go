package handler

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/tetran/go-web-app-example/entity"
	"github.com/tetran/go-web-app-example/testutil"
)

func TestRegisterUser(t *testing.T) {
	t.Parallel()

	type want struct {
		status  int
		rspFile string
	}
	tests := map[string]struct {
		reqFile string
		want    want
	}{
		"ok": {
			reqFile: "testdata/register_user/ok_req.json.golden",
			want: want{
				status:  http.StatusCreated,
				rspFile: "testdata/register_user/ok_rsp.json.golden",
			},
		},
		"badRequest": {
			reqFile: "testdata/register_user/bad_request_req.json.golden",
			want: want{
				status:  http.StatusBadRequest,
				rspFile: "testdata/register_user/bad_request_rsp.json.golden",
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
				"/users",
				bytes.NewReader(testutil.LoadFile(t, tt.reqFile)),
			)

			moq := &RegisterUserServiceMock{}
			moq.RegisterUserFunc = func(ctx context.Context, name, password, role string) (*entity.User, error) {
				if tt.want.status == http.StatusCreated {
					return &entity.User{ID: 1}, nil
				}
				return nil, errors.New("error from mock")
			}

			sut := RegisterUser{
				Service:   moq,
				Validator: validator.New(),
			}
			sut.ServeHTTP(w, r)

			resp := w.Result()
			testutil.AssertResponse(t, resp, tt.want.status, testutil.LoadFile(t, tt.want.rspFile))
		})
	}
}
