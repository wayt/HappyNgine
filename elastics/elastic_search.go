package elastics

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/wayt/happyngine/env"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
)

var Config *ESConfig

type ESConfig struct {
	Hosts    []string
	Index    string
	Username string
	Password string
}

func init() {

	Config = &ESConfig{
		Hosts:    []string{env.Get("ELASTICSEARCH_PORT_9200_TCP_ADDR") + ":" + env.Get("ELASTICSEARCH_PORT_9200_TCP_PORT")},
		Index:    env.Get("HAPPY_ELASTICSEARCH_INDEX"),
		Username: env.Get("HAPPY_ELASTICSEARCH_USERNAME"),
		Password: env.Get("HAPPY_ELASTICSEARCH_PASSWORD"),
	}
}

func randomHost() string {

	// Set default host target
	host := Config.Hosts[0]

	// Select a random host if possible
	if l := len(Config.Hosts); l > 1 {
		host = Config.Hosts[rand.Intn(l)]
	}

	return host
}

func newRequest(method, path string, body []byte) (req *http.Request, err error) {

	host := &url.URL{
		Scheme: "http",
		Host:   randomHost(),
		Path:   path,
	}

	if body != nil {

		req, err = http.NewRequest(method, host.String(), bytes.NewReader(body))
	} else {

		req, err = http.NewRequest(method, host.String(), nil)
	}

	if err != nil {
		return
	}

	if Config.Username != "" {

		req.SetBasicAuth(Config.Username, Config.Password)
	}

	return req, err
}

func do(req *http.Request, respStruct interface{}) error {

	log.Debugln("elastic.do:", req.URL.String())
	resp, err := http.DefaultClient.Do(req)
	log.Debugln("elastic.do: done:", req.URL.String())
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if respStruct != nil {

		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		if err := json.Unmarshal(respBody, &respStruct); err != nil {
			return err
		}
	}

	return nil
}

type IndexResult struct {
	Index   string `json:"_index"`
	Type    string `json:"_type"`
	Id      string `json:"_id"`
	Version int    `json:"_version"`
	Created bool   `json:"created"`
}

// Put a document at an index
func Index(_type, id string, obj interface{}) (*IndexResult, error) {

	body, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}

	req, err := newRequest("PUT", fmt.Sprintf("/%s/%s/%s", Config.Index, _type, id), body)

	var resp IndexResult
	if err := do(req, &resp); err != nil {

		return nil, err
	}

	return &resp, nil
}

type SearchHitHighlight map[string][]string

type SearchExplanation struct {
	Value       float64             `json:"value"`             // e.g. 1.0
	Description string              `json:"description"`       // e.g. "boost" or "ConstantScore(*:*), product of:"
	Details     []SearchExplanation `json:"details,omitempty"` // recursive details
}

type SearchHit struct {
	Score       *float64               `json:"_score"`       // computed score
	Index       string                 `json:"_index"`       // index name
	Id          string                 `json:"_id"`          // external or internal
	Type        string                 `json:"_type"`        // type
	Version     *int64                 `json:"_version"`     // version number, when Version is set to true in SearchService
	Sort        []interface{}          `json:"sort"`         // sort information
	Highlight   SearchHitHighlight     `json:"highlight"`    // highlighter information
	Source      *json.RawMessage       `json:"_source"`      // stored document source
	Fields      map[string]interface{} `json:"fields"`       // returned fields
	Explanation *SearchExplanation     `json:"_explanation"` // explains how the score was computed
}

type SearchHits struct {
	TotalHits int64        `json:"total"`     // total number of hits found
	MaxScore  *float64     `json:"max_score"` // maximum score of all hits
	Hits      []*SearchHit `json:"hits"`      // the actual hits returned
}

type SearchSuggestionOption struct {
	Text    string      `json:"text"`
	Score   float32     `json:"score"`
	Freq    int         `json:"freq"`
	Payload interface{} `json:"payload"`
}

type SearchSuggestion struct {
	Text    string                   `json:"text"`
	Offset  int                      `json:"offset"`
	Length  int                      `json:"length"`
	Options []SearchSuggestionOption `json:"options"`
}

type SearchSuggest map[string][]SearchSuggestion

// searchFacetTerm is the result of a terms facet.
// See http://www.elasticsearch.org/guide/en/elasticsearch/reference/current/search-facets-terms-facet.html.
type searchFacetTerm struct {
	Term  string `json:"term"`
	Count int    `json:"count"`
}

// searchFacetRange is the result of a range facet.
// See http://www.elasticsearch.org/guide/en/elasticsearch/reference/current/search-facets-range-facet.html.
type searchFacetRange struct {
	From       *float64 `json:"from"`
	FromStr    *string  `json:"from_str"`
	To         *float64 `json:"to"`
	ToStr      *string  `json:"to_str"`
	Count      int      `json:"count"`
	Min        *float64 `json:"min"`
	Max        *float64 `json:"max"`
	TotalCount int      `json:"total_count"`
	Total      *float64 `json:"total"`
	Mean       *float64 `json:"mean"`
}

