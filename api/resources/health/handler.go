package health

import "net/http"

func Read(w http.ResponseWriter, r *http.Request) {
	_, e := w.Write([]byte("."))
	if e != nil {
		panic(e)
	}
}
