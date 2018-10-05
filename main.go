package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/encoding/wkb"
	"github.com/paulmach/orb/encoding/wkt"
)

var schema_feed = `
create table if not exists feed (
	id int(10) auto_increment primary key,
	title varchar(20) null,
	url varchar(50) null
)
`

var schema_feed_items = `
create table if not exists feed_items (
	id int(10) auto_increment primary key,
	title varchar(20) null,
	description varchar(20) null,
	url varchar(50) null,
	feed_id int(10) not null
)
`

var schema_gis = `
create table if not exists gis_test (
  id int(10) not null,
  name varchar(10) not null,
  gis geometry not null,
  geohash varchar(20) generated always as (st_geohash(gis, 10)) virtual,
  primary key (id),
  unique key id (id),
  spatial key idx_gis (gis),
  key idx_geohash (geohash)
)
`

type Feed struct {
	ID    int            `json:"id"`
	Title sql.NullString `json:"title"`
	URL   sql.NullString `json:"url"`
}

type FeedItem struct {
	ID          int            `json:"id"`
	Title       sql.NullString `json:"title"`
	Description sql.NullString `json:"description"`
	URL         sql.NullString `json:"url"`
	FeedID      int            `json:"feed_id" db:"feed_id"`
	Feed        `db:"feed"`
}

func main() {
	mysqlConnectString := "foo:bar@tcp(127.0.0.1)/robot?charset=utf8mb4"
	db, err := sqlx.Connect("mysql", mysqlConnectString)
	if err != nil {
		log.Fatalln(err)
	}
	db.MustExec(schema_feed)
	db.MustExec(schema_feed_items)
	db.MustExec(schema_gis)
	/*tx := db.MustBegin()
	tx.MustExec("insert into feed (id, title, url) values (?, ?, ?)", 1, "a", "a.com")
	tx.MustExec("insert into feed (id, title, url) values (?, ?, ?)", 2, "b", "b.com")
	tx.MustExec("insert into feed_items (id, title, description, url, feed_id) values (?, ?, ?, ?, ?)", 1, "a.item1", nil, nil, 1)
	tx.MustExec("insert into feed_items (id, title, description, url, feed_id) values (?, ?, ?, ?, ?)", 2, "a.item2", nil, nil, 1)
	tx.MustExec("insert into feed_items (id, title, description, url, feed_id) values (?, ?, ?, ?, ?)", 3, "b.item1", nil, nil, 2)
	tx.MustExec("insert into gis_test (id, name, gis) values (?, ?, st_geomfromtext('point(118.803664 32.079682)'))", 1, "a")
	tx.MustExec("insert into gis_test (id, name, gis) values (?, ?, st_geomfromtext('point(118.863329 32.060038)'))", 2, "b")
	tx.MustExec("insert into gis_test (id, name, gis) values (?, ?, st_geomfromtext('point(119.094445 32.13495)'))", 3, "c")
	tx.MustExec("insert into gis_test (id, name, gis) values (?, ?, st_geomfromtext('point(121.483507 31.234282)'))", 4, "d")
	tx.Commit()*/

	var feeds []Feed
	err = db.Select(&feeds, "SELECT * FROM feed")
	fmt.Printf("%#v\n", feeds)
	var items []FeedItem
	sqls := `SELECT
      feed_items.*,
      feed.id "feed.id",
      feed.title "feed.title",
      feed.url "feed.url"
    FROM
      feed_items JOIN feed ON feed_items.feed_id = feed.id`
	err = db.Select(&items, sqls)
	fmt.Printf("%#v\n", items)

	var polygon orb.Point
	query := "SELECT ST_AsBinary(gis) FROM gis_test WHERE name = '张三'"
	row := db.QueryRow(query)
	err = row.Scan(wkb.Scanner(&polygon))
	fmt.Printf("%#v\n", polygon)
	fmt.Printf("%#v\n", wkt.MarshalString(polygon))
	fmt.Printf("%#v\n", wkb.Value(polygon))
	/*p := orb.Point{
		108.9498710632,
		34.2588125935,
	}*/
	_, err = db.Exec(`INSERT INTO gis_test (id, name, gis) VALUES (10, 'e', st_geomfromwkb(?))`, wkb.Value(polygon))
	fmt.Println(err)
}
