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

type EnvironmentTask struct {
	Db services.DatastoreTask
}

func(env *EnvironmentTask) GetTasksHandler(w http.ResponseWriter, r *http.Request) {
	paramFromURL := mux.Vars(r)
	userId, err := strconv.Atoi(paramFromURL["id"])
	if err != nil {
		utils.Respond(w, utils.Message(false,"Invalid id","Bad Request"))
		return
	}

	var strSlice []string
	idStr := r.URL.Query().Get("id")
	assigneeIdStr := r.URL.Query().Get("assignee_id")
	groupIdStr := r.URL.Query().Get("group_id")
	strSlice = append(strSlice, idStr, assigneeIdStr, groupIdStr)
	title := r.URL.Query().Get("title")
	var idSlice []int

	for _, k := range strSlice {
		if k != "" {
			tmp, err := strconv.Atoi(k)
			if err != nil {
				utils.Respond(w, utils.Message(false,"bad parameters", "Bad Request"))
				return
			}
			idSlice = append(idSlice, tmp)
		} else {
			idSlice = append(idSlice, 0)
		}
	}

	tasks, err := env.Db.GetTasks(idSlice, title, userId)
	if err != nil {
		utils.Respond(w, utils.Message(false, err.Error(), "Internal Server Error"))
		return
	}

	resp := utils.Message(true,"Get tasks", "")
	resp["tasks"]= tasks
	utils.Respond(w, resp)
}

func (env *EnvironmentTask)GetTaskHandler(w http.ResponseWriter, r *http.Request) {
	paramFromURL := mux.Vars(r)
	userId, err := strconv.Atoi(paramFromURL["id"])
	if err != nil {
		utils.Respond(w, utils.Message(false,"Invalid id","Bad Request"))
		return
	}

	taskId, err := strconv.Atoi(paramFromURL["task_id"])
	if err != nil {
		utils.Respond(w, utils.Message(false,"Invalid id","Bad Request"))
		return
	}

	var labels []models.Label
	task, labels, err := env.Db.GetTaskById(taskId)
	if err != nil {
		utils.Respond(w, utils.Message(false, err.Error(), "Internal Server Error"))
		return
	}

	if task.CreatorId != userId {
		utils.Respond(w, utils.Message(false, "Id don't match", "Unauthorized"))
		return
	}

	resp := utils.Message(true, "Get task", "")
	resp["task"] = task
	resp["task_labels"] = labels
	utils.Respond(w, resp)
}

func (env *EnvironmentTask)CreateTask(w http.ResponseWriter, r *http.Request) {
	paramFromURL := mux.Vars(r)
	id, err := strconv.Atoi(paramFromURL["id"])
	if err != nil {
		utils.Respond(w, utils.Message(false,"Invalid id","Bad Request"))
		return
	}

	groupId, _ := strconv.Atoi(paramFromURL["group_id"])

	task := models.Task{}
	err = json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		utils.Respond(w, utils.Message(false,"Invalid body", "Bad Request"))
		return
	}

	task.GroupId = groupId
	task, err = env.Db.CreateTask(task, id)
	if err != nil {
		utils.Respond(w, utils.Message(false,err.Error(), "Internal Server Error"))
		return
	}

	resp := utils.Message(true,"Create task", "")
	resp["task"] = task
	utils.Respond(w, resp)
}

func (env *EnvironmentTask)UpdateTaskHandler(w http.ResponseWriter, r *http.Request) {
	paramFromURL := mux.Vars(r)
	userId, err := strconv.Atoi(paramFromURL["id"])
	if err != nil {
		utils.Respond(w, utils.Message(false,"Invalid id","Bad Request"))
		return
	}

	taskId, err := strconv.Atoi(paramFromURL["task_id"])
	if err != nil {
		utils.Respond(w, utils.Message(false,"Invalid id","Bad Request"))
		return
	}

	task := models.Task{}
	err = json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		utils.Respond(w, utils.Message(false,"Invalid body", "Bad Request"))
		return
	}

	if  task.Title == "" || task.Description == "" ||
		task.Deadline == 0 || task.State == "" ||
		task.Priority == 0 || task.AssigneeId == 0 || task.Duration == 0 {
		utils.Respond(w, utils.Message(false,"Invalid body", "Bad Request"))
		return
	}

	task, err = env.Db.UpdateTask(task, taskId, userId)
	if err != nil {
		utils.Respond(w, utils.Message(false, err.Error(),"Internal Server Error"))
		return
	}

	resp := utils.Message(true, "Update task", "")
	resp["task"] = task
	utils.Respond(w, resp)
}
