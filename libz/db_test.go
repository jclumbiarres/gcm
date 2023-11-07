package libz_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	"github.com/jclumbiarres/gcm/libz"
)

func TestInitDB(t *testing.T) {
	libz.InitDB()
	defer libz.CloseDB()

	assert.NotNil(t, libz.DB, "DB should not be nil")
	assert.IsType(t, &gorm.DB{}, libz.DB, "DB should be a *gorm.DB")
}

func TestCloseDB(t *testing.T) {
	libz.InitDB()
	defer libz.CloseDB()

	sqlDB, err := libz.DB.DB()
	assert.NoError(t, err, "should not return an error")

	libz.CloseDB()

	err = sqlDB.Ping()
	assert.Error(t, err, "should return an error after closing DB")
}
