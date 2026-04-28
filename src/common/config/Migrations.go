package config

import (
	"fmt"
	"gestrym-ai/src/common/models"
	"gestrym-ai/src/common/utils"
)

var logger = utils.NewLogger()

func MigrateDB() (IDatabaseConnection, error) {
	connection := NewPostgresConnection()
	db := connection.GetDB()

	// Register all models for AutoMigrate
	err := db.AutoMigrate(
		&models.AIRecommendation{},
	)

	if err != nil {
		logger.Error(fmt.Sprintf("[ERROR] Error al migrar las entidades: %s", err.Error()))
		return nil, err
	}

	logger.Info("[OK] Todas las migraciones completadas exitosamente")
	return connection, nil
}
