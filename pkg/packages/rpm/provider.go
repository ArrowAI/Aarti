package rpm

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	"aarti/pkg/utils/logger"

	"github.com/gorilla/mux"

	"aarti/pkg/packages"
	"aarti/pkg/storage"
)

const Name = "rpm"

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

func (p *provider) config(repo string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		var name string
		if repo != "" {
			name = strings.NewReplacer("/", "-").Replace(repo)
		} else {
			name = strings.NewReplacer("/", "-", ".", "-").Replace(strings.TrimPrefix(strings.Split(r.Host, ":")[0], Name+"."))
		}
		if _, err := storage.FromContext(ctx).Stat(ctx, RepositoryPublicKey); err != nil {
			storage.Error(w, err)
			return
		}
		host := strings.TrimSuffix(r.Host, "/")
		user, pass, _ := r.BasicAuth()
		url := strings.TrimSuffix(fmt.Sprintf("%s://%s/%s", packages.Scheme(r), host, strings.TrimPrefix(strings.TrimSuffix(r.URL.Path, ".repo"), "/")), "/")

		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		if err := repoDefinition(w, name, url, RepositoryPublicKey, user, pass); err != nil {
			logger.C(r.Context()).WithError(err).Error("failed to execute template")
		}
	}
}

func (p *provider) setup(repo string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		if _, err := storage.FromContext(ctx).Stat(ctx, RepositoryPublicKey); err != nil {
			storage.Error(w, err)
			return
		}
		var name string
		if repo != "" {
			name = strings.NewReplacer("/", "-").Replace(repo)
		} else {
			name = strings.NewReplacer("/", "-", ".", "-").Replace(strings.TrimPrefix(strings.Split(r.Host, ":")[0], Name+"."))
		}
		user, pass, _ := r.BasicAuth()
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
			Method:  http.MethodGet,
			Path:    ".repo",
			Handler: p.config,
		},
		{
			Method:  http.MethodGet,
			Path:    "/setup",
			Handler: p.setup,
		},
		{
			Method: http.MethodPut,
			Path:   "/push",
			Handler: packages.Push(func(r *http.Request, reader io.Reader, size int64, key string) (storage.Artifact, error) {
				return NewPackage(reader, size, key)
			}),
		},
		{
			Method: http.MethodGet,
			Path:   "/repodata/{filename}",
			Handler: packages.Pull(func(r *http.Request) string {
				return mux.Vars(r)["filename"]
			}),
		},
		{
			Method: http.MethodGet,
			Path:   "/{filename}",
			Handler: packages.Pull(func(r *http.Request) string {
				return mux.Vars(r)["filename"]
			}),
		},
		{
			Method: http.MethodDelete,
			Path:   "/{filename}",
			Handler: packages.Delete(func(r *http.Request) string {
				return mux.Vars(r)["filename"]
			}),
		},
	}
}
