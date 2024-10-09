package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"regexp"
	"sync"
)

// Structures de données
type Student struct {
	FirstName string
	LastName  string
	Age       int
	Gender    string
}

type Class struct {
	Name         string
	Department   string
	Level        string
	StudentCount int
	Students     []Student
}

type UserData struct {
	FirstName string
	LastName  string
	BirthDate string
	Gender    string
}

var (
	viewCount = 0
	userData  UserData
	mu        sync.Mutex
	templates *template.Template
)

func main() {
	// Initialisation des templates
	var err error
	templates, err = template.ParseGlob("templates/*.html")
	if err != nil {
		log.Fatal(err)
	}

	// Routes
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/promo", promoHandler)
	http.HandleFunc("/change", changeHandler)
	http.HandleFunc("/user/form", userFormHandler)
	http.HandleFunc("/user/treatment", userTreatmentHandler)
	http.HandleFunc("/user/display", userDisplayHandler)

	// Servir les fichiers statiques
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	fmt.Println("Serveur démarré sur http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	err := templates.ExecuteTemplate(w, "home.html", nil)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func promoHandler(w http.ResponseWriter, r *http.Request) {
	class := Class{
		Name:         "B1 Informatique",
		Department:   "Informatique",
		Level:        "Bachelor 1",
		StudentCount: 3,
		Students: []Student{
			{FirstName: "Jean", LastName: "Dupont", Age: 20, Gender: "M"},
			{FirstName: "Marie", LastName: "Martin", Age: 19, Gender: "F"},
			{FirstName: "Paul", LastName: "Bernard", Age: 21, Gender: "M"},
		},
	}

	err := templates.ExecuteTemplate(w, "promo.html", class)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func changeHandler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	viewCount++
	currentCount := viewCount
	mu.Unlock()

	data := struct {
		Count  int
		IsPair bool
	}{
		Count:  currentCount,
		IsPair: currentCount%2 == 0,
	}

	err := templates.ExecuteTemplate(w, "change.html", data)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func userFormHandler(w http.ResponseWriter, r *http.Request) {
	err := templates.ExecuteTemplate(w, "user_form.html", nil)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func userTreatmentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/user/form", http.StatusSeeOther)
		return
	}

	firstName := r.FormValue("firstName")
	lastName := r.FormValue("lastName")
	birthDate := r.FormValue("birthDate")
	gender := r.FormValue("gender")

	// Validation
	nameRegex := regexp.MustCompile("^[a-zA-Z]{1,32}$")
	if !nameRegex.MatchString(firstName) || !nameRegex.MatchString(lastName) {
		http.Redirect(w, r, "/user/form", http.StatusSeeOther)
		return
	}

	if gender != "masculin" && gender != "féminin" && gender != "autre" {
		http.Redirect(w, r, "/user/form", http.StatusSeeOther)
		return
	}

	userData = UserData{
		FirstName: firstName,
		LastName:  lastName,
		BirthDate: birthDate,
		Gender:    gender,
	}

	http.Redirect(w, r, "/user/display", http.StatusSeeOther)
}

func userDisplayHandler(w http.ResponseWriter, r *http.Request) {
	err := templates.ExecuteTemplate(w, "user_display.html", userData)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
