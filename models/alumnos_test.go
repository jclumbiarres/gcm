package models_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/jclumbiarres/gcm/models"
)

func TestAlumnoModel(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("../test.sqlite"))
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}
	db.AutoMigrate(&models.Alumno{})

	alumno := models.Alumno{
		Nombre:      "Juan",
		FechaInicio: time.Now(),
		FechaFin:    time.Now().AddDate(0, 0, 7),
		Activo:      true,
		Comentario:  "Comentario de prueba",
	}

	// Create
	result := db.Create(&alumno)
	assert.NoError(t, result.Error, "should not return an error")

	// Read
	var alumnoDB models.Alumno
	result = db.First(&alumnoDB, alumno.ID)
	assert.NoError(t, result.Error, "should not return an error")
	assert.Equal(t, alumno.Nombre, alumnoDB.Nombre, "should be equal")
	assert.Equal(t, alumno.FechaInicio.Unix(), alumnoDB.FechaInicio.Unix(), "should be equal")
	assert.Equal(t, alumno.FechaFin.Unix(), alumnoDB.FechaFin.Unix(), "should be equal")
	assert.Equal(t, alumno.Activo, alumnoDB.Activo, "should be equal")
	assert.Equal(t, alumno.Comentario, alumnoDB.Comentario, "should be equal")

	// Update
	alumno.Nombre = "Pedro"
	result = db.Save(&alumno)
	assert.NoError(t, result.Error, "should not return an error")

	// Delete
	result = db.Delete(&alumno)
	assert.NoError(t, result.Error, "should not return an error")
}
