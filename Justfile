generate:
  ./scripts/generate.sh

lint:
  ./scripts/lint.sh

reset-db:
  if test -f "server/internal/store/shelter.db"; then \
    rm server/internal/store/shelter.db ; \
  fi

start:
  cd server && go run cmd/main.go

test:
  cd tests && pytest test.py
