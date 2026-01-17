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

main: main.go
	@echo "Building backend.."
	@go build $<

clean:
	@echo "Cleaning.."
	@rm -f main
