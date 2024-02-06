package helm

import (
	"context"
	"encoding/json"

	hrepo "helm.sh/helm/v3/pkg/repo"
	"sigs.k8s.io/yaml"

	"aarti/pkg/codec"
	"aarti/pkg/crypt/openpgp"
	"aarti/pkg/storage"
)

const (
	RepositoryPublicKey  = "repository.key"
	RepositoryPrivateKey = "private.key"
)

var _ storage.Repository = (*repo)(nil)

type repo struct{}

func (r *repo) Index(_ context.Context, _ string, artifacts ...storage.Artifact) ([]storage.Artifact, error) {
	cs := storage.MustAs[*Package](artifacts)
	i := hrepo.NewIndexFile()
	for _, v := range cs {
		if err := i.MustAdd(v.Metadata, v.Path(), "", v.PkgDigest); err != nil {
			return nil, err
		}
	}
	b, err := yaml.Marshal(i)
	if err != nil {
		return nil, err
	}
	return []storage.Artifact{storage.NewFile("index.yaml", b)}, nil
}

func (r *repo) GenerateKeypair() (string, string, error) {
	return openpgp.GenerateKeypair("Artifact Registry", "Helm Registry", "")
}

func (r *repo) KeyNames() (string, string) {
	return RepositoryPublicKey, RepositoryPrivateKey
}

func (r *repo) Codec() storage.Codec {
	return codec.Funcs[storage.Artifact]{
		Format: "json",
		EncodeFunc: func(v storage.Artifact) ([]byte, error) {
			return json.Marshal(v)
		},
		DecodeFunc: func(b []byte) (storage.Artifact, error) {
			var a Package
			return &a, json.Unmarshal(b, &a)
		},
	}
}

func (r *repo) Name() string {
	return "helm"
}
