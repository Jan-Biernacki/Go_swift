# Go Swift API

This project implements a RESTful API for managing SWIFT (BIC) codes using Go, Gin, and PostgreSQL. It includes functionality for parsing SWIFT code data from a CSV file, storing it in a PostgreSQL database, and exposing endpoints for retrieving, creating, and deleting records.

## Project Structure
    go_swift/
    ├── cmd/
    │   └── server/
    │       └── main.go  # Application entry point
    ├── data/
    │   └── swift_codes.csv  # CSV file with SWIFT code data
    ├── internal/
    │   ├── controllers/
    │   │   └── swift_controller.go  # REST API endpoints
    │   ├── models/
    │   │   └── swift_code.go  # SWIFT code model (GORM)
    │   ├── parsing/
    │   │   └── parse.go  # CSV parsing logic
    │   └── repositories/
    │       └── database.go  # Database initialization and migrations
    ├── go.mod
    ├── go.sum
    ├── Dockerfile
    ├── docker-compose.yml
    └── README.md



## Requirements

- Docker & Docker Compose installed
- Go 1.23 or later

## How to Run (Docker Compose)

1. Clone this repository

    git clone <repository-url>
    cd go_swift

2. Ensure Docker is running, then execute:

    docker-compose up --build

This command will build the application image, start the PostgreSQL container, and run the API on port 8080.

3. Check running containers : 

    docker ps



4. Endpoints 

- **GET** `/v1/swift-codes/{swiftCode}`
  - Returns JSON. If HQ, also includes "branches".
- **GET** `/v1/swift-codes/country/{iso2}`
  - Returns all codes for a country.
- **POST** `/v1/swift-codes`
  - Create a new SWIFT code (JSON body).
- **DELETE** `/v1/swift-codes/{swiftCode}`
  - Delete a SWIFT code if bankName and countryISO2 match.


5. To stop the application and database, run:

    docker compose down

6. Database veryfication (manual way to check whether the database contains records) : 

    docker exec -it go_swift-db-1 psql -U postgres -d swift_codes

        SELECT COUNT(*) FROM swift_codes;

7.  Endpoint usage example: 

    Look for an existing record :

        http://localhost:8080/v1/swift-codes/AAISALTRXXX

    Should return : 
      {
    "swiftCode": "NEWTEST33XXX",
    "bankName": "New Test Bank",
    "address": "456 New St",
    "countryISO2": "US",
    "countryName": "UNITED STATES",
    "isHeadquarter": true
      }


    Look for non existing record : 

        curl http://localhost:8080/v1/swift-codes/MRINDJJDXXX

    Load new data to the databasefor convinienve, JSON data regarding this bank is stored in a payload.json file :

        curl.exe -X POST -H "Content-Type: application/json" -d "@payload.json" http://localhost:8080/v1/swift-codes

    Check that the new record exists : 

        curl http://localhost:8080/v1/swift-codes/MRINDJJDXXX

    Remove new record

        curl.exe -X DELETE http://localhost:8080/v1/swift-codes/MRINDJJDXXX
    
    Verify that the record no longer exists : 

        curl http://localhost:8080/v1/swift-codes/MRINDJJDXXX

8. Unit test : 
    go test -v ./internal/parsing
9. Integration test :  
    go test -v

10. Remarks 

    - Lack of variety of unit tests is caused by not using TDD since the beggining of the project, I will remember to use it in my future projects. Good lesson.
    - There is a minor failsafe in regard to SWIFT codes, the program expects them to be 8 or 11 chars long. Eight char long code case is for a missinput in csv file assuming that somebody could copy the first 8 chars of a headquaters code negelecting the triple X at the end, hence if there is a code having just 8 chars it will be considered missinputed headquaters.

