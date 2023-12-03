# Introduction: Creating a Car-Pooling Platform 🚗
Welcome to my Car-Pooling Platform, a microservices-based solution for efficient and convenient carpooling service for users! This platform caters to two primary user groups - passengers and car owners. 

During the creation of a user account, a default passenger profile is created where first name, last name, mobile number, and email address are required. For a user who is also a car owner, the default passenger profile can be changed to a car owner profile. 
The user is required to provide a driver’s license number and a car plate number. Subsequently, users can update any information in their account. Users are able to delete their accounts after 1-year if the car-pooling platform is no longer relevant to them. The 1-year data retention is for audit purposes.
Users who are car owners publish car-pooling trips with addresses of pick-up locations, alternative pick-up locations, start traveling time, address of destination and number of passengers their car can accommodate. 
When trips are published, car owners wait for passengers to select and enrol for respective trips.
The platform assigns seats based on a first-come-first-serve basis. The car owners will be able to start trips or cancel them 30 minutes before the scheduled time. 
Passengers can search and browse through a listing of published trips and enrol to any published trips as long as there are vacancies, and no date and time conflicts with enrolled trips. 
Users can retrieve all trips that have taken before in reverse chronological order. 

## Objectives 🎯
1. Demonstrate Ability to Develop REST APIs: Showcase your proficiency in designing and implementing RESTful APIs.

2. Conscientious Consideration in Designing Microservices: Develop microservices with careful consideration, adhering to best practices in design and architecture.

## Design considerations

* For my design considerations,  after reading through the requirements, I had to implement a **2 tiered microservice** architecture with implementation of a persistant storage. I identified key entities and business capabilities in the system that I have created, such as Users and Trips. I have created 2 separate microservices for each entity or business capability to achieve a modular and scalable architecture.
  
* For **User Service**, it is responsible for managing user information and will be exposed to endpoints for creating, updating, deleting, and retrieving user details. Ensured proper validation of user data during creation and update. Likewise for **Trip Service**, it manages information related to carpooling trips.It will be exposed to endpoints for creating, updating, deleting, and retrieving trip details. Implemented logic to enforce business rules, such as ensuring a car owner is valid and that trip times are in the future.
  
* **Data Storage:** For my database, I used a singluar mysql database that has 2 tables, for users and trips.
  
* **Communication between Microservices:** Used RESTful APIs for communication between microservices.Implemented proper error handling and response codes.
  
* **Event Sourcing and Logging:** Consider implementing event sourcing for tracking changes in state. Implement logging for each microservice to capture important events and errors. 
  
* **Error Handling and Resilience:** Implement proper error-handling mechanisms in each microservice. Use Circuit Breaker patterns to handle faults gracefully and prevent cascading failures. Implement retries for transient failures and timeouts so that code will still be able to function even though a functionality cannot be used.


## Architectural Diagram 📐


<!-- GETTING STARTED -->
## Getting Started 

This is an example of how you may give instructions on setting up your project locally.
To get a local copy up and running follow these simple example steps.
1. Import the necessary packages in main.go:
* gorilla mux
```sh
go get -u github.com/gorilla/mux
```

* mysql
```sh
go get -u github.com/go-sql-driver/mysql
```

2. Setting up your database:
   
Create user:
```sql
CREATE USER 'user'@'localhost' IDENTIFIED BY 'password';
GRANT ALL ON *.* TO 'user'@'localhost'
```

Create database:
```sql
CREATE database carpooling;
USE carpooling;
CREATE TABLE IF NOT EXISTS users (
    id VARCHAR(255) PRIMARY KEY,
    first_name VARCHAR(255),programmes
    last_name VARCHAR(255),
    mobile_number VARCHAR(20),
    email VARCHAR(255),
    driver_license VARCHAR(255),
    car_plate_number VARCHAR(20),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- Create the 'trips' table
CREATE TABLE IF NOT EXISTS trips (
    id VARCHAR(255) PRIMARY KEY,
    car_owner_id VARCHAR(255),
    pickup_location VARCHAR(255),
    alternative_pickup VARCHAR(255),
    start_travel_time TIMESTAMP,
    destination VARCHAR(255),
    available_seats INT,
    creation_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    cancellation_time TIMESTAMP
);
```

2. Run main.go using the following command
```sh
go run main.go
```

3. Run console.go using the following command
```sh
go run main.go
```

4. Now you will be able to access the console application.

## Prerequisites

* You should be able to set up your go lang project locally by using the **"go mod init" command**.

* If you face any issues, please head over [this website](https://go.dev/doc/tutorial/getting-started) on how to create a local Go program.

* **mysql workbench** should be installed


<!-- USAGE EXAMPLES -->
## Usage

Use this space to show useful examples of how a project can be used. 





## Contributing

Sole contributor: Zacharia Aslam

## License

[MIT](https://choosealicense.com/licenses/mit/)
