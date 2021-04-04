package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/Smirrra/just-to-do-it/src/models"
	"github.com/Smirrra/just-to-do-it/src/services"
	"github.com/Smirrra/just-to-do-it/src/utils"
)

type EnvironmentGroup struct {
	Db services.DatastoreGroup
}

func (env *EnvironmentGroup)CreateGroupHandler(w http.ResponseWriter, r *http.Request) {
	group := models.Group{}
	err := json.NewDecoder(r.Body).Decode(&group)
	if err != nil {
		utils.Respond(w, utils.Message(false,"Invalid body", "Bad Request"))
		return
	}
	if group.Title == "" || group.Description == "" {
		utils.Respond(w, utils.Message(false,"Invalid body", "Bad Request"))
		return
	}

	paramFromURL := mux.Vars(r)
	userId, err := strconv.Atoi(paramFromURL["id"])
	if err != nil {
		utils.Respond(w, utils.Message(false,"Invalid id", "Bad Request"))
		return
	}

	group, err = env.Db.CreateGroup(group, userId)
	if err != nil {
		utils.Respond(w, utils.Message(false, err.Error(), "Internal Server Error"))
		return
	}

	resp := utils.Message(true, "Created group", "")
	resp["group"] = group
	utils.Respond(w, resp)
}

func (env *EnvironmentGroup)GetGroupHandler(w http.ResponseWriter, r *http.Request) {
	paramFromURL := mux.Vars(r)
	groupId, err := strconv.Atoi(paramFromURL["group_id"])
	if err != nil {
		utils.Respond(w, utils.Message(false,"Invalid id","Bad Request"))
		return
	}

	group, err := env.Db.GetGroup(groupId)
	if err != nil {
		utils.Respond(w, utils.Message(false, err.Error(), "Internal Server Error"))
		return
	}

	resp := utils.Message(true, "Get group", "")
	resp["group"] = group
	utils.Respond(w, resp)
}

func (env *EnvironmentGroup)UpdateGroupHandler(w http.ResponseWriter, r *http.Request) {
	paramFromURL := mux.Vars(r)
	groupId, err := strconv.Atoi(paramFromURL["group_id"])
	if err != nil {
		utils.Respond(w, utils.Message(false,"Invalid id","Bad Request"))
		return
	}

	group := models.Group{}
	err = json.NewDecoder(r.Body).Decode(&group)
	if err != nil {
		utils.Respond(w, utils.Message(false, "Invalid id", "Bad Request"))
		return
	}
	if group.Title == "" || group.Description == "" {
		utils.Respond(w, utils.Message(false, "Invalid body", "Bad Request"))
		return
	}

	group, err = env.Db.UpdateGroup(groupId, group)
	if err != nil {
		utils.Respond(w, utils.Message(false, err.Error(), "Internal Server Error"))
		return
	}

	resp := utils.Message(true, "Update group", "")
	group.Id = groupId
	resp["group"] = group
	utils.Respond(w, resp)
}

func (env *EnvironmentGroup)DeleteGroupHandler(w http.ResponseWriter, r *http.Request) {
	paramFromURL := mux.Vars(r)
	groupId, err := strconv.Atoi(paramFromURL["group_id"])
	if err != nil {
		utils.Respond(w, utils.Message(false,"Invalid id","Bad Request"))
		return
	}

	err = env.Db.DeleteGroup(groupId)
	if err != nil {
		utils.Respond(w, utils.Message(false, err.Error(), "Internal Server Error"))
		return
	}

	resp := utils.Message(true, "Deleted group", "")
	utils.Respond(w, resp)
}

func (env *EnvironmentGroup)GetGroupsByUserId(w http.ResponseWriter, r *http.Request) {
	paramFromURL := mux.Vars(r)
	userId, err := strconv.Atoi(paramFromURL["id"])
	if err != nil {
		utils.Respond(w, utils.Message(false,"Invalid id","Bad Request"))
		return
	}

	groups, err := env.Db.GetGroups(userId)
	if err != nil {
		utils.Respond(w, utils.Message(false, err.Error(),"Internal Server Error"))
		return
	}

	resp := utils.Message(true, "Get groups", "")
	resp["groups"] = groups
	utils.Respond(w, resp)
}