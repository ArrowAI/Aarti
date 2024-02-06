package packages

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/gorilla/mux"

	"aarti/pkg/storage"
)

var ErrUnknownProvider = errors.New("unknown provider")

type Provider interface {
	Routes() []*Route
	Repository() storage.Repository
}

type ProviderFactory func(ctx context.Context) (Provider, error)

var providers = map[string]ProviderFactory{}

func Register(name string, factory ProviderFactory) {
	providers[name] = factory
}

func Providers() []string {
	var ret []string
	for k := range providers {
		ret = append(ret, k)
	}
	return ret
}

func Names() []string {
	var ret []string
	for k := range providers {
		ret = append(ret, k)
	}
	return ret
}

func New(ctx context.Context, name string) (Provider, error) {
	f, ok := providers[name]
	if !ok {
		return nil, fmt.Errorf("%s: %w", name, ErrUnknownProvider)
	}
	return f(ctx)
}

func Init(ctx context.Context, r *mux.Router, domain, repo string) error {
	for k, v := range providers {
		p, err := v(ctx)
		if err != nil {
			return err
		}
		mdlw := storage.Middleware(p.Repository())("repo")
		subs := []*mux.Router{r.PathPrefix("/" + k).Subrouter()}
		if domain != "" {
			subs = append(subs, r.Host(k+"."+domain).Subrouter())
		}
		for _, v := range subs {
			v.Use(mdlw)
			for _, vv := range p.Routes() {
				p := vv.Path
				if !strings.HasPrefix(p, "/") {
					p = "/" + p
				}
				if repo == "" {
					p = "/{repo:.+}" + vv.Path
				}
				if err := v.Path(p).Methods(vv.Method).HandlerFunc(makeHandler("", vv.Handler)).GetError(); err != nil {
					return fmt.Errorf("%s: %q: %w", k, vv.Path, err)
				}
			}
		}
	}
	return nil
}
