package databases

import (
	"fmt"
  "errors"

	"github.com/korbajan/sqlson/internal/configs"
	"github.com/korbajan/sqlson/pkg/databases/dberrors"
	"github.com/korbajan/sqlson/pkg/databases/mysql"
	"github.com/korbajan/sqlson/pkg/databases/postgres"
)

type ExecutorConnection interface {
  GetVersion() string
  RunQuery(string) (string, error)
}

func NewExecutorConnection(databaseConfig *configs.Database) (ExecutorConnection, error) { // or perhaps a generic instead of interface?
  
  var getExecutorError dberrors.NewExecutorError
  var c ExecutorConnection
  // Try connecting to PostgreSQL
  c, err := postgres.NewExecutor(databaseConfig) 
  if err == nil {
    return c, nil
  }
  if errors.As(err, &getExecutorError) {
    return nil, err
  }

  // If it fails, try connecting to MySQL/MariaDB
  c, err = mysql.NewExecutor(databaseConfig) 
  if err == nil {
    return c, nil
  }
  return nil, fmt.Errorf("could not determine database type: %v", err)
}

