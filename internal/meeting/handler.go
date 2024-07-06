package meeting

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

// CreateMeetingHandler handles the request to create a new meeting
func CreateMeetingHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse request body
		var meeting Meeting
		err := json.NewDecoder(r.Body).Decode(&meeting)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Implement logic to create a new meeting in the database
		// Example: Insert meeting record into database
		insertQuery := `
			INSERT INTO meetings (host_id, code, start_time, active)
			VALUES ($1, $2, $3, $4)
			RETURNING meeting_id
		`
		var meetingID uint
		err = db.QueryRow(insertQuery, meeting.HostID, meeting.Code, time.Now(), true).Scan(&meetingID)
		if err != nil {
			http.Error(w, "Failed to create meeting", http.StatusInternalServerError)
			return
		}

		// Set HTTP status code and response
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"meeting_id": meetingID,
			"message":    "Meeting created successfully",
		})
	}
}

// GetAllMeetingsHandler handles the request to retrieve all meetings
func GetAllMeetingsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var meetings []Meeting

		query := `
			SELECT meeting_id, host_id, code, start_time, end_time, active
			FROM meetings
		`
		rows, err := db.Query(query)
		if err != nil {
			http.Error(w, "Failed to fetch meetings", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		for rows.Next() {
			var meeting Meeting
			err := rows.Scan(&meeting.MeetingID, &meeting.HostID, &meeting.Code, &meeting.StartTime, &meeting.EndTime, &meeting.Active)
			if err != nil {
				http.Error(w, "Failed to scan meeting row", http.StatusInternalServerError)
				return
			}
			meetings = append(meetings, meeting)
		}
		if err := rows.Err(); err != nil {
			http.Error(w, "Error reading meeting rows", http.StatusInternalServerError)
			return
		}

		// Encode meetings to JSON
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(meetings); err != nil {
			http.Error(w, "Failed to encode meetings to JSON", http.StatusInternalServerError)
			return
		}
	}
}

// GetMeetingByIDHandler handles the request to retrieve a meeting by ID
func GetMeetingByIDHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract meeting ID from request parameters
		vars := mux.Vars(r)
		meetingIDStr := vars["id"]
		meetingID, err := strconv.ParseUint(meetingIDStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid meeting ID", http.StatusBadRequest)
			return
		}

		var meeting Meeting

		query := `
			SELECT meeting_id, host_id, code, start_time, end_time, active
			FROM meetings
			WHERE meeting_id = $1
		`
		err = db.QueryRow(query, meetingID).Scan(&meeting.MeetingID, &meeting.HostID, &meeting.Code, &meeting.StartTime, &meeting.EndTime, &meeting.Active)
		if err != nil {
			http.Error(w, "Meeting not found", http.StatusNotFound)
			return
		}

		// Encode meeting to JSON
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(meeting); err != nil {
			http.Error(w, "Failed to encode meeting to JSON", http.StatusInternalServerError)
			return
		}
	}
}

// EndMeetingHandler handles the request to end a meeting
func EndMeetingHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract meeting ID from request parameters
		vars := mux.Vars(r)
		meetingIDStr := vars["id"]
		meetingID, err := strconv.ParseUint(meetingIDStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid meeting ID", http.StatusBadRequest)
			return
		}

		// Implement logic to end the meeting in the database
		// Example: Update end_time and set active to false
		updateQuery := `
			UPDATE meetings
			SET end_time = $1, active = false
			WHERE meeting_id = $2
		`
		_, err = db.Exec(updateQuery, time.Now(), meetingID)
		if err != nil {
			http.Error(w, "Failed to end meeting", http.StatusInternalServerError)
			return
		}

		// Return success response
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Meeting ended successfully",
		})
	}
}
