package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/workspace/evoting/ev-webservice/internal/config"
	"github.com/workspace/evoting/ev-webservice/internal/consensusgroup"
	"github.com/workspace/evoting/ev-webservice/internal/country"
	"github.com/workspace/evoting/ev-webservice/internal/politicalparty"
	"github.com/workspace/evoting/ev-webservice/internal/user"
	"github.com/workspace/evoting/ev-webservice/pkg/log"
	"github.com/workspace/evoting/ev-webservice/pkg/token"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"gopkg.in/olahol/melody.v1"
)

// Server serves HTTP requests for our image service.
type Server struct {
	db         *mongo.Database
	router     *gin.Engine
	tokenMaker token.Maker
	config     config.Config
	logger     log.Logger
}

var (
	flagConfig = flag.String("config", "./config/local.yml", "path to the config file")
)

// NewServer creates a new HTTP server and set up routing.
func NewServer(db *mongo.Database, config *config.Config, logger log.Logger) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSecretKey)
	if err != nil {
		return nil, fmt.Errorf("Could not create token maker instance: %w", err)
	}

	server := &Server{
		db:         db,
		tokenMaker: tokenMaker,
		config:     *config,
		logger:     logger,
	}
	router := gin.Default()

	router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"API": "Web service"})
	})

	server.router = router

	//Build Router handler
	server.buildHandler()

	//Build Socket handler
	server.buildSocketHandler()

	return server, nil
}

//  buildHandler sets up the HTTP routing and builds an HTTP handler.
func (server *Server) buildHandler() {
	// Middlewares
	{
		//recovery middleware
		server.router.Use(gin.Recovery())
		//middleware which injects a 'RequestID' into the context and header of each request.
		server.router.Use(requestid.New())
	}
	v1 := server.router.Group("/api/v1")

	//Register user handlers
	user.RegisterHandlers(
		v1,
		user.NewUserService(
			user.NewMongoUserRepository(server.db, server.logger),
			server.logger,
		),
		server.logger,
	)

	//Register country handlers and service
	countryService := country.NewCountryService(
		country.NewMongoCountryRepository(server.db, server.logger),
		server.logger,
	)
	country.RegisterHandlers(v1, countryService, server.logger)

	//Register Political Party handlers
	politicalparty.RegisterHandlers(
		v1,
		politicalparty.NewPoliticalPartyService(
			politicalparty.NewMongoPoliticalPartyRepository(server.db, server.logger),
			server.logger,
		),
		countryService,
		server.logger,
	)

	//Register consensus group handlers
	consensusgroup.RegisterHandlers(
		v1,
		consensusgroup.NewGroupService(
			consensusgroup.NewMongoGroupRepository(server.db, server.logger),
			server.logger,
		),
		countryService,
		server.logger,
	)
}

// buildSocketHandler handles socket
func (server *Server) buildSocketHandler() {
	mrouter := melody.New()

	server.router.GET("/ws", func(c *gin.Context) {
		mrouter.HandleRequest(c.Writer, c.Request)
	})

	mrouter.HandleMessage(func(s *melody.Session, msg []byte) {
		mrouter.Broadcast(msg)
	})
}

// Start runs the HTTP server on a specific address.
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func main() {
	flag.Parse()

	// load application configurations'
	cfg, err := config.Load(*flagConfig)

	// create root logger tagged with server version
	logger := log.New(cfg.LogFile)

	if err != nil {
		logger.Errorf("failed to load application configuration: %s", err)
		os.Exit(-1)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.NewClient(options.Client().ApplyURI(cfg.DBSource))

	if err != nil {
		logger.Error(err)
		os.Exit(-1)
	}

	err = client.Connect(ctx)
	if err != nil {
		logger.Error(err)
		os.Exit(-1)
	}

	// Check the connection
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		logger.Error(err)
		os.Exit(-1)
	}

	defer client.Disconnect(ctx)

	database := client.Database(cfg.Database.Name)

	server, err := NewServer(
		database,
		cfg,
		logger,
	)
	if err != nil {
		logger.Error(err)
		os.Exit(-1)
	}
	address := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)

	logger.Infof("server %v is running at %v", cfg.Version, address)
	err = server.Start(address)
	if err != nil {
		logger.Error(err)
		os.Exit(-1)
	}
}
