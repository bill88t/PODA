PORT=5137

.PHONY: breadport
breadport:
	$(eval PORT = 80)
	echo $(PORT)

.PHONY: frontend
frontend:
	cd frontend; npm run build

build: main frontend

run: build
	PORT=$(PORT) ./main

breadrun: breadport | run

main: main.go
	PORT=$(PORT) go build $<

clean:
	rm main
