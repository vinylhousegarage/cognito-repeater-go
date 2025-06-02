package rootroot

import (
	"net/http"
)

func NewRootHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/login", http.StatusFound)
}
