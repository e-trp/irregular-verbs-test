package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	_ "github.com/mattn/go-sqlite3"
)

const dbname = "verbs.db"

func openDB(dbname string) *sql.DB {
	db, err := sql.Open("sqlite3", dbname)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return db
}

var db = openDB(dbname)

var query = map[string]string{
	"select_verbs": "select infinit, simple, participl from verbs order by id",
	"insert_verb":  "insert into verbs(infinit, simple, participl) values(?, ?, ?)",
	"insert_error": "insert into errors(test_id, source, user_guess) values(?, ?, ?)",
	"start_test":   "insert into tests(datetime_stamp, user_name) values(CURRENT_TIMESTAMP, 'test')",
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

func select_verbs() []Verbs {
	var verbs []Verbs = make([]Verbs, 0, 50)

	rows, err := db.Query(query["select_verbs"])
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

func insert_verbs(verbs []Verbs) int {

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

	reader, err_count := bufio.NewReader(os.Stdin), 0

	ret, err := db.Exec(query["start_test"])

	if err != nil {
		fmt.Println(err)
		return
	}

	test_id, _ := ret.LastInsertId()

	insert_error, err := db.Prepare(query["insert_error"])
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, i := range rand.Perm(len(verbs)) {
		fmt.Println("Infinitive ", verbs[i].infinit)
		fmt.Print("Enter past simple form:")
		simple_guess, _ := reader.ReadString('\n')
		fmt.Println("Enter past participle form:")
		participl_guess, _ := reader.ReadString('\n')
		fmt.Printf("Past simple verb %s your choice %s", verbs[i].simple, simple_guess)
		fmt.Printf("Past simple verb %s your choice %s", verbs[i].participl, participl_guess)

		if strings.EqualFold(verbs[i].simple, simple_guess) {
			insert_error.Exec(test_id, verbs[i].simple, simple_guess)
		}

		if strings.EqualFold(verbs[i].participl, participl_guess) {
			insert_error.Exec(test_id, verbs[i].participl, participl_guess)
		}
	}
	return verbs, 0
}

func main() {

	run_test(select_verbs())

	db.Close()

}
