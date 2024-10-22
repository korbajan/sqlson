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


type PostgresExecutor struct {
  host string
  port int
  user string
  password string
  databaseName string
  db *gorm.DB
  dbVersion string
}


const (
  DefaultHost = "localhost"
  DefaultPort = 5432
  DefaultUser = "postgres"
)

func NewExecutor(databaseConfig *configs.Database) *PostgresExecutor {
  host := databaseConfig.Host
  if databaseConfig.Host == "" {
    host = DefaultHost
  }
  port := databaseConfig.Port
  if databaseConfig.Port == 0 {
    port = DefaultPort
  }
  user := databaseConfig.User
  if databaseConfig.User == "" {
    host = DefaultHost
  }
  
  return &PostgresExecutor{
    host: host,
    port: port,
    user: user,
    databaseName: databaseConfig.Name,
    password: databaseConfig.Password,
  }
  
}

func (pg *PostgresExecutor) GetDSN() string {
  return fmt.Sprintf(
    "host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
    pg.host,
    pg.user,
    pg.password,
    pg.databaseName,
    pg.port,
  )
}

func (pg *PostgresExecutor) PrepareDBConnection() error {
  db, err := gorm.Open(postgres.Open(pg.GetDSN()), &gorm.Config{
    Logger: logger.Discard, // Silences all gorm logs
  })
  if err != nil {
    if strings.Contains(err.Error(), "SQLSTATE 3D000") {
      return dberrors.NewDBCheckTypeError(err)
    }
    if strings.Contains(err.Error(), "SQLSTATE 28P01") {
      return dberrors.NewDBCheckTypeError(err)
    }
    return err
  }
  pg.db = db
  return nil
}

func ParseToJsonAggQuery(sqlQuery string) string {
  return fmt.Sprintf("SELECT json_agg(t) FROM (%s) t;", sqlQuery)
}

func (pg *PostgresExecutor) Execute(sqlQuery string) (string, error) {
  if pg.db == nil {
    err := pg.PrepareDBConnection()
    if err != nil {
      return "", err
    }
  }
  var jsonResult string // Use Raw to execute the query and scan the result
  if err := pg.db.Raw(ParseToJsonAggQuery(sqlQuery)).Scan(&jsonResult).Error; err != nil {
    return "", err
  }
  return jsonResult, nil
}

func (pg *PostgresExecutor) GetVersion() string {
  pg.db.Raw("SELECT version();").Scan(&pg.dbVersion)
  return pg.dbVersion
}

