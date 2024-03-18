package handlers

import (
	"encoding/json"
	"html/template"
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/lectures/6-databases/05_crudapp_db_tests/pkg/items"
	"github.com/go-park-mail-ru/lectures/6-databases/05_crudapp_db_tests/pkg/session"

	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"go.uber.org/zap"
)

// mockgen -source=items.go -destination=items_mock.go -package=handlers ItemRepositoryInterface
type ItemRepositoryInterface interface {
	GetAll() ([]*items.Item, error)
	GetByID(int64) (*items.Item, error)
	Add(*items.Item) (int64, error)
	Update(*items.Item) (int64, error)
	Delete(int64) (int64, error)
}

type ItemsHandler struct {
	Tmpl      *template.Template
	ItemsRepo ItemRepositoryInterface
	Logger    *zap.SugaredLogger
}

func (h *ItemsHandler) List(w http.ResponseWriter, r *http.Request) {
	elems, err := h.ItemsRepo.GetAll()
	if err != nil {
		h.Logger.Error("GetAll err", err)
		http.Error(w, `DB err`, http.StatusInternalServerError)
		return
	}

	err = h.Tmpl.ExecuteTemplate(w, "index.html", struct {
		Items []*items.Item
	}{
		Items: elems,
	})
	if err != nil {
		h.Logger.Error("ExecuteTemplate err", err)
		http.Error(w, `Template errror`, http.StatusInternalServerError)
		return
	}
}

func (h *ItemsHandler) AddForm(w http.ResponseWriter, r *http.Request) {
	err := h.Tmpl.ExecuteTemplate(w, "create.html", nil)
	if err != nil {
		h.Logger.Error("ExecuteTemplate err", err)
		http.Error(w, `Template errror`, http.StatusInternalServerError)
		return
	}
}

// type ItemsAddInput struct {

// }

func (h *ItemsHandler) Add(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	item := new(items.Item)
	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)
	err := decoder.Decode(item, r.PostForm)
	if err != nil {
		h.Logger.Error("Form err", err)
		http.Error(w, `Bad form`, http.StatusBadRequest)
		return
	}

	sess, _ := session.SessionFromContext(r.Context())
	// item.UserID = sess.UserID

	lastID, err := h.ItemsRepo.Add(item)
	if err != nil {
		h.Logger.Error("Db err", err)
		http.Error(w, `DB err`, http.StatusInternalServerError)
		return
	}
	h.Logger.Infof("Insert with id LastInsertId: %v %v", lastID, sess)
	http.Redirect(w, r, "/", http.StatusFound)
}

func (h *ItemsHandler) Edit(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, `{"error": "bad id"}`, http.StatusBadGateway)
		return
	}

	item, err := h.ItemsRepo.GetByID(int64(id))
	if err != nil {
		h.Logger.Error("Db err", err)
		http.Error(w, `DB err`, http.StatusInternalServerError)
		return
	}
	if item == nil {
		http.Error(w, `no item`, http.StatusNotFound)
		return
	}

	err = h.Tmpl.ExecuteTemplate(w, "edit.html", item)
	if err != nil {
		h.Logger.Error("ExecuteTemplate err", err)
		http.Error(w, `Template errror`, http.StatusInternalServerError)
		return
	}
}

func (h *ItemsHandler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, `Bad id`, http.StatusBadRequest)
		return
	}

	r.ParseForm()
	item := new(items.Item)
	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)
	err = decoder.Decode(item, r.PostForm)
	if err != nil {
		h.Logger.Error("Form err", err)
		http.Error(w, `Bad form`, http.StatusBadRequest)
		return
	}
	item.ID = uint32(id)

	sess, _ := session.SessionFromContext(r.Context())
	item.SetUpdated(sess.UserID)

	ok, err := h.ItemsRepo.Update(item)
	if err != nil {
		h.Logger.Error("Db err", err)
		http.Error(w, `db error`, http.StatusInternalServerError)
		return
	}

	h.Logger.Infof("update: %v %v", item, ok)

	http.Redirect(w, r, "/", http.StatusFound)
}

func (h *ItemsHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, `{"error": "bad id"}`, http.StatusBadGateway)
		return
	}

	affected, err := h.ItemsRepo.Delete(int64(id))
	if err != nil {
		h.Logger.Error("Db err", err)
		http.Error(w, `{"error": "db error"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-type", "application/json")
	respJSON, _ := json.Marshal(map[string]int64{
		"updated": affected,
	})
	w.Write(respJSON)
}
