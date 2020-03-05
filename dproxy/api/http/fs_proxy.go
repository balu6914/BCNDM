package http

import (
	"bufio"
	"datapace/dproxy"
	log "datapace/logger"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var errUnauthorized = errors.New("unauthorized")

type FsProxy struct {
	svc         dproxy.Service
	localFsRoot string
	logger      log.Logger
	logPrefix   string
}

func NewFsProxy(svc dproxy.Service, localFsRoot string, logger log.Logger) *FsProxy {
	return &FsProxy{svc: svc, localFsRoot: localFsRoot, logger: logger, logPrefix: "fs"}
}

func (f *FsProxy) GetFile(w http.ResponseWriter, r *http.Request) {
	fp, err := f.prepareFilePath(r.Header.Get("Authorization"))
	if err != nil {
		f.logger.Error(fmt.Sprintf("%s: failed to prepare file path %s with error %s", f.logPrefix, fp, err))
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(http.StatusText(http.StatusUnauthorized)))
		return
	}
	f.logger.Info(fmt.Sprintf("%s: received request from %s to %s", f.logPrefix, r.RemoteAddr, fp))
	file, err := os.Open(fp)
	if err != nil {
		f.logger.Error(fmt.Sprintf("%s: failed to open file %s with error %s", f.logPrefix, fp, err))
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(http.StatusText(http.StatusNotFound)))
		return
	}
	defer file.Close()
	b := bufio.NewReader(file)
	io.Copy(w, b)
}

func (f *FsProxy) prepareFilePath(token string) (string, error) {
	targetURL, err := f.svc.GetTargetURL(token)
	if err != nil {
		return "", err
	}
	fp := filepath.Join(f.localFsRoot, targetURL)
	if strings.Contains(fp, "..") {
		return "", errUnauthorized
	}
	return fp, nil
}
