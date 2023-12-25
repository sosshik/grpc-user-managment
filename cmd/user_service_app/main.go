package main

import (
	"fmt"
	"net"

	"git.foxminded.ua/foxstudent106264/task-4.1/internal/api"
	"git.foxminded.ua/foxstudent106264/task-4.1/internal/database"
	"git.foxminded.ua/foxstudent106264/task-4.1/pkg/config"
	proto "git.foxminded.ua/foxstudent106264/task-4.1/protos/gen/go/user_service"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Warn("No .env file")
	}

	cfg := config.GetConfig()

	level, err := log.ParseLevel(cfg.LogLevel)
	if err != nil {
		fmt.Printf("Error parsing log level: %v, setting log level to info\n", err)
		log.SetLevel(log.InfoLevel)
	} else {
		log.SetLevel(level)
		fmt.Printf("log level was set to %s\n", cfg.LogLevel)
	}
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
	fmt.Printf("config initialized\n")
}
func main() {
	cfg := config.GetConfig()

	db, err := database.NewDatabase(cfg)
	if err != nil {
		log.Warn(err)
	}
	defer db.DB.Close()

	s := grpc.NewServer()
	srv := &api.ServerAPI{DB: db}
	proto.RegisterUserServiceServer(s, srv)

	l, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Warn(err)
	}
	if err := s.Serve(l); err != nil {
		log.Warn(err)
	}
}
