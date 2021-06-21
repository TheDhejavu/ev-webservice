package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/workspace/evoting/ev-webservice/internal/config"
	"github.com/workspace/evoting/ev-webservice/pkg/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/olahol/melody.v1"
)

// Server serves HTTP requests for our image service.
type Server struct {
	db     *mongo.Database
	router *gin.Engine
}

var (
	flagConfig = flag.String("config", "./config/local.yml", "path to the config file")
)

func main() {
	flag.Parse()

	// load application configurations'
	cfg, err := config.Load(*flagConfig)

	// create root logger tagged with server version
	logger := log.New(cfg.LogFile).With(nil, "version", cfg.Version)

	if err != nil {
		logger.Errorf("failed to load application configuration: %s", err)
		os.Exit(-1)
	}
	client, err := mongo.NewClient(options.Client().ApplyURI(cfg.DBSource))
	if err != nil {
		logger.Error(err)
		os.Exit(-1)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		logger.Error(err)
		os.Exit(-1)
	}
	defer client.Disconnect(ctx)

	database := client.Database(cfg.Database.Name)

	server := NewServer(database)
	address := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)

	logger.Infof("server %v is running at %v", cfg.Version, address)
	err = server.Start(address)
	if err != nil {
		logger.Error(err)
		os.Exit(-1)
	}
}

//  buildHandler sets up the HTTP routing and builds an HTTP handler.
func buildHandler(logger log.Logger, cfg *config.Config) {}

// buildSocketHandler handles socket
func buildSocketHandler(router *gin.Engine) {
	mrouter := melody.New()

	router.GET("/ws", func(c *gin.Context) {
		mrouter.HandleRequest(c.Writer, c.Request)
	})

	mrouter.HandleMessage(func(s *melody.Session, msg []byte) {
		mrouter.Broadcast(msg)
	})
}

// NewServer creates a new HTTP server and set up routing.
func NewServer(db *mongo.Database) *Server {
	server := &Server{db: db}
	router := gin.Default()

	router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"API": "Web service"})
	})

	buildSocketHandler(router)

	server.router = router
	return server
}

// Start runs the HTTP server on a specific address.
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}
