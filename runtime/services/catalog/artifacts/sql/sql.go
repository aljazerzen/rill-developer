package sql

import (
	"context"
	"errors"
	"path"
	"strings"

	runtimev1 "github.com/rilldata/rill/proto/gen/rill/runtime/v1"
	"github.com/rilldata/rill/runtime/drivers"
	"github.com/rilldata/rill/runtime/pkg/fileutil"
	"github.com/rilldata/rill/runtime/services/catalog/artifacts"
)

/**
 * this package contains code to map an sql file to a catalog object
 */

type artifact struct{}

var NotSupported = errors.New("only model supported for sql")

func init() {
	artifacts.Register(".sql", &artifact{})
}

func (r *artifact) DeSerialise(ctx context.Context, filePath string, blob string) (*drivers.CatalogEntry, error) {
	ext := fileutil.FullExt(filePath)
	fileName := path.Base(filePath)
	name := strings.TrimSuffix(fileName, ext)
	return &drivers.CatalogEntry{
		Type: drivers.ObjectTypeModel,
		Object: &runtimev1.Model{
			Name:    name,
			Sql:     blob,
			Dialect: runtimev1.Model_DIALECT_DUCKDB,
		},
		Name: name,
		Path: filePath,
	}, nil
}

func (r *artifact) Serialise(ctx context.Context, catalogObject *drivers.CatalogEntry) (string, error) {
	if catalogObject.Type != drivers.ObjectTypeModel {
		return "", NotSupported
	}
	return catalogObject.GetModel().Sql, nil
}
