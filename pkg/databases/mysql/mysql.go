package mysql

import (
  "encoding/json"
	"fmt"
  "strings"
  
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/driver/mysql"

	"github.com/korbajan/sqlson/internal/configs"
	"github.com/korbajan/sqlson/pkg/databases/dberrors"
)

type MysqlExecutorConnection struct {
  db *gorm.DB
}

const (
  DefaultHost = "localhost"
  DefaultPort = 3306
  DefaultUser = "root"
)

func NewExecutor(databaseConfig *configs.Database) (*MysqlExecutorConnection, error) {
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
    "%s:%s@tcp(%s:%d)/%s",
    databaseConfig.User,
    databaseConfig.Password,
    databaseConfig.Host,
    databaseConfig.Port,
    databaseConfig.Name,
  )

  db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
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
  
  return &MysqlExecutorConnection{db: db}, nil
}

// Function to execute a raw SQL query and return results in JSON format
func (c *MysqlExecutorConnection) RunQuery(sqlQuery string) (string, error) {
  // Execute the raw SQL query
  rows, err := c.db.Raw(sqlQuery).Rows()
  if err != nil {
    return "", err
  }
  defer rows.Close()

  // Get column names
  columns, err := rows.Columns()
  if err != nil {
    return "", err
  }

  // Prepare to scan the data into a slice of maps
  results := make([]map[string]interface{}, 0)

  // Iterate through the rows
  for rows.Next() {
    // Create a slice to hold the values for each row
    values := make([]interface{}, len(columns))
    for i := range values {
      values[i] = new(interface{}) // Store pointer to interface
    }

    // Scan the row into the values slice
    if err := rows.Scan(values...); err != nil {
      return "", err
    }

    // Create a map for the current row
    rowMap := make(map[string]interface{})
    for i, colName := range columns {
      value := *(values[i].(*interface{})) // Dereference pointer

      // Handle specific types if needed (e.g., converting byte slices)
      switch v := value.(type) {
      case []byte:
        rowMap[colName] = string(v) // Convert byte slice to string if needed
      default:
        rowMap[colName] = v
      }

    }
    results = append(results, rowMap)
  }

  // Convert results to JSON format
  jsonData, err := json.Marshal(results)
  if err != nil {
    return "", err
  }

  // return jsonData, nil or perhaps Writer could be better?
  return string(jsonData), nil
}
    
func (c *MysqlExecutorConnection) GetVersion() string {
  var version string
  c.db.Raw("SELECT version();").Scan(&version)
  return version
}