// searchFacetEntry is a general facet entry.
// See http://www.elasticsearch.org/guide/en/elasticsearch/reference/current/search-facets.html
type searchFacetEntry struct {
	// Key for this facet, e.g. in histograms
	Key interface{} `json:"key"`
	// Date histograms contain the number of milliseconds as date:
	// If e.Time = 1293840000000, then: Time.at(1293840000000/1000) => 2011-01-01
	Time int64 `json:"time"`
	// Number of hits for this facet
	Count int `json:"count"`
	// Min is either a string like "Infinity" or a float64.
	// This is returned with some DateHistogram facets.
	Min interface{} `json:"min,omitempty"`
	// Max is either a string like "-Infinity" or a float64
	// This is returned with some DateHistogram facets.
	Max interface{} `json:"max,omitempty"`
	// Total is the sum of all entries on the recorded Time
	// This is returned with some DateHistogram facets.
	Total float64 `json:"total,omitempty"`
	// TotalCount is the number of entries for Total
	// This is returned with some DateHistogram facets.
	TotalCount int `json:"total_count,omitempty"`
	// Mean is the mean value
	// This is returned with some DateHistogram facets.
	Mean float64 `json:"mean,omitempty"`
}

type SearchFacet struct {
	Type    string             `json:"_type"`
	Missing int                `json:"missing"`
	Total   int                `json:"total"`
	Other   int                `json:"other"`
	Terms   []searchFacetTerm  `json:"terms"`
	Ranges  []searchFacetRange `json:"ranges"`
	Entries []searchFacetEntry `json:"entries"`
}

type SearchFacets map[string]*SearchFacet

type Aggregations map[string]json.RawMessage

type SearchResult struct {
	TookInMillis int64         `json:"took"`            // search time in milliseconds
	ScrollId     string        `json:"_scroll_id"`      // only used with Scroll and Scan operations
	Hits         *SearchHits   `json:"hits"`            // the actual search hits
	Suggest      SearchSuggest `json:"suggest"`         // results from suggesters
	Facets       SearchFacets  `json:"facets"`          // results from facets
	Aggregations Aggregations  `json:"aggregations"`    // results from aggregations
	TimedOut     bool          `json:"timed_out"`       // true if the search timed out
	Error        string        `json:"error,omitempty"` // used in MultiSearch only
}

func Search(_type, query string) (*SearchResult, error) {

	req, err := newRequest("POST", fmt.Sprintf("/%s/%s/_search", Config.Index, _type), []byte(query))
	if err != nil {
		return nil, err
	}

	resp := new(SearchResult)
	if err := do(req, resp); err != nil {

		return nil, err
	}

	if len(resp.Error) > 0 {
		return resp, errors.New(resp.Error)
	}

	return resp, nil
}

type GetResult struct {
	Index   string           `json:"_index"`
	Type    string           `json:"_type"`
	Id      string           `json:"_id"`
	Version int64            `json:"_version,omitempty"`
	Source  *json.RawMessage `json:"_source,omitempty"`
	Found   bool             `json:"found,omitempty"`
	Fields  []string         `json:"fields,omitempty"`
	Error   string           `json:"error,omitempty"` // used only in MultiGet
}

// Get a document
func Get(_type, id string) (*GetResult, error) {

	req, err := newRequest("GET", fmt.Sprintf("/%s/%s/%s", Config.Index, _type, id), nil)
	if err != nil {
		return nil, err
	}

	var resp GetResult
	if err := do(req, &resp); err != nil {

		return nil, err
	}

	return &resp, nil

}

type DeleteResult struct {
	Found   bool   `json:"found"`
	Index   string `json:"_index"`
	Type    string `json:"_type"`
	Id      string `json:"_id"`
	Version int    `json:"_version"`
}

// Delete a document
func Delete(_type, id string) (*DeleteResult, error) {

	req, err := newRequest("DELETE", fmt.Sprintf("/%s/%s/%s", Config.Index, _type, id), nil)
	if err != nil {
		return nil, err
	}

	resp := new(DeleteResult)
	if err := do(req, resp); err != nil {

		return nil, err
	}

	return resp, nil

}

type PingResult struct {
	Status      int    `json:"status"`
	Name        string `json:"name"`
	ClusterName string `json:"cluster_name"`
	Version     struct {
		Number         string `json:"number"`
		BuildHash      string `json:"build_hash"`
		BuildTimestamp string `json:"build_timestamp"`
		BuildSnapshot  bool   `json:"build_snapshot"`
		LuceneVersion  string `json:"lucene_version"`
	} `json:"version"`
	TagLine string `json:"tagline"`
}

func Ping() (*PingResult, error) {

	req, err := newRequest("GET", "", nil)
	if err != nil {
		return nil, err
	}

	var resp PingResult
	if err := do(req, &resp); err != nil {

		return nil, err
	}

	return &resp, nil
}
