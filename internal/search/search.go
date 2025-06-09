package search

import (
	"banditsecret/internal/parser"
	"bytes"
	"encoding/json"
	"fmt"
	"log"

	"github.com/elastic/go-elasticsearch"
)

type CaptionEntry = parser.CaptionEntry

func StoreCaptionsToSearchEngine(captions []CaptionEntry) {
	es, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{"http://localhost:9200"},
	})

	if err != nil {
		log.Fatalf("Error creating client %s", err)
	}

	if _, err := es.Info(); err != nil {
		log.Fatalf("Error testing connection: %s", err)
	}

	indexName := "captions"
	_, err = es.Indices.Create(indexName)
	if err != nil {
		fmt.Println("Index creation may have failed or already exists.")
	}

	err = bulkUpload(es, indexName, captions)
	if err != nil {
		log.Fatalf("Bulk upload failed: %s", err)
	}
}

func bulkUpload(es *elasticsearch.Client, indexName string, captions []CaptionEntry) error {

	var buf bytes.Buffer

	for _, caption := range captions {

		docID := fmt.Sprintf("%s-%d", caption.VideoId, caption.Start)
		meta := fmt.Appendf(nil, `{ "index" : { "_index" : "%s", "_id" : "%s" } }%s`, indexName, docID, "\n")

		data, err := json.Marshal(caption)
		if err != nil {
			log.Fatalf("Error marshaling caption: %s", err)
		}
		buf.Write(meta)
		buf.Write(data)
		buf.Write([]byte("\n"))
	}

	res, err := es.Bulk(bytes.NewReader(buf.Bytes()), es.Bulk.WithIndex(indexName))
	if err != nil {
		return fmt.Errorf("bulk indexing error: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("bulk response error: %s", res.String())
	}
	fmt.Println("Bulk upload successful.")
	return nil
}
