make the MYSQL database first:
CREATE DATABASE players;
USE players;

CREATE TABLE userdata (
  userID INT NOT NULL AUTO_INCREMENT,
  username VARCHAR(25) UNIQUE NOT NULL,
  PRIMARY KEY (userID)
);

CREATE TABLE userscores (
  scoreID INT NOT NULL AUTO_INCREMENT UNIQUE,
  userID INT NOT NULL,
  score INT NOT NULL,
  scoreDate TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP(),
  PRIMARY KEY (scoreID),
  FOREIGN KEY (userID) REFERENCES userdata(userID)
);
change the env file based on your database info

go build the server then run the server.exe
then do the same in the client folder
