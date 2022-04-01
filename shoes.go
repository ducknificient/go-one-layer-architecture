package main

import (
	"context"
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/jackc/pgx/v4/pgxpool"
	jsoniter "github.com/json-iterator/go"
	"github.com/julienschmidt/httprouter"
)

type Shoes struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (h *Handler) sAppShoesCreate(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	b, err := ioutil.ReadAll(r.Body)

	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}

	defer r.Body.Close()

	var request Shoes

	err = jsoniter.Unmarshal(b, &request)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}
	rName := strings.TrimSpace(request.Name)

	if len(rName) == 0 {
		log.Println(err)
		http.Error(w, "empty name", 200)
		return
	}

	lastSeq, err := GetShoesLastSeq(h.dbpool)
	if err != nil {
		//log.Println(err)
		fmt.Println(err.Error())
		http.Error(w, err.Error(), 500)
		return
	}

	newSeq := lastSeq + 1

	err = InsertShoes(h.dbpool, request, newSeq)
	if err != nil {
		log.Println(err)
		// http.Error(w, http.StatusText(500), 500)
		return
	}

	handleDefault(w, http.StatusOK, true, "ok")
}

func (h *Handler) sAppShoesEdit(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	b, err := ioutil.ReadAll(r.Body)

	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}

	defer r.Body.Close()

	var request Shoes

	err = jsoniter.Unmarshal(b, &request)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}

	rID := strings.TrimSpace(request.ID)
	rName := strings.TrimSpace(request.Name)

	if len(rID) == 0 {
		log.Println(err)
		http.Error(w, "empty id", 200)
		return
	}

	if len(rName) == 0 {
		log.Println(err)
		http.Error(w, "empty name", 200)
		return
	}

	err = EditShoes(h.dbpool, request)
	if err != nil {
		log.Println(err)
		// http.Error(w, http.StatusText(500), 500)
		return
	}

	handleDefault(w, http.StatusOK, true, "ok")
}

func (h *Handler) sAppShoesList(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	// We can now access the connection pool directly in our handlers.
	shoes, err := AllShoes(h.dbpool)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}

	fmt.Println(shoes)

	handleDefault(w, http.StatusOK, true, "ok")
}

func (h *Handler) sAppShoesDelete(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	b, err := ioutil.ReadAll(r.Body)

	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}

	defer r.Body.Close()

	var request Shoes

	err = jsoniter.Unmarshal(b, &request)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}

	rID := strings.TrimSpace(request.ID)

	if len(rID) == 0 {
		log.Println(err)
		http.Error(w, "empty id", 200)
		return
	}

	err = DeleteShoes(h.dbpool, request)
	if err != nil {
		log.Println(err)
		// http.Error(w, http.StatusText(500), 500)
		return
	}

	handleDefault(w, http.StatusOK, true, "ok")
}

func GetShoesLastSeq(db *pgxpool.Pool) (int, error) {
	lastSeq := 0
	q := `SELECT seq FROM smaster.shoes WHERE flag != 'delete' ORDER BY seq DESC LIMIT 1`
	err := db.QueryRow(context.Background(), q).Scan(&lastSeq)
	fmt.Println("error ")
	if err != nil && err != sql.ErrNoRows && err.Error() != "no rows in result set" {
		fmt.Println(err.Error())
		return lastSeq, err
	}

	return lastSeq, nil

}

func InsertShoes(db *pgxpool.Pool, req Shoes, seq int) error {

	q := `INSERT INTO smaster.shoes(seq,name,description,createdby,createdip,updatedby,updatedip) VALUES($1,$2,$3,$4,$5,$4,$5)`
	_, err := db.Exec(context.Background(), q, seq, req.Name, req.Description, `backoffice`, ``)
	if err != nil {
		// log.Println(err)
		fmt.Println(err.Error())
	}

	return nil
}

func AllShoes(db *pgxpool.Pool) ([]Shoes, error) {

	var shoes []Shoes

	q := `SELECT id,name,description,flag::text FROM smaster.shoes WHERE flag != 'delete'`
	rows, err := db.Query(context.Background(), q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {

		var (
			ID          sql.NullString
			Name        sql.NullString
			Description sql.NullString
			Flag        sql.NullString
		)

		err := rows.Scan(&ID, &Name, &Description, &Flag)
		if err != nil {
			return nil, err
		}

		shoe := Shoes{
			ID:          ID.String,
			Name:        Name.String,
			Description: Description.String,
		}

		shoes = append(shoes, shoe)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return shoes, nil
}

func EditShoes(db *pgxpool.Pool, req Shoes) error {

	q := `UPDATE smaster.shoes SET name=$2,description=$3,updatedby=$4,updatedip=$5,flag='update' WHERE id = $1`
	_, err := db.Exec(context.Background(), q, req.ID, req.Name, req.Description, `backoffice`, ``)
	if err != nil {
		// log.Println(err)
		fmt.Println(err.Error())
	}

	return nil
}

func DeleteShoes(db *pgxpool.Pool, req Shoes) error {

	q := `UPDATE smaster.shoes SET flag='delete',updatedby=$2,updatedip=$3 WHERE id = $1`
	_, err := db.Exec(context.Background(), q, req.ID, `backoffice`, ``)
	if err != nil {
		// log.Println(err)
		fmt.Println(err.Error())
	}

	return nil
}
