package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/m-szalik/pact-framework-example/go-api-server/model"
	"io"
	"log/slog"
	"net/http"
	"os/signal"
	"strconv"
	"syscall"
)

type server struct {
	books []model.ServerBook
}

func (s *server) getList(writer http.ResponseWriter, _ *http.Request) {
	slog.Info("http request 'get list'")
	handleResponse(s.books, http.StatusOK, writer)
}

func (s *server) getByID(writer http.ResponseWriter, request *http.Request) {
	paramID, ok := mux.Vars(request)["id"]
	slog.Info(fmt.Sprintf("http request 'byId' with id=%s", paramID))
	if !ok {
		handleResponse(errors.New("missing id"), http.StatusBadRequest, writer)
		return
	}
	id, err := strconv.Atoi(paramID)
	if err != nil {
		handleResponse(err, http.StatusBadRequest, writer)
		return
	}
	for _, book := range s.books {
		if book.ID == id {
			handleResponse(book, http.StatusOK, writer)
			return
		}
	}
	handleResponse(nil, http.StatusNotFound, writer)
}

func (s *server) deleteById(writer http.ResponseWriter, request *http.Request) {
	// TODO implement, Pact does not check if the logic is ok but if the request and response match expectations
}

func handleResponse(output any, httpStatus int, writer http.ResponseWriter) {
	switch x := output.(type) {
	case nil:
		writer.WriteHeader(httpStatus)
	case error:
		writer.WriteHeader(httpStatus)
		str := fmt.Sprintf("Error: %s\n", x.Error())
		_, _ = writer.Write([]byte(str))
	default:
		buff, err := json.MarshalIndent(output, "", "  ")
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			_, _ = writer.Write([]byte(err.Error()))
			return
		}
		writer.Header().Set("Content-Type", "application/json")
		_, _ = writer.Write(buff)
	}
}

func (s *server) httpHandler() http.Handler {
	router := mux.NewRouter()
	router.Methods("GET").Path("/books").HandlerFunc(s.getList)
	router.Methods("POST").Path("/books").HandlerFunc(s.createBook)
	router.Methods("GET").Path("/books/{id}").HandlerFunc(s.getByID)
	router.Methods("DELETE").Path("/books/{id}").HandlerFunc(s.deleteById)
	return router
}

func (s *server) serverStart(ctx context.Context) {
	listenOn := ":8080"
	httpServer := &http.Server{
		Addr:    listenOn,
		Handler: s.httpHandler(),
	}
	go func() {
		<-ctx.Done()
		slog.Info("Closing http server.")
		_ = httpServer.Close()
	}()
	slog.Info(fmt.Sprintf("Server started. Listening on %s", listenOn))
	err := httpServer.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		panic(err)
	}
}

func (s *server) createBook(writer http.ResponseWriter, request *http.Request) {
	buff, _ := io.ReadAll(request.Body)
	writer.Header().Add("Content-Type", "application/json")
	writer.WriteHeader(http.StatusCreated)
	_, _ = writer.Write(buff)
}

func main() {
	ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	s := newServer()
	s.serverStart(ctx)
}

func newServer() *server {
	return &server{
		books: []model.ServerBook{
			{ID: 0, Title: "The Go Programming Language"},
			{ID: 1, Title: "Clean Code"},
			{ID: 5, Title: "Effective Java"},
		},
	}
}
