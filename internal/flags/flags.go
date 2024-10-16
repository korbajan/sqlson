package flags

import (
  "flag"

)
  
func Parse(host *string, port *int, name, user, password *string, query *string) {
  flag.StringVar(host, "host", "localhost", "db host name/ip")
  flag.IntVar(port, "port", 0, "db host name/ip")
  flag.StringVar(password, "password", "your-secret-pw", "users password")
  flag.StringVar(user, "user", "root", "user name")
  flag.StringVar(name, "database", "database", "database name")
  flag.StringVar(query, "query", "", "sql query to execute")
  flag.Parse()
}
