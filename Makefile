PORT=5173

.PHONY: all
all: clean build

.PHONY: check-go
check-go:
	@command -v go >/dev/null 2>&1 || { echo >&2 "Error: go is not installed."; exit 1; }

# tsc already is running from vite
# tsc is not global always but it is always local to node nodules
# .PHONY: check-tsc
# check-tsc:
# 	@command -v frontend/node_modules/typescript/bin/tsc >/dev/null 2>&1 || { echo >&2 "Error: typescript is not installed."; exit 1; }

.PHONY: breadport
breadport:
	@$(eval PORT = 80)

.PHONY: frontend
frontend: # check-tsc - already from vite
	@echo "Building frontend.."
	@cd frontend; npm run build

build: check-go main frontend

npm:
	@echo "Refreshing NPM packages"
	@cd frontend; npm install

run:
	@echo "Running on port $PORT"
	@PORT=$(PORT) ./main

breadrun: breadport | run

main:
	@echo "Building backend.."
	@go build -o main .

clean:
	@echo "Cleaning build.."
	@rm -f main

full_clean: rmdb clean

rmdb:
	@echo "Wiping database.."
	@rm -v poda.db-wal 2>/dev/null || true
	@rm -v poda.db-shm 2>/dev/null || true
	@rm -v poda.db 2>/dev/null || true

test_bill_reg:
	@echo "Registering bill88t.."
	curl -X POST http://localhost:5173/api/v1/users/signup \
	  -H "Content-Type: application/json" \
	  -d '{"fname": "Bill", "lname": "Sideris", "email": "bill88t@feline.gr", "password": "securepassword123", "birthday": "2002-12-01" }'
	@echo
	@echo

test_bill_login:
	@echo "Testing login.."
	curl -X POST http://localhost:5173/api/v1/users/login \
	  -H "Content-Type: application/json" \
	  -d '{ "email":  "bill88t@feline.gr", "password": "securepassword123" }'
	@echo
	@echo

runtests: test_bill_reg test_bill_login
