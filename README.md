# Form3 account API client #

### by Aleksandr Kulik

## Description

The project is split into 3 packages:
* **f3** - tiny package containing F3-related headers, endpoints and a bit of logic, which is not to be shared with client users
* **accounts** - the API client itself covered with unit tests
* **test** - separated package for everything related to e2e testing
 
### Client configuration

The only required parameter for client creation is **api url**, but I've also added several functional options for additional client configuration:
* WithTimeout - timeout applied to every client request
* WithRateLimiter - although there is a rate limiting functionality on the server side, I still decided to add it also to the client, as it might be useful in some cases in order to prevent restrictions on the API side in the first place
* WithBasePath - path added to the url, might be useful for setting api version

## Testing

One of the main requirements for the task was to be able to run all the tests  with docker-compose up command.
The solution was implemented accordingly: all the tests are run as soon as all the related containers (i.e. database and api for e2e testing) are alive. There is also a Makefile in place for simplifying the testing, so all the testing can be performed only using make commands. Docker compose can be managed with make commands as well.

In order to start the containers (will also start the tests in the container):
```
make up
```
When you're done with the testing you can finish with 
```
make down
```
For running all the tests locally:
```
make test
```

Running all the BDD scenarios (requires either API container to be alive or to set up non-local API url):
```
make godog
```

Running only unit tests with coverage report:
```
make cover
```
