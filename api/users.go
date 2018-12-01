package api

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type UserWithoutPassword struct {
	Uid   int64  `json:"uid"`
	Email string `json:"email"`
	Name  string `json:"name"`
}
type User struct {
	UserWithoutPassword
	Password string `json:"password"`
}

type Link struct {
	Source int64 `json:"from"`
	Dest   int64 `json:"to"`
}

type UserRepo interface {
	Add(*User) error
	GetFromUid(int64) (User, error)
	GetFromEmail(string) (User, error)
	AddLink(from, to int64) error
	GetIncomingLinks(int64) ([]UserWithoutPassword, error)
	GetOutgoingLinks(int64) ([]UserWithoutPassword, error)
}

func userRouter(repo UserRepo) *mux.Router {
	router := mux.NewRouter()
	router.Handle("/users/{uid}", getUserInfo(repo)).Methods(http.MethodGet)
	router.Handle("/users/links/in", getIncomingLinks(repo)).Methods(http.MethodGet)
	router.Handle("/users/links/out", getOutgoingLinks(repo)).Methods(http.MethodGet)
	router.Handle("/users/links/out", link(repo)).Methods(http.MethodPost)
	return router
}

func getUserInfo(userRepo UserRepo) HttpErrorHandler {
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

func link(userRepo UserRepo) HttpErrorHandler {
	type email struct {
		Email string `json:"email"`
	}
	return func(w http.ResponseWriter, r *http.Request) error {
		// get uid from token, added to context by authMiddleware
		tokenUid := r.Context().Value("uid").(int64)

		var e email
		if err := json.NewDecoder(r.Body).Decode(&e); err != nil {
			return HandlerError{err, http.StatusBadRequest}
		}

		otherUser, err := userRepo.GetFromEmail(e.Email)
		if err != nil {
			switch err.(type) {
			case ErrUserDoesNotExist:
				return HandlerError{err, http.StatusBadRequest}
			default:
				return err
			}
		}

		return userRepo.AddLink(tokenUid, otherUser.Uid)
	}
}

func getIncomingLinks(userRepo UserRepo) HttpErrorHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		// get uid from token, added to context by authMiddleware
		tokenUid := r.Context().Value("uid").(int64)

		users, err := userRepo.GetIncomingLinks(tokenUid)
		if err != nil {
			return err
		}

		json.NewEncoder(w).Encode(users)
		return nil
	}
}

func getOutgoingLinks(userRepo UserRepo) HttpErrorHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		// get uid from token, added to context by authMiddleware
		tokenUid := r.Context().Value("uid").(int64)

		users, err := userRepo.GetOutgoingLinks(tokenUid)
		if err != nil {
			return err
		}

		json.NewEncoder(w).Encode(users)
		return nil
	}
}
