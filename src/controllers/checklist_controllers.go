package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/Smirrra/just-to-do-it/src/models"
	"github.com/Smirrra/just-to-do-it/src/utils"
)

func (env *EnvironmentTask)CreateChecklistHandler(w http.ResponseWriter, r *http.Request) {
	paramFromURL := mux.Vars(r)
	taskId, err := strconv.Atoi(paramFromURL["task_id"])
	if err != nil {
		utils.Respond(w, utils.Message(false, err.Error(), "Bad Request"))
		return
	}

	checklist := models.Checklist{}
	err = json.NewDecoder(r.Body).Decode(&checklist)
	if err != nil {
		utils.Respond(w, utils.Message(false, err.Error(), "Bad Request"))
		return
	}

	checklist, err = env.Db.CreateChecklist(checklist, taskId)
	if err != nil {
		utils.Respond(w, utils.Message(false, err.Error(), "Internal Server Error"))
		return
	}

	resp := utils.Message(true, "Create checklist", "")
	resp["checklist"] = checklist
	utils.Respond(w, resp)
}

func (env *EnvironmentTask)CreateItemHandler(w http.ResponseWriter, r *http.Request) {
	paramFromURL := mux.Vars(r)
	checklistId, err := strconv.Atoi(paramFromURL["checklist_id"])
	if err != nil {
		utils.Respond(w, utils.Message(false, err.Error(), "Bad Request"))
		return
	}

	item := models.ChecklistItem{}
	err = json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		utils.Respond(w, utils.Message(false, err.Error(), "Bad Request"))
		return
	}

	item, err = env.Db.CreateChecklistItem(item, checklistId)
	if err != nil {
		utils.Respond(w, utils.Message(false, err.Error(), "Internal Server Error"))
		return
	}

	resp := utils.Message(true, "Create item", "")
	resp["item"] = item
	utils.Respond(w, resp)
}

func (env *EnvironmentTask)GetChecklistHandler(w http.ResponseWriter, r *http.Request) {
	paramFromURL := mux.Vars(r)
	checklistId, err := strconv.Atoi(paramFromURL["checklist_id"])
	if err != nil {
		utils.Respond(w, utils.Message(false, err.Error(), "Bad Request"))
		return
	}

	checklist, item, err := env.Db.GetChecklist(checklistId)
	if err != nil {
		utils.Respond(w, utils.Message(false, err.Error(), "Internal Server Error"))
		return
	}

	resp := utils.Message(true, "Get checklist", "")
	resp["checklist"] = checklist
	resp["item"] = item
	utils.Respond(w, resp)
}

func (env *EnvironmentTask)UpdateChecklistHandler(w http.ResponseWriter, r *http.Request) {
	paramFromURL := mux.Vars(r)
	checklistId, err := strconv.Atoi(paramFromURL["checklist_id"])
	if err != nil {
		utils.Respond(w, utils.Message(false, err.Error(), "Bad Request"))
		return
	}

	checklist := models.Checklist{}
	err = json.NewDecoder(r.Body).Decode(&checklist)
	if err != nil {
		utils.Respond(w, utils.Message(false, err.Error(), "Bad Request"))
		return
	}
	if checklist.Name == "" {
		utils.Respond(w, utils.Message(false, "Invalid body", "Bad Request"))
		return
	}

	checklist, err = env.Db.UpdateChecklist(checklistId, checklist)
	if err != nil {
		utils.Respond(w, utils.Message(false, err.Error(), "Internal Server Error"))
		return
	}

	resp := utils.Message(false, "Update Checklist", "")
	resp["checklist"] = checklist
	utils.Respond(w, resp)
}

func (env *EnvironmentTask)DeleteChecklistHandler(w http.ResponseWriter, r *http.Request) {
	paramFromURL := mux.Vars(r)
	checklistId, err := strconv.Atoi(paramFromURL["checklist_id"])
	if err != nil {
		utils.Respond(w, utils.Message(false, err.Error(), "Bad Request"))
		return
	}

	err = env.Db.DeleteChecklist(checklistId)
	if err != nil {
		utils.Respond(w, utils.Message(false, err.Error(), "Internal Server Error"))
		return
	}
	utils.Respond(w, utils.Message(true, "Delete checklist", ""))
}

func (env *EnvironmentTask)GetChecklistItems(w http.ResponseWriter, r *http.Request) {
	paramFromURL := mux.Vars(r)
	checklistId, err := strconv.Atoi(paramFromURL["checklist_id"])
	if err != nil {
		utils.Respond(w, utils.Message(false, err.Error(), "Bad Request"))
		return
	}

	items, err := env.Db.GetChecklistItems(checklistId)
	if err != nil {
		utils.Respond(w, utils.Message(false, err.Error(), "Internal Server Error"))
		return
	}

	resp := utils.Message(true, "Get items", "")
	resp["items"] = items
	utils.Respond(w, resp)
}

func (env *EnvironmentTask)UpdateItemHandler(w http.ResponseWriter, r *http.Request) {
	paramFromURL := mux.Vars(r)
	checklistId, err := strconv.Atoi(paramFromURL["checklist_id"])
	if err != nil {
		utils.Respond(w, utils.Message(false, err.Error(), "Bad Request"))
		return
	}
	itemId, err := strconv.Atoi(paramFromURL["item_id"])
	if err != nil {
		utils.Respond(w, utils.Message(false, err.Error(), "Bad Request"))
		return
	}

	item := models.ChecklistItem{}
	err = json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		utils.Respond(w, utils.Message(false, err.Error(), "Bad Request"))
		return
	}

	item, err = env.Db.UpdateChecklistItem(itemId, checklistId, item)
	if err != nil {
		utils.Respond(w, utils.Message(false, err.Error(), "Internal Server Error"))
		return
	}

	resp := utils.Message(true, "Update item", "")
	resp["item"] = item
	utils.Respond(w, resp)
}

func (env *EnvironmentTask)DeleteItemHandler(w http.ResponseWriter, r *http.Request) {
	paramFromURL := mux.Vars(r)
	itemId, err := strconv.Atoi(paramFromURL["item_id"])
	if err != nil {
		utils.Respond(w, utils.Message(false, err.Error(), "Bad Request"))
		return
	}

	err = env.Db.DeleteChecklistItem(itemId)
	if err != nil {
		utils.Respond(w, utils.Message(false, err.Error(), "Internal Server Error"))
		return
	}
	resp := utils.Message(true, "Delete item", "")
	utils.Respond(w, resp)
}
