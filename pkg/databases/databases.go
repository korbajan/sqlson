package databases

import (
	"fmt"
	"log"
  "errors"

	"github.com/korbajan/sqlson/internal/configs"
	"github.com/korbajan/sqlson/pkg/databases/dberrors"
	"github.com/korbajan/sqlson/pkg/databases/mysql"
	"github.com/korbajan/sqlson/pkg/databases/postgres"
)

// DatabaseType represents the type of database
type DatabaseType int

const (
  Unknown DatabaseType = iota
  PostgreSQL
  MySQL
)

type QueryExecutor interface {
  GetDSN() string
  PrepareDBConnection() error
  GetVersion() string
  Execute(string) (string, error)
}

func CheckDatabaseType(postgresExecuter QueryExecutor, mysqlExecuter QueryExecutor) (DatabaseType, string, error) {
  var checkTypeError dberrors.DBCheckTypeError
  // Try connecting to PostgreSQL
  err := postgresExecuter.PrepareDBConnection()
  if err == nil {
    return PostgreSQL, postgresExecuter.GetVersion(), nil
  }
  if errors.As(err, &checkTypeError) {
    return PostgreSQL, "", err
  }

  // If it fails, try connecting to MySQL/MariaDB
  err = mysqlExecuter.PrepareDBConnection()
  if err == nil {
    return MySQL, mysqlExecuter.GetVersion(), nil
  }
  return Unknown, "", fmt.Errorf("could not determine database type: %v", err)
}

func Execute(databaseConfig *configs.Database, sqlQuery string) (string, error) {
  
  var executor QueryExecutor

  postgresExecuter := postgres.NewExecutor(databaseConfig)
  mysqlExecuter := mysql.NewExecutor(databaseConfig)

  dbType, _, err := CheckDatabaseType(postgresExecuter, mysqlExecuter)
  if err != nil {
    log.Fatal(err)
    return "", err
  }
 
  switch dbType {
    case PostgreSQL:
      executor = postgresExecuter
    case MySQL:
      executor = mysqlExecuter
    default:
      log.Fatal("Unknown database type.")
    }

  return executor.Execute(sqlQuery)
}
