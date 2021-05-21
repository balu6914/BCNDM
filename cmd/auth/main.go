package main

import (
	"fmt"
	"github.com/datapace/datapace/auth/mail"
	"github.com/datapace/datapace/auth/recovery"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
	mgo "gopkg.in/mgo.v2"

	"github.com/datapace/datapace"

	accessapi "github.com/datapace/datapace/access-control/api/grpc"
	"github.com/datapace/datapace/auth"
	"github.com/datapace/datapace/auth/access"
	"github.com/datapace/datapace/auth/aes"
	"github.com/datapace/datapace/auth/api"
	grpcapi "github.com/datapace/datapace/auth/api/grpc"
	httpapi "github.com/datapace/datapace/auth/api/http"
	"github.com/datapace/datapace/auth/bcrypt"
	"github.com/datapace/datapace/auth/jwt"
	"github.com/datapace/datapace/auth/mongo"
	"github.com/datapace/datapace/auth/transactions"
	log "github.com/datapace/datapace/logger"
	accessproto "github.com/datapace/datapace/proto/access"
	authproto "github.com/datapace/datapace/proto/auth"
	transactionsproto "github.com/datapace/datapace/proto/transactions"
	transactionsapi "github.com/datapace/datapace/transactions/api/grpc"
)

const (
	envHTTPPort         = "DATAPACE_AUTH_HTTP_PORT"
	envGRPCPort         = "DATAPACE_AUTH_GRPC_PORT"
	envDBURL            = "DATAPACE_AUTH_DB_URL"
	envDBUser           = "DATAPACE_AUTH_DB_USER"
	envDBPass           = "DATAPACE_AUTH_DB_PASS"
	envDBName           = "DATAPACE_AUTH_DB_NAME"
	envTransactionsURL  = "DATAPACE_TRANSACTIONS_URL"
	envAccessControlURL = "DATAPACE_ACCESS_CONTROL_URL"
	envSecret           = "DATAPACE_AUTH_SECRET"
	envEncryptionKey    = "DATAPACE_ENCRYPTION_KEY"
	envAdminEmail       = "DATAPACE_ADMIN_EMAIL"
	envAdminPassword    = "DATAPACE_ADMIN_PASSWORD"
	envSmtpIdentity     = "DATAPACE_SMTP_IDENTITY"
	envSmtpURL          = "DATAPACE_SMTP_URL"
	envSmtpHost         = "DATAPACE_SMTP_HOST"
	envSmtpUser         = "DATAPACE_SMTP_USER"
	envSmtpPassword     = "DATAPACE_SMTP_PASSWORD"
	envSmtpFrom         = "DATAPACE_SMTP_FROM"
	envFrontendURL      = "DATAPACE_FRONTEND_URL"
	envPassRecoveryTpl  = "DATAPACE_PASSWORD_RECOVERY_TPL"

	defHTTPPort         = "8080"
	defGRPCPort         = "8081"
	defDBURL            = "0.0.0.0"
	defDBUser           = ""
	defDBPass           = ""
	defDBName           = "auth"
	defTransactionsURL  = "localhost:8081"
	defAccessControlURL = "localhost:8081"
	defSecret           = "github.com/datapace/datapace"
	defEncryptionKey    = "AES256Key-32Characters1234567890"
	defAdminEmail       = "admin@datapace.localhost"
	defAdminPassword    = "datapaceadmin"
	defSmtpIdentity     = ""
	defSmtpURL          = "smtp.mailtrap.io:25"
	defSmtpHost         = "smtp.mailtrap.io"
	defSmtpUser         = "3b29d66d776ccc"
	defSmtpPassword     = "8bfabd687f207b"
	defSmtpFrom         = "noreply@datapace.io"
	defFrontendURL      = "https://datapace.io"
	defPassRecoveryTpl  = "auth/mail/templates/passwordRecovery.html"

	dbConnectTimeout = 5000
	dbSocketTimeout  = 5000
)

type config struct {
	httpPort         string
	grpcPort         string
	dbURL            string
	dbUser           string
	dbPass           string
	dbName           string
	dbConnectTimeout int
	dbSocketTimeout  int
	transactionsURL  string
	accessControlURL string
	secret           string
	encryptionKey    string
	adminEmail       string
	adminPassword    string
	smtpIdentity     string
	smtpURL          string
	smtpHost         string
	smtpUser         string
	smtpPassword     string
	smtpFrom         string
	frontendURL      string
	passRecoveryTpl  string
}

