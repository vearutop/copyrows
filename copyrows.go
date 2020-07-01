package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/url"

	"github.com/Masterminds/squirrel"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

func main() {

	src := flag.String("src", "", "Source DB DSN, example: mysql://user:pass@host.io/dbname, required")
	dst := flag.String("dst", "", "Destination DB DSN, example: postgres://user:pass@host.io:5432/dbname?sslmode=disable, required")
	query := flag.String("query", "", "Query to SELECT of source database, required")
	table := flag.String("table", "", "Destination table dstDrv, required")
	pageSize := flag.Int("page-size", 1000, "Number of rows to INSERT in single statement, required")

	flag.Parse()

	if *src == "" || *dst == "" || *query == "" || *table == "" {
		fmt.Println("Copyrows copies data between databases (MySQL and PostgreSQL are supported).")

		flag.Usage()
		return
	}

	srcConn, err := sql.Open(dsn(*src))
	if err != nil {
		log.Panic(err)
	}

	dstDrv, dstDsn := dsn(*dst)
	dstConn, err := sql.Open(dstDrv, dstDsn)
	if err != nil {
		log.Panic(err)
	}

	var ph squirrel.PlaceholderFormat = squirrel.Dollar
	if dstDrv == "mysql" {
		ph = squirrel.Question
	}

	err = srcConn.Ping()
	if err != nil {
		log.Panic(err)
	}

	err = dstConn.Ping()
	if err != nil {
		log.Panic(err)
	}

	rows, err := srcConn.Query(*query)
	if err != nil {
		log.Panic(err)
	}

	cols, err := rows.Columns()
	if err != nil {
		log.Panic(err)
	}

	ins := squirrel.StatementBuilder.PlaceholderFormat(ph).Insert(*table)
	ins = ins.Columns(cols...)

	total := 0

	flush := func() {
		q, a, err := ins.ToSql()
		if err != nil {
			log.Panic(err)
		}

		res, err := dstConn.Exec(q, a...)
		if err != nil {
			log.Panic(err)
		}

		aff, err := res.RowsAffected()
		if err != nil {
			log.Panic(err)
		}

		log.Printf("rows affected: %d, total inserted: %d\n", aff, total)

		ins = squirrel.StatementBuilder.PlaceholderFormat(ph).Insert(*table)
		ins = ins.Columns(cols...)
	}

	item := 0
	for rows.Next() {
		if item >= *pageSize {
			flush()
			item = 0
		}

		total++
		item++

		// Create a slice of interface{}'s to represent each column,
		// and a second slice to contain pointers to each item in the columns slice.
		columns := make([]interface{}, len(cols))
		columnPointers := make([]interface{}, len(cols))
		for i, _ := range columns {
			columnPointers[i] = &columns[i]
		}

		// Scan the result into the column pointers...
		if err := rows.Scan(columnPointers...); err != nil {
			log.Panic(err)
		}

		vals := make([]interface{}, len(cols))
		for i, _ := range cols {
			val := columnPointers[i].(*interface{})
			vals[i] = *val
		}

		ins = ins.Values(vals...)
	}

	flush()
}

func dsn(d string) (string, string) {
	u, err := url.Parse(d)
	if err != nil {
		log.Panic(err)
	}

	if u.Scheme == "mysql" {
		d := u.User.String()

		if d != "" {
			d += "@"
		}

		d += "tcp(" + u.Host + ")/" + u.RequestURI()
	}

	return u.Scheme, d
}
