# Jobsity Stock Chat


## Requirements:

### Docker
Download and install Docker.
If you are on a Mac (see https://docs.docker.com/docker-for-mac/install).
If you are on Ubuntu (https://docs.docker.com/install/linux/docker-ce/ubuntu/).

## Running the development environment

The local development environment consists of 6 docker containers:
- `postgres`: the database service.
- `migrate`: in charge of running the migrations to create the tables for storing users and posts. It also seeds the `StockBot` and two more users.
- `rabbitmq`: the message broker service.
- `srv`: Go http server. Handles users and posts.
- `bot`: Go worker. Handles stocks processing.
- `vue`: Vue.js frontend application.

#### Running the application
The recommended way to start all the services together is executing `docker-compose up` in the root project folder.

Please wait until all services are fully up, it may take a while especially for the rabbitmq container.

- Access the front-end site with the following URLs: `http://localhost:3000`.
Two users are already created: `Alice` and `Bob`.
You can access with any of them using the password `12345`, or you can create new users using the signup endpoint
 -  `POST http://localhost:5000/signup` to create a user
  <pre>Request: <br>{<br>"username": "UserThree",<br>"password": "12345"<br>}</pre>
  <pre>Response: <br>{<br>"id": "4e6bff32-ec75-4996-8721-03bf9bc5b785", <br>"username": "UserThree"<br>}</pre>

#### Running Separately

To run the `srv` or `bot` services locally (outside of docker)

Stop the container
```
docker stop stockchat-srv-1 (or stockchat-bot-1)
```

Change the env variables specified in the `.env` file: 
```
POSTGRES_HOST=localhost:5432
RABBITMQ_HOST=localhost:5672
```

Run the local server (you should have Go 1.20 installed):
```
cd server (or cd bot)
go run cmd/main.go
```

To run the `vue` client service locally (outside of docker), stop the container and run (you should have Node.js v.16.17):
```
docker stop stockchat-vue-1
cd client
npm install
npm run dev
```

#### Stoping the application
To stop all the services execute `docker-compose down` in the root project folder.