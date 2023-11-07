package models

import (
	"time"

	"github.com/jclumbiarres/gcm/libz"
	"gorm.io/gorm"
)

type Alumno struct {
	gorm.Model
	Nombre      string
	FechaInicio time.Time
	FechaFin    time.Time
	FIniStr     string `gorm:"-"`
	FFinStr     string `gorm:"-"`
	Activo      bool
	Comentario  string
	Sprints     []Sprint
}

type Sprint struct {
	gorm.Model
	AlumnoID    uint
	Numero      int
	Nombre      string
	FechaInicio time.Time
	FechaFin    time.Time
	FechaCorre  time.Time
	Nota        float32
	Activo      bool
	Completado  bool
	Comentario  string
	FIniStr     string `gorm:"-"`
	FFinStr     string `gorm:"-"`
	FCorStr     string `gorm:"-"`
}

// GetAlumnos devuelve una lista de todos los alumnos en la base de datos,
// incluyendo información sobre los sprints en los que están inscritos.
func (a Alumno) GetAlumnos() (alumnos []Alumno) {
	libz.DB.Preload("Sprints").Find(&alumnos)
	return alumnos
}

// GetSprintById devuelve un sprint por su id y el id del alumno al que pertenece.
func (a Alumno) GetSprintById(id int, alumnoid int) (sprint Sprint, err error) {
	err = libz.DB.Where("numero = ? AND alumno_id = ?", id, alumnoid).First(&sprint).Error
	if err != nil {
		return sprint, err
	}
	return sprint, nil
}

// ModificarAlumno actualiza los datos de un alumno en la base de datos.
func (a Alumno) ModificarAlumno() (err error) {
	libz.DB.Preload("Sprints").Save(&a)
	return
}

// GetAlumno devuelve un objeto Alumno con el id especificado y sus sprints asociados.
// Si no se encuentra ningún alumno con el id especificado, se devuelve un objeto Alumno vacío.
func (a *Alumno) GetAlumno(id int) (alumno Alumno) {
	libz.DB.Preload("Sprints").Find(&alumno, id)
	return
}

// CreateAlumno crea un nuevo registro de alumno en la base de datos.
// Recibe un puntero a la estructura Alumno y devuelve un error en caso de haberlo.
func (a *Alumno) CreateAlumno(alumno Alumno) (err error) {
	libz.DB.Create(&alumno)
	return
}

// BuscarAlumno busca un alumno por su ID en la base de datos y devuelve el resultado.
// Parámetros:
// - id: el ID del alumno a buscar.
// Retorna:
// - alumno: el alumno encontrado en la base de datos.
func BuscarAlumno(id int) (alumno Alumno) {
	libz.DB.Preload("Sprints").Find(&alumno, id)
	return
}

// UpdateSprint actualiza un sprint del alumno en la base de datos.
// Recibe como parámetro un objeto Sprint y lo busca en la lista de sprints del alumno.
// Si lo encuentra, lo actualiza y guarda los cambios en la base de datos.
// Retorna un error en caso de que ocurra algún problema al guardar los cambios.
func (a *Alumno) UpdateSprint(sprint Sprint) (err error) {
	for i, s := range a.Sprints {
		if s.ID == sprint.ID {
			a.Sprints[i] = sprint
			err = libz.DB.Save(&a).Error
			return
		}
	}
	return err
}

func UpdateSingleSprint(alumnoID uint, sprintID uint, sprintData Sprint) (err error) {
	var alumno Alumno
	if err = libz.DB.Preload("Sprints").First(&alumno, alumnoID).Error; err != nil {
		return
	}

	for i, s := range alumno.Sprints {
		if s.ID == sprintID {
			alumno.Sprints[i].Nombre = sprintData.Nombre
			alumno.Sprints[i].FechaInicio = sprintData.FechaInicio
			alumno.Sprints[i].FechaFin = sprintData.FechaFin
			alumno.Sprints[i].FechaCorre = sprintData.FechaCorre
			alumno.Sprints[i].Nota = sprintData.Nota
			alumno.Sprints[i].Activo = sprintData.Activo
			alumno.Sprints[i].Completado = sprintData.Completado
			alumno.Sprints[i].Comentario = sprintData.Comentario
			return
		}
	}
	err = libz.DB.Save(&alumno).Error
	return
}

