package controllers

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jclumbiarres/gcm/libz"
	"github.com/jclumbiarres/gcm/models"
)

func AlumnoInfo(w http.ResponseWriter, r *http.Request) {
	var salida models.Alumno
	tmpl, err := template.ParseFiles("views/sprints/listado.html")
	if err != nil {
		log.Fatal(err)
	}

	numero := chi.URLParam(r, "numero")
	numeroInt, err := strconv.Atoi(numero)
	if err != nil {
		log.Fatal(err)
	}

	// Loop through all sprints and format the date strings
	alumno := salida.GetAlumno(numeroInt)
	for i := range alumno.Sprints {
		alumno.Sprints[i].FCorStr = formatAsDate(alumno.Sprints[i].FechaCorre)
		alumno.Sprints[i].FFinStr = formatAsDate(alumno.Sprints[i].FechaFin)
	}

	data := struct {
		Alumno models.Alumno
		Merda  int
	}{
		Alumno: alumno,
		Merda:  numeroInt,
	}

	tmpl.Execute(w, data)
}

func EditarAlumno(w http.ResponseWriter, r *http.Request) {
	var salida models.Alumno
	var sprints models.Sprint
	tmpl, err := template.ParseFiles("views/sprints/editar.html")
	if err != nil {
		log.Fatal(err)
	}

	// Get the `numero` and `sprint` variables from the URL path
	numero := chi.URLParam(r, "numero")
	sprint := chi.URLParam(r, "sprint")

	// Convert the `numero` and `sprint` variables to integers
	numeroInt, err := strconv.Atoi(numero)
	if err != nil {
		log.Fatal(err)
	}
	sprintInt, err := strconv.Atoi(sprint)
	if err != nil {
		log.Fatal(err)
	}

	// Query the database for the corresponding sprint
	sprints, errs := salida.GetSprintById(sprintInt, numeroInt)
	if errs != nil {
		log.Fatal(err)
	}
	sprints.FFinStr = formatAsDate(sprints.FechaFin)
	sprints.FCorStr = formatAsDate(sprints.FechaCorre)
	// Create a data structure to pass to the template
	data := struct {
		Alumno    models.Alumno
		Sprint    models.Sprint
		Numero    int
		SprintNum int
	}{
		Alumno:    salida.GetAlumno(numeroInt),
		Sprint:    sprints,
		Numero:    numeroInt,
		SprintNum: sprintInt,
	}

	// Execute the template with the data
	tmpl.Execute(w, data)

}

func AlumnosHandler(w http.ResponseWriter, r *http.Request) {

	tmpl, err := template.ParseFiles("views/alumnos/listado.html")
	if err != nil {
		log.Fatal(err)
	}
	salida := models.Alumno{}.GetAlumnos()
	// range salida sprints
	for i := range salida {
		for j := range salida[i].Sprints {
			salida[i].Sprints[j].FIniStr = formatAsDate(salida[i].Sprints[j].FechaInicio)
			salida[i].Sprints[j].FCorStr = formatAsDate(salida[i].Sprints[j].FechaCorre)
			salida[i].Sprints[j].FFinStr = formatAsDate(salida[i].Sprints[j].FechaFin)
		}
	}

	tmpl.Execute(w, salida)
}

func AnadirAlumno(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("views/alumnos/anadir.html")
	if err != nil {
		log.Fatal(err)
	}
	tmpl.Execute(w, nil)
}

