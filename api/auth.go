package api

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
)

func authRouter(repo UserRepo) *mux.Router {
	router := mux.NewRouter()
	router.Handle("/auth/signup", signup(repo)).Methods(http.MethodPost)
	router.Handle("/auth/signin", signin(repo)).Methods(http.MethodPost)
	return router
}

func (user User) Valid() error {
	if len(user.Email) < 3 {
		return errors.New("invalid email")
	}
	if len(user.Password) < 6 {
		return errors.New("password must be at least 6 characters")
	}
	if len(user.Name) < 1 {
		return errors.New("user must have a name")
	}
	return nil
}

type ErrUserExists error
type ErrUserDoesNotExist error

func keyFunc(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, errors.New("invalid signing method")
	}
	return []byte("Very Secret Key, Shhh....."), nil
}

func signup(userRepo UserRepo) HttpErrorHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		// decode user details from request body
		var user User
		json.NewDecoder(r.Body).Decode(&user)
		if err := user.Valid(); err != nil {
			return HandlerError{err, http.StatusBadRequest}
		}

		// generate a bcrypt hash from the user's plaintext password
		hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
		if err != nil {
			return err
		}
		// replace the password in the user struct with the bcrypt hash before sending to the database
		user.Password = string(hash)

		// add the new user to the database
		if err = userRepo.Add(&user); err != nil {
			switch err.(type) {
			case ErrUserExists:
				return HandlerError{errors.New("user with this email already exists"), http.StatusConflict}
			default:
				return err
			}
		}
		return nil
	}
}

func signin(userRepo UserRepo) HttpErrorHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		// decode user details from request body
		var user User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			return HandlerError{err, http.StatusBadRequest}
		}

		// get stored user details
		storedUser, err := userRepo.GetFromEmail(user.Email)
		if err != nil {
			switch err.(type) {
			case ErrUserDoesNotExist:
				return HandlerError{errors.New("incorrect email or password"), http.StatusUnauthorized}
			default:
				return err
			}
		}

		// compare the user's plaintext password with the stored hash
		hash := []byte(storedUser.Password)
		pwd := []byte(user.Password)
		if err = bcrypt.CompareHashAndPassword(hash, pwd); err != nil {
			return HandlerError{errors.New("incorrect email or password"), http.StatusUnauthorized}
		}

		// return a signed JSON Web Token which embeds the user's UID
		// this token should be used in the Authorization header of requests to authenticated endpoints
		// the token is all that the server needs in order to verify the user's UID
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"uid": storedUser.Uid})
		key, _ := keyFunc(token)
		tokenString, err := token.SignedString(key)
		if err != nil {
			return err
		}

		w.Write([]byte(tokenString))
		return nil
	}
}

func authMiddleware(next http.Handler) HttpErrorHandler {
	parser := jwt.Parser{
		ValidMethods:  []string{jwt.SigningMethodHS256.Alg()},
		UseJSONNumber: true,
	}
	return func(w http.ResponseWriter, r *http.Request) error {
		// get jwt token with claims
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			status := http.StatusUnauthorized
			return HandlerError{errors.New(http.StatusText(status)), status}
		}
		var claims jwt.MapClaims
		if _, err := parser.ParseWithClaims(tokenString, &claims, keyFunc); err != nil {
			return HandlerError{err, http.StatusUnauthorized}
		}

		// parse uid from jwt claims
		uidJSON, ok := claims["uid"].(json.Number)
		if !ok {
			log.Print(uidJSON)
			return errors.New("failed to parse uid from jwt")
		}
		uid, err := uidJSON.Int64()
		if err != nil {
			return err
		}

		// set the uid in the request context and pass on the request
		ctx := context.WithValue(r.Context(), "uid", uid)
		next.ServeHTTP(w, r.WithContext(ctx))
		return nil
	}
}
