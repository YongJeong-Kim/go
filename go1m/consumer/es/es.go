package es

import (
	"fmt"
	elasticsearch7 "github.com/elastic/go-elasticsearch/v7"
	"log"
	"strings"
)

func main() {
	es, err := elasticsearch7.NewClient(elasticsearch7.Config{
		Addresses: []string{"http://localhost:9200"},
	})
	if err != nil {
		panic(err)
	}

	//CreateDoc(es, "products")
	createIndex(es)
}

func deleteIndex(es *elasticsearch7.Client, name string) {
}

func deleteAllDocs(es *elasticsearch7.Client, index string) {

}

func CreateDoc(es *elasticsearch7.Client, index string, firstname, lastname, email, gender, ipAddress, createdAt string) {
	q := fmt.Sprintf(`
		{
			"first_name": "%s",
			"last_name": "%s",
			"email": "%s",
			"gender": "%s",
			"ip_address": "%s",
			"created_at": "%s"
		}
	`, firstname, lastname, email, gender, ipAddress, createdAt)
	res, err := es.Index(index, strings.NewReader(q),
		//es.Index.WithDocumentID("1"),
		es.Index.WithPretty(),
	)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	//log.Println(res)
}

func createIndex(es *elasticsearch7.Client) {
	index := "products"
	mapping := `
    {
      "settings": {
        "number_of_shards": 3
      },
      "mappings": {
        "properties": {
          "first_name": {
            "type": "text"
          },
					"last_name": {
            "type": "text"
          },
					"email": {
            "type": "text"
          },
					"gender": {
            "type": "text"
          },
					"ip_address": {
            "type": "text"
          },
					"created_at": {
            "type": "text"
          }
        }
      }
    }`

	res, err := es.Indices.Create(
		index,
		es.Indices.Create.WithBody(strings.NewReader(mapping)),
	)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(res)
}
