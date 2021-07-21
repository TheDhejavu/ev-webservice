package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/go-echarts/statsview"
	newrelic "github.com/newrelic/go-agent"
	"github.com/newrelic/go-agent/_integrations/nrgin/v1"
	"github.com/thedhejavu/ev-blockchain-protocol/rpc"

	"github.com/workspace/evoting/ev-webservice/internal/accrediation"
	"github.com/workspace/evoting/ev-webservice/internal/auth"
	"github.com/workspace/evoting/ev-webservice/internal/chain"
	"github.com/workspace/evoting/ev-webservice/internal/config"
	"github.com/workspace/evoting/ev-webservice/internal/consensusgroup"
	"github.com/workspace/evoting/ev-webservice/internal/country"
	"github.com/workspace/evoting/ev-webservice/internal/election"
	"github.com/workspace/evoting/ev-webservice/internal/identity"
	"github.com/workspace/evoting/ev-webservice/internal/politicalparty"
	"github.com/workspace/evoting/ev-webservice/internal/user"
	"github.com/workspace/evoting/ev-webservice/internal/voting"
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
type key int

const keyNrID key = iota

var (
	nrapp      newrelic.Application
	flagConfig = flag.String("config", "./config/local.yml", "path to the config file")
)

func initNewRelic(conf *config.Config) {
	var err error
	nrConfig := newrelic.NewConfig("app", conf.NewrelicKey)
	nrapp, err = newrelic.NewApplication(nrConfig)
	if err != nil {
		panic("Failed to setup NewRelic: " + err.Error())
	}
}

//populateNewRelicInContext get the request context populated
func setNewRelicInContext() gin.HandlerFunc {

	return func(c *gin.Context) {
		//Setup context
		ctx := c.Request.Context()

		//Set newrelic context
		var txn newrelic.Transaction
		//newRelicTransaction is the key populated by nrgin Middleware
		value, exists := c.Get("newRelicTransaction")
		if exists {
			if v, ok := value.(newrelic.Transaction); ok {
				txn = v
			}
			ctx = context.WithValue(ctx, keyNrID, txn)
		}
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

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
	router.Use(static.Serve("/assets", static.LocalFile("storage", false)))

	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// // Set a lower memory limit for multipart forms (default is 32 MiB)
	// router.MaxMultipartMemory = 8 << 20 // 8 MiB
	initNewRelic(config)

	router.Use(nrgin.Middleware(nrapp))
	router.Use(setNewRelicInContext())

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
	// General Middlewares
	{
		//recovery middleware
		server.router.Use(gin.Recovery())
		//middleware which injects a 'RequestID' into the context and header of each request.
		server.router.Use(requestid.New())
	}
	v1 := server.router.Group("/api/v1")

	//Register user handlers
	userService := user.NewUserService(
		user.NewMongoUserRepository(server.db, server.logger),
		server.logger,
	)

	//Register country handlers and service
	countryService := country.NewCountryService(
		country.NewMongoCountryRepository(server.db, server.logger),
		server.logger,
	)
	country.RegisterHandlers(v1, countryService, server.logger)

	user.RegisterHandlers(
		v1,
		userService,
		server.logger,
	)

	//Register identity handlers
	identityService := identity.NewIdentityService(
		identity.NewMongoIdentityRepository(server.db, server.logger, server.config),
		server.logger,
	)

	authMiddleware := auth.NewAuthMiddleware(
		userService,
		identityService,
		server.tokenMaker,
		server.logger,
	)

	identity.RegisterHandlers(
		v1,
		identityService,
		countryService,
		identityService,
		authMiddleware,
		server.config,
		server.logger,
	)

	//Register auth handlers
	auth.RegisterHandlers(
		v1,
		auth.NewIdentityService(
			identityService,
			userService,
			server.logger,
			server.config,
			server.tokenMaker,
		),
		server.logger,
	)

	//Register Political Party handlers
	politicalPartyService := politicalparty.NewPoliticalPartyService(
		politicalparty.NewMongoPoliticalPartyRepository(server.db, server.logger),
		server.logger,
	)
	politicalparty.RegisterHandlers(
		v1,
		politicalPartyService,
		countryService,
		server.logger,
	)

	//Register consensus group handlers
	consensusGroupRepo := consensusgroup.NewMongoGroupRepository(server.db, server.logger)
	consensusGroupService := consensusgroup.NewGroupService(
		consensusGroupRepo,
		server.logger,
	)
	consensusgroup.RegisterHandlers(
		v1,
		consensusGroupService,
		countryService,
		authMiddleware,
		server.logger,
	)

	//Register blockchain handler
	blockchainService := chain.NewBlockchainService(
		rpc.NewClient(server.config.RPCServerURL),
	)

	electionRepo := election.NewMongoElectionRepository(server.db, server.logger)
	//Register election handlers
	election.RegisterHandlers(
		v1,
		election.NewElectionService(
			electionRepo,
			blockchainService,
			consensusGroupService,
			server.logger,
		),
		countryService,
		politicalPartyService,
		identityService,
		authMiddleware,
		server.logger,
	)

	chain.RegisterHandlers(
		v1,
		authMiddleware,
		blockchainService,
		server.logger,
	)

	//Register accreditation handler
	accrediation.RegisterHandlers(
		v1,
		authMiddleware,
		accrediation.NewAccreditationService(
			blockchainService,
			electionRepo,
			consensusGroupService,
			server.logger,
		),
		server.logger,
	)

	//Register voting handler
	voting.RegisterHandlers(
		v1,
		authMiddleware,
		voting.NewVotingService(
			blockchainService,
			electionRepo,
			consensusGroupService,
			server.logger,
		),
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

	mgr := statsview.New()

	// Start() runs a HTTP server at `localhost:18066` by default.
	go mgr.Start()
	err = server.Start(address)
	if err != nil {
		logger.Error(err)
		os.Exit(-1)
	}
}
