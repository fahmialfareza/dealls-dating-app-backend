.PHONY: run-http generate-mocks test coverage force

lint:
	@go mod vendor
	@echo "Running golang lint"
	@golangci-lint run 
	@rm -rf vendor

run-http:
	@echo "Running http server..."
	@go run cmd/main.go http

generate-mocks:
	@mockgen -package=mocks_redis -destination=configs/cache/redis/mocks/mock_redis.go github.com/fahmialfareza/deals-dating-app-backend/configs/cache/redis IRedis
	@mockgen -package=mocks_postgres_repository -destination=internal/repository/postgres/mocks/mocks_postgres.go github.com/fahmialfareza/deals-dating-app-backend/internal/repository/postgres IPostgresRepository,IUserPostgresRepository,ISwipePostgresRepository,IPurchasePostgresRepository
	@mockgen -package=mocks_redis_respository -destination=internal/repository/redis/mocks/mocks_redis.go github.com/fahmialfareza/deals-dating-app-backend/internal/repository/redis IRedisRepository,IUserRedisRepository
	@mockgen -package=mocks_repository -destination=internal/repository/mocks/mocks_repository.go github.com/fahmialfareza/deals-dating-app-backend/internal/repository IRepository,IUserRepository,ISwiperRepository,IPurchaseRepository

test:
	@echo "Running tests, excluding mocks..."
	@mkdir -p coverage
	@find ./internal -type d ! -path '*/mocks' | while read -r dir; do \
	    go test -coverprofile="coverage/$$(basename $$dir).out" "$$dir"; \
	done
	@echo "Combining coverage reports..."
	@echo "mode: set" > coverage/coverage.out
	@find coverage -name '*.out' ! -name 'coverage.out' | xargs -I {} tail -n +2 {} >> coverage/coverage.out
	@go tool cover -func=coverage/coverage.out | grep total | awk '{print "Total Coverage:", $$3}'

coverage: force
	@echo "Combining coverage reports..."
	@echo "mode: set" > coverage/coverage.out
	@find coverage -name '*.out' ! -name 'coverage.out' | xargs -I {} tail -n +2 {} >> coverage/coverage.out
	@go tool cover -html=coverage/coverage.out
	@echo "Coverage report available in coverage/coverage.out"
	@go tool cover -func=coverage/coverage.out | grep total | awk '{print "Total Coverage:", $$3}'

force:
