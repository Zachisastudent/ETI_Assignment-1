package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
)

const baseURL = "http://localhost:8222/api/v1"

// User represents a user in the car-pooling platform
type User struct {
	ID             string    `json:"id"`
	FirstName      string    `json:"first_name"`
	LastName       string    `json:"last_name"`
	MobileNumber   string    `json:"mobile_number"`
	Email          string    `json:"email"`
	DriverLicense  string    `json:"driver_license,omitempty"`
	CarPlateNumber string    `json:"car_plate_number,omitempty"`
	IsCarOwner     bool      `json:"is_car_owner"`
	CreatedAt      time.Time `json:"created_at"`
}

// Trip represents a car-pooling trip published by a car owner
type Trip struct {
	ID                 string    `json:"id"`
	CarOwnerID         string    `json:"car_owner_id"`
	PickupLocation     string    `json:"pickup_location"`
	AltPickupLocation  string    `json:"alt_pickup_location,omitempty"`
	StartTime          time.Time `json:"start_time"`
	Destination        string    `json:"destination"`
	AvailableSeats     int       `json:"available_seats"`
	EnrolledPassengers []string  `json:"enrolled_passengers,omitempty"`
	TotalSeats         int       `json:"total_seats"`
	Started            bool      `json:"started"` // New field
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		printMenu()

		fmt.Print("Enter an option: ")
		scanner.Scan()
		option := scanner.Text()

		switch option {
		case "1":
			listAllUsers()
		case "2":
			createNewUser(scanner)
		case "3":
			updateUser(scanner)
		case "4":
			deleteUser(scanner)
		case "5":
			listAllTrips()
		case "6":
			createNewTrip(scanner)
		case "7":
			enrollPassenger(scanner)
		case "8":
			startTrip(scanner)
		case "9":
			cancelTrip(scanner)
		case "10":
			listTripStatus(scanner)
		case "11":
			fmt.Println("Exiting the program.")
			return
		default:
			fmt.Println("Invalid option. Please try again.")
		}
	}
}

func printMenu() {
	fmt.Println("1. List all users")
	fmt.Println("2. Create new user")
	fmt.Println("3. Update user")
	fmt.Println("4. Delete user")
	fmt.Println("5. List all trips")
	fmt.Println("6. Create new trip")
	fmt.Println("7. Enroll passenger in a trip")
	fmt.Println("8. Start a trip")
	fmt.Println("9. Delete/cancel a trip")
	fmt.Println("10.List trip status")
	fmt.Println("11. Quit")
}

func listAllUsers() {
	getData("users")
}

func createNewUser(scanner *bufio.Scanner) {
	fmt.Print("Enter the ID of the user to be created: ")
	scanner.Scan()
	userID := scanner.Text()

	// Check if the user already exists
	if userExists(userID) {
		fmt.Println("Error - User already exists")
		return
	}

	fmt.Print("Enter the first name: ")
	scanner.Scan()
	firstName := scanner.Text()

	fmt.Print("Enter the last name: ")
	scanner.Scan()
	lastName := scanner.Text()

	fmt.Print("Enter the mobile number: ")
	scanner.Scan()
	mobileNumber := scanner.Text()

	fmt.Print("Enter the email address: ")
	scanner.Scan()
	email := scanner.Text()

	fmt.Print("Is the user also a car owner? (true/false): ")
	scanner.Scan()
	isCarOwnerStr := scanner.Text()
	isCarOwner, err := strconv.ParseBool(isCarOwnerStr)
	if err != nil {
		fmt.Println("Invalid input for car owner. Please enter true or false.")
		return
	}

	newUser := map[string]interface{}{
		"id":            userID,
		"first_name":    firstName,
		"last_name":     lastName,
		"mobile_number": mobileNumber,
		"email":         email,
		"is_car_owner":  isCarOwner,
	}

	if isCarOwner {
		fmt.Print("Enter the driver's license number: ")
		scanner.Scan()
		driverLicense := scanner.Text()

		fmt.Print("Enter the car plate number: ")
		scanner.Scan()
		carPlateNumber := scanner.Text()

		newUser["driver_license"] = driverLicense
		newUser["car_plate_number"] = carPlateNumber
	}

	createOrUpdateUser("POST", userID, newUser)
}

