package archive // import "github.com/ory/dockertest/docker/pkg/archive"

import (
	"syscall"
	"time"
)

func timeToTimespec(time time.Time) (ts syscall.Timespec) {
	if time.IsZero() {
		// Return UTIME_OMIT special value
		ts.Sec = 0
		ts.Nsec = ((1 << 30) - 2)
		return
	}
	return syscall.NsecToTimespec(time.UnixNano())
}
