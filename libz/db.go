package libz

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// DB es una variable global que representa una conexión a la base de datos utilizando GORM.
var DB *gorm.DB

// InitDB abre una conexión a la base de datos SQLite y la asigna a la variable global DB.
// Si no se puede establecer la conexión, se lanza un error.
func InitDB() {
	var err error
	DB, err = gorm.Open(sqlite.Open("datos.sqlite"))
	if err != nil {
		panic("Fallo al conectarse a la base de datos")
	}

}

// CloseDB cierra la conexión a la base de datos.
// Si hay un error al cerrar la conexión, se lanza un pánico.
func CloseDB() {
	sqlDB, err := DB.DB()
	if err != nil {
		panic("Fallo al cerrar la base de datos")
	}
	sqlDB.Close()
}
