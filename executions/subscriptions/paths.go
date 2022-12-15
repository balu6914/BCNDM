package subscriptions

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/datapace/datapace/executions"
)

type pathRepository struct {
	subscriptionsURL string
}

// New returns path repository instance.
func New(subscriptionsURL string) executions.PathRepository {
	return pathRepository{
		subscriptionsURL: subscriptionsURL,
	}
}

func (repo pathRepository) Current(owner, data string) (string, error) {
	url := fmt.Sprintf("%s/owner/%s/stream/%s/subscriptions", repo.subscriptionsURL, owner, data)
	resp, _ := http.Get(url)

	if resp.StatusCode == 200 {
		defer resp.Body.Close()

		var sub viewSubRes
		if err := json.NewDecoder(resp.Body).Decode(&sub); err != nil {
			return "", err
		}

		return sub.URL, nil
	}

	return "", nil

}
