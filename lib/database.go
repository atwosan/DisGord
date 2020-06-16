package lib

import (
	"database/sql"
	"database/sql/driver"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

type recv struct {
	come string
	send string
}

type Database struct {
	Conn *sql.DB
	Msgs map[string]string
	Tx   driver.Tx
}

func (db Database) Add_msg(f, s, tp string) {
	db.Msgs[f] = s
	sl := "insert into " + tp + " values (" + "'" + f + "'" + "," + "'" + s + "'" + ")"
	db.Conn.Exec(sl)
	db.Tx.Commit()
}

func (db Database) Delete_msg(come, tp, t string) {
	delete(db.Msgs, come)
	db.Conn.Exec("delete from " + tp + " where " + t + " = " + "'" + come + "'")
	db.Tx.Commit()
}

func SetupDB() (Database, Database) {
	db, err := sql.Open("sqlite3", "./data.db")
	if err != nil {
		log.Fatal(err)
	}
	var msg []recv
	var ngword []recv
	rows, _ := db.Query("select * from msg")
	rows_ng, _ := db.Query("select * from ngword")
	defer rows.Close()
	defer rows_ng.Close()
	for rows.Next() {
		var m recv
		if err := rows.Scan(&m.come, &m.send); err != nil {
			log.Println(err)
		}
		msg = append(msg, m)
	}
	for rows.Next() {
		var m recv
		if err := rows_ng.Scan(&m.come, &m.send); err != nil {
			log.Println(err)
		}
		ngword = append(ngword, m)
	}
	ngwords := make(map[string]string)
	for _, key := range ngword {
		ngwords[key.come] = key.send
	}
	tx, _ := db.Begin()
	msgs := make(map[string]string)
	for _, key := range msg {
		msgs[key.come] = key.send
	}
	return Database{Conn: db, Msgs: msgs, Tx: tx}, Database{Conn: db, Msgs: ngwords, Tx: tx}

}
