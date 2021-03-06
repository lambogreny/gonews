package app

import (
	"context"
	"net/http"

	"github.com/acoshift/gonews/pkg/model"
)

// Mount mounts handlers to mux
func Mount(mux *http.ServeMux) {
	mux.Handle("/", fetchUser(http.HandlerFunc(index))) // list all news
	mux.Handle("/upload/", http.StripPrefix("/upload", http.FileServer(http.Dir("upload"))))
	mux.Handle("/news/", http.StripPrefix("/news", http.HandlerFunc(newsView))) // /news/:id

	mux.HandleFunc("/register", adminRegister)
	mux.HandleFunc("/login", adminLogin)

	adminMux := http.NewServeMux()
	adminMux.HandleFunc("/logout", adminLogout)
	adminMux.HandleFunc("/list", adminList)
	adminMux.HandleFunc("/create", adminCreate)
	adminMux.HandleFunc("/edit", adminEdit)

	mux.Handle("/admin/", http.StripPrefix("/admin", onlyAdmin(adminMux)))
}

func onlyAdmin(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sess := model.GetSession(r)
		ok, err := model.CheckUserID(sess.UserID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if !ok {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
		h.ServeHTTP(w, r)
	})
}

func fetchUser(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sess := model.GetSession(r)
		if sess.UserID == "" {
			h.ServeHTTP(w, r)
			return
		}
		username, err := model.GetUsernameFromID(sess.UserID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		ctx := context.WithValue(r.Context(), "username", username)
		nr := r.WithContext(ctx)
		h.ServeHTTP(w, nr)
	})
}
