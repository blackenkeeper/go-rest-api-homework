package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Task - структура для создания типизированных задач
type Task struct {
	ID           string   `json:"id"`
	Description  string   `json:"description"`
	Note         string   `json:"note"`
	Applications []string `json:"applications"`
}

// Список задач к исполнению или TODO-лист
var tasks = map[string]Task{
	"1": {
		ID:          "1",
		Description: "Сделать финальное задание темы REST API",
		Note:        "Если сегодня сделаю, то завтра будет свободный день. Ура!",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
		},
	},
	"2": {
		ID:          "2",
		Description: "Протестировать финальное задание с помощью Postman",
		Note:        "Лучше это делать в процессе разработки, каждый раз, когда запускаешь сервер и проверяешь хендлер",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
			"Postman",
		},
	},
}

// Обработчик для GET-запроса с эндпоинтом "/tasks". Используется для вывода полного списка задач
func getTasks(w http.ResponseWriter, r *http.Request) {
	tasksSerializedInJson, err := json.MarshalIndent(&tasks, "", "    ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(tasksSerializedInJson)
}

// Обработчик для GET-запроса с эндпоинтом "/tasks/{id}". Используется для вывода конкретной задачи по {id}
func getTask(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	task, ok := tasks[id]
	if !ok {
		http.Error(w, "Задачи с таким ID нет в вашем списке", http.StatusNoContent)
		return
	}

	taskSerializedInJson, err := json.MarshalIndent(task, "", "    ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(taskSerializedInJson)
}

// Обработчик для POST-запроса с эндпоинтом "/tasks". Используется для добавления новой задачи в список
func postTask(w http.ResponseWriter, r *http.Request) {
	var taskDeserilializedFromRequestBody Task
	var bufferForRequestBody bytes.Buffer

	_, err := bufferForRequestBody.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(bufferForRequestBody.Bytes(), &taskDeserilializedFromRequestBody); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tasks[taskDeserilializedFromRequestBody.ID] = taskDeserilializedFromRequestBody

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

// Обработчик для DELETE-запроса с эндпоинтом "/tasks/{id}". Используется для удаления задачи из списка по {id}
func deleteTask(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if _, ok := tasks[id]; ok {
		delete(tasks, id)
	} else {
		http.Error(w, "Задачи с таким ID нет в вашем списке", http.StatusNoContent)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func main() {
	r := chi.NewRouter()

	r.Get("/tasks", getTasks)
	r.Get("/tasks/{id}", getTask)
	r.Post("/tasks", postTask)
	r.Delete("/tasks/{id}", deleteTask)

	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}
