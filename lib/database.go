package lib

import (
    "database/sql"
    "log"
    "database/sql/driver"
    _ "github.com/mattn/go-sqlite3"
)

type Msg struct {
    come string
    send string
}

type Database struct {
    Conn *sql.DB
    Msgs map[string]string
    Tx driver.Tx
}

func (db Database) Registration_msg(come string, send string) {
    db.Msgs[come] = send
    s := "insert into msg values (" + "'" + come + "'" + "," + "'" + send + "'" + ")"
    db.Conn.Exec(s)
    db.Tx.Commit()
}

func (db Database) Delete_msg(come string) {
    delete(db.Msgs, come)
    db.Conn.Exec("delete from msg where come_msg = " + "'" + come + "'")
    db.Tx.Commit()
}

func SetupDB() (Database) {
    db, err := sql.Open("sqlite3", "./messages.db")
    if err != nil {
        panic(err)
    }
    var msg []Msg
    rows,_ := db.Query("select * from msg")
    defer rows.Close()
    for rows.Next() {
        var m Msg
        err := rows.Scan(&m.come, &m.send)
        if err != nil {
            log.Println(err)
        }
        msg = append(msg, m)
    }
    msgs := make(map[string]string)
    for _,key := range msg {
        msgs[key.come] = key.send
    }
    tx, err := db.Begin()
    return Database{Conn: db, Msgs: msgs, Tx: tx}

}
