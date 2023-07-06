package main

//here we are starting the game and sending the score to the webapp to show the highscores
// run this in command prompt for a better ui
import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	//use these if you want to save to a txt
	//"io"
	//"os"
	//"bufio"

	snake "projectSnakeClient/game"
)

// using string for points to make sure the server reads it correctly
type playerScore struct {
	Name   string `json:"name"`
	Points string `json:"points"`
	Date   string `json:"date"`
}

// start the game and get the data in the correct format
func startGame() playerScore {
	points, name, time := snake.Start("Test")
	playerData := playerScore{name, strconv.Itoa(points), time}
	return playerData
}

func main() {
	playerData := startGame()
	err := postScore(playerData)
	if err != nil {
		log.Fatal(err)
	}
	//the save to file
	/*
		// Open the file in append mode or create a new file if it doesn't exist
		file, err := os.OpenFile("scores.txt", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		err = writeToFile(file, playerData)
		if err != nil {
			panic(err)
		}

		scanFile(file)
	*/

}

// post the scores to /score
func postScore(score playerScore) error {
	scoreJSON, err := json.Marshal(score)
	if err != nil {
		return err
	}

	// Log the score JSON before making the request to make sure it is sent
	log.Printf("Sending score: %s\n", string(scoreJSON))
	//change the host if you are using another one
	resp, err := http.Post("http://localhost:8089/score", "application/json", bytes.NewBuffer(scoreJSON))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	//just checking if all is well
	if resp.StatusCode != http.StatusOK {
		errMsg := fmt.Sprintf("Failed to post score: %s", resp.Status)
		log.Println(errMsg)
		return errors.New(errMsg)
	}

	return nil
}

//the save to file again
/*
var scores []playerScore

	func scanFile(file *os.File) {
		scanner := bufio.NewScanner(file)
		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
		for scanner.Scan() {
			var score playerScore
			if err := json.Unmarshal(scanner.Bytes(), &score); err != nil {
				log.Println("Failed to parse score:", err)
				continue
			}
			scores = append(scores, score)
		}
	}

	func writeToFile(file *os.File, playerData playerScore) error {
		// Serialize the playerData to JSON
		jsonData, err := json.Marshal(playerData)
		if err != nil {
			return err
		}

		// Seek to the end of the file
		_, err = file.Seek(0, io.SeekEnd)
		if err != nil {
			return err
		}

		// Write the JSON data to the file in append mode
		_, err = file.Write(jsonData)
		if err != nil {
			return err
		}

		return nil

}
*/
