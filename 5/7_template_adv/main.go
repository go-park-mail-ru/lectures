package main

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/go-park-mail-ru/lectures/5/7_template_adv/item"
	"github.com/go-park-mail-ru/lectures/5/7_template_adv/template"
)

//go:generate hero -source=./template/

var ExampleItems = []*item.Item{
	&item.Item{1, "rvasily", "Mail.ru Group -> Mail & Portal -> Mail"},
	&item.Item{2, "dmitrydorofeev", "Mail.ru Group -> Mail & Portal -> BeepCar"},
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		buffer := new(bytes.Buffer)
		template.Index(ExampleItems, buffer)
		w.Write(buffer.Bytes())
	})

	fmt.Println("starting server at :8080")
	http.ListenAndServe(":8080", nil)
}
