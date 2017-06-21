## How to start working with API?


1. Install [docker](https://docs.docker.com/engine/installation/) and [docker-compose](https://docs.docker.com/compose/install/)

2. From projects root build containers

`$ docker-compose build`

3. Start a project

`$ docker-compose up`

4. Bash into Golang container

`$ docker exec -it golang_echo bash`

5. Init SQL schema with fallowing:

`$ cd /db && go run init.go ./init.sql`


### You are ready to GO! Server works on port 8001