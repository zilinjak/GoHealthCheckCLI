# Go Health Checker
> A simple health checker for your web applications.
> 
> Go Health Checker checks your website specified in CLI args every couple of seconds. 
> This is repeated until CTRL+C is pressed. 
> Once that is done, the program will print stats about analysis.

## Run the app
...

## Architecture

The application is built using these components:
 - View - representation of the data to user
   - CLIView - command line interface
 - Model - data structure
   - HealthCheckResult - data structure for health check result
 - Service - business logic
   - HealthCheckService - service for health check, pinging the websites
 - Store - data access layer
   - InMemoryStore - in-memory store for health check results
 - Controller - putting it all together
   - HealthCheckController - controller for health check, responsible for starting and stopping the health check



