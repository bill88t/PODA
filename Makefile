PORT=5173

.PHONY: breadport
breadport:
	@$(eval PORT = 80)

.PHONY: frontend
frontend:
	cd frontend; npm run build

build: main frontend

run:
	PORT=$(PORT) ./main

breadrun: breadport | run

main: main.go
	go build $<

clean:
	rm main