func updateUser(scanner *bufio.Scanner) {
	fmt.Print("Enter the ID of the user to be updated: ")
	scanner.Scan()
	userID := scanner.Text()

	// Check if the user exists
	if !userExists(userID) {
		fmt.Println("Error - User does not exist")
		return
	}

	fmt.Print("Enter the first name: ")
	scanner.Scan()
	firstName := scanner.Text()

	fmt.Print("Enter the last name: ")
	scanner.Scan()
	lastName := scanner.Text()

	fmt.Print("Enter the mobile number: ")
	scanner.Scan()
	mobileNumber := scanner.Text()

	fmt.Print("Enter the email address: ")
	scanner.Scan()
	email := scanner.Text()

	fmt.Print("Is the user also a car owner? (true/false): ")
	scanner.Scan()
	isCarOwnerStr := scanner.Text()
	isCarOwner, err := strconv.ParseBool(isCarOwnerStr)
	if err != nil {
		fmt.Println("Invalid input for car owner. Please enter true or false.")
		return
	}

	updatedUser := map[string]interface{}{
		"id":            userID,
		"first_name":    firstName,
		"last_name":     lastName,
		"mobile_number": mobileNumber,
		"email":         email,
		"is_car_owner":  isCarOwner,
	}

	if isCarOwner {
		fmt.Print("Enter the driver's license number: ")
		scanner.Scan()
		driverLicense := scanner.Text()

		fmt.Print("Enter the car plate number: ")
		scanner.Scan()
		carPlateNumber := scanner.Text()

		updatedUser["driver_license"] = driverLicense
		updatedUser["car_plate_number"] = carPlateNumber
	}

	createOrUpdateUser("PUT", userID, updatedUser)
}

func deleteUser(scanner *bufio.Scanner) {
	fmt.Print("Enter the ID of the user to be deleted: ")
	scanner.Scan()
	userID := scanner.Text()

	// Check if the user exists
	if !userExists(userID) {
		fmt.Println("Error - User does not exist")
		return
	}

	// Retrieve user information from the server
	userResp, err := http.Get(baseURL + "/users/" + userID)
	if err != nil {
		fmt.Println("Error retrieving user information:", err)
		return
	}
	defer userResp.Body.Close()

	if userResp.StatusCode != http.StatusOK {
		fmt.Println("Error retrieving user information:", userResp.Status)
		return
	}

	var user User
	if err := json.NewDecoder(userResp.Body).Decode(&user); err != nil {
		fmt.Println("Error decoding user information:", err)
		return
	}

	// Check if the account has been active for at least 1 year
	oneYearAgo := time.Now().AddDate(-1, 0, 0)
	if user.CreatedAt.After(oneYearAgo) {
		fmt.Println("Error - Account cannot be deleted before 1 year")
		return
	}

	deleteUserByID(userID)
}

func listAllTrips() {
	getData("trips")
}

func createNewTrip(scanner *bufio.Scanner) {
	fmt.Print("Enter the ID of the trip to be created: ")
	scanner.Scan()
	tripID := scanner.Text()

	// Check if the trip already exists
	if tripExists(tripID) {
		fmt.Println("Error - Trip already exists")
		return
	}

	fmt.Print("Enter the ID of the car owner: ")
	scanner.Scan()
	carOwnerID := scanner.Text()

	// Check if the car owner exists
	if !userExists(carOwnerID) {
		fmt.Println("Error - Car owner does not exist")
		return
	}

	// Retrieve car owner information from the server
	carOwnerResp, err := http.Get(baseURL + "/users/" + carOwnerID)
	if err != nil {
		fmt.Println("Error retrieving car owner information:", err)
		return
	}
	defer carOwnerResp.Body.Close()

	if carOwnerResp.StatusCode != http.StatusOK {
		fmt.Println("Error retrieving car owner information:", carOwnerResp.Status)
		return
	}

	var carOwner User
	if err := json.NewDecoder(carOwnerResp.Body).Decode(&carOwner); err != nil {
		fmt.Println("Error decoding car owner information:", err)
		return
	}

	// Check if the car owner is a car owner
	if !carOwner.IsCarOwner {
		fmt.Println("Error - Only car owners can create trips")
		return
	}

	// Check if the car owner has the required fields if they are a car owner
	if carOwner.IsCarOwner {
		if carOwner.DriverLicense == "" || carOwner.CarPlateNumber == "" {
			fmt.Println("Error - Car owner profile incomplete")
			return
		}
	}

	fmt.Print("Enter the pickup location: ")
	scanner.Scan()
	pickupLocation := scanner.Text()

	fmt.Print("Enter the alternative pickup location (optional, press enter to skip): ")
	scanner.Scan()
	altPickupLocation := scanner.Text()

	fmt.Print("Enter the start time (e.g., '15:04' for 3:04 PM): ")
	scanner.Scan()
	startTimeStr := scanner.Text()

	// Convert the user input to a complete time string
	completeTimeStr := time.Now().Format("2006-01-02") + " " + startTimeStr + ":00"

	// Parse the complete time string
	startTime, err := time.Parse("2006-01-02 15:04:05", completeTimeStr)
	if err != nil {
		fmt.Println("Invalid input for start time. Please enter a valid time in '15:04' format.")
		return
	}

	// Validate that the start time is at least 30 minutes in the future
	currentTime := time.Now()
	if startTime.Before(currentTime.Add(30 * time.Minute)) {
		fmt.Println("Error - Trips must be scheduled at least 30 minutes in the future")
		return
	}

	fmt.Print("Enter the destination: ")
	scanner.Scan()
	destination := scanner.Text()

	// Enter the number of available seats
	fmt.Print("Enter the number of total seats in the car: ")
	scanner.Scan()
	totalSeatsStr := scanner.Text()
	totalSeats, err := strconv.Atoi(totalSeatsStr)
	if err != nil {
		fmt.Println("Invalid input for total seats. Please enter a valid number.")
		return
	}

	newTrip := map[string]interface{}{
		"id":                  tripID,
		"car_owner_id":        carOwnerID,
		"pickup_location":     pickupLocation,
		"alt_pickup_location": altPickupLocation,
		"start_time":          startTime,
		"destination":         destination,
		//"available_seats":    availableSeats,
		"total_seats": totalSeats,
		"started":     false, // Set to false initially
	}

	createOrUpdateTrip("POST", tripID, newTrip)
}

