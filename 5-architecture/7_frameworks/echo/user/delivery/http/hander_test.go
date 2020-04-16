package http

import (
	"echo_example/model"
	"echo_example/user"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo"
)

func TestGetUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	usecase := user.NewMockUsecase(ctrl)

	usecase.EXPECT().GetUser(gomock.Eq("ivan")).Return(model.User{
		Username: "Ivan",
	}, nil)

	handler := userHandler{usecase: usecase}
	e := echo.New()

	req := httptest.NewRequest(http.MethodGet, "/user/ivan", nil)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)
	c.SetPath("/user/:username")
	c.SetParamNames("username")
	c.SetParamValues("ivan")

	err := handler.GetUser(c)

	if err != nil {
		t.Errorf("err is not nil: %s", err)
	}

	body, _ := ioutil.ReadAll(rec.Body)

	if strings.Trim(string(body), "\n") != `{"Username":"Ivan"}` {
		t.Errorf("Expected: %s, got: %s", `{"Username":"Ivan"}`, string(body))
	}

}
