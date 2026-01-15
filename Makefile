# Makefile - one-command setup & run

.PHONY: all build run dashboard clean

all: build

build:
	go build -o athletelog-cli
	cd ts-dashboard && npm install && npm run build

run: build
	./athletelog-cli

dashboard: build
	./athletelog-cli dashboard

# Start static server in background
serve:
	npx serve . &

clean:
	rm -f athletelog-cli
	rm -f ts-dashboard/main.js
	cd ts-dashboard && npm run clean || true
