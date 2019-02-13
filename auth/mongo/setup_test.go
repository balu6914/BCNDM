package mongo_test

import (
	"datapace/auth/mongo"
	log "datapace/logger"
	"fmt"
	"os"
	"testing"

	"github.com/ory/dockertest"
)

const (
	testDB          = "datapace"
	userDB          = ""
	passDB          = ""
	timeoutDB       = 5000
	socketTimeoutDB = 5000
)

var testLog = log.New(os.Stdout)

func TestMain(m *testing.M) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		testLog.Error(fmt.Sprintf("Could not connect to docker: %s", err))
	}

	cfg := []string{
		fmt.Sprintf("MONGO_INITDB_DATABASE=%s", testDB),
	}

	container, err := pool.Run("mongo", "3.6", cfg)
	if err != nil {
		testLog.Error(fmt.Sprintf("Could not start container: %s", err))
	}

	port := container.GetPort("27017/tcp")
	addr := fmt.Sprintf("localhost:%s", port)

	err = pool.Retry(func() error {
		db, err = mongo.Connect(addr, timeoutDB, socketTimeoutDB, testDB, userDB, passDB)
		return err
	})
	if err != nil {
		testLog.Error(fmt.Sprintf("Could not connect to docker: %s", err))
	}
	defer db.Close()

	code := m.Run()

	if err := pool.Purge(container); err != nil {
		testLog.Error(fmt.Sprintf("Could not purge container: %s", err))
	}

	os.Exit(code)
}
