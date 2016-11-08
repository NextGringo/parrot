package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/anthonynsimon/parrot/errors"
	"github.com/anthonynsimon/parrot/model"
	"github.com/anthonynsimon/parrot/render"
	"golang.org/x/crypto/bcrypt"
)

func createUser(w http.ResponseWriter, r *http.Request) {
	// TODO(anthonynsimon): handle user already exists
	user := model.User{}
	errs := decodeAndValidate(r.Body, &user)
	if errs != nil {
		render.ErrorWithStatus(w, http.StatusBadRequest, errs)
		return
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		render.Error(w, errors.ErrInternal)
		return
	}

	user.Password = string(hashed)

	err = store.CreateUser(&user)
	if err != nil {
		render.Error(w, err)
		return
	}

	render.JSON(w, http.StatusCreated, map[string]interface{}{
		"message": fmt.Sprintf("created user with email: %s", user.Email),
	})
}

func getUserIDFromContext(ctx context.Context) (int, error) {
	v := ctx.Value("userID")
	if v == nil {
		return -1, errors.ErrInternal
	}
	str := v.(string)
	if v == "" {
		return -1, errors.ErrInternal
	}
	id, err := strconv.Atoi(str)
	if err != nil {
		return -1, errors.ErrInternal
	}
	return id, nil
}

func decodeAndValidate(r io.Reader, m model.Validatable) error {
	if err := json.NewDecoder(r).Decode(m); err != nil {
		return errors.ErrBadRequest
	}
	return m.Validate()
}
