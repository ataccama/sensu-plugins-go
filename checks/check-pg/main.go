package main

import (
	"database/sql"
	"fmt"
	"sensuatc/sensuutil"

	_ "github.com/lib/pq"
	"github.com/spf13/cobra"
)

// Argument variables and Cobra rootCmd
//
var (
	cmd        *cobra.Command
	host       string
	port       uint
	user       string
	password   string
	database   string
	threshWarn uint64
	threshCrit uint64
)

// Initialize CLI arguments
func init() {
	cmd = sensuutil.Cmd("check-pg", check)
	cmd.Flags().Uint64VarP(&threshWarn, "warning", "w", 100, "Connection count warning threshold")
	cmd.Flags().Uint64VarP(&threshCrit, "critical", "c", 200, "Connection count critical threshold")
	cmd.Flags().StringVarP(&host, "host", "H", "/var/run/postgresql", "Hostname/IP of the Postgres")
	cmd.Flags().UintVarP(&port, "port", "P", 5432, "Postgres TCP port")
	cmd.Flags().StringVarP(&user, "user", "u", "monitoring", "Username to connect as")
	cmd.Flags().StringVarP(&password, "password", "p", "", "Password for the user")
	cmd.Flags().StringVarP(&database, "database", "n", "postgres", "Database name to connect to")
}

// Execute the check
func main() {
	cmd.Execute()
}

// Check function
//
func check(c *cobra.Command, args []string) {
	var connStr string

	if host[0] == '/' {
		connStr = fmt.Sprintf("user=%s host=%s dbname=%s", user, host, database)
	} else {
		connStr = fmt.Sprintf("user=%s passowrd=%s host=%s port=%d dbname=%s sslmode=disable", user, password, host, port, database)
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		sensuutil.Exit("CRITICAL", err.Error())
	}

	var data uint64
	row := db.QueryRow("SELECT COUNT(pid) FROM pg_stat_activity WHERE usename NOT LIKE 'monitoring';")

	err = row.Scan(&data)
	if err != nil {
		sensuutil.Exit("CRITICAL", err.Error())
	}

	exitStatus := "OK"
	if data >= threshCrit {
		exitStatus = "CRITICAL"
	} else if data >= threshWarn {
		exitStatus = "WARNING"
	}
	sensuutil.Exit(exitStatus, fmt.Sprintf("%v connections", data))
}
