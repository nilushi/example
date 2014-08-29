package main

import (
	"fmt"
	"strings"

	"github.com/wojka/arc64/solr"
)

func main() {

	s, err := solr.NewSolr(`http://192.168.1.100:8983/solr`, `wojka`, 2)
	if err != nil {
		fmt.Println(err)
		return
	}
	//	addDocuments(s)

	fmt.Println(`================================== Fetch by Distance =====================================`)

	result, err := FetchByDistance(s, 7.3, 79.5, []string{`SIDEJOB`}, 0, 100, 100000)

	if err != nil {
		fmt.Println(err)
	}

	for i, v := range result.Documents {
		fmt.Printf("result[%v] : Title (%v), Skills (%v)\n", i, v[`keywords`], v[`skills`])
	}

	fmt.Println(`================================== Fetch by Relevance and Distance  =====================================`)

	result, err = FetchByRelevanceAndDistance(s, `trainer`, 7.3, 79.5, []string{`SIDEJOB`}, 0, 100, 100000)
	if err != nil {
		fmt.Println(err)
	}

	for i, v := range result.Documents {
		fmt.Printf("result[%v] : Title (%v), Skills (%v)\n", i, v[`keywords`], v[`skills`])
	}

	fmt.Println(`================================== Fetch Sample Records by Relevance =====================================`)

	result, err = FetchSampleRecordsByRelevance(s, `tainer`, 0, 100)
	if err != nil {
		fmt.Println(err)
	}

	for i, v := range result.Documents {
		fmt.Printf("result[%v] : Title (%v), Skills (%v)\n", i, v[`keywords`], v[`skills`])
	}

}

func FetchSampleRecordsByRelevance(s *solr.Solr, keywords string, startIndex int, limit int) (result *solr.Result, err error) {

	query := solr.NewQuery(solr.AddFuzzyLogic(keywords), startIndex, limit)
	query.DefType(`edismax`)
	query.QueryFields(`keywords^0.2 skills^0.7`)
	query.FilterQuery(`accounttype:EXAMPLE`)
	query.Sort(`score desc`)
	query.SpellcheckQuery(keywords)
	result, err = s.Search(query)

	return

}

func FetchByDistance(s *solr.Solr, lat float64, long float64, labels []string, startIndex int, limit int, maxDistance int) (result *solr.Result, err error) {

	var filterQuery string
	if len(labels) > 0 {
		filterQuery = "labels: (" + strings.Join(labels, " AND ") + ")"
	}

	query := solr.NewQuery(`*:*`, startIndex, limit)
	query.FilterQuery(filterQuery)
	query.Geofilt(lat, long, `location`, maxDistance)
	query.FilterQuery(`accounttype:EXAMPLE`)
	query.Sort(`geodist() asc`)
	result, err = s.Search(query)

	return

}

func FetchByRelevanceAndDistance(s *solr.Solr, keywords string, lat float64, long float64, labels []string, startIndex int, limit int, maxDistance int) (result *solr.Result, err error) {

	var filterQuery string
	if len(labels) > 0 {
		filterQuery = "labels: (" + strings.Join(labels, " AND ") + ")"
	}
	query := solr.NewQuery(solr.AddFuzzyLogic(keywords), startIndex, limit)
	query.DefType(`edismax`)
	query.FilterQuery(filterQuery)
	query.QueryFields(`keywords^0.2 skills^0.7`)
	query.FilterQuery(`accounttype:FREE`)
	query.SpacialParam(lat, long, `location`, maxDistance)
	query.Sort(`score desc`)
	query.BoostFunction(solr.Recip(`geodist()`, 1, 1, 1, 4))
	result, err = s.Search(query)
	return

}

func addDocuments(s *solr.Solr) {
	doc1 := map[string]interface{}{`id`: `doc3`, `keywords`: `baby sitter`, `location`: `43.697225, -79.404949`, `lables`: []string{`SIDEJOB`}, `skills`: []string{`nanny`, `baby sitter`}, `accounttype`: `EXAMPLE`}
	err := s.AddDocuments(doc1)
	fmt.Println(err)

}
