package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/m-szalik/pact-framework-example/go-api-server/model"
	"log/slog"
	"net/http"
	"os/signal"
	"strconv"
	"syscall"
)

type server struct {
	books []model.ServerBook
}

func (s *server) list(writer http.ResponseWriter, _ *http.Request) {
	slog.Info("http request 'list'")
	handleResponse(s.books, http.StatusOK, writer)
}

func (s *server) byID(writer http.ResponseWriter, request *http.Request) {
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
	handleResponse(nil, http.StatusNoContent, writer)
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
		_, _ = writer.Write(buff)
	}
}

func (s *server) httpHandler() http.Handler {
	router := mux.NewRouter()
	router.Path("/books").HandlerFunc(s.list)
	router.Path("/books/{id}").HandlerFunc(s.byID)
	return router
}

func main() {
	ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	s := server{
		books: prepareData(),
	}

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

func prepareData() []model.ServerBook {
	return []model.ServerBook{
		{ID: 0, Title: "The Go Programming Language"},
		{ID: 1, Title: "Clean Code"},
		{ID: 2, Title: "Introduction to Algorithms"},
		{ID: 3, Title: "The Pragmatic Programmer"},
		{ID: 4, Title: "Design Patterns"},
		{ID: 5, Title: "Effective Java"},
		{ID: 6, Title: "Refactoring"},
		{ID: 7, Title: "Structure and Interpretation of Computer Programs"},
		{ID: 8, Title: "Domain-Driven Design"},
		{ID: 9, Title: "Code Complete"},
	}
}
