package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
)

func getCakeHandler(w http.ResponseWriter, _ *http.Request, u User, _ UserRepository) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("[" + u.Email + "], your favourite cake is " + u.FavoriteCake))
}

func wrapJwt(
	jwt *JWTService,
	f func(http.ResponseWriter, *http.Request, *JWTService),
) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		f(rw, r, jwt)
	}
}

func newRouter(u *UserService, jwtService *JWTService) *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/user/register", u.Register).Methods(http.MethodPost)
	r.HandleFunc("/user/jwt", wrapJwt(jwtService, u.JWT)).Methods(http.MethodPost)
	r.HandleFunc("/user/me", jwtService.jwtAuth(u.repository, getCakeHandler)).Methods(http.MethodGet)
	r.HandleFunc("/user/favorite_cake", jwtService.jwtAuth(u.repository, updateCakeHandler)).Methods(http.MethodPut)
	r.HandleFunc("/user/email", jwtService.jwtAuth(u.repository, updateEmailHandler)).Methods(http.MethodPut)
	r.HandleFunc("/user/password", jwtService.jwtAuth(u.repository, updatePasswordHandler)).Methods(http.MethodPut)

	return r
}

func newLoggingRouter(u *UserService, jwtService *JWTService) *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/cake", logRequest(jwtService.jwtAuth(u.repository, getCakeHandler))).Methods(http.MethodGet)
	r.HandleFunc("/user/register", logRequest(u.Register)).Methods(http.MethodPost)
	r.HandleFunc("/user/jwt", logRequest(wrapJwt(jwtService, u.JWT))).Methods(http.MethodPost)
	r.HandleFunc("/user/me", logRequest(jwtService.jwtAuth(u.repository, getCakeHandler))).Methods(http.MethodGet)
	r.HandleFunc("/user/favorite_cake", logRequest(jwtService.jwtAuth(u.repository, updateCakeHandler))).Methods(http.MethodPut)
	r.HandleFunc("/user/email", jwtService.jwtAuth(u.repository, updateEmailHandler)).Methods(http.MethodPut)
	r.HandleFunc("/user/password", logRequest(jwtService.jwtAuth(u.repository, updatePasswordHandler))).Methods(http.MethodPut)

	return r
}

func main() {
	users := NewInMemoryUserStorage()
	userService := UserService{repository: users}
	jwtService, err := NewJWTService("pubkey.rsa", "privkey.rsa")
	if err != nil {
		panic(err)
	}

	r := newLoggingRouter(&userService, jwtService)

	srv := http.Server{
		Addr:    ":8080",
		Handler: r,
	}
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	go func() {
		<-interrupt
		ctx, cancel := context.WithTimeout(context.Background(),
			5*time.Second)
		defer cancel()
		srv.Shutdown(ctx)
	}()
	log.Println("Server started, hit Ctrl+C to stop")
	err1 := srv.ListenAndServe()
	if err1 != nil {
		log.Println("Server exited with error:", err)
	}
	log.Println("Good bye :)")
}
