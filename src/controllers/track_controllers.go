package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/Smirrra/just-to-do-it/src/models"
	"github.com/Smirrra/just-to-do-it/src/utils"
)

func (env *EnvironmentGroup)CreateTrackHandler(w http.ResponseWriter, r *http.Request) {
	paramFromURL := mux.Vars(r)
	groupId, err := strconv.Atoi(paramFromURL["group_id"])
	if err != nil {
		utils.Respond(w, utils.Message(false,"Invalid id","Bad Request"))
		return
	}

	track := models.Track{}
	err = json.NewDecoder(r.Body).Decode(&track)
	if err != nil {
		utils.Respond(w, utils.Message(false,"Invalid body","Bad Request"))
		return
	}

	track, err = env.Db.CreateTrack(groupId, track)
	if err != nil {
		utils.Respond(w, utils.Message(false, err.Error(), "Internal Server Error"))
		return
	}

	resp := utils.Message(true, "Create track", "")
	resp["track"] = track
	utils.Respond(w, resp)
}

func (env *EnvironmentGroup)GetTrackHandler(w http.ResponseWriter, r *http.Request) {
	paramFromURL := mux.Vars(r)
	trackId, err := strconv.Atoi(paramFromURL["track_id"])
	if err != nil {
		utils.Respond(w, utils.Message(false,"Invalid id","Bad Request"))
		return
	}

	track, tasks, err := env.Db.GetTrack(trackId)
	if err != nil {
		utils.Respond(w, utils.Message(false, err.Error(), "Internal Server Error"))
		return
	}

	resp := utils.Message(true, "Get track info", "")
	resp["track"] = track
	resp["tasks"] = tasks
	utils.Respond(w, resp)
}

func (env *EnvironmentGroup)UpdateTrackHandler(w http.ResponseWriter, r *http.Request) {
	paramFromURL := mux.Vars(r)
	trackId, err := strconv.Atoi(paramFromURL["track_id"])
	if err != nil {
		utils.Respond(w, utils.Message(false,"Invalid id","Bad Request"))
		return
	}

	var updateTrack models.Track
	err = json.NewDecoder(r.Body).Decode(&updateTrack)
	if err != nil {
		utils.Respond(w, utils.Message(false,"Invalid body","Bad Request"))
		return
	}

	track, err := env.Db.UpdateTrack(trackId, updateTrack)
	if err != nil {
		utils.Respond(w, utils.Message(false, err.Error(), "Internal Server Error"))
		return
	}

	resp := utils.Message(true, "Updated track", "")
	resp["track"] = track
	utils.Respond(w, resp)
}

func (env *EnvironmentGroup)DeleteTrackHandler(w http.ResponseWriter, r *http.Request) {
	paramFromURL := mux.Vars(r)
	trackId, err := strconv.Atoi(paramFromURL["track_id"])
	if err != nil {
		utils.Respond(w, utils.Message(false,"Invalid id","Bad Request"))
		return
	}

	err = env.Db.DeleteTrack(trackId)
	if err != nil {
		utils.Respond(w, utils.Message(false, err.Error(), "Internal Server Error"))
		return
	}

	resp := utils.Message(true, "Deleted track and all task in track", "")
	utils.Respond(w, resp)
}

func (env *EnvironmentGroup)AddTaskInTrackHandler(w http.ResponseWriter, r *http.Request) {
	paramFromURL := mux.Vars(r)
	taskId, err := strconv.Atoi(paramFromURL["task_id"])
	if err != nil {
		utils.Respond(w, utils.Message(false,"Invalid id","Bad Request"))
		return
	}
	trackId, err := strconv.Atoi(paramFromURL["track_id"])
	if err != nil {
		utils.Respond(w, utils.Message(false,"Invalid id","Bad Request"))
		return
	}

	taskInTrack, err := env.Db.AddTaskInTrack(taskId, trackId)
	if err != nil {
		utils.Respond(w, utils.Message(false, err.Error(), "Internal Server Error"))
		return
	}

	resp := utils.Message(true, "Added task in track", "")
	resp["struct"] = taskInTrack
	utils.Respond(w, resp)
}

func (env *EnvironmentGroup)CreateTaskInTrackHandler(w http.ResponseWriter, r *http.Request) {
	paramFromURL := mux.Vars(r)
	trackId, err := strconv.Atoi(paramFromURL["track_id"])
	if err != nil {
		utils.Respond(w, utils.Message(false,"Invalid id","Bad Request"))
		return
	}
	userId, err := strconv.Atoi(paramFromURL["id"])
	if err != nil {
		utils.Respond(w, utils.Message(false,"Invalid id","Bad Request"))
		return
	}

	var task models.Task
	err = json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		utils.Respond(w, utils.Message(false,"Invalid body","Bad Request"))
		return
	}

	taskInTrack, task, err := env.Db.CreateTaskInTrack(userId, trackId, task)
	if err != nil {
		utils.Respond(w, utils.Message(false, err.Error(), "Internal Server Error"))
		return
	}

	resp := utils.Message(true, "Create task in track", "")
	resp["task"] = task
	resp["track"] = taskInTrack
	utils.Respond(w, resp)
}

func (env *EnvironmentGroup)DeleteTaskInTrack(w http.ResponseWriter, r *http.Request) {
	paramFromURL := mux.Vars(r)
	taskId, err := strconv.Atoi(paramFromURL["task_id"])
	if err != nil {
		utils.Respond(w, utils.Message(false, "Invalid id", "Bad Request"))
		return
	}
	trackId, err := strconv.Atoi(paramFromURL["track_id"])
	if err != nil {
		utils.Respond(w, utils.Message(false, "Invalid id", "Bad Request"))
		return
	}

	err = env.Db.DeleteTaskInTrack(trackId, taskId)
	if err != nil {
		utils.Respond(w, utils.Message(false, err.Error(), "Internal Server Error"))
		return
	}

	resp := utils.Message(true, "Deleted task in track", "")
	utils.Respond(w, resp)
}
