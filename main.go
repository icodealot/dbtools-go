package main

import (
	"fmt"
	"log"

	"github.com/icodealot/dbtools-go/example"
)

var sql string = `
select 'What is the ultimate answer?' as "QUESTION" from dual;
select 42 as "ANSWER" from dual;`

func main() {
	fmt.Println("Hello, DBTOOLS!")
	cfg := example.DBToolsConfig{
		ConnectionId: "ocid1.databasetoolsconnection.oc1.phx.change-me",
		ContentType:  "application/sql",
		Payload:      sql,
	}

	rawBytes, err := example.ExecuteDBToolsConnection(cfg)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(string(rawBytes))
}
