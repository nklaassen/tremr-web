package api

import (
	"github.com/gorilla/mux"
	"net/http"
)

type DataStore struct {
	TremorRepo
	MedicineRepo
	ExerciseRepo
	UserRepo
}
type Env struct {
	DataStore
	Reboot chan struct{}
}

func NewRouter(env *Env) *mux.Router {
	r := mux.NewRouter()
	r.PathPrefix("/tremors").Handler(authMiddleware(tremorsRouter(env.DataStore.TremorRepo)))
	r.PathPrefix("/meds").Handler(authMiddleware(medsRouter(env.DataStore.MedicineRepo)))
	r.PathPrefix("/exercises").Handler(authMiddleware(exercisesRouter(env.DataStore.ExerciseRepo)))
	r.PathPrefix("/users").Handler(authMiddleware(userRouter(env.DataStore.UserRepo)))
	r.PathPrefix("/auth").Handler(authRouter(env.DataStore.UserRepo))
	r.Handle("/update", update(env.Reboot)).Methods(http.MethodPost)
	return r
}
