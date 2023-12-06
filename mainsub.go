package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

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

var (
	users = map[string]User{}
	trips = map[string]Trip{}
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/api/v1/users/{id}", getUser).Methods("GET", "DELETE")
	r.HandleFunc("/api/v1/users", getAllUsers).Methods("GET")
	r.HandleFunc("/api/v1/users/{id}", createOrUpdateUser).Methods("POST", "PUT")

	r.HandleFunc("/api/v1/trips/{id}", getTrip).Methods("GET", "DELETE")
	r.HandleFunc("/api/v1/trips", getAllTrips).Methods("GET")
	r.HandleFunc("/api/v1/trips/{id}", createOrUpdateTrip).Methods("POST", "PUT")
	// Add a new route for enrolling passengers
	r.HandleFunc("/api/v1/trips/{id}/enroll", enrollPassenger).Methods("PUT")
	r.HandleFunc("/api/v1/trips/{id}/start", startTrip).Methods("PUT")

	fmt.Println("Starting car-pooling server on port 8222")
	http.ListenAndServe(":8222", r)
}

// getUser handles GET and DELETE requests for a specific user
func getUser(w http.ResponseWriter, r *http.Request) {
	userID := mux.Vars(r)["id"]

	if user, ok := users[userID]; ok {
		if r.Method == "GET" {
			json.NewEncoder(w).Encode(user)
		} else if r.Method == "DELETE" {
			delete(users, userID)
			fmt.Fprintf(w, "User %s deleted", userID)
		}
	} else {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "Invalid user ID")
	}
}

// getAllUsers handles GET requests to retrieve all users
func getAllUsers(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(users)
}

// createOrUpdateUser handles POST and PUT requests to create or update a user
func createOrUpdateUser(w http.ResponseWriter, r *http.Request) {
	userID := mux.Vars(r)["id"]
	var user User

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&user); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Invalid request payload")
		return
	}

	// Set the creation time if the user is being created
	if r.Method == "POST" {
		user.CreatedAt = time.Now()
	}

	// Check if the user is also a car owner
	if user.IsCarOwner {
		if user.DriverLicense == "" || user.CarPlateNumber == "" {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "Driver's license and car plate number are required for car owners")
			return
		}
	}

	users[userID] = user
	w.WriteHeader(http.StatusAccepted)
	fmt.Fprintf(w, "User %s %s successfully", r.Method, userID)
}

// getTrip handles GET and DELETE requests for a specific trip
func getTrip(w http.ResponseWriter, r *http.Request) {
	tripID := mux.Vars(r)["id"]

	if trip, ok := trips[tripID]; ok {
		if r.Method == "GET" {
			json.NewEncoder(w).Encode(trip)
		} else if r.Method == "DELETE" {
			delete(trips, tripID)
			fmt.Fprintf(w, "Trip %s deleted", tripID)
		}
	} else {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "Invalid trip ID")
	}
}

// getAllTrips handles GET requests to retrieve all trips
func getAllTrips(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(trips)
}

// createOrUpdateTrip handles POST and PUT requests to create or update a trip
func createOrUpdateTrip(w http.ResponseWriter, r *http.Request) {
	tripID := mux.Vars(r)["id"]
	var trip Trip

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&trip); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Invalid request payload")
		return
	}

	// Check if the car owner exists
	carOwner, carOwnerExists := users[trip.CarOwnerID]
	if !carOwnerExists {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "Error - Car owner does not exist")
		return
	}

	// Check if the car owner is a car owner
	if !carOwner.IsCarOwner {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Error - Only car owners can create trips")
		return
	}

	// Check if the car owner has the required fields if they are a car owner
	if carOwner.IsCarOwner {
		if carOwner.DriverLicense == "" || carOwner.CarPlateNumber == "" {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "Error - Car owner profile incomplete")
			return
		}
	}
// Check if the start time is at least 30 minutes in the future

	if time.Until(trip.StartTime) < 30*time.Minute {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Error - Trips must be scheduled at least 30 minutes in the future")
		return
	}

	// Update the EnrolledPassengers in the Trips map
	if existingTrip, ok := trips[tripID]; ok {
		trip.EnrolledPassengers = existingTrip.EnrolledPassengers
	}

	// Set available seats to the difference between total seats and enrolled passengers
	trip.AvailableSeats = trip.TotalSeats - len(trip.EnrolledPassengers)
	if trip.AvailableSeats < 0 {
		trip.AvailableSeats = 0
	}

	trips[tripID] = trip
	w.WriteHeader(http.StatusAccepted)
	fmt.Fprintf(w, "Trip %s %s successfully", r.Method, tripID)
}

// enrollPassenger handles the enrollment of passengers in a trip
func enrollPassenger(w http.ResponseWriter, r *http.Request) {
    tripID := mux.Vars(r)["id"]

    var enrollmentData map[string]string
    decoder := json.NewDecoder(r.Body)
    if err := decoder.Decode(&enrollmentData); err != nil {
        w.WriteHeader(http.StatusBadRequest)
        fmt.Fprint(w, "Invalid request payload")
        return
    }

    userID, exists := enrollmentData["user_id"]
    if !exists {
        w.WriteHeader(http.StatusBadRequest)
        fmt.Fprint(w, "Error - User ID is required in the request payload")
        return
    }

    if trip, ok := trips[tripID]; ok {
        // Check if the user already enrolled
        for _, passengerID := range trip.EnrolledPassengers {
            if passengerID == userID {
                w.WriteHeader(http.StatusBadRequest)
                fmt.Fprint(w, "Error - User already enrolled in this trip")
                return
            }
        }

        // Update the EnrolledPassengers field
        trip.EnrolledPassengers = append(trip.EnrolledPassengers, userID)

        // Update the trip in the trips map
        trips[tripID] = trip

        w.WriteHeader(http.StatusAccepted)
        fmt.Fprintf(w, "User %s enrolled in trip %s successfully", userID, tripID)
    } else {
        w.WriteHeader(http.StatusNotFound)
        fmt.Fprint(w, "Invalid trip ID")
    }
}


// startTrip handles the starting of a trip
func startTrip(w http.ResponseWriter, r *http.Request) {
	tripID := mux.Vars(r)["id"]

	// Retrieve the car owner ID from the request header
	carOwnerID := r.Header.Get("car-owner-id")

	// Retrieve the trip based on the tripID
	trip, ok := trips[tripID]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "Invalid trip ID")
		return
	}

	// Check if the user starting the trip is the car owner
	if trip.CarOwnerID != carOwnerID {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, "Error - Only the car owner can start the trip")
		return
	}

	// Check if the trip is already started
	if trip.Started {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Error - Trip is already started")
		return
	}

	// Check if the trip has at least one enrolled passenger
	if len(trip.EnrolledPassengers) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Error - Trip cannot start without any enrolled passengers")
		return
	}

	// Check if the start time is within the allowed window
	if time.Until(trip.StartTime) < 30*time.Minute {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Error - Trip cannot be started more than 30 minutes after scheduled time")
		return
	}

	// Mark the trip as started
	trip.Started = true
	trips[tripID] = trip

	w.WriteHeader(http.StatusAccepted)
	fmt.Fprintf(w, "Trip %s started successfully", tripID)
}


