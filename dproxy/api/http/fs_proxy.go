package http

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/datapace/datapace/dproxy"
	log "github.com/datapace/datapace/logger"
	"github.com/go-zoo/bone"
)

var errUnauthorized = errors.New("unauthorized")

type FsProxy struct {
	svc         dproxy.Service
	localFsRoot string
	logger      log.Logger
	logPrefix   string
	PathPrefix  string
}

func NewFsProxy(svc dproxy.Service, localFsRoot, pathPrefix string, logger log.Logger) *FsProxy {
	return &FsProxy{svc: svc, localFsRoot: localFsRoot, PathPrefix: pathPrefix, logger: logger, logPrefix: "fs"}
}

func (f *FsProxy) GetFile(w http.ResponseWriter, r *http.Request) {
	fp, err := f.prepareFilePath(r)
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

func (f *FsProxy) PutFile(w http.ResponseWriter, r *http.Request) {
	fp, err := f.prepareFilePath(r)
	if err != nil {
		f.logger.Error(fmt.Sprintf("%s: failed to prepare file path %s with error %s", f.logPrefix, fp, err))
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(http.StatusText(http.StatusUnauthorized)))
		return
	}
	f.logger.Info(fmt.Sprintf("%s: received PUT request from %s to %s", f.logPrefix, r.RemoteAddr, fp))
	f.saveFile(fp, w, r)
}

func (f *FsProxy) PostFile(w http.ResponseWriter, r *http.Request) {
	fp, err := f.prepareFilePath(r)
	if err != nil {
		f.logger.Error(fmt.Sprintf("%s: failed to prepare file path %s with error %s", f.logPrefix, fp, err))
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(http.StatusText(http.StatusUnauthorized)))
		return
	}
	fp = fp + time.Now().Format(time.RFC3339Nano)
	f.logger.Info(fmt.Sprintf("%s: received POST request from %s to %s", f.logPrefix, r.RemoteAddr, fp))
	f.saveFile(fp, w, r)
}

func (f *FsProxy) prepareFilePath(r *http.Request) (string, error) {
	t := r.Header.Get("Authorization")
	//if there is no token in authorization header, try token in the url
	if t == "" {
		t = bone.GetValue(r, "token")
	}
	targetURL, err := f.svc.GetTargetURL(t)
	if err != nil {
		return "", err
	}
	fp := filepath.Join(f.localFsRoot, targetURL)
	if strings.Contains(fp, "..") {
		return "", errUnauthorized
	}
	return fp, nil
}

func (f *FsProxy) saveFile(fp string, w http.ResponseWriter, r *http.Request) {
	file, err := os.Create(fp)
	if err != nil {
		f.logger.Error(fmt.Sprintf("%s: failed to create file %s with error %s", f.logPrefix, fp, err))
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(http.StatusText(http.StatusInternalServerError)))
		return
	}
	defer file.Close()
	defer r.Body.Close()
	io.Copy(file, r.Body)
}
