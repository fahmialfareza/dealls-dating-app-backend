package postgres

import "github.com/fahmialfareza/deals-dating-app-backend/configs/database/postgres"

func NewPostgresUtil(postgres postgres.IPostgres) {
	PostgresClient = postgres
}
