package main

// @title Expense Tracker API
// @version 1.0
// @description Backend API for expense tracking.
// @host localhost:8080
// @BasePath /

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "expense-tracker/expense-service/docs"

	httpSwagger "github.com/swaggo/http-swagger"

	_ "github.com/lib/pq"
)

var db *sql.DB

type Expense struct {
	Id          int     `json:"id"`
	UserId      int     `json:"user_id"`
	Category    string  `json:"category"`
	Amount      float64 `json:"amount"`
	Description string  `json:"description"`
	Date        string  `json:"date"`
}

// @Summary Get expenses by user ID
// @Produce json
// @Param user_id query int true "User ID"
// @Success 200 {array} Expense
// @Failure 400 {string} string "User ID is required"
// @Failure 404 {string} string "No expenses found"
// @Router /get_expenses [get]
func getExpenses(w http.ResponseWriter, r *http.Request) {
	userId := r.URL.Query().Get("user_id")
	if userId == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	rows, err := db.Query(`SELECT id, user_id, category, amount, description, date FROM expenses WHERE user_id = $1`, userId)
	if err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var expenses []Expense
	for rows.Next() {
		var e Expense
		if err := rows.Scan(&e.Id, &e.UserId, &e.Category, &e.Amount, &e.Description, &e.Date); err != nil {
			log.Println("❌ QUERY ERROR:", err)
			return
		}
		expenses = append(expenses, e)
	}

	if len(expenses) == 0 {
		http.Error(w, "No expenses found for the given user", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(expenses)
}

// @Summary Add a new expense
// @Accept json
// @Produce json
// @Param expense body Expense true "Expense to add"
// @Success 201 {string} string "Expense added successfully"
// @Router /add_expense [post]
func addExpense(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	var exp Expense
	err := json.NewDecoder(r.Body).Decode(&exp)
	if err != nil {
		http.Error(w, "Bad request: "+err.Error(), http.StatusBadRequest)
		return
	}

	query := `INSERT INTO expenses (user_id, category, amount, description, date)
	          VALUES ($1, $2, $3, $4, $5)`
	_, err = db.Exec(query, exp.UserId, exp.Category, exp.Amount, exp.Description, exp.Date)
	if err != nil {
		log.Println("❌ INSERT ERROR:", err)
http.Error(w, "Error inserting data: "+err.Error(), http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintln(w, "Expense added successfully")
}

// @Summary Delete expense by ID
// @Param id query int true "Expense ID"
// @Success 204 {string} string "Deleted"
// @Router /delete_expense [delete]
func deleteExpense(w http.ResponseWriter, r *http.Request) {
	expenseID := r.URL.Query().Get("id")
	if expenseID == "" {
		http.Error(w, "Expense ID is required", http.StatusBadRequest)
		return
	}

	_, err := db.Exec(`DELETE FROM expenses WHERE id = $1`, expenseID)
	if err != nil {
		http.Error(w, "Failed to delete expense: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// @Summary Update expense
// @Accept json
// @Produce json
// @Param expense body Expense true "Expense to update"
// @Success 200 {string} string "Updated"
// @Router /update_expense [put]
func updateExpense(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	var exp Expense
	err := json.NewDecoder(r.Body).Decode(&exp)
	if err != nil {
		http.Error(w, "Bad request: "+err.Error(), http.StatusBadRequest)
		return
	}

	_, err = db.Exec(`
		UPDATE expenses
		SET user_id=$1, category=$2, amount=$3, description=$4, date=$5
		WHERE id=$6`,
		exp.UserId, exp.Category, exp.Amount, exp.Description, exp.Date, exp.Id)

	if err != nil {
		http.Error(w, "Failed to update expense: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Expense updated successfully")
}

func main() {
	var err error
	db, err = sql.Open("postgres", "host=db port=5432 user=aknurturakhan dbname=expense_db sslmode=disable")
	if err != nil {
		log.Fatal("Error connecting to the database:", err)
	}

	http.HandleFunc("/get_expenses", getExpenses)
	http.HandleFunc("/add_expense", addExpense)
	http.HandleFunc("/delete_expense", deleteExpense)
	http.HandleFunc("/update_expense", updateExpense)


	http.Handle("/swagger/", http.StripPrefix("/swagger/", httpSwagger.WrapHandler))
	http.Handle("/", http.FileServer(http.Dir("./")))

	fmt.Println("🚀 Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
	fmt.Println()
}

// ~/Downloads/expense-tracker/expense-service