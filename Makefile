APP_NAME=momentum-app
PORT=8080

build:	## Build the container
				docker build . --tag $(APP_NAME)

run:		## Run container on port
				docker run -i -t --rm -p=$(PORT):$(PORT) --name="$(APP_NAME)" $(APP_NAME)

stop:		## Stop and remoove running container
				docker stop $(APP_NAME); docker rm $(APP_NAME)
