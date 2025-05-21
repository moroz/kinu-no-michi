package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "modernc.org/sqlite"
)

func main() {
	conn, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		log.Fatal(err)
	}

	var result string
	if err := conn.QueryRow("select unixepoch()").Scan(&result); err != nil {
		log.Fatal(err)
	}
	fmt.Println(result)

	fmt.Println("Hello, world!")
}
