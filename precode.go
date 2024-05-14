package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Структура объекта задачи
type Task struct {
	ID           string   `json:"id"`
	Description  string   `json:"description"`
	Note         string   `json:"note"`
	Applications []string `json:"applications"`
}

// Тестовая мапа для сущностей задач
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
		Description: "Протестировать финальное задание с помощью Postmen",
		Note:        "Лучше это делать в процессе разработки, каждый раз, когда запускаешь сервер и проверяешь хендлер",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
			"Postman",
		},
	},
}

// Получить все задачи
func getAllTasks(w http.ResponseWriter, r *http.Request) {
	// Сериализация данных
	resp, err := json.Marshal(tasks)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Установка заголовков
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	// Запись данных в тело ответа
	w.Write(resp)
}

// Добавить новую задачу
func addTask(w http.ResponseWriter, r *http.Request) {
	var task Task
	var buf bytes.Buffer

	// Получаем данные из тела запроса
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Десериализуем данные задачи
	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Проверка на наличие задачи с таким же id
	_, ok := tasks[task.ID]
	if ok {
		http.Error(w, "Задача с таким id уже существует", http.StatusConflict)
		return
	}

	// Записываем новую задачу в мапу
	tasks[task.ID] = task

	// Формируем заголовки ответа
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

}

// Получить задачу по id
func getTask(w http.ResponseWriter, r *http.Request) {
	// Поиск задачи по id из параметров запроса
	id := chi.URLParam(r, "id")
	task, ok := tasks[id]
	if !ok {
		http.Error(w, "Задача отсутствует", http.StatusBadRequest)
		return
	}

	// Сериализация объекта задачи
	taskJSON, err := json.Marshal(task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Установка заголовков
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	// Запись данных в тело ответа
	w.Write(taskJSON)
}

// Удалить задачу по id
func deleteTask(w http.ResponseWriter, r *http.Request) {
	// Получаем id задачи из параметров
	id := chi.URLParam(r, "id")

	// Проверка на наличие задачи с таким id
	_, ok := tasks[id]
	if !ok {
		http.Error(w, "Задача отсутствует", http.StatusBadRequest)
		return
	}

	// Удаляем значение из мапы по id
	delete(tasks, id)

	// Формируем заголовки ответа
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

}

func main() {
	r := chi.NewRouter()

	r.Get("/tasks", getAllTasks)
	r.Post("/tasks", addTask)
	r.Get("/tasks/{id}", getTask)
	r.Delete("/tasks/{id}", deleteTask)

	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}