/* Refactorizados */
/* ---------------------------------------------------------------------------*/
func (a *Alumno) ModificarSprint(sprint *Sprint) {
	if err := libz.DB.Where("nombre = ?", a.Nombre).Preload("Sprints").First(a).Error; err != nil {
		panic(err)
	}
	err := libz.DB.Transaction(func(tx *gorm.DB) error {
		for i := range a.Sprints {
			if a.Sprints[i].Numero == sprint.Numero {
				a.Sprints[i].Completado = sprint.Completado
				a.Sprints[i].Activo = sprint.Activo
				a.Sprints[i].Comentario = sprint.Comentario
				a.Sprints[i].Nota = sprint.Nota
				a.Sprints[i].FechaCorre = sprint.FechaCorre
				if err := tx.Save(&a.Sprints[i]).Error; err != nil {
					return err
				}
				break
			}
		}
		return nil
	})

	if err != nil {
		panic(err)
	}
}

func (a *Alumno) ModificarComentario(comentario string) (err error) {
	err = libz.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&a).Update("comentario", comentario).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

/* ---------------------------------------------------------------------------*/
/* Fin refactorizados */

// CrearAlumno agrega un nuevo alumno a la base de datos con los sprints iniciales.
// Retorna un error en caso de que la creación del alumno falle.
func (a Alumno) CrearAlumno() (err error) {
	a.Sprints = append(a.Sprints, Sprint{
		AlumnoID:    a.ID,
		Numero:      1,
		Nombre:      "Sprint 1",
		FechaInicio: time.Now(),
		FechaFin:    time.Now().Add(24 * time.Hour * 7),
		FechaCorre:  time.Now().Add(24 * time.Hour * 7),
		Nota:        0,
		Activo:      true,
		Completado:  false,
		Comentario:  "En curso",
	}, Sprint{
		AlumnoID:   a.ID,
		Numero:     2,
		Nombre:     "Sprint 2",
		Nota:       0,
		Activo:     false,
		Completado: false,
		Comentario: "",
	}, Sprint{
		AlumnoID: a.ID,

		Numero:     3,
		Nombre:     "Sprint 3",
		Nota:       0,
		Activo:     false,
		Completado: false,
		Comentario: "",
	}, Sprint{
		AlumnoID: a.ID,

		Numero:     9,
		Nombre:     "Prueba Nivel",
		Nota:       0,
		Activo:     false,
		Completado: false,
		Comentario: "",
	}, Sprint{
		AlumnoID: a.ID,

		Numero:     4,
		Nombre:     "Sprint 4",
		Nota:       0,
		Activo:     false,
		Completado: false,
		Comentario: "",
	}, Sprint{
		AlumnoID: a.ID,

		Numero:     5,
		Nombre:     "Sprint 5",
		Nota:       0,
		Activo:     false,
		Completado: false,
		Comentario: "",
	}, Sprint{
		AlumnoID: a.ID,

		Numero:     6,
		Nombre:     "Sprint 6",
		Nota:       0,
		Activo:     false,
		Completado: false,
		Comentario: "",
	}, Sprint{
		AlumnoID: a.ID,

		Numero:     7,
		Nombre:     "Sprint 7",
		Nota:       0,
		Activo:     false,
		Completado: false,
		Comentario: "",
	}, Sprint{
		AlumnoID: a.ID,

		Numero:     8,
		Nombre:     "Sprint 8",
		Nota:       0,
		Activo:     false,
		Completado: false,
		Comentario: "",
	})

	libz.DB.Create(&a)
	return
}
