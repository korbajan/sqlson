package postgres

import (
	"fmt"
	"strings"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/korbajan/sqlson/internal/configs"
	"github.com/korbajan/sqlson/pkg/databases/dberrors"
)


type PostgresExecutorConnection struct {
  db *gorm.DB
}


const (
  DefaultHost = "localhost"
  DefaultPort = 5432
  DefaultUser = "postgres"
)

func NewExecutor(databaseConfig *configs.Database) (*PostgresExecutorConnection, error) {
  if databaseConfig.Host == "" {
    databaseConfig.Host = DefaultHost
  }
  if databaseConfig.Port == 0 {
    databaseConfig.Port = DefaultPort
  }
  if databaseConfig.User == "" {
    databaseConfig.User = DefaultUser
  }
  
  dsn := fmt.Sprintf(
    "host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
    databaseConfig.Host,
    databaseConfig.User,
    databaseConfig.Password,
    databaseConfig.Name,
    databaseConfig.Port,
  )
  
  db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
    Logger: logger.Discard, // Silences all gorm logs
  })

  //TODO: switch-case in separate common/generic function???
  if err != nil {
    if strings.Contains(err.Error(), "SQLSTATE 3D000") {
      return nil, dberrors.CreateNewExecutorError(err)
    }
    if strings.Contains(err.Error(), "SQLSTATE 28P01") {
      return nil, dberrors.CreateNewExecutorError(err)
    }
    return nil, err
  }
  
  return &PostgresExecutorConnection{db: db}, nil
  
}

func ParseToJsonAggQuery(sqlQuery string) string {
  return fmt.Sprintf("SELECT json_agg(t) FROM (%s) t;", sqlQuery)
}

func (c *PostgresExecutorConnection) RunQuery(sqlQuery string) (string, error) {
  var jsonResult string // Use Raw to execute the query and scan the result
  if err := c.db.Raw(ParseToJsonAggQuery(sqlQuery)).Scan(&jsonResult).Error; err != nil {
    return "", err
  }
  return jsonResult, nil
}

func (c *PostgresExecutorConnection) GetVersion() string {
  var version string
  c.db.Raw("SELECT version();").Scan(&version)
  return version
}

