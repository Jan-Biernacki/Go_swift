# Go_Swift

This project is a REST API for managing SWIFT codes using Go, Gin, and PostgreSQL. The API allows fetching, creating, and deleting SWIFT codes stored in a PostgreSQL database. The application is containerized using Docker and runs on localhost:8080.
This project works based on the provided csv file, other source of data must follow the format within swift_codes.csv in order to run properly.


## Requirements

- Docker & Docker Compose installed
- (Alternatively) Go 1.20+ if you want to run locally

## How to Run (Docker Compose)

1. Clone this repository

    git clone <repository-url>
    cd go_swift

2. Ensure Docker is running, then execute:

    docker-compose up --build

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

6. Database veryfication : 

    docker exec -it go_swift-db-1 psql -U postgres -d swift_codes

        SELECT COUNT(*) FROM swift_codes;

7. 