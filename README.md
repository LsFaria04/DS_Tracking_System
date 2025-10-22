# Tracking system backend

The tracking system backend is developed using the Go programing language

To start developing use the following command to create the docker containers:

```shell
#Use the dev composer file to start postgres, pgadmin and the go containers
docker compose -f compose.dev.yml up --watch
```

To stop the containers and remove them use:

```shell
docker compose -f compose.dev.yml  down
```

Whenever changes are made to a Dockerfile or Compose file, use the following commands to rebuild and ensure those updates  are applied to your containers:

```shell
docker compose -f compose.dev.yml down -v
docker compose -f compose.dev.yml build --no-cache
docker compose -f compose.dev.yml up --watch
```

To run the production Compose and Dockerfile (i.e., the ones without the .dev suffix), make sure to create a .env file containing the required environment variables. Place this file in the same directory as your Compose configuration to ensure proper loading.

**DO NOT COMMIT THE .ENV FILE!**