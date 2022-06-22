package json_test

import (
	"encoding/json"
	"testing"

	jsoniter "github.com/json-iterator/go"
	"github.com/ohler55/ojg/oj"
)

const JSON = `{
   "min_position": 4,
   "has_more_items": false,
   "items_html": "Bike",
   "new_latent_count": 6,
   "data": {
      "length": 22,
      "text": "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat."
   },
   "numericalArray": [
      31,
      24,
      30,
      31,
      33
   ],
   "StringArray": [
      "Oxygen",
      "Carbon",
      "Carbon",
      "Carbon"
   ],
   "multipleTypesArray": true,
   "objArray": [
      {
         "class": "upper",
         "age": 1
      },
      {
         "class": "middle",
         "age": 6
      },
      {
         "class": "upper",
         "age": 0
      },
      {
         "class": "upper",
         "age": 0
      },
      {
         "class": "lower",
         "age": 0
      }
   ]
}`

var bytesJSON = []byte(JSON)

func BenchmarkOjg(b *testing.B) {
	b.ReportAllocs()
	b.SetBytes(int64(len(JSON)))
	for i := 0; i < b.N; i++ {
		_, err := oj.ParseString(JSON)
		if err != nil {
			panic(err)
		}
	}
}

func BenchmarkJsonIter(b *testing.B) {
	var v map[string]interface{}
	b.ReportAllocs()
	b.SetBytes(int64(len(JSON)))
	for i := 0; i < b.N; i++ {
		err := jsoniter.Unmarshal(bytesJSON, &v)
		if err != nil {
			panic(err)
		}
	}
}

func BenchmarkStd(b *testing.B) {
	var v map[string]interface{}
	b.ReportAllocs()
	b.SetBytes(int64(len(JSON)))
	for i := 0; i < b.N; i++ {
		err := json.Unmarshal(bytesJSON, &v)
		if err != nil {
			panic(err)
		}
	}
}