func main() {
	cfg := loadConfig()

	logger := log.New(os.Stdout)

	ms := connectToDB(cfg, logger)
	defer ms.Close()

	tconn := newGRPCConn(cfg.transactionsURL, logger)
	defer tconn.Close()

	acconn := newGRPCConn(cfg.accessControlURL, logger)
	defer acconn.Close()

	tc := transactionsapi.NewClient(tconn)
	ac := accessapi.NewClient(acconn)

	svc := newService(cfg, ms, tc, ac, logger)
	initAdmin(svc, cfg.adminEmail, cfg.adminPassword, logger)

	errs := make(chan error, 2)

	go startHTTPServer(svc, cfg.httpPort, logger, errs)

	go startGRPCServer(svc, cfg.grpcPort, logger, errs)

	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	err := <-errs
	logger.Error(fmt.Sprintf("Auth service terminated: %s", err))
}

func loadConfig() config {
	return config{
		httpPort:         datapace.Env(envHTTPPort, defHTTPPort),
		grpcPort:         datapace.Env(envGRPCPort, defGRPCPort),
		dbURL:            datapace.Env(envDBURL, defDBURL),
		dbUser:           datapace.Env(envDBUser, defDBUser),
		dbPass:           datapace.Env(envDBPass, defDBPass),
		dbName:           datapace.Env(envDBName, defDBName),
		dbConnectTimeout: dbConnectTimeout,
		dbSocketTimeout:  dbSocketTimeout,
		transactionsURL:  datapace.Env(envTransactionsURL, defTransactionsURL),
		accessControlURL: datapace.Env(envAccessControlURL, defAccessControlURL),
		secret:           datapace.Env(envSecret, defSecret),
		encryptionKey:    datapace.Env(envEncryptionKey, defEncryptionKey),
		adminEmail:       datapace.Env(envAdminEmail, defAdminEmail),
		adminPassword:    datapace.Env(envAdminPassword, defAdminPassword),
		smtpIdentity:     datapace.Env(envSmtpIdentity, defSmtpIdentity),
		smtpURL:          datapace.Env(envSmtpURL, defSmtpURL),
		smtpHost:         datapace.Env(envSmtpHost, defSmtpHost),
		smtpUser:         datapace.Env(envSmtpUser, defSmtpUser),
		smtpPassword:     datapace.Env(envSmtpPassword, defSmtpPassword),
		smtpFrom:         datapace.Env(envSmtpFrom, defSmtpFrom),
		frontendURL:      datapace.Env(envFrontendURL, defFrontendURL),
		passRecoveryTpl:  datapace.Env(envPassRecoveryTpl, defPassRecoveryTpl),
	}
}

func connectToDB(cfg config, logger log.Logger) *mgo.Session {
	ms, err := mongo.Connect(
		cfg.dbURL,
		cfg.dbConnectTimeout,
		cfg.dbSocketTimeout,
		cfg.dbName,
		cfg.dbUser,
		cfg.dbPass,
	)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to connect to Mongo: %s", err))
		os.Exit(1)
	}

	return ms
}

