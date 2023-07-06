package main

// the db interactions are here
import (
	"database/sql"
	"fmt"
	"strconv"
)

type Store interface {
	CreateScore(score *Score) error
	GetScores() ([]*Score, error)
}

type dbStore struct {
	db *sql.DB
}

// sending the score data to the db
func (store *dbStore) CreateScore(score *Score) error {
	// Retrieve the player ID based on the player name
	playerID, err := store.getPlayerID(score.Name)
	if err != nil {
		// If player not found, create a new player and retrieve the generated player ID
		playerID, err = store.createPlayer(score.Name)
		if err != nil {
			return err
		}
	}
	// storing the player data and preventing sql injections
	stmt, err := store.db.Prepare("INSERT INTO userscores (score, scoreDate, userID) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	// changing the points to int
	points, err := strconv.Atoi(score.Points)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(points, score.Date, playerID)
	return err
}

// getting the id of the user
func (store *dbStore) getPlayerID(playerName string) (int, error) {
	query := "SELECT userID FROM userdata WHERE username = ?"
	var playerID int
	err := store.db.QueryRow(query, playerName).Scan(&playerID)
	if err != nil {
		if err == sql.ErrNoRows {
			// Player not found, return an error with a custom message
			return 0, fmt.Errorf("player not found")
		}
		return 0, err
	}
	return playerID, nil
}

// if the user does not exist create one
func (store *dbStore) createPlayer(playerName string) (int, error) {
	stmt, err := store.db.Prepare("INSERT INTO userdata (username) VALUES (?)")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	_, err = stmt.Exec(playerName)
	if err != nil {
		return 0, err
	}

	return store.getPlayerID(playerName)
}

func (store *dbStore) GetScores() ([]*Score, error) {
	// Query the database for all highscores, and return the highst score but only one score from each person is taken
	rows, err := store.db.Query("SELECT u.username, us.score, us.scoreDate FROM userscores us JOIN userdata u ON u.userID = us.userID WHERE (us.userID, us.score) IN (SELECT us2.userID, MAX(us2.score) FROM userscores us2 GROUP BY us2.userID)")
	// Return in case of an error, and defer the closing of the row structure
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Create the data structure that is returned from the function and given to the html
	// By default, this will be an empty array of scores unless the database has data
	scores := []*Score{}
	for rows.Next() {
		// For each row returned by the table, create a pointer to a score
		score := &Score{}
		// add the info
		if err := rows.Scan(&score.Name, &score.Points, &score.Date); err != nil {
			return nil, err
		}
		//append the result to the returned array, and repeat for the next row untill every score is taken
		scores = append(scores, score)
	}
	return scores, nil
}

// i don"t have a clue
var store Store

func InitStore(s Store) {
	store = s
}
