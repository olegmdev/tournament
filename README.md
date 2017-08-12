# Social Tournament on Go

Specific tournament application which provides isolated transaction locks using `SELECT ... FOR UPDATE` statements.
Its built using gin-gonic for routing, gorm as ORM and viper for configuration stuff.

### Getting started

1. Clone this repository
2. Configure the project
2.1. Make sure to update `config/development.yaml` to run this app manually via CLI (outside of docker container)
2.2. Make sure to update `config/local.yaml` to run this app within docker container
2.3. You should have following env variables for docker: `POSTGRES_USER`, `POSTGRES_PASSWORD`, `POSTGRES_DB`
2. Build services `docker-compose build`
4. Run the app with following command `docker-compose up -d`