func enrollPassenger(scanner *bufio.Scanner) {
	fmt.Print("Enter the ID of the trip to enroll in: ")
	scanner.Scan()
	tripID := scanner.Text()

	// Check if the trip exists
	if !tripExists(tripID) {
		fmt.Println("Error - Trip does not exist")
		return
	}

	fmt.Print("Enter the ID of the user to enroll in the trip: ")
	scanner.Scan()
	userID := scanner.Text()

	enrollData := map[string]interface{}{
		"user_id": userID,
	}

	createOrUpdateTrip("PUT", tripID+"/enroll", enrollData)
}

// startTrip handles the starting of a trip
func startTrip(scanner *bufio.Scanner) {
	fmt.Print("Enter the ID of the trip to start: ")
	scanner.Scan()
	tripID := scanner.Text()

	// Check if the trip exists
	if !tripExists(tripID) {
		fmt.Println("Error - Trip does not exist")
		return
	}

	// Retrieve trip information from the server
	tripResp, err := http.Get(baseURL + "/trips/" + tripID)
	if err != nil {
		fmt.Println("Error retrieving trip information:", err)
		return
	}
	defer tripResp.Body.Close()

	if tripResp.StatusCode != http.StatusOK {
		fmt.Println("Error retrieving trip information:", tripResp.Status)
		return
	}

	var trip Trip
	if err := json.NewDecoder(tripResp.Body).Decode(&trip); err != nil {
		fmt.Println("Error decoding trip information:", err)
		return
	}

	// Display trip information
	fmt.Printf("Trip %s Information:\n", tripID)
	fmt.Printf(" - Started: %v\n", trip.Started)
	fmt.Printf(" - Enrolled Passengers: %v\n", trip.EnrolledPassengers)

	// Check if the user is the car owner
	fmt.Print("Enter your user ID as the car owner: ")
	scanner.Scan()
	carOwnerID := scanner.Text()

	// Check if the user starting the trip is the car owner
	if trip.CarOwnerID != carOwnerID {
		fmt.Println("Error - Only the car owner can start the trip")
		return
	}

	// Check if the trip is already started
	if trip.Started {
		fmt.Println("Error - Trip is already started")
		return
	}

	// Check if the trip has at least one enrolled passenger
	if len(trip.EnrolledPassengers) == 0 {
		fmt.Println("Error - Trip cannot start without any enrolled passengers")
		return
	}

	// Check if the start time is within the allowed window
	if time.Until(trip.StartTime) < -30*time.Minute {
		fmt.Println("Error - Trip cannot be started more than 30 minutes after scheduled time")
		return
	}

	// Mark the trip as started on the server
	startTripOnServer(tripID, carOwnerID)
	fmt.Println("Trip started successfully.")
}

