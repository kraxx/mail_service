NAME =	mail_service

DEPS =	github.com/gorilla/context \
		github.com/gorilla/mux \
		github.com/joho/godotenv

OUTPUT_DIR =	./bin/
OUTPUT_BINARY =	${OUTPUT_DIR}${NAME}

RED =				\033[31m
GREEN =				\033[32m
BLUE =				\033[34m
YELLOW =			\033[33m
MAGENTA =			\033[35m
GREY =				\033[37m
GREEN_LIGHT =		\033[92m
YELLOW_LIGHT =		\033[93m
YELLOW_BOLD =		\033[1;33m
YELLOW_LIGHT_BOLD =	\033[1;93m
MAGENTA_LIGHT =		\033[95m
BLINK =				\033[5m
GREEN_LIGHT_BLINK =	\033[5;92m
END_COLOUR =		\033[0m

TAG? =	$(shell git rev-list HEAD --max-count=1 --abbrev-commit)

.PHONY: install build serve test clean heroku docker re # invalidate these commands if they exist outside this script
.SILENT: # Prepends everything with @ (command executed without printing to stdout)
all: install build serve
install:
	echo "${YELLOW_LIGHT_BOLD}Installing dependencies${END_COLOUR}"
	go get ${DEPS}
build:
	echo "${YELLOW_LIGHT_BOLD}Building binary${END_COLOUR}"
	go build -o ${OUTPUT_BINARY} -ldflags "-X main.version=$(TAG)" .
serve:
	echo "${YELLOW_LIGHT_BOLD}Executing binary${END_COLOUR}"
	${OUTPUT_BINARY}
test:
	echo "${YELLOW_LIGHT}Testing ${BLUE}${NAME}${END_COLOUR}"
	go run main.go
clean:
	echo "${YELLOW_LIGHT_BOLD}Cleaning installations and binary${END_COLOUR}"
	go clean
	rm -f ${OUTPUT_BINARY}
heroku:
	echo "${GREEN_LIGHT}Setting up Godep files for Heroku${END_COLOUR}"
	godep save
	echo "${GREEN_LIGHT}Git commits before push${END_COLOUR}"
	git add .
	git commit -m "Ran \"make heroku\""
	echo "${GREEN_LIGHT}Pushing to Heroku${END_COLOUR}"
	git push heroku master
	rm -rf Godeps vendor
	echo "${GREEN_LIGHT}Removed Godep files because we don't need them${END_COLOUR}"
docker:
	echo "${GREEN_LIGHT}Building Docker image on Heroku registry${END_COLOUR}"
	docker build -t registry.heroku.com/${NAME}/web .
	echo "${GREEN_LIGHT}Deploying Docker image to Heroku${END_COLOUR}"
	docker push registry.heroku.com/${NAME}/web
re: clean all
# pack:
#	GOOS=linux make build
#	sudo docker build -t us.gcr.io/itsashopchallenge/${NAME}:$(TAG) .

# upload:
# 	gcloud docker -- us.gcr.io/itsashopchallenge/${NAME}:$(TAG)

# deploy:
# 	envsubst < k8s/deployment.yml | kubectl apply -f -

# ship: test pack upload deploy clean
