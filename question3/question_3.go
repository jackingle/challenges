package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
	"time"

	"github.com/olekukonko/tablewriter"
)

type NarutoList struct {
	RequestHash        string  `json:"request_hash"`
	RequestCached      bool    `json:"request_cached"`
	RequestCacheExpiry int     `json:"request_cache_expiry"`
	Results            Results `json:"results"`
	LastPage           int     `json:"last_page"`
}
type Results []struct {
	MalID     int       `json:"mal_id"`
	URL       string    `json:"url"`
	ImageURL  string    `json:"image_url"`
	Title     string    `json:"title"`
	Airing    bool      `json:"airing"`
	Synopsis  string    `json:"synopsis"`
	Type      string    `json:"type"`
	Episodes  int       `json:"episodes"`
	Score     float64   `json:"score"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	Members   int       `json:"members"`
	Rated     string    `json:"rated"`
}

var record = []string{"MalID", "URL", "ImageURL", "Title", "Airing", "Synopsis", "Type", "Episodes", "Score", "StartDate", "EndDate", "Members", "Rated"}

// ByScore implements sort.Interface based on the Score field.
type ByScore Results

func (a ByScore) Len() int { return len(a) }

func (a ByScore) Less(i, j int) bool { return a[i].Score < a[j].Score }
func (a ByScore) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

// ByRated implements sort.Interface based on the Rated field.
type ByRated Results

func (a ByRated) Len() int { return len(a) }

func (a ByRated) Less(i, j int) bool { return a[i].Rated < a[j].Rated }
func (a ByRated) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

// Write a program that sorts all of the results from this GET query https://api.jikan.moe/v3/search/anime?q=naruto by score, then rated.
// Output the data in CSV and also display this in an ASCII table on the console.

const (
	narutoURL             = "https://api.jikan.moe/v3/search/anime?q=naruto"
	sortedListByRated     = "sortedListByRated"
	sortedListByScore     = "sortedListByScore"
	sortedByScoreAndRated = "sortedListByScoreAndRated"
	original              = "original"
)

func main() {
	if err := getScore(); err != nil {
		log.Fatal(err)
	}
}

func getScore() error {
	client := http.DefaultClient

	resp, err := client.Get(narutoURL)
	if err != nil {
		return err
	}
	list := NarutoList{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(body, &list); err != nil {
		return err
	}
	if err := outputAll(list); err != nil {
		return err
	}
	if err := makeAllASCIITables(); err != nil {
		return err
	}
	return nil
}

func sortByScoreAndRated(list NarutoList) NarutoList {

	sort.Sort(ByScore(list.Results))
	sort.Sort(ByRated(list.Results))
	return list
}

func sortByScore(list NarutoList) NarutoList {

	sort.Sort(ByScore(list.Results))
	return list
}

func sortByRated(list NarutoList) NarutoList {

	sort.Sort(ByRated(list.Results))
	return list
}

func outputAll(list NarutoList) error {
	listByScoreAndRated := sortByScoreAndRated(list)
	listByScore := sortByScore(list)
	listByRated := sortByRated(list)
	if err := outputCSV(listByScoreAndRated, sortedByScoreAndRated); err != nil {
		return err
	}

	if err := outputCSV(listByScore, sortedListByScore); err != nil {
		return err
	}

	if err := outputCSV(listByRated, sortedListByRated); err != nil {
		return err
	}

	if err := outputCSV(list, original); err != nil {
		return err
	}

	return nil
}

func outputCSV(list NarutoList, filename string) error {
	newFile := fmt.Sprintf("%v.csv", filename)
	file, err := os.Create(newFile)
	if err != nil {
		return err
	}

	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()
	if err := writer.Write(record); err != nil {
		return err
	}

	for _, value := range list.Results {

		s := []string{fmt.Sprintf("%v", value.MalID), value.URL, value.ImageURL, value.Title, fmt.Sprintf("%v", value.Airing), value.Synopsis, value.Type, fmt.Sprintf("%v", value.Episodes), fmt.Sprintf("%f", value.Score), value.StartDate.String(), value.EndDate.String(), fmt.Sprintf("%v", value.Members), value.Rated}
		err = writer.Write(s)
		if err != nil {
			return err
		}
	}

	return nil
}

func makeAllASCIITables() error {
	if err := makeASCIITable(sortedByScoreAndRated); err != nil {
		return err
	}

	if err := makeASCIITable(sortedListByScore); err != nil {
		return err
	}

	if err := makeASCIITable(sortedListByRated); err != nil {
		return err
	}

	if err := makeASCIITable(original); err != nil {
		return err
	}
	return nil
}

func makeASCIITable(filename string) error {
	newFile := fmt.Sprintf("%v.csv", filename)
	table, err := tablewriter.NewCSV(os.Stdout, newFile, true)
	if err != nil {
		panic(err)
	}
	table.SetAlignment(tablewriter.ALIGN_CENTER) // Set Alignment
	table.Render()
	return nil
}

// var record = []string{"MalID", "URL", "ImageURL", "Title", "Airing", "Synopsis", "Type", "Episodes", "Score", "StartDate", "EndDate", "Members", "Rated"}
//

// func (*NarutoList) CSVheader(w io.Writer) {
// 	cw := csv.NewWriter(w)
// 	cw.Write([]string{"A Key", "B Key", "C Key", "D Key"})
// 	cw.Flush()
// }

// func (rm *NarutoList) CSVrow(w io.Writer) {
// 	cw := csv.NewWriter(w)
// 	cw.Write([]string{rm.Results, rm.B, rm.C, rm.D.F})
// 	cw.Write([]string{rm.A, rm.B, rm.C, rm.D.G})

// 	is, _ := json.Marshal(rm.D.H)
// 	cw.Write([]string{rm.A, rm.B, rm.C, string(is)})
// 	cw.Flush()
// }

// func writeCSV(list NarutoList) {
// 	cw := csv.NewWriter(w)

// 	for i, j := range list.Results {
// 		cw.Write(j)
// 	}
// }