// startTripOnServer marks the trip as started on the server
func startTripOnServer(tripID, carOwnerID string) {
	url := baseURL + "/trips/" + tripID + "/start"
	request, err := http.NewRequest("PUT", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	// Set the car owner ID in the request headers
	request.Header.Set("car-owner-id", carOwnerID)

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		fmt.Println("Error executing request:", err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		fmt.Println("Error starting trip:", response.Status)
		return
	}

	fmt.Println("Trip started successfully on the server.")
}

func userExists(userID string) bool {
	response, err := http.Get(baseURL + "/users/" + userID)
	if err != nil {
		fmt.Println("Error checking if user exists:", err)
		return false
	}
	defer response.Body.Close()

	return response.StatusCode == http.StatusOK
}

func tripExists(tripID string) bool {
	response, err := http.Get(baseURL + "/trips/" + tripID)
	if err != nil {
		fmt.Println("Error checking if trip exists:", err)
		return false
	}
	defer response.Body.Close()

	return response.StatusCode == http.StatusOK
}

func createOrUpdateUser(method, userID string, user map[string]interface{}) {
	jsonBody, err := json.Marshal(user)
	if err != nil {
		fmt.Println("Error encoding user JSON:", err)
		return
	}

	url := baseURL + "/users/" + userID
	request, err := http.NewRequest(method, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		fmt.Println("Error executing request:", err)
		return
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}

	fmt.Println(string(body))
}

func createOrUpdateTrip(method, tripID string, trip map[string]interface{}) {
	jsonBody, err := json.Marshal(trip)
	if err != nil {
		fmt.Println("Error encoding trip JSON:", err)
		return
	}

	url := baseURL + "/trips/" + tripID
	request, err := http.NewRequest(method, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		fmt.Println("Error executing request:", err)
		return
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}

	fmt.Println(string(body))
}

func getData(endpoint string) {
	url := baseURL + "/" + endpoint

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	fmt.Println("Response:", string(body))
}

func deleteUserByID(userID string) {
	url := baseURL + "/users/" + userID
	request, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		fmt.Println("Error executing request:", err)
		return
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}

	fmt.Println(string(body))
}

func cancelTrip(scanner *bufio.Scanner) {
	fmt.Print("Enter the ID of the trip to cancel: ")
	scanner.Scan()
	tripID := scanner.Text()

	// Check if the trip exists
	if !tripExists(tripID) {
		fmt.Println("Error - Trip does not exist")
		return
	}

	// Retrieve trip information from the server
	tripResp, err := http.Get(baseURL + "/trips/" + tripID)
	if err != nil {
		fmt.Println("Error retrieving trip information:", err)
		return
	}
	defer tripResp.Body.Close()

	if tripResp.StatusCode != http.StatusOK {
		fmt.Println("Error retrieving trip information:", tripResp.Status)
		return
	}

	var trip Trip
	if err := json.NewDecoder(tripResp.Body).Decode(&trip); err != nil {
		fmt.Println("Error decoding trip information:", err)
		return
	}

	// Check if the trip is already started
	if trip.Started {
		fmt.Println("Error - Trip is already started and cannot be canceled")
		return
	}

	// Check if the trip is within the cancellation window
	if time.Until(trip.StartTime) < -30*time.Minute {
		fmt.Println("Error - Trip cannot be canceled more than 30 minutes after scheduled time")
		return
	}

	// Perform the trip cancellation
	createOrUpdateTrip("DELETE", tripID, nil)
}

// listTripStatus prints out the status of the trip, including whether it has started
func listTripStatus(scanner *bufio.Scanner) {
	fmt.Print("Enter the ID of the trip to check status: ")
	scanner.Scan()
	tripID := scanner.Text()

	// Check if the trip exists
	if !tripExists(tripID) {
		fmt.Println("Error - Trip does not exist")
		return
	}

	// Retrieve trip information from the server
	tripResp, err := http.Get(baseURL + "/trips/" + tripID)
	if err != nil {
		fmt.Println("Error retrieving trip information:", err)
		return
	}
	defer tripResp.Body.Close()

	if tripResp.StatusCode != http.StatusOK {
		fmt.Println("Error retrieving trip information:", tripResp.Status)
		return
	}

	var trip Trip
	if err := json.NewDecoder(tripResp.Body).Decode(&trip); err != nil {
		fmt.Println("Error decoding trip information:", err)
		return
	}

	// Display updated trip status
	fmt.Printf("Trip %s Status:\n", tripID)
	fmt.Printf(" - Started: %v\n", trip.Started) // Updated trip status from the server
	fmt.Printf(" - Enrolled Passengers: %v\n", trip.EnrolledPassengers)
}
