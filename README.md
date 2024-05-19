# Exchange rate service
This is a service with APIs that will allow you to:
- find out the current dollar (USD) to hryvnia (UAH) exchange rate;
- sign an email to receive information on exchange rate changes.

Undone task:  
- send email with the currency rate to subscribed emails;
- validation for email format.  

## Folder structure
/api - holds the swagger file
/cmd - contains main.go file and config file
/db - holds the liquibase migrations
/deployments - contains docker-compose.yaml and .env files
/internal - contains all internal service logics, e.g. services, api handler, store etc.

## Run app
1. Navigate to deployments directory:
```
cd ./deployments
``` 
2. Run the command:
```
docker-compose up --build
```
3. When app service is running test the API request. If the .env file isn't changed, the service API should be available at port **9000**.
