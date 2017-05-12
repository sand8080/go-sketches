package db

import (
	"database/sql"
	"fmt"
	"time"
	"github.com/lib/pq"
	"github.com/sand8080/go-sketches/utils"
	"github.com/sand8080/go-sketches/search"
)

var dumpEvery = 10

func GetDBConnection(host string, port int, user, password, dbname, sslmode string) (*sql.DB, error) {
	connString := fmt.Sprintf("host=%s port=%d user=%s password=%s "+
		"dbname=%s sslmode=%s", host, port, user, password, dbname, sslmode)
	conn, err := sql.Open("postgres", connString)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func DropTables(conn *sql.DB) error {
	fmt.Println("Dropping tables")
	tx, err := conn.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec("DROP TABLE IF EXISTS components")
	if err != nil {
		return err
	}
	_, err = tx.Exec("DROP TABLE IF EXISTS relations")
	if err != nil {
		return err
	}
	_, err = tx.Exec("DROP TABLE IF EXISTS objects")
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	fmt.Println("Tables are dropped")
	return nil
}

func CreateTables(conn *sql.DB) error {
	fmt.Println("Creating tables")
	tx, err := conn.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec("CREATE TABLE IF NOT EXISTS objects (id SERIAL PRIMARY KEY)")
	if err != nil {
		return err
	}
	_, err = tx.Exec("CREATE TABLE IF NOT EXISTS relations (id SERIAL PRIMARY KEY, " +
		"object_id integer NOT NULL, relative_ids integer[])")
	if err != nil {
		return err
	}
	_, err = tx.Exec("CREATE TABLE IF NOT EXISTS components (id SERIAL PRIMARY KEY, " +
		"object_ids integer[])")
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	fmt.Println("Tables are created")
	return nil
}

func get_min_obj_id(conn *sql.DB) (int, error) {
	var min_id int
	err := conn.QueryRow("SELECT MIN(id) FROM objects").Scan(&min_id)
	if err != nil {
		return 0, err
	}
	return min_id, nil
}
func get_max_obj_id(conn *sql.DB) (int, error) {
	var max_id int
	err := conn.QueryRow("SELECT MAX(id) FROM objects").Scan(&max_id)
	if err != nil {
		return 0, err
	}
	return max_id, nil
}

func insertObjects(conn *sql.DB, objs_num int) error {
	fmt.Printf("Inserting %d objects\n", objs_num)
	start := time.Now()

	txn, err := conn.Begin()
	if err != nil {
		return err
	}
	defer txn.Rollback()

	stmt, err := txn.Prepare(pq.CopyIn("objects", "id"))
	if err != nil {
		return err
	}

	start_chunk := time.Now()
	for i := 1; i <= objs_num; i++ {
		_, err = stmt.Exec(i)
		if err != nil {
			return err
		}
		if i % dumpEvery == 0 {
			fmt.Printf("Inserted %d objects in %v\n",
				i, time.Since(start_chunk))
			start_chunk = time.Now()
		}
	}

	if err := stmt.Close(); err != nil {
		return err
	}
	if err := txn.Commit(); err != nil {
		return err
	}
	fmt.Printf("%d objects are inserted in %v\n", objs_num, time.Since(start))
	return nil
}

func min(l, r int) int {
	if r < l {
		return r
	}
	return l
}

func insertRelations(conn *sql.DB, rels_num, min_opr, max_opr int) error {
	fmt.Printf("Inserting %d relations\n", rels_num)
	start := time.Now()

	txn, err := conn.Begin()
	if err != nil {
		return err
	}
	defer txn.Rollback()

	min_obj_id, err := get_min_obj_id(conn)
	if err != nil {
		return nil
	}
	max_obj_id, err := get_max_obj_id(conn)
	if err != nil {
		return nil
	}
	max_opr = min(max_opr, max_obj_id - min_obj_id)

	stmt, err := txn.Prepare(pq.CopyIn("relations", "object_id",
		"relative_ids"))
	if err != nil {
		return err
	}

	start_chunk := time.Now()
	utils.InitRandom()

	for i := 1; i <= rels_num; i++ {
		// Generating relative_ids
		obj_id := utils.Random(min_obj_id, max_obj_id + 1)
		relative_size := utils.Random(min_opr, max_opr + 1)
		relative_ids := make(map[int]bool, relative_size)
		for i := 0; i < relative_size; i++ {
			relative_ids[utils.Random(min_obj_id, max_obj_id + 1)] = true
		}
		delete(relative_ids, obj_id)
		ids := make([]int, 0, len(relative_ids))
		for k := range relative_ids {
			ids = append(ids, k)
		}

		_, err = stmt.Exec(obj_id, pq.Array(ids))
		if err != nil {
			return err
		}

		if i % dumpEvery == 0 {
			fmt.Printf("Inserted %d relations in %v\n",
				i, time.Since(start_chunk))
			start_chunk = time.Now()
		}
	}

	if err := stmt.Close(); err != nil {
		return err
	}
	if err := txn.Commit(); err != nil {
		return err
	}
	fmt.Printf("%d relations are inserted in %v\n",
		rels_num, time.Since(start))
	return nil
}

func createIndexesOnRelations(conn *sql.DB) error {
	fmt.Println("Creating indexes on the relations table")
	start := time.Now()

	txn, err := conn.Begin()
	if err != nil {
		return err
	}
	defer txn.Rollback()

	_, err = txn.Exec("ALTER TABLE relations ADD CONSTRAINT " +
		"relations_object_id_fkey FOREIGN KEY (object_id) " +
		"REFERENCES objects(id)")
	if err != nil {
		return err
	}
	fmt.Printf("Creating FK on relation object_id finished in: %v\n",
		time.Since(start))

	start_obj_id := time.Now()
	fmt.Println("Creating index on relation object_id started")
	_, err = txn.Exec("CREATE INDEX relations_object_id_idx " +
		"ON relations(object_id)")
	if err != nil {
		return err
	}
	fmt.Printf("Creating index on relation object_id finished in: %v\n",
		time.Since(start_obj_id))


	if err := txn.Commit(); err != nil {
		return err
	}
	fmt.Printf("All indexes creation on 'relations' performed in: %v\n",
		time.Since(start))
	return nil
}

func FillTables(conn *sql.DB, objs_num, rels_num, min_opr, max_opr int) error {
	fmt.Println("Filling data into tables")
	if err := insertObjects(conn, objs_num); err != nil {
		return err
	}

	if err := insertRelations(conn, rels_num, min_opr, max_opr); err != nil {
		return err
	}

	if err := createIndexesOnRelations(conn); err != nil {
		return err
	}
	fmt.Println("Data is filled into tables")
	return nil
}

func makeUnion(txn *sql.Tx) (*search.DisjointSetInt, error) {
	union := search.NewDisjointSetInt(0)

	var objs_num int
	err := txn.QueryRow("SELECT COUNT(*) FROM objects").Scan(&objs_num)
	if err != nil {
		return nil, err
	}

	rows, err := txn.Query("SELECT id FROM objects")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	obj_ids_chunks := make([][]int, 0, objs_num / dumpEvery)
	var obj_id, counter int
	chunk := make([]int, 0, dumpEvery)
	for rows.Next() {
		err := rows.Scan(&obj_id)
		if err != nil {
			return nil, err
		}
		chunk = append(chunk, obj_id)
		counter++

		if counter % dumpEvery == 0 {
			obj_ids_chunks = append(obj_ids_chunks, chunk)
			chunk = make([]int, 0, dumpEvery)
		}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	if len(chunk) > 0 {
		obj_ids_chunks = append(obj_ids_chunks, chunk)
	}

	//Fetching relations by ids chunks
	rels_stmt, err := txn.Prepare("SELECT object_id, relative_ids " +
		"FROM relations " +
		"WHERE object_id = ANY($1)")
	if err != nil {
		return nil, err
	}
	defer rels_stmt.Close()

	var object_id int
	var relative_ids []int
	for _, ids := range obj_ids_chunks {
		rows, err := rels_stmt.Query(pq.Array(ids))
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		for rows.Next() {
			err := rows.Scan(&object_id, &relative_ids)
			if err != nil {
				return nil, err
			}
			fmt.Println("data:", object_id, relative_ids)
		}
		if err := rows.Err(); err != nil {
			return nil, err
		}
	}

	//counter = 0
	//tick_every = 100000
	//s = datetime.datetime.utcnow()
	//curs.execute(select_objs)
	//ids_chunks = split_every(curs.fetchall(), 10000)
	//for i, ids in enumerate(ids_chunks):
	//obj_ids = []
	//for row in ids:
	//obj_id = row[0]
	//obj_ids.append(obj_id)
	//union.union((obj_id,))
	//
	//curs.execute(select_rels, {'object_ids': tuple(obj_ids)})
	//for rel_row in curs.fetchall():
	//obj_id = rel_row[0]
	//relative_ids = rel_row[1]
	//union.union([obj_id] + relative_ids)
	//counter += 1
	//if counter % tick_every == 0:
	//print("Processed records num: {0} in {1}".format(
	//	counter, datetime.datetime.utcnow() - s))
	//s = datetime.datetime.utcnow()
	return union, nil
}

func saveDisjointSets(conn *sql.DB, union *search.DisjointSetInt) error {
	start := time.Now()
	fmt.Println("Disjoint sets saving started")
	fmt.Printf("Disjoint sets saving finished at: %v\n",
		time.Since(start))
	return nil
}


func cleanComponents(conn *sql.DB) error {
	fmt.Println("Cleaning components table")
	if _, err := conn.Exec("DELETE FROM components"); err != nil {
		return err
	}
	fmt.Println("Components table is cleaned")
	return nil
}

func calculateDisjoints(conn *sql.DB) (*search.DisjointSetInt, error) {
	fmt.Println("Disjoint sets calculation started")
	start := time.Now()
	txn, err := conn.Begin()
	if err != nil {
		return nil, err
	}
	defer txn.Rollback()

	union, err := makeUnion(txn)
	if err != nil {
		return nil, err
	}
	if err = txn.Commit(); err != nil {
		return nil, err
	}
	fmt.Printf("Disjoint sets calculation finished in: %v\n",
		time.Since(start))
	return union, err
}

func RecalculateDisjoints(conn *sql.DB) error {
	if err := cleanComponents(conn); err != nil {
		return err
	}

	union, err := calculateDisjoints(conn)
	if err != nil {
		return err
	}

	if err := saveDisjointSets(conn, union); err != nil {
		return err
	}
	return nil
}