func AnadirAlumnoDb(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		println(err.Error())
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}
	nombre := r.FormValue("nombre")
	fechaInicio := r.FormValue("fechaInicio")
	fechaFin := r.FormValue("fechaFin")
	fechaFinSprint1 := r.FormValue("fechaFinSprint1")
	fechaFinSprint2 := r.FormValue("fechaFinSprint2")
	fechaFinSprint3 := r.FormValue("fechaFinSprint3")
	fechaFinSprint4 := r.FormValue("fechaFinSprint4")
	fechaFinSprint5 := r.FormValue("fechaFinSprint5")
	fechaFinSprint6 := r.FormValue("fechaFinSprint6")
	fechaFinSprint7 := r.FormValue("fechaFinSprint7")
	fechaFinSprint8 := r.FormValue("fechaFinSprint8")
	dates := []string{fechaFinSprint1, fechaFinSprint2, fechaFinSprint3, fechaFinSprint4, fechaFinSprint5, fechaFinSprint6, fechaFinSprint7, fechaFinSprint8, fechaInicio, fechaFin}
	var fInicio, fFin time.Time
	var ffechaFinSprints []time.Time

	for _, date := range dates {
		ffechaFinSprint, err := time.Parse("2006-01-02", date)
		if err != nil {
			println(err.Error())
			http.Error(w, "Error parsing form", http.StatusBadRequest)
			return
		}
		ffechaFinSprints = append(ffechaFinSprints, ffechaFinSprint)
	}

	fInicio = ffechaFinSprints[len(ffechaFinSprints)-2]
	fFin = ffechaFinSprints[len(ffechaFinSprints)-1]
	ffechaFinSprints = ffechaFinSprints[:len(ffechaFinSprints)-2]

	creacion := models.Alumno{
		Nombre:      nombre,
		FechaInicio: fInicio,
		FechaFin:    fFin,
		Comentario:  "Recien iniciado",
		Activo:      true,
	}
	// Hay que hacer un loop para esto también, en la refactorización se hará, por ahora el proof of concept.
	creacion.Sprints = append(creacion.Sprints, models.Sprint{
		Numero:      1,
		Nombre:      "Sprint 1",
		FechaInicio: time.Now(),
		FechaFin:    ffechaFinSprints[0],
		Nota:        0,
		Activo:      true,
		Completado:  false,
		Comentario:  "En curso",
	}, models.Sprint{
		Numero:     2,
		Nombre:     "Sprint 2",
		Nota:       0,
		Activo:     false,
		FechaFin:   ffechaFinSprints[1],
		Completado: false,
		Comentario: "",
	}, models.Sprint{
		Numero:     3,
		Nombre:     "Sprint 3",
		Nota:       0,
		Activo:     false,
		FechaFin:   ffechaFinSprints[2],
		Completado: false,
		Comentario: "",
	}, models.Sprint{
		Numero:     9,
		Nombre:     "Prueba Nivel",
		Nota:       0,
		FechaFin:   time.Now().AddDate(10, 0, 0),
		Activo:     false,
		Completado: false,
		Comentario: "",
	}, models.Sprint{
		Numero:     4,
		Nombre:     "Sprint 4",
		Nota:       0,
		Activo:     false,
		FechaFin:   ffechaFinSprints[3],
		Completado: false,
		Comentario: "",
	}, models.Sprint{
		Numero:     5,
		Nombre:     "Sprint 5",
		Nota:       0,
		Activo:     false,
		FechaFin:   ffechaFinSprints[4],
		Completado: false,
		Comentario: "",
	}, models.Sprint{
		Numero:     6,
		Nombre:     "Sprint 6",
		Nota:       0,
		Activo:     false,
		FechaFin:   ffechaFinSprints[5],
		Completado: false,
		Comentario: "",
	}, models.Sprint{
		Numero:     7,
		Nombre:     "Sprint 7",
		Nota:       0,
		FechaFin:   ffechaFinSprints[6],
		Activo:     false,
		Completado: false,
		Comentario: "",
	}, models.Sprint{
		Numero:     8,
		Nombre:     "Sprint 8",
		Nota:       0,
		Activo:     false,
		FechaFin:   ffechaFinSprints[7],
		Completado: false,
		Comentario: "",
	})
	libz.DB.Create(&creacion)
	w.Write([]byte("Alumno añadido"))
}
