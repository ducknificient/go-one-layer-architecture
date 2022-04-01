package main

import (
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/julienschmidt/httprouter"
)

type Handler struct {
	dbpool *pgxpool.Pool
}

func getRouter(db *pgxpool.Pool) *httprouter.Router {
	h := Handler{
		dbpool: db,
	}
	router := httprouter.New()

	router.POST("/app/shoes/create", h.sAppShoesCreate)
	router.GET("/app/shoes/list", h.sAppShoesList)
	router.POST("/app/shoes/edit", h.sAppShoesEdit)
	router.POST("/app/shoes/delete", h.sAppShoesDelete)

	return router
}
