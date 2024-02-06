package apk

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"

	"aarti/pkg/utils/logger"

	"github.com/gorilla/mux"

	"aarti/pkg/crypt/rsa"
	"aarti/pkg/packages"
	"aarti/pkg/storage"
)

const Name = "apk"

var _ packages.Provider = (*provider)(nil)

func init() {
	packages.Register(Name, newProvider)
}

func newProvider(_ context.Context) (packages.Provider, error) {
	return &provider{}, nil
}

type provider struct{}

func (p *provider) Repository() storage.Repository {
	return &repo{}
}

func (p *provider) downloadKey(_ string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		s := storage.FromContext(ctx)
		pub, f, err := rsa.PublicKeyAndFingerprintFromPrivateKey(s.Key())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/x-pem-file")
		w.Header().Set("Content-Length", fmt.Sprintf("%d", len(pub)))
		w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename=%s`, fmt.Sprintf("aarticlient@%s.rsa.pub", hex.EncodeToString(f))))
		w.Header().Set("Access-Control-Expose-Headers", "Content-Disposition")
		io.Copy(w, bytes.NewReader(pub))
	}
}

func (p *provider) setup(_ string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		s := storage.FromContext(ctx)
		if _, err := s.Stat(ctx, RepositoryPublicKey); err != nil {
			storage.Error(w, err)
			return
		}
		branch, repository := mux.Vars(r)["branch"], mux.Vars(r)["repository"]
		user, pass, _ := r.BasicAuth()
		args := SetupArgs{
			User:       user,
			Password:   pass,
			Scheme:     packages.Scheme(r),
			Host:       r.Host,
			Path:       strings.TrimSuffix(r.URL.Path, fmt.Sprintf("/%s/%s/setup", branch, repository)),
			Branch:     branch,
			Repository: repository,
		}
		if err := scriptTemplate.Execute(w, args); err != nil {
			logger.C(r.Context()).WithError(err).Error("failed to execute template")
		}
	}
}

func (p *provider) Routes() []*packages.Route {
	return []*packages.Route{
		{
			Method:  http.MethodGet,
			Path:    "/{branch}/{repository}/key",
			Handler: p.downloadKey,
		},
		{
			Method:  http.MethodGet,
			Path:    "/{branch}/{repository}/setup",
			Handler: p.setup,
		},
		{
			Method: http.MethodPut,
			Path:   "/{branch}/{repository}/push",
			Handler: packages.Push(func(r *http.Request, reader io.Reader, size int64, key string) (storage.Artifact, error) {
				branch, repo := mux.Vars(r)["branch"], mux.Vars(r)["repository"]
				return NewPackage(reader, branch, repo, size)
			}),
		},
		{
			Method: http.MethodGet,
			Path:   "/{branch}/{repository}/{architecture}/{filename}",
			Handler: packages.Pull(func(r *http.Request) string {
				branch, repo, arch, filename := mux.Vars(r)["branch"], mux.Vars(r)["repository"], mux.Vars(r)["architecture"], mux.Vars(r)["filename"]
				return filepath.Join(branch, repo, arch, filename)
			}),
		},
		{
			Method: http.MethodDelete,
			Path:   "/{branch}/{repository}/{architecture}/{filename}",
			Handler: packages.Delete(func(r *http.Request) string {
				branch, repo, arch, filename := mux.Vars(r)["branch"], mux.Vars(r)["repository"], mux.Vars(r)["architecture"], mux.Vars(r)["filename"]
				return filepath.Join(branch, repo, arch, filename)
			}),
		},
	}
}
