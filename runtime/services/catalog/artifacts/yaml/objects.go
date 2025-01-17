package yaml

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/jinzhu/copier"
	runtimev1 "github.com/rilldata/rill/proto/gen/rill/runtime/v1"
	"github.com/rilldata/rill/runtime/drivers"
	"github.com/rilldata/rill/runtime/pkg/fileutil"
	"google.golang.org/protobuf/types/known/structpb"
)

/**
 * This file contains the mapping from CatalogObject to Yaml files
 */

const Version = "0.0.1"

type Source struct {
	Version string
	Type    string
	URI     string `yaml:"uri,omitempty"`
	Path    string `yaml:"path,omitempty"`
	Region  string `yaml:"region,omitempty"`
}

type MetricsView struct {
	Version          string
	DisplayName      string `yaml:"display_name"`
	Description      string
	From             string
	TimeDimension    string `yaml:"time_dimension"`
	TimeGrains       []string
	DefaultTimeGrain string `yaml:"default_timegrain"`
	Dimensions       []*Dimension
	Measures         []*Measure
}

type Measure struct {
	Label       string
	Expression  string
	Description string
	Format      string `yaml:"format_preset"`
}

type Dimension struct {
	Label       string
	Property    string `copier:"Name"`
	Description string
}

func toSourceArtifact(catalog *drivers.CatalogEntry) (*Source, error) {
	source := &Source{
		Version: Version,
		Type:    catalog.GetSource().Connector,
	}

	props := catalog.GetSource().Properties.AsMap()
	path, ok := props["path"].(string)
	if ok {
		if catalog.GetSource().Connector == "file" {
			source.Path = path
		} else {
			source.URI = path
		}
	}
	region, ok := props["aws.region"].(string)
	if ok {
		source.Region = region
	}

	return source, nil
}

func toMetricsViewArtifact(catalog *drivers.CatalogEntry) (*MetricsView, error) {
	metricsArtifact := &MetricsView{}
	err := copier.Copy(metricsArtifact, catalog.Object)
	if err != nil {
		return nil, err
	}

	metricsArtifact.Version = Version
	return metricsArtifact, nil
}

func fromSourceArtifact(source *Source, path string) (*drivers.CatalogEntry, error) {
	props := map[string]interface{}{}
	if source.Type == "file" {
		props["path"] = source.Path
	} else {
		props["path"] = source.URI
	}
	if source.Region != "" {
		props["aws.region"] = source.Region
	}
	propsPB, err := structpb.NewStruct(props)
	if err != nil {
		return nil, err
	}

	name := strings.TrimSuffix(filepath.Base(path), fileutil.FullExt(path))
	return &drivers.CatalogEntry{
		Name: name,
		Type: drivers.ObjectTypeSource,
		Path: path,
		Object: &runtimev1.Source{
			Name:       name,
			Connector:  source.Type,
			Properties: propsPB,
		},
	}, nil
}

func fromMetricsViewArtifact(metrics *MetricsView, path string) (*drivers.CatalogEntry, error) {
	apiMetrics := &runtimev1.MetricsView{}
	err := copier.Copy(apiMetrics, metrics)
	if err != nil {
		return nil, err
	}

	// this is needed since measure names are not given by the user
	for i, measure := range apiMetrics.Measures {
		measure.Name = fmt.Sprintf("measure_%d", i)
	}

	name := strings.TrimSuffix(filepath.Base(path), fileutil.FullExt(path))
	apiMetrics.Name = name
	return &drivers.CatalogEntry{
		Name:   name,
		Type:   drivers.ObjectTypeMetricsView,
		Path:   path,
		Object: apiMetrics,
	}, nil
}
