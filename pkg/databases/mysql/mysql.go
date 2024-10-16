package mysql

import (
  "encoding/json"
	"fmt"
  
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/driver/mysql"

	"github.com/korbajan/sqlson/internal/configs"
)

type MysqlExecutor struct {
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
  DefaultPort = 3306
  DefaultUser = "root"
)

func NewExecutor(databaseConfig *configs.Database) *MysqlExecutor {
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
  
  return &MysqlExecutor{
    host: host,
    port: port,
    user: user,
    databaseName: databaseConfig.Name,
    password: databaseConfig.Password,
  }
  
}
func (msql *MysqlExecutor) GetDSN() string {
  return fmt.Sprintf(
    "%s:%s@tcp(%s:%d)/%s",
    msql.user,
    msql.password,
    msql.host,
    msql.port,
    msql.databaseName,
  )
}

func (msql *MysqlExecutor) PrepareDBConnection() error {
  db, err := gorm.Open(mysql.Open(msql.GetDSN()), &gorm.Config{
    Logger: logger.Discard, // Silences all gorm logs
  })
  if err != nil {
    return err
  }
  msql.db = db
  return nil
}
// Function to execute a raw SQL query and return results in JSON format
func (mysql *MysqlExecutor) Execute(sqlQuery string) (string, error) {
  if mysql.db == nil {
    err := mysql.PrepareDBConnection()
    if err != nil {
      return "", err
    }
  }
  // Execute the raw SQL query
  rows, err := mysql.db.Raw(sqlQuery).Rows()
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

  return string(jsonData), nil
}
    
func (msql *MysqlExecutor) GetVersion() string {
  msql.db.Raw("SELECT version();").Scan(&msql.dbVersion)
  return msql.dbVersion
}
