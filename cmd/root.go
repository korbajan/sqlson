package cmd

import (
  "fmt"
  "os"

  "github.com/korbajan/sqlson/internal/configs"
  "github.com/korbajan/sqlson/internal/flags"
  "github.com/korbajan/sqlson/pkg/databases"
)

func parseFromArgs(databaseConfig *configs.Database, query *string) () {
  flags.Parse(&databaseConfig.Host, &databaseConfig.Port, &databaseConfig.Name, &databaseConfig.User, &databaseConfig.Password, query)
}

func Execute() {
  databaseConfig := configs.NewDatabaseConfig()
  var sqlQuery string
  parseFromArgs(&databaseConfig, &sqlQuery)
  jsonString, err := databases.Execute(&databaseConfig, sqlQuery)
  if err != nil {
    fmt.Println(err)
    os.Exit(1)
  }
  fmt.Println(jsonString)
}
