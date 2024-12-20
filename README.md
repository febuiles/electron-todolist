## pooltasks

### Installation

#### Server

You can run the backend using `docker-compose`:

```sh
docker-compose up --build
```

If you don't have `docker-compose` available, you can manually build
the Docker image and run it (`docker build -t pooltasks . && docker
run -p 8000:8080 pooltasks`), or build it with Go (`go mod download && go build .`).

You can verify it's working by hitting the ping endpoint:

```
curl http://localhost:8000/_ping
# => PONG
```

#### Client 

You can run the Electron app 

### Design 

### Considerations

Given the nature of the task I took some shortcuts:

* I decided to use SQLite instead of something like Postgres or MySQL
  to make my life easier. If running inside Docker, the database will
  be wiped every time the container is stopped.
* The client stores sync data locally in the appData folder. In order
  to run more than one instance of the application, they need to be
  run as different users.
* Given the requirement to have the server running locally, additional
  instances of the app running on different machines forces the user
  to enter the backend hostname by hand. Ideally this would be handled
  by sharing a deep link instead of only a slug.
* The server code is very light on validations, especifically: it has
  no authentication/authorization. All the endpoints are open to
  manipulation.
* The Electron application needs more tests.
* I iterated very quickly on this, and am pretty certain there's dead
  code inside the project.

