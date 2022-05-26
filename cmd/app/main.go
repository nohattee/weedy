package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	ctrl_v1 "weedy/internal/controller/http/v1"
	"weedy/internal/repository"
	"weedy/pkg/config"
	"weedy/pkg/httpserver"

	_ "github.com/lib/pq"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "app",
	Short: "App is a very fast static site generator",
	Long: `A Fast and Flexible Static Site Generator built with
				  love by spf13 and friends in Go.
				  Complete documentation is available at http://hugo.spf13.com`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.LoadConfigFromEnv()

		var dns string
		if cfg.DB.Connection != "" {
			dns = cfg.DB.Connection
		} else {
			dns = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", cfg.DB.Host, cfg.DB.Port, cfg.DB.Username, cfg.DB.Password, cfg.DB.Name)
		}
		db, err := sql.Open("postgres", dns)
		if err != nil {
			log.Fatalf("sql.Open err: %v", err)
		}
		err = db.Ping()
		if err != nil {
			log.Fatalf("db.Ping err: %v", err)
		}

		userRepo := &repository.UserRepo{DB: db}
		server := httpserver.NewHttpServer(cfg.HTTP)
		server.AddController(ctrl_v1.NewUserController(userRepo))
		server.Run()
	},
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
