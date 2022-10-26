package mongo

import (
	"fmt"
	"os"
	"testing"

	"github.com/ory/dockertest"
	"gopkg.in/mgo.v2"
)

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
	addr := fmt.Sprintf("mongodb://localhost:%s", port)

	if err := pool.Retry(func() error {
		db, err = mgo.Dial(addr)
		return err
	}); err != nil {
		testLog.Error(fmt.Sprintf("Could not connect to docker: %s", err))
	}

	code := m.Run()

	if err := pool.Purge(container); err != nil {
		testLog.Error(fmt.Sprintf("Could not purge container: %s", err))
	}

	os.Exit(code)
}
