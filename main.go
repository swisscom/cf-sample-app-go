package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	cfenv "github.com/cloudfoundry-community/go-cfenv"
	_ "github.com/go-sql-driver/mysql"
)

var (
	buildstamp string
	githash    string
)

func main() {
	log.SetOutput(os.Stdout)
	var port string
	if port = os.Getenv("PORT"); len(port) == 0 {
		port = "4000"
	}

	cStr := getConnectionStr("mydb")
	log.Printf("Connect str %s ", cStr)
	db, err := sql.Open("mysql", cStr)
	if err != nil {
		log.Fatal("Can not open DB:", err)
	}
	defer db.Close()

	var version string
	db.QueryRow("SELECT VERSION()").Scan(&version)
	log.Printf("Connected to :%s\n", version)

	if err := db.Ping(); err != nil {
		log.Fatal("DB Error Ping", err.Error())
	}

	http.HandleFunc("/", defaultHandler)
	http.HandleFunc("/info", infoHandler)
	log.Printf(fmt.Sprintf("Listening at %s", port))
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}

}

func getConnectionStr(n string) string {
	appEnv, err := cfenv.Current()
	if err != nil {
		log.Fatal("hr")
	}
	mysqlService, err := appEnv.Services.WithName(n)
	if err != nil {
		log.Fatal(err)
	}
	uri, ok := mysqlService.Credentials["uri"].(string)
	if !ok {
		log.Fatal("No valid MariabDB uri\n")
	}
	u, err := url.Parse(uri)
	if err != nil {
		log.Fatal("No valid MariabDB uri\n")
	}
	return fmt.Sprintf("%s@tcp(%s)%s", u.User.String(), u.Host, u.Path)
}

func defaultHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(w, "hello Swisscom cloud!")
}

func infoHandler(w http.ResponseWriter, req *http.Request) {
	r := "Binary INFO:\n"
	r += fmt.Sprintf("buildstamp %s\n", buildstamp)
	r += fmt.Sprintf("githash %s\n", githash)
	r += fmt.Sprintf("\n\nENV Variables\n")
	for _, e := range os.Environ() {
		r += fmt.Sprintf("%s\n", e)
	}
	fmt.Fprintln(w, r)

	return
}
