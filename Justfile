generate:
  ./scripts/generate.sh

lint:
  ./scripts/lint.sh

reset-db:
  rm server/internal/store/shelter.db

start:
  cd server && go run cmd/main.go
