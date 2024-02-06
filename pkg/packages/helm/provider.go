package helm

import (
	"context"
	"io"
	"net/http"
	"strings"

	"aarti/pkg/utils/logger"

	"github.com/gorilla/mux"

	"aarti/pkg/packages"
	"aarti/pkg/storage"
)

const Name = "helm"

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

func (p *provider) setup(_ string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		if _, err := storage.FromContext(ctx).Stat(ctx, RepositoryPublicKey); err != nil {
			storage.Error(w, err)
			return
		}
		user, pass, _ := r.BasicAuth()
		repo := mux.Vars(r)["repo"]
		var name string
		if repo != "" {
			name = strings.NewReplacer("/", "-").Replace(repo)
		} else {
			name = strings.NewReplacer("/", "-", ".", "-").Replace(strings.TrimPrefix(strings.Split(r.Host, ":")[0], Name+"."))
		}
		args := SetupArgs{
			Name:     name,
			User:     user,
			Password: pass,
			Scheme:   packages.Scheme(r),
			Host:     r.Host,
			Path:     strings.TrimSuffix(r.URL.Path, "/setup"),
		}
		if err := scriptTemplate.Execute(w, args); err != nil {
			logger.C(r.Context()).WithError(err).Error("failed to execute template")
		}
	}
}

func (p *provider) Routes() []*packages.Route {
	return []*packages.Route{
		{
			Path:    "/setup",
			Method:  http.MethodGet,
			Handler: p.setup,
		},
		{
			Path:   "/{filename}",
			Method: http.MethodGet,
			Handler: packages.Pull(func(r *http.Request) string {
				return mux.Vars(r)["filename"]
			}),
		},
		{
			Path:   "/push",
			Method: http.MethodPut,
			Handler: packages.Push(func(r *http.Request, reader io.Reader, size int64, key string) (storage.Artifact, error) {
				return NewPackage(reader)
			}),
		},
		{
			Path:   "/{filename}",
			Method: http.MethodDelete,
			Handler: packages.Delete(func(r *http.Request) string {
				return mux.Vars(r)["filename"]
			}),
		},
	}
}
