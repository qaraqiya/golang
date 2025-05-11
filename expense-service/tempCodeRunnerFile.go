package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type Expense struct {
	Id          int     `json:"id"`
	UserId      int     `json:"user_id"`
	Category    string  `json:"category"`
	Amount      float64 `json:"amount"`
	Description string  `json:"description"`
	Date        string  `json:"date"`
}

var expenses = []Expense{
	{Id: 1, UserId: 1, Category: "Food", Amount: 20.5, Description: "Lunch", Date: "2025-05-10"},
	{Id: 2, UserId: 1, Category: "Transport", Amount: 10.0, Description: "Taxi", Date: "2025-05-10"},
}

func getExpenses(w http.ResponseWriter, r *http.Request) {
	userId := r.URL.Query().Get("user_id")
	if userId == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	var filteredExpenses []Expense
	for _, expense := range expenses {
		if fmt.Sprintf("%d", expense.UserId) == userId {
			filteredExpenses = append(filteredExpenses, expense)
		}
	}

	if len(filteredExpenses) == 0 {
		http.Error(w, "No expenses found for the given user", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(filteredExpenses)
}

func createExpense(w http.ResponseWriter, r *http.Request) {
	var expense Expense
	if err := json.NewDecoder(r.Body).Decode(&expense); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Генерация ID (в реальном приложении это будет делать база данных)
	expense.Id = len(expenses) + 1
	expenses = append(expenses, expense)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(expense)
}

func deleteExpense(w http.ResponseWriter, r *http.Request) {
	expenseID := r.URL.Query().Get("id")
	if expenseID == "" {
		http.Error(w, "Expense ID is required", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(expenseID)
	if err != nil {
		http.Error(w, "Invalid Expense ID", http.StatusBadRequest)
		return
	}

	// Поиск расхода и его удаление
	for i, expense := range expenses {
		if expense.Id == id {
			expenses = append(expenses[:i], expenses[i+1:]...)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}

	http.Error(w, "Expense not found", http.StatusNotFound)
}

func main() {
	http.HandleFunc("/get_expenses", getExpenses)
	http.HandleFunc("/create_expense", createExpense)
	http.HandleFunc("/delete_expense", deleteExpense)

	fmt.Println("🚀 Starting HTTP server on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
