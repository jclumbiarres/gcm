package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jclumbiarres/gcm/controllers"
	"github.com/jclumbiarres/gcm/libz"
	"github.com/jclumbiarres/gcm/models"
)

func main() {
	libz.InitDB()
	libz.DB.AutoMigrate(&models.Alumno{}, &models.Sprint{})
	// bananas := models.Alumno{
	// 	Nombre:      "Joe Bananas",
	// 	FechaInicio: time.Now(),
	// 	FechaFin:    time.Now().Add(24 * time.Hour * 7),
	// 	Activo:      true,
	// 	Comentario:  "Abuelete",
	// }
	// bananas.CrearAlumno()

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	fs := http.FileServer(http.Dir("assets"))
	r.Handle("/assets/*", http.StripPrefix("/assets/", fs))
	r.Get("/", controllers.RootHandler)
	r.Mount("/debug", middleware.Profiler())
	r.Get("/alumnos", controllers.AlumnosHandler)
	r.Get("/alumnos/add", controllers.AnadirAlumno)
	r.Post("/alumnos/add", controllers.AnadirAlumnoDb)
	r.Get("/alumnos/{numero}", controllers.AlumnoInfo)
	r.Get("/alumnos/{numero}/edit/{sprint}", controllers.EditarAlumno)
	r.Put("/alumnos/{numero}/edit/{sprint}", controllers.ActualizarSprint)
	r.Get("/home", controllers.HomeHandler)
	r.NotFound(controllers.NotFoundHandler)
	http.ListenAndServe(":9000", r)
}
