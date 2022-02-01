package main

import (
	"context"
	"crypto/md5"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
)

func createEnvVars() {
	_ = os.Setenv("CAKE_SUPERADMIN_EMAIL", "supadmin@gmail.com")
	_ = os.Setenv("CAKE_SUPERADMIN_PASSWORD", "myNameIsSuperadmin")
	_ = os.Setenv("CAKE_SUPERADMIN_CAKE", "bestCake")
}

func processEnvVars(u *UserService) {
	passwordDigest := md5.New().Sum([]byte(os.Getenv("CAKE_SUPERADMIN_PASSWORD")))
	supadmin := User{
		Email:          os.Getenv("CAKE_SUPERADMIN_EMAIL"),
		Role:           SuperAdminRole,
		Banned:         false,
		PasswordDigest: string(passwordDigest),
		FavoriteCake:   os.Getenv("CAKE_SUPERADMIN_CAKE"),
		BanHistory:     BanHistory{},
	}
	_ = u.repository.Add(supadmin.Email, supadmin)
}

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
	createEnvVars()
	r := mux.NewRouter()

	r.HandleFunc("/user/register", u.Register).Methods(http.MethodPost)
	r.HandleFunc("/user/jwt", wrapJwt(jwtService, u.JWT)).Methods(http.MethodPost)
	r.HandleFunc("/user/me", jwtService.jwtAuth(u.repository, getCakeHandler)).Methods(http.MethodGet)
	r.HandleFunc("/user/favorite_cake", jwtService.jwtAuth(u.repository, updateCakeHandler)).Methods(http.MethodPut)
	r.HandleFunc("/user/email", jwtService.jwtAuth(u.repository, updateEmailHandler)).Methods(http.MethodPut)
	r.HandleFunc("/user/password", jwtService.jwtAuth(u.repository, updatePasswordHandler)).Methods(http.MethodPut)
	r.HandleFunc("/admin/ban", jwtService.jwtAuthAdmin(u.repository, banHandler)).Methods(http.MethodPost)
	r.HandleFunc("/admin/unban", jwtService.jwtAuthAdmin(u.repository, unbanHandler)).Methods(http.MethodPost)
	r.HandleFunc("/admin/inspect", jwtService.jwtAuthAdmin(u.repository, inspectHandler)).Methods(http.MethodGet)
	r.HandleFunc("/admin/promote", jwtService.jwtAuthSuperAdmin(u.repository, promoteHandler)).Methods(http.MethodPost)
	r.HandleFunc("/admin/fire", jwtService.jwtAuthSuperAdmin(u.repository, fireHandler)).Methods(http.MethodPost)

	processEnvVars(u)

	return r
}

func newLoggingRouter(u *UserService, jwtService *JWTService) *mux.Router {
	createEnvVars()
	r := mux.NewRouter()

	r.HandleFunc("/cake", logRequest(jwtService.jwtAuth(u.repository, getCakeHandler))).Methods(http.MethodGet)
	r.HandleFunc("/user/register", logRequest(u.Register)).Methods(http.MethodPost)
	r.HandleFunc("/user/jwt", logRequest(wrapJwt(jwtService, u.JWT))).Methods(http.MethodPost)
	r.HandleFunc("/user/me", logRequest(jwtService.jwtAuth(u.repository, getCakeHandler))).Methods(http.MethodGet)
	r.HandleFunc("/user/favorite_cake", logRequest(jwtService.jwtAuth(u.repository, updateCakeHandler))).Methods(http.MethodPut)
	r.HandleFunc("/user/email", jwtService.jwtAuth(u.repository, updateEmailHandler)).Methods(http.MethodPut)
	r.HandleFunc("/user/password", logRequest(jwtService.jwtAuth(u.repository, updatePasswordHandler))).Methods(http.MethodPut)
	r.HandleFunc("/admin/ban", logRequest(jwtService.jwtAuthAdmin(u.repository, banHandler))).Methods(http.MethodPost)
	r.HandleFunc("/admin/unban", logRequest(jwtService.jwtAuthAdmin(u.repository, unbanHandler))).Methods(http.MethodPost)
	r.HandleFunc("/admin/inspect", logRequest(jwtService.jwtAuthAdmin(u.repository, inspectHandler))).Methods(http.MethodGet)
	r.HandleFunc("/admin/promote", logRequest(jwtService.jwtAuthSuperAdmin(u.repository, inspectHandler))).Methods(http.MethodPost)
	r.HandleFunc("/admin/fire", logRequest(jwtService.jwtAuthSuperAdmin(u.repository, fireHandler))).Methods(http.MethodPost)

	processEnvVars(u)

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
