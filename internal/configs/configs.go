package configs 

type Database struct {
  Name string
  Host string
  Port int
  User string
  Password string
}

func NewDatabaseConfig() Database {
  var config Database
  return config
}
