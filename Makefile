FRONT_END_BINARY=frontApp

up:
	@echo "Starting Docker images..."
	docker-compose up -d
	@echo "Docker images started!"

up_build: build_broker build_auth build_logger build_mail
	@echo "Stopping Docker images (if running...)"
	docker-compose down
	@echo "Building (when required) and starting Docker images..."
	docker-compose up -d --build
	@echo "Docker images started!"

down:
	@echo "Stopping Docker images..."
	docker-compose down
	@echo "Docker images stopped!"

build_broker:
	@echo "Building broker binary..."
	cd ./broker && env GOOS=linux CGO_ENABLED=0 go build -o brokerService ./api
	@echo "Broker image built!"

build_logger:
	@echo "Building logger binary..."
	cd ./logger && env GOOS=linux CGO_ENABLED=0 go build -o loggerService ./api
	@echo "Logger image built!"

build_auth:
	@echo "Building auth binary..."
	cd ./auth && env GOOS=linux CGO_ENABLED=0 go build -o authService ./api
	@echo "Auth image built!"

build_mail:
	@echo "Building mail binary..."
	cd ./mail && env GOOS=linux CGO_ENABLED=0 go build -o mailService ./api
	@echo "Mail image built!"

build_front:
	@echo "Building front end binary..."
	cd ./go-micro && env CGO_ENABLED=0 go build -o ${FRONT_END_BINARY} ./cmd/web
	@echo "Done!"

start: build_front
	@echo "Starting front end"
	cd ./go-micro && ./${FRONT_END_BINARY} &

stop:
	@echo "Stopping front end"
	@-pkill -SIGTERM -f "./${FRONT_END_BINARY}"
	@echo "Front end stopped"
