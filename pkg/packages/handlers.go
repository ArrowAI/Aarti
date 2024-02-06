package packages

import (
	"io"
	"net/http"

	"aarti/pkg/utils/logger"

	"aarti/pkg/storage"
)

type ArtifactFactory func(r *http.Request, reader io.Reader, size int64, key string) (storage.Artifact, error)

func Push(fn ArtifactFactory) HandlerFunc {
	return func(_ string) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			var (
				reader io.ReadCloser
				size   int64
			)
			if file, header, err := r.FormFile("file"); err == nil {
				reader, size = file, header.Size
			} else {
				reader, size = r.Body, r.ContentLength
			}
			defer reader.Close()
			logger.C(ctx).Debugf("ensuring storage is initialized")
			s := storage.FromContext(ctx)
			if err := s.Init(ctx); err != nil {
				storage.Error(w, err)
				return
			}
			logger.C(ctx).Debugf("parsing artifact")
			pkg, err := fn(r, reader, size, s.Key())
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer pkg.Close()
			logger.C(ctx).WithFields("name", pkg.Name(), "filepath", pkg.Path(), "arch", pkg.Arch()).Infof("uploading artifact")
			if err := s.Write(ctx, pkg); err != nil {
				storage.Error(w, err)
				return
			}
			w.WriteHeader(http.StatusCreated)
		}
	}
}

func Pull(fn func(r *http.Request) string) HandlerFunc {
	return func(_ string) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if err := storage.FromContext(r.Context()).ServeFile(w, r, fn(r)); err != nil {
				storage.Error(w, err)
				return
			}
		}
	}
}

func Delete(fn func(t *http.Request) string) HandlerFunc {
	return func(_ string) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			if err := storage.FromContext(ctx).Delete(ctx, fn(r)); err != nil {
				storage.Error(w, err)
				return
			}
		}
	}
}

func NotFound(w http.ResponseWriter, _ *http.Request) {
	http.Error(w, "404 not found", http.StatusNotFound)
}
