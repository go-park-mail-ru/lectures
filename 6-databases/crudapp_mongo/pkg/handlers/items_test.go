package handlers

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http/httptest"
	"testing"

	"crudapp/pkg/items"

	"github.com/golang/mock/gomock"
	"go.uber.org/zap"
)

/*

repo
	db -> mock

handler
	repo -> mock


*/

func TestItemsHandlerList(t *testing.T) {

	// мы передаём t сюда, это надо чтобы получить корректное сообщение если тесты не пройдут
	ctrl := gomock.NewController(t)

	// Finish сравнит последовательсноть вызовов и выведет ошибку если последовательность другая
	defer ctrl.Finish()

	st := NewMockItemRepository(ctrl)
	service := &ItemsHandler{
		ItemsRepo: st,
		Logger:    zap.NewNop().Sugar(), // не пишет логи
		Tmpl:      template.Must(template.ParseGlob("../../templates/*")),
	}

	resultItems := []*items.Item{
		{ID: 1, Title: "database/sql"},
	}

	// тут мы записываем последовтаельность вызовов и результат
	st.EXPECT().GetAll().Return(resultItems, nil)

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	service.List(w, req)

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)

	img := `database/sql`
	if !bytes.Contains(body, []byte(img)) {
		t.Errorf("no text found")
		return
	}

	// GetPhotos error
	// тут мы записываем последовтаельность вызовов и результат
	st.EXPECT().GetAll().Return(nil, fmt.Errorf("no results"))

	req = httptest.NewRequest("GET", "/", nil)
	w = httptest.NewRecorder()

	service.List(w, req)

	resp = w.Result()
	if resp.StatusCode != 500 {
		t.Errorf("expected resp status 500, got %d", resp.StatusCode)
		return
	}

	// template expand error
	service.Tmpl, _ = template.New("tmplError").Parse("{{.NotExist}}")

	st.EXPECT().GetAll().Return(resultItems, nil)

	req = httptest.NewRequest("GET", "/", nil)
	w = httptest.NewRecorder()

	service.List(w, req)

	resp = w.Result()
	if resp.StatusCode != 500 {
		t.Errorf("expected resp status 500, got %d", resp.StatusCode)
		return
	}

}
