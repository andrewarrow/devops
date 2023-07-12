package main

import "fmt"

func Backup(db, ip string) {
	s := `"pg_dump postgres://fred:fred@localhost:5432/%s > %s.sql"`
	Run("aa", ip, fmt.Sprintf(s, db, db))
	s = `"gzip %s.sql"`
	Run("aa", ip, fmt.Sprintf(s, db))
}
