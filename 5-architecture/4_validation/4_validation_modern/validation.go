package main

import (
	"errors"
	"fmt"
	"net/http"
	"unicode/utf8"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/schema"
)

// curl -v 'http://127.0.0.1:8080/?to=v.romanov@corp.mail.ru&priority=low&subject=Hello!&inner=ignored&id=12&flag=23'
// curl -i 'http://127.0.0.1:8080/?to=bad&priority=urgent&subject=HelloWorld!'
type SendMessage struct {
	Id        int    `validate:"omitempty"`
	Priority  string `validate:"required,oneof=low normal high"`
	Recipient string `schema:"to" validate:"required,email"`
	Subject   string `validate:"required,msgSubject"`
	Inner     string `schema:"-" validate:"-"`
	flag      int
}

var validate = validator.New(validator.WithRequiredStructEnabled())

func msgSubject(field validator.FieldLevel) bool {
	subject := field.Field().String()
	length := utf8.RuneCountInString(subject)
	return length >= 1 && length <= 10
}

func handler(w http.ResponseWriter, r *http.Request) {
	msg := &SendMessage{}

	// Парсинг параметров остался тем же.
	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)
	if err := decoder.Decode(msg, r.URL.Query()); err != nil {
		http.Error(w, "invalid query parameters", http.StatusBadRequest)
		return
	}

	// Затем проверяем правила, заданные тегами validate.
	if err := validate.Struct(msg); err != nil {
		var allErrs validator.ValidationErrors
		if !errors.As(err, &allErrs) {
			fmt.Println(err)
			http.Error(w, "internal", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "request %s\n\nMsg: %#v\n\n", r.URL, msg)
		for _, field := range allErrs {
			fmt.Fprintf(w, "field: %s, rule: %s, param: %s\n\n",
				field.Field(), field.Tag(), field.Param())
		}
		fmt.Fprintf(w, "error: %s\n\n", err)
		return
	}

	fmt.Fprintf(w, "request %s\n\nMsg: %#v\n\nmsg is correct\n\n", r.URL, msg)
}

func main() {
	if err := validate.RegisterValidation("msgSubject", msgSubject); err != nil {
		fmt.Println(err)
		return
	}

	http.HandleFunc("/", handler)
	fmt.Println("starting server at :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println(err)
	}
}
