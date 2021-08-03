package main

import (
	"database/sql"
	"fmt"
	"math/rand"
	"net/http"

	"github.com/PuerkitoBio/goquery"
	_ "github.com/mattn/go-sqlite3"
)

var query = map[string]string{
	"fetch_verbs": "select infinit, simple, participl from verbs order by id",
	"insert_verb": "insert into verbs(infinit, simple, participl) values(?,?,?)",
}

var host_map = map[string]map[string]string{
	"www.worddy.co": {
		"url":            "https://www.worddy.co/en/list-of-irregular-verbs-english",
		"table_selector": "table.table.table-striped.table-v._fs-mob-14 tbody tr",
	},
}

type Verbs struct {
	infinit   string
	simple    string
	participl string
}

func fetch_verbs(dbname string) []Verbs {
	var verbs []Verbs = make([]Verbs, 0, 50)

	db, err := sql.Open("sqlite3", dbname)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	rows, err := db.Query(query["fetch_verbs"])
	if err != nil {
		fmt.Println(err)
		return nil
	}

	for rows.Next() {
		verb := Verbs{}
		rows.Scan(&verb.infinit, &verb.simple, &verb.participl)
		verbs = append(verbs, verb)
	}

	rows.Close()

	return verbs
}

func pars_site(host string) []Verbs {
	var verbs []Verbs = make([]Verbs, 0, 50)

	res, err := http.Get(host_map[host]["url"])
	if err != nil {
		return nil
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil
	}

	doc.Find(host_map[host]["table_selector"]).Each(func(i int, sel *goquery.Selection) {
		cols := sel.Find("td").Map(func(i int, s *goquery.Selection) string { return s.Text() })
		verbs = append(verbs, Verbs{cols[0], cols[1], cols[2]})
	})

	return verbs
}

func insert_verbs(dbname string, verbs []Verbs) int {

	db, err := sql.Open("sqlite3", dbname)
	if err != nil {
		fmt.Println(err)
		return 1
	}

	stmt, err := db.Prepare(query["insert_verb"])
	if err != nil {
		fmt.Println(err)
		return 1
	}

	for _, verb := range verbs {
		affect, err := stmt.Exec(verb.infinit, verb.simple, verb.participl)
		if err != nil {
			return 1
		}
		fmt.Println("row affected ", affect)
	}

	return 0
}

func run_test(verbs []Verbs) (wrong []Verbs, err_count int) {

	for _, index := range rand.Perm(len(verbs)) {
		var simple, participl string
		fmt.Println("Infinitive ", verbs[index].infinit)

		fmt.Scanln(&simple)
		fmt.Scanln(&participl)
		fmt.Printf("Past simple verb %s your choice %s", verbs[index].simple, simple)
		fmt.Printf("Past simple verb %s your choice %s", verbs[index].participl, participl)
	}
	return verbs, 0
}

func main() {

	// insert_verbs("verbs.db", pars_site("www.worddy.co"))
	// fmt.Println(fetch_verbs("verbs.db"))
	// fmt.Println("test")
	run_test(fetch_verbs("verbs.db"))
}
