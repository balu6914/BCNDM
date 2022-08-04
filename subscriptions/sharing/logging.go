package sharing

import (
	"fmt"
	log "github.com/datapace/datapace/logger"
	"github.com/datapace/sharing"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

var _ Service = (*loggingMiddleware)(nil)

type loggingMiddleware struct {
	logger log.Logger
	svc    Service
}

// NewLoggingMiddleware adds logging facilities to the core service.
func NewLoggingMiddleware(svc Service, logger log.Logger) Service {
	return &loggingMiddleware{logger, svc}
}

func (l loggingMiddleware) GetSharings(rcvUserId string, rcvGroupIds []string) (sharings []sharing.Sharing, err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method GetSharings(%s, %s) took %s to complete", rcvUserId, rcvGroupIds, time.Since(begin))
		if err != nil {
			st, _ := status.FromError(err)
			if st != nil && st.Code() == codes.Unavailable {
				l.logger.Warn("Failed to resolve shared streams: sharing service is unavailable (not supported).")
				return
			}
			l.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		l.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())

	return l.svc.GetSharings(rcvUserId, rcvGroupIds)
}
