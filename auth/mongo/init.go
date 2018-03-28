package postgres

import (
	"fmt"

	"gitlab.com/blocksense/monetasa/auth"
)

const errDuplicate string = "unique_violation"

type connection struct {
	ClientID  string `gorm:"primary_key"`
	ChannelID string `gorm:"primary_key"`
}

func (c connection) TableName() string {
	return "channel_clients"
}

// Connect creates a connection to the PostgreSQL instance. A non-nil error
// is returned to indicate failure.
func Connect(host, port, name, user, pass string) (*gorm.DB, error) {
	t := "host=%s port=%s user=%s dbname=%s password=%s sslmode=disable"
	url := fmt.Sprintf(t, host, port, user, name, pass)

  cfg := config{
		AuthPort: 					 p,
		AuthHost: 					 env(envAuthHost, authHost),
		MongoURL:            env(envMongoURL, defMongoURL),
		MongoUser:           defMongoUsername,
		MongoPass:           defMongoPassword,
		MongoDatabase:       defMongoDatabase,
		MongoPort:           defMongoPort,
		MongoConnectTimeout: defMongoConnectTimeout,
		MongoSocketTimeout:  defMongoSocketTimeout,
	}

  ms, err := connectToMongo(cfg)
	if err != nil {
		logger.Error("Failed to connect to Mongo.", zap.Error(err))
		return
	}
	defer ms.Close()

}
