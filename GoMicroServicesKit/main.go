package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome")
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	t1 := time.Now()
	fmt.Fprintf(w, "Welcome!")
	t2 := time.Now()
	log.Printf("[%s] %q %v\n", r.Method, r.URL.String(), t2.Sub(t1))
}

func loggingHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		t1 := time.Now()
		next.ServeHTTP(w, r)
		t2 := time.Now()
		log.Printf("[%s] %q %v\n", r.Method, r.URL.String(), t2.Sub(t1))
	}

	return http.HandlerFunc(fn)
}

func recoverHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("panic: %+v", err)
				http.Error(w, http.StatusText(500), 500)
			}
		}()

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

type appContext struct {
	db *sql.DB
}

func (c *appContext) authHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		authToken := r.Header.Get("Authorization")
		user, err := getUser(c.db, authToken)

		if err != nil {
			http.Error(w, http.StatusText(401), 401)
			return
		}

		context.Set(r, "user", user)
		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func (c *appContext) adminHandler(w http.ResponseWriter, r *http.Request) {
	user := context.Get(r, "user")
	// Maybe other operations on the database
	json.NewEncoder(w).Encode(user)
}

func wrapHandler(h http.Handler) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		context.Set(r, "params", ps)
		h.ServeHTTP(w, r)
	}
}
func (c *appContext) teaHandler(w http.ResponseWriter, r *http.Request) {
	params := context.Get(r, "params").(httprouter.Params)
	tea := getTea(c.db, params.ByName("id"))
	json.NewEncoder(w).Encode(tea)
}
func main() {
	db := sql.Open("postgres", "...")
	appC := appContext{db}
	srv := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	commonHandlers := chain.New(context.ClearHandler, loggingHandler, recoverHandler)

	http.HandleFunc("/", wrapHandler(commonHandlers.ThenFunc(aboutHandler)))
	http.HandleFunc("/index", wrapHandler(commonHandlers.ThenFunc(indexHandler)))
	http.Handle("/admin", wrapHandler(commonHandlers.Append(appC.authHandler).ThenFunc(adminHandler)))
	router.GET("/teas/:id", wrapHandler(commonHandlers.ThenFunc(appC.teaHandler)))

	srv.ListenAndServe()
}
