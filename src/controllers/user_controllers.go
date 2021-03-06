package controllers

import (
	"encoding/json"
	"github.com/Smirrra/just-to-do-it/src/auth"
	"github.com/Smirrra/just-to-do-it/src/models"
	"github.com/Smirrra/just-to-do-it/src/services"
	"github.com/Smirrra/just-to-do-it/src/utils"
	"net/http"
)

type EnvironmentUser struct {
	Db services.DatastoreUser
}

func (env *EnvironmentUser) ResponseLoginHandler(w http.ResponseWriter, r *http.Request) {
	user := models.User{}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		utils.Respond(w, utils.Message(false,"Invalid body","Bad Request"))
		return
	}

	user, err = env.Db.Login(user.Login, user.Password)
	if err != nil {
		utils.Respond(w, utils.Message(false,"Invalid login or password","Unauthorized"))
		return
	}

	accToken, err := auth.CreateAccessToken(user.Id)
	if err != nil {
		utils.Respond(w, utils.Message(false, err.Error(), "Unauthorized"))
		return
	}
	auth.SetCookieForAccToken(w, accToken)

	refToken, err := auth.CreateRefreshToken(user.Id)
	if err != nil {
		utils.Respond(w, utils.Message(false, err.Error(), "Unauthorized"))
		return
	}
	auth.SetCookieForRefToken(w, refToken)

	resp := utils.Message(true, "Logged In", "")
	resp["user"] = user
	utils.Respond(w, resp)
}

func (env *EnvironmentUser) ResponseRegisterHandler (w http.ResponseWriter, r *http.Request) {
	user := models.User{}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		utils.Respond(w, utils.Message(false,"Invalid body","Bad Request"))
		return
	}

	if user.Password == ""  || user.Login == "" || user.Email == "" {
		utils.Respond(w, utils.Message(false,"Invalid body","Bad Request"))
		return
	}

	user, msg, errStr := env.Db.Register(user)
	if msg != "" {
		utils.Respond(w, utils.Message(false, msg, errStr))
		return
	}

	accToken, err := auth.CreateAccessToken(user.Id)
	if err != nil {
		utils.Respond(w, utils.Message(false, err.Error(), "Unauthorized"))
		return
	}
	auth.SetCookieForAccToken(w, accToken)

	refToken, err := auth.CreateRefreshToken(user.Id)
	if err != nil {
		utils.Respond(w, utils.Message(false, err.Error(), "Unauthorized"))
		return
	}
	auth.SetCookieForRefToken(w, refToken)

	resp := utils.Message(true, "User created", "")
	resp["user"] = user
	utils.Respond(w, resp)
}

func (env *EnvironmentUser) ConfirmEmailHandler (w http.ResponseWriter, r *http.Request) {
	hash := r.URL.Query().Get("hash")
	err := env.Db.Confirm(hash)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (env *EnvironmentUser) GetUserHandler (w http.ResponseWriter, r *http.Request) {
	id, err := auth.CheckUser(w, r)
	if err != nil {
		utils.Respond(w, utils.Message(false, err.Error(), "Unauthorized"))
		return
	}

	user, err := env.Db.GetUser(int(id))
	if err != nil {
		utils.Respond(w, utils.Message(false,"Not found user in db","Internal Server Error"))
		return
	}

	resp := utils.Message(true, "Get user", "")
	resp["user"] = user
	utils.Respond(w, resp)
}

func (env *EnvironmentUser) UpdateUserHandler (w http.ResponseWriter, r *http.Request) {
	id, err := auth.CheckUser(w, r)
	if err != nil {
		utils.Respond(w, utils.Message(false, err.Error(), "Unauthorized"))
		return
	}

	user := models.User{}
	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		utils.Respond(w, utils.Message(false,"Invalid body","Bad Request"))
		return
	}

	if  user.Email == "" || user.Fullname == "" ||
		user.Login == "" || user.Password == "" {
		utils.Respond(w, utils.Message(false,"Invalid body","Bad Request"))
		return
	}

	user, err = env.Db.UpdateUser(int(id), user)
	if err != nil {
		utils.Respond(w, utils.Message(false, "Database error", "Internal Server Error"))
		return
	}

	resp := utils.Message(true, "Update user", "")
	resp["user"] = user
	utils.Respond(w, resp)
}

func (env *EnvironmentUser) DeleteUserHandler (w http.ResponseWriter, r *http.Request) {
	id, err := auth.CheckUser(w, r)
	if err != nil {
		utils.Respond(w, utils.Message(false, err.Error(), "Unauthorized"))
		return
	}

	err = env.Db.DeleteUser(id)
	if err != nil {
		utils.Respond(w, utils.Message(false,err.Error(),"Internal Server Error"))
		return
	}

	resp := utils.Message(true, "User deleted", "")
	utils.Respond(w, resp)
}