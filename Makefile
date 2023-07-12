APP_NAME=ib-ui

## build: builds all binaries
build: clean build_app
	@printf "All binaries built!\n"

## clean: cleans all binaries and runs go clean
clean:
	@echo "Cleaning..."
	@- rm -f dist/*
	@go clean
	@echo "Cleaned!"

build: build_app

#https://pkg.go.dev/cmd/go#hdr-Compile_packages_and_dependencies
## build_app: builds the front end
build_app:
	@echo "Building application ${APP_NAME}..."
	@go build -o dist/${APP_NAME} ./cmd/web
	@go build -o dist/ib-cli ./cmd/cli
	@echo "Application '${APP_NAME}' and 'ib-cli' built!"

build_app_amd:
	@echo "Building application ${APP_NAME}..."
	@GOARCH=amd64;GOOS=darwin;go build -o dist/amd/${APP_NAME} ./cmd/web
	@echo "Application 'amd/${APP_NAME}' and 'ib-cli' built!"

## start: all applications in project
start: build start_app

## starts the application
start_app:
	@echo "Starting the application ${APP_NAME}..."
	./dist/${APP_NAME} &
	@echo "Application ${APP_NAME} running!"

## stop: stops application
stop: stop_app
	@echo "All applications stopped"

## stop_app: stops the application
stop_app:
	@echo "Stopping the application ${APP_NAME} ..."
	@-pkill -SIGTERM -f "${APP_NAME}"
	@echo "Stopped ${APP_NAME} application"

restart: stop start

