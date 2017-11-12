package elastic

import (
	"context"
	"github.com/ohsu-comp-bio/funnel/config"
	elastic "gopkg.in/olivere/elastic.v5"
	"time"
)

// Elastic provides an elasticsearch database server backend.
type Elastic struct {
	client    *elastic.Client
	conf      config.Elastic
	taskIndex string
	nodeIndex string
}

// NewElastic returns a new Elastic instance.
func NewElastic(conf config.Elastic) (*Elastic, error) {
	client, err := elastic.NewClient(
		elastic.SetURL(conf.URL),
		elastic.SetSniff(false),
		elastic.SetRetrier(
			elastic.NewBackoffRetrier(
				elastic.NewExponentialBackoff(time.Millisecond*50, time.Minute),
			),
		),
	)
	if err != nil {
		return nil, err
	}
	return &Elastic{
		client,
		conf,
		conf.IndexPrefix + "-tasks",
		conf.IndexPrefix + "-nodes",
	}, nil
}

// Close closes the database client.
func (es *Elastic) Close() error {
	es.client.Stop()
	return nil
}

func (es *Elastic) initIndex(ctx context.Context, name, body string) error {
	exists, err := es.client.
		IndexExists(name).
		Do(ctx)

	if err != nil {
		return err
	} else if !exists {
		if _, err := es.client.CreateIndex(name).Body(body).Do(ctx); err != nil {
			return err
		}
	}
	return nil
}

// Init initializing the Elasticsearch indices.
func (es *Elastic) Init(ctx context.Context) error {
	taskMappings := `{
    "mappings": {
      "task":{
        "properties":{
          "id": {
            "type": "keyword"
          },
          "state": {
            "type": "keyword"
          },
          "inputs": {
            "type": "nested"
          },
          "logs": {
            "type": "nested",
            "properties": {
              "logs": {
                "type": "nested"
              }
            }
          }
        }
      }
    }
  }`
	if err := es.initIndex(ctx, es.taskIndex, taskMappings); err != nil {
		return err
	}
	if err := es.initIndex(ctx, es.nodeIndex, ""); err != nil {
		return err
	}
	return nil
}

var idFieldSort = elastic.NewFieldSort("id").
	Desc().
	// Handles the case where there are no documents in the index.
	UnmappedType("keyword")

var minimal = elastic.NewFetchSourceContext(true).Include("id", "state")
var basic = elastic.NewFetchSourceContext(true).
	Exclude("logs.logs.stderr", "logs.logs.stdout", "inputs.content")
