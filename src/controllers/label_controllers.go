package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/Smirrra/just-to-do-it/src/models"
	"github.com/Smirrra/just-to-do-it/src/utils"
)

func (env *EnvironmentTask) CreateLabelHandler (w http.ResponseWriter, r *http.Request) {
	paramFromURL := mux.Vars(r)
	taskId, err := strconv.Atoi(paramFromURL["task_id"])
	if err != nil {
		utils.Respond(w, utils.Message(false,"Invalid id","Bad Request"))
		return
	}

	label := models.Label{}
	err = json.NewDecoder(r.Body).Decode(&label)
	if err != nil {
		utils.Respond(w, utils.Message(false,"Invalid body", "Bad Request"))
		return
	}
	label.TaskId = taskId

	label, err = env.Db.CreateLabel(label, taskId)
	if err != nil {
		utils.Respond(w, utils.Message(false, err.Error(), "Internal Server Error"))
		return
	}

	resp := utils.Message(true, "Create label", "")
	resp["task_label"] = label
	utils.Respond(w, resp)
}

func (env *EnvironmentTask)GetLabelHandler(w http.ResponseWriter, r *http.Request) {
	paramFromURL := mux.Vars(r)
	labelId, err := strconv.Atoi(paramFromURL["label_id"])
	if err != nil {
		utils.Respond(w, utils.Message(false,"Invalid id","Bad Request"))
		return
	}

	label, err := env.Db.GetLabel(labelId)
	if err != nil {
		utils.Respond(w, utils.Message(false, err.Error(), "Internal Server Error"))
		return
	}

	resp := utils.Message(true, "Get task", "")
	resp["task_label"] = label
	utils.Respond(w, resp)
}

func (env *EnvironmentTask)GetLabelsByTaskIdHandler(w http.ResponseWriter, r *http.Request) {
	paramFromURL := mux.Vars(r)
	taskId, err := strconv.Atoi(paramFromURL["task_id"])
	if err != nil {
		utils.Respond(w, utils.Message(false,"Invalid id","Bad Request"))
		return
	}

	labels, err := env.Db.GetLabelsByTaskId(taskId)
	if err != nil {
		utils.Respond(w, utils.Message(false, err.Error(), "Internal Server Error"))
		return
	}

	resp := utils.Message(true, "Get labels", "")
	resp["task_labels"] = labels
	utils.Respond(w, resp)
}

func (env *EnvironmentTask)UpdateLabelColorHandler(w http.ResponseWriter, r *http.Request) {
	paramFromURL := mux.Vars(r)
	labelId, err := strconv.Atoi(paramFromURL["label_id"])
	if err != nil {
		utils.Respond(w, utils.Message(false,"Invalid id","Bad Request"))
		return
	}

	label := models.Label{}
	err = json.NewDecoder(r.Body).Decode(&label)
	if err != nil {
		utils.Respond(w, utils.Message(false,"Invalid body", "Bad Request"))
		return
	}

	label, err = env.Db.UpdateLabelColor(labelId, label.Color)
	if err != nil {
		utils.Respond(w, utils.Message(false, err.Error(), "Internal Server Error"))
		return
	}

	resp := utils.Message(true, "Update label color", "")
	resp["task_label"] = label
	utils.Respond(w, resp)
}

func (env *EnvironmentTask)UpdateLabelTitleHandler(w http.ResponseWriter, r *http.Request) {
	paramFromURL := mux.Vars(r)
	labelId, err := strconv.Atoi(paramFromURL["label_id"])
	if err != nil {
		utils.Respond(w, utils.Message(false,"Invalid id","Bad Request"))
		return
	}

	label := models.Label{}
	err = json.NewDecoder(r.Body).Decode(&label)
	if err != nil {
		utils.Respond(w, utils.Message(false,"Invalid body", "Bad Request"))
		return
	}

	label, err = env.Db.UpdateLabelTitle(labelId, label.Title)
	if err != nil {
		utils.Respond(w, utils.Message(false, err.Error(), "Internal Server Error"))
		return
	}

	resp := utils.Message(true, "Update label color", "")
	resp["task_label"] = label
	utils.Respond(w, resp)
}

func (env *EnvironmentTask)DeleteLabelHandler(w http.ResponseWriter, r *http.Request) {
	paramFromURL := mux.Vars(r)
	labelId, err := strconv.Atoi(paramFromURL["label_id"])
	if err != nil {
		utils.Respond(w, utils.Message(false,"Invalid id","Bad Request"))
		return
	}

	err = env.Db.DeleteLabel(labelId)
	if err != nil {
		utils.Respond(w, utils.Message(false, err.Error(), "Internal Server Error"))
		return
	}

	resp := utils.Message(true, "Label deleted", "")
	utils.Respond(w, resp)
}
