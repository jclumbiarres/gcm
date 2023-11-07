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

func ActualizarSprint(w http.ResponseWriter, r *http.Request) {
	var salida models.Alumno
	err := r.ParseForm()
	if err != nil {
		println(err.Error())
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}
	correccion := r.FormValue("fechaCorre")
	nota := r.FormValue("nota")
	comentario := r.FormValue("comentario")
	completado := r.FormValue("completado")
	fechaCorre, err := time.Parse("2006-01-02", correccion)
	if err != nil {
		log.Fatal(err)
	}
	notaFloat, err := strconv.ParseFloat(nota, 32)
	if err != nil {
		log.Fatal(err)
	}
	tmpl, err := template.ParseFiles("views/sprints/actualizado.html")
	if err != nil {
		log.Fatal(err)
	}

	// Get the `numero` and `sprint` variables from the URL path
	numero := chi.URLParam(r, "numero")
	sprint := chi.URLParam(r, "sprint")
	numero_int, err := strconv.Atoi(numero)
	if err != nil {
		log.Fatal(err)
	}
	salida = salida.GetAlumno(numero_int)
	// procesa := salida.GetAlumno(numero_int)
	sprint_int, err := strconv.Atoi(sprint)
	if err != nil {
		log.Fatal(err)
	}
	var completadoSprint bool
	log.Println(completado)
	if completado == "on" {
		completadoSprint = true
	} else {
		completadoSprint = false
	}
	log.Println(completadoSprint)
	var alumno models.Alumno
	if err := libz.DB.Preload("Sprints").First(&alumno, uint(numero_int)).Error; err != nil {
		panic(err)
	}
	sprint_modificado := &models.Sprint{
		Numero:     sprint_int,
		Completado: completadoSprint,
		Comentario: comentario,
		Nota:       float32(notaFloat),
		FechaCorre: fechaCorre,
	}
	alumno.ModificarSprint(sprint_modificado)
	tmpl.Execute(w, nil)
}
