package api

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type User struct {
	Uid      int64  `json:"uid"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

type UserRepo interface {
	Add(*User) error
	GetFromUid(int64) (User, error)
	GetFromEmail(string) (User, error)
}

func userRouter(repo UserRepo) *mux.Router {
	router := mux.NewRouter()
	router.Handle("/users/{uid}", getUser(repo)).Methods(http.MethodGet)
	return router
}

func getUser(userRepo UserRepo) HttpErrorHandler {
	type UserWithoutPassword struct {
		Uid   int64  `json:"uid"`
		Email string `json:"email"`
		Name  string `json:"name"`
	}
	return func(w http.ResponseWriter, r *http.Request) error {
		// get uid from token, added to context by authMiddleware
		tokenUid := r.Context().Value("uid").(int64)

		// get uid from url
		vars := mux.Vars(r)
		urlUid, err := strconv.ParseInt(vars["uid"], 10, 64)
		if err != nil {
			return HandlerError{err, http.StatusBadRequest}
		}

		// make sure only the logged-in user can see their profile
		if tokenUid != urlUid {
			status := http.StatusUnauthorized
			return HandlerError{errors.New(http.StatusText(status)), status}
		}

		// get stored user details
		user, err := userRepo.GetFromUid(urlUid)
		if err != nil {
			switch err.(type) {
			case ErrUserDoesNotExist:
				return HandlerError{errors.New("user does not exist"), http.StatusBadRequest}
			default:
				return err
			}
		}

		// don't want to send password hash from database
		userWithoutPassword := UserWithoutPassword{user.Uid, user.Email, user.Name}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(userWithoutPassword)
		return nil
	}
}
