package util

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"reflect"
)

func NewDB() *sql.DB {
	db, err := sql.Open("sqlite3", "./identifier.sqlite")
	if err != nil {
		log.Fatal("Failed geetting db", err)
	}
	return db
}

//type QueryDB interface {
//	fetch()
//	insert()
//	update()
//}

type QueryDB struct {
	Db *sql.DB
}

func ds2slice(dataStruct interface{}) []any {
	var ds []any

	dsType := reflect.TypeOf(dataStruct)
	for i := 0; i < dsType.NumField(); i++ {
		ds = append(ds, dsType.Field(i).Type.String())
	}
	for i := 0; i < len(ds); i++ {
		switch ds[i] {
		case "int":
			ds[i] = *new(int)
		case "string":
			ds[i] = *new(string)
		}
	}

	return ds
}

func (r *QueryDB) Fetch(table string, columns []string, constraints []string) (*sql.Rows, error) {
	template := "SELECT "
	for _, val := range columns {
		template += val + ","
	}
	template = template[:len(template)-1] + " FROM " + table
	if constraints != nil {
		template += " WHERE "
		for _, val := range constraints {
			template += val + " AND "
		}
		template = template[:len(template)-5]
	}

	rows, err := r.Db.Query(template)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (r *QueryDB) Insert(table string, datas map[string]any) (map[string]int64, error) {
	template := fmt.Sprintf("INSERT INTO %s ", table)
	cols := "("
	vals := "VALUES ("
	for key, val := range datas {
		cols += key + ","
		switch reflect.TypeOf(val).String() {
		case "bool":
			if val == true {
				vals += fmt.Sprintf(`"%v"`, "0") + ","
			} else {
				vals += fmt.Sprintf(`"%v"`, "1") + ","
			}
		default:
			vals += fmt.Sprintf(`"%v"`, val) + ","
		}
	}
	template = template + cols[:len(cols)-1] + ") " + vals[:len(vals)-1] + ")"
	ret, err := r.Db.Exec(template)
	if err != nil {
		log.Fatal("When insert into database, the following error occured: ", err)
		return nil, err
	}
	var stat map[string]int64 = map[string]int64{"rowsAffected": 0, "lastInsertID": 0}
	stat["rowsAffected"], _ = ret.RowsAffected()
	stat["lastInsertID"], _ = ret.LastInsertId()
	return stat, nil
}

func (r *QueryDB) Update(table string, datas map[string]any, constraints []string) (map[string]int64, error) {
	template := fmt.Sprintf("UPDATE %s SET ", table)
	for key, val := range datas {
		template += key + "="
		switch reflect.TypeOf(val).String() {
		case "bool":
			if val == true {
				template += fmt.Sprintf(`"%v"`, "0") + ","
			} else {
				template += fmt.Sprintf(`"%v"`, "1") + ","
			}
		default:
			template += fmt.Sprintf(`"%v"`, val) + ","
		}
	}
	template = template[:len(template)-1]
	if constraints != nil {
		template += " WHERE "

		for _, val := range constraints {
			template += val + " AND "
		}
		template = template[:len(template)-5]
	}

	ret, err := r.Db.Exec(template)
	if err != nil {
		log.Fatal("When update data, the following error occured: ", err)
		return nil, err
	}
	var stat map[string]int64 = map[string]int64{"rowsAffected": 0, "lastInsertID": 0}
	stat["rowsAffected"], _ = ret.RowsAffected()
	stat["lastInsertID"], _ = ret.LastInsertId()
	return stat, nil
}

//func main() {
//	db := newDB()
//
//	r := queryDB{
//		db: db,
//		dataStruct: struct {
//			id   int
//			text string
//			done int
//		}{},
//	}
//	rows, _ := r.fetch("todo", []string{
//		"*",
//	}, []string{
//		`id=3`,
//	})
//
//	ds := ds2slice(r.dataStruct)
//	for rows.Next() {
//		rows.Scan(&ds[0], &ds[1], &ds[2])
//		fmt.Println(ds)
//	}
//
//	a, _ := r.insert("todo", map[string]any{
//		"id":   3,
//		"text": "testing",
//		"done": false,
//	})
//	fmt.Println(a)
//
//	ab, _ := r.update("todo", map[string]any{
//		"text": "hhhh",
//	}, []string{
//		`done="1"`,
//		`id="1"`,
//		`text="a"`,
//	})
//	fmt.Println(ab)
//}
