package main

import (
	"encoding/json"
	"fmt"
	"github.com/microcosm-cc/bluemonday"
	"net/http"
)

// для санитайзинга на сторое фронта используйте https://github.com/cure53/DOMPurify

func main() {
	sanitizer := bluemonday.UGCPolicy()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		comment := `<a onblur="alert(document.сookie)" href="https://www.mail.ru">Mail.ru</a>`
		comment = sanitizer.Sanitize(comment)
		resp, _ := json.Marshal(map[string]interface{}{
			"comment": comment,
		})
		w.Write(resp)
	})

	fmt.Println("starting server at :8080")
	http.ListenAndServe(":8080", nil)
}
