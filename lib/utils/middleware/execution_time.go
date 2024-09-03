package middleware

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
)


type CreateLogExecutionTimeParams struct {
	Environment   string `json:"environment"`
	Status		  int 	 `json:"status"`
	Service       string `json:"service"`
	Route         string `json:"route"`
	Method        string `json:"method"`
	ExecutionTime int64  `json:"execution_time" validate:"required"`
}

func ExecutionTimeMiddleware(db *sql.DB, environment string, service string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.Context()
		startTime := time.Now()

		defer func() {
			duration := time.Since(startTime)

			params := CreateLogExecutionTimeParams{
				Environment:   environment,
				Status: c.Response().StatusCode(),
				Service:       service,
				Route:         c.OriginalURL(),
				Method:        c.Method(),
				ExecutionTime: int64(duration / time.Millisecond),
			}

			_, err := CreateLogExecutionTimeQuery(db, ctx, params)
			if err != nil {
				log.Printf("execution time errors: %v", err.Error())
			}

			log.Printf("Status: %d, Route: %s, Method: %s, Execution Time: %vms", params.Status, params.Route, params.Method, params.ExecutionTime)
		}()

		return c.Next()
	}
}

func CreateLogExecutionTimeQuery(db *sql.DB, ctx context.Context, body CreateLogExecutionTimeParams) (_ int64, err error) {
	query := `
		INSERT INTO 
			log_execution_time (environment, status, service, route, method, execution_time)
		VALUES
			(?, ?, ?, ?, ?, ?)
	`

	last_id := int64(0)

	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		return last_id, err
	}

	// Close Statement on end function
	defer stmt.Close()

	values := []interface{}{
		body.Environment,
		body.Status,
		body.Service,
		body.Route,
		body.Method,
		body.ExecutionTime,
	}

	// Exec query with values
	query_result, err := stmt.ExecContext(ctx, values...)
	if err != nil {
		return last_id, err
	}

	last_id, err = query_result.LastInsertId()
	if err != nil {
		return last_id, err
	}

	return last_id, err
}