func initAdmin(svc auth.Service, adminEmail, adminPassword string, logger log.Logger) {
	user := auth.User{
		Email:        adminEmail,
		ContactEmail: adminEmail,
		Password:     adminPassword,
		ID:           "",
		FirstName:    "",
		LastName:     "",
		Company:      "",
		Address:      "",
		Phone:        "",
		Role:         auth.AdminRole,
	}

	policies := map[string]auth.Policy{
		auth.AdminRole: {
			Name:    auth.AdminRole,
			Owner:   user.ID,
			Version: "1.0.0",
			Rules: []auth.Rule{
				{
					Action: auth.Any,
					Type:   "user",
				},
				{
					Action: auth.Any,
					Type:   "stream",
				},
				{
					Action: auth.Any,
					Type:   "subscription",
				},
				{
					Action: auth.Any,
					Type:   "policy",
				},
				{
					Action: auth.Any,
					Type:   "contract",
				},
				{
					Action: auth.Read,
					Type:   "token",
				},
				{
					Action: auth.List,
					Type:   "token",
				},
				{
					Action: auth.Buy,
					Type:   "token",
				},
				{
					Action: auth.Withdraw,
					Type:   "token",
				},
			},
		},
		auth.AdminUserRole: {
			Name:    auth.AdminUserRole,
			Owner:   user.ID,
			Version: "1.0.0",
			Rules: []auth.Rule{
				{
					Action: auth.Any,
					Type:   "user",
				},
				{
					Action: auth.Any,
					Type:   "stream",
				},
				{
					Action: auth.Any,
					Type:   "subscription",
				},
				{
					Action: auth.Any,
					Type:   "policy",
				},
				{
					Action: auth.Any,
					Type:   "contract",
				},
				{
					Action: auth.Read,
					Type:   "token",
				},
				{
					Action: auth.List,
					Type:   "token",
				},
			},
		},
		auth.AdminWalletRole: {
			Name:    auth.AdminWalletRole,
			Owner:   user.ID,
			Version: "1.0.0",
			Rules: []auth.Rule{
				{
					Action: auth.Read,
					Type:   "user",
				},
				{
					Action: auth.List,
					Type:   "user",
				},
				{
					Action: auth.Any,
					Type:   "stream",
				},
				{
					Action: auth.Any,
					Type:   "subscription",
				},
				{
					Action: auth.Any,
					Type:   "policy",
				},
				{
					Action: auth.Any,
					Type:   "contract",
				},
				{
					Action: auth.Read,
					Type:   "token",
				},
				{
					Action: auth.List,
					Type:   "token",
				},
				{
					Action: auth.Buy,
					Type:   "token",
				},
				{
					Action: auth.Withdraw,
					Type:   "token",
				},
			},
		},
		"user": {
			Name:    "user",
			Owner:   user.ID,
			Version: "1.0.0",
			Rules: []auth.Rule{
				{
					Action: auth.Create,
					Type:   "stream",
				},
				{
					Action: auth.Create,
					Type:   "contract",
				},
				{
					Action: auth.CreateBulk,
					Type:   "stream",
				},
				{
					Action: auth.CreateBulk,
					Type:   "subscription",
				},
				{
					Action: auth.List,
					Type:   "stream",
				},
				{
					Action: auth.List,
					Type:   "contract",
				},
				{
					Action: auth.List,
					Type:   "subscription",
				},
				{
					Action: auth.Sign,
					Type:   "contract",
				},
				{
					Action: auth.Read,
					Type:   "stream",
				},
				{
					Action: auth.Read,
					Type:   "subscription",
				},
				{
					Action: auth.Any,
					Type:   "stream",
					Condition: auth.SimpleCondition{
						Key: "ownerID",
					},
				},
				{
					Action: auth.Any,
					Type:   "contract",
					Condition: auth.SimpleCondition{
						Key: "ownerID",
					},
				},

				{
					Action: auth.Any,
					Type:   "subscription",
					Condition: auth.SimpleCondition{
						Key: "ownerID",
					},
				},
				{
					Action: auth.Any,
					Type:   "user",
					Condition: auth.SimpleCondition{
						Key: "id",
					},
				},
				{
					Action: auth.Any,
					Type:   "token",
				},
			},
		},
	}

	if err := svc.InitAdmin(user, policies); err != nil {
		logger.Error(fmt.Sprintf("Failed to create admin: %s", err))
		os.Exit(1)
	}
}

func newGRPCConn(addr string, logger log.Logger) *grpc.ClientConn {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to connect to GRPC service on address %s: %s", addr, err))
		os.Exit(1)
	}

	return conn
}

func newService(cfg config, ms *mgo.Session, tc transactionsproto.TransactionsServiceClient, asc accessproto.AccessServiceClient, logger log.Logger) auth.Service {
	cipher := aes.NewCipher([]byte(cfg.encryptionKey))
	users := mongo.NewUserRepository(ms, cipher)
	policies := mongo.NewPolicyRepository(ms)
	hasher := bcrypt.New()
	idp := jwt.New(cfg.secret)
	ts := transactions.NewService(tc)
	ac := access.New(asc)
	rc := recovery.New()
	mailsvc := mail.New(cfg.smtpIdentity, cfg.smtpURL, cfg.smtpHost, cfg.smtpUser, cfg.smtpPassword, cfg.smtpFrom, cfg.frontendURL, cfg.passRecoveryTpl)

	svc := auth.New(users, policies, hasher, idp, ts, ac, rc, mailsvc)
	svc = api.LoggingMiddleware(svc, logger)
	svc = api.MetricsMiddleware(
		svc,
		kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
			Namespace: "auth",
			Subsystem: "api",
			Name:      "request_count",
			Help:      "Number of requests received.",
		}, []string{"method"}),
		kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: "auth",
			Subsystem: "api",
			Name:      "request_latency_microseconds",
			Help:      "Total duration of requests in microseconds.",
		}, []string{"method"}),
	)

	return svc
}

func startHTTPServer(svc auth.Service, port string, logger log.Logger, errs chan error) {
	p := fmt.Sprintf(":%s", port)
	logger.Info(fmt.Sprintf("Users HTTP service started, exposed port %s", port))
	errs <- http.ListenAndServe(p, httpapi.MakeHandler(svc))
}

func startGRPCServer(svc auth.Service, port string, logger log.Logger, errs chan error) {
	p := fmt.Sprintf(":%s", port)
	listener, err := net.Listen("tcp", p)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to listen on port %s: %s", port, err))
	}

	server := grpc.NewServer()
	authproto.RegisterAuthServiceServer(server, grpcapi.NewServer(svc))
	logger.Info(fmt.Sprintf("Auth gRPC service started, exposed port %s", port))
	errs <- server.Serve(listener)
}
