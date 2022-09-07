package mongo

import (
	"fmt"
	"github.com/datapace/datapace/access-control"
	log "github.com/datapace/datapace/logger"
	"github.com/ory/dockertest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/mgo.v2"
	"os"
	"testing"
)

var (
	db      *mgo.Session
	testLog = log.New(os.Stdout)
)

func TestMain(m *testing.M) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		testLog.Error(fmt.Sprintf("Could not connect to docker: %s", err))
	}

	cfg := []string{
		fmt.Sprintf("MONGO_INITDB_DATABASE=%s", dbName),
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

func TestRequestAccessMultipleTimes(t *testing.T) {
	db.DB(dbName).DropDatabase()
	repo := NewAccessRequestRepository(db)
	acId0, err := repo.GrantAccess("user0", "user1")
	require.NotEmpty(t, acId0)
	require.Nil(t, err)
	ac, err := repo.One(acId0)
	require.Nil(t, err)
	require.Equal(t, "user0", ac.Receiver)
	require.Equal(t, "user1", ac.Sender)
	require.Equal(t, access.State("approved"), ac.State)
	err = repo.Revoke("user0", acId0)
	require.Nil(t, err)
	acId1, err := repo.RequestAccess("user1", "user0")
	assert.Nil(t, err)
	assert.Equal(t, acId0, acId1)
	acId2, err := repo.RequestAccess("user1", "user0")
	assert.Nil(t, err)
	assert.Equal(t, acId0, acId2)
	acId3, err := repo.RequestAccess("user2", "user0")
	assert.Nil(t, err)
	assert.NotEqual(t, acId0, acId3)
}
