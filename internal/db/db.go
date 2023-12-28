package db

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/lib/pq"
)

// Prompt represents the structure of our resource
type Prompt struct {
	ID        int
	Prompt    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// DB is a database handle
var DB *sql.DB

// Connect initializes the database connection
func Connect() error {
	var err error
	connectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))

	DB, err = sql.Open("postgres", connectionString)
	if err != nil {
		return err
	}

	return DB.Ping()
}

// GetPrompts retrieves all prompts from the database
func GetPrompts() ([]Prompt, error) {
	rows, err := DB.Query("SELECT id, prompt, created_at, updated_at FROM prompts")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var prompts []Prompt
	for rows.Next() {
		var p Prompt
		if err := rows.Scan(&p.ID, &p.Prompt, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		prompts = append(prompts, p)
	}

	return prompts, nil
}

// CreatePrompt inserts a new prompt into the database
func CreatePrompt(promptText string) (int, error) {
	var promptID int
	err := DB.QueryRow("INSERT INTO prompts(prompt, created_at, updated_at) VALUES($1, $2, $3) RETURNING id", promptText, time.Now(), time.Now()).Scan(&promptID)
	if err != nil {
		return 0, err
	}

	return promptID, nil
}

// GetPrompt retrieves a single prompt by its ID
func GetPrompt(id int) (*Prompt, error) {
	var p Prompt
	err := DB.QueryRow("SELECT id, prompt, created_at, updated_at FROM prompts WHERE id = $1", id).Scan(&p.ID, &p.Prompt, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &p, nil
}

// DeletePrompt removes a prompt from the database
func DeletePrompt(id int) error {
	_, err := DB.Exec("DELETE FROM prompts WHERE id = $1", id)
	return err
}
