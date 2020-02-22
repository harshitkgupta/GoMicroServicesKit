package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"text/template"
	"time"

	"github.com/harshitkgupta/welcome-app/chain"
)

func getPort() string {
	p := os.Getenv("PORT")
	if p != "" {
		return ":" + p
	}
	return ":8080"
}

type Welcome struct {
	Name string
	Time string
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

func welcomeHandler(w http.ResponseWriter, r *http.Request) {
	welcome := Welcome{"Anonymous", time.Now().Format(time.Stamp)}
	if name := r.FormValue("name"); name != "" {
		welcome.Name = name
	}
	templates := template.Must(template.ParseFiles("templates/welcome-template.html"))

	if err := templates.ExecuteTemplate(w, "welcome-template.html", welcome); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

type RadioButton struct {
	Name       string
	Value      string
	IsDisabled bool
	IsChecked  bool
	Text       string
}

type PageVariables struct {
	PageTitle        string
	PageRadioButtons []RadioButton
	Answer           string
}

func GetViolinHandler(w http.ResponseWriter, r *http.Request) {
	// Display some radio buttons to the user

	Title := "Which do you prefer?"
	MyRadioButtons := []RadioButton{
		RadioButton{"animalselect", "cats", false, false, "Cats"},
		RadioButton{"animalselect", "dogs", false, false, "Dogs"},
	}

	MyPageVariables := PageVariables{
		PageTitle:        Title,
		PageRadioButtons: MyRadioButtons,
	}

	t, err := template.ParseFiles("templates/violin.html") //parse the html file homepage.html
	if err != nil {                                        // if there is an error
		log.Print("template parsing error: ", err) // log it
	}

	err = t.Execute(w, MyPageVariables) //execute the template and pass it the HomePageVars struct to fill in the gaps
	if err != nil {                     // if there is an error
		log.Print("template executing error: ", err) //log it
	}

}

func ShowViolinHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	// r.Form is now either
	// map[animalselect:[cats]] OR
	// map[animalselect:[dogs]]
	// so get the animal which has been selected
	youranimal := r.Form.Get("animalselect")

	Title := "Your preferred animal"
	MyPageVariables := PageVariables{
		PageTitle: Title,
		Answer:    youranimal,
	}

	// generate page by passing page variables into template
	t, err := template.ParseFiles("templates/violin.html") //parse the html file homepage.html
	if err != nil {                                        // if there is an error
		log.Print("template parsing error: ", err) // log it
	}

	err = t.Execute(w, MyPageVariables) //execute the template and pass it the HomePageVars struct to fill in the gaps
	if err != nil {                     // if there is an error
		log.Print("template executing error: ", err) //log it
	}
}

func wrapHandler(h http.Handler) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)
	}
}
func main() {
	srv := &http.Server{
		Addr:         "0.0.0.0" + getPort(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	commonHandlers := chain.New(loggingHandler, recoverHandler)
	//http.HandleFunc("/", commonHandlers.ThenFunc(welcomeHandler)))
	http.Handle("/static/", http.StripPrefix("/static/",
		http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", wrapHandler(commonHandlers.ThenFunc(welcomeHandler)))
	http.HandleFunc("/getViolin", wrapHandler(commonHandlers.ThenFunc(GetViolinHandler)))
	http.HandleFunc("/showViolin", wrapHandler(commonHandlers.ThenFunc(ShowViolinHandler)))

	fmt.Println("Listening")
	fmt.Println(srv.ListenAndServe())

}
