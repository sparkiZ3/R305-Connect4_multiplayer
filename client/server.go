package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"time"
)

var connection net.Conn

// cette fonction permet de se connecter au serveur. Une fois la connexion établie, elle lance un listener via une goroutine.
func connect() bool {
	if !isConnected {
		conn, err := net.Dial("tcp", adress+":"+port)
		if err != nil {
			log.Println("Dial error:", err)
			time.Sleep(1 * time.Second)
			return false
		}
		log.Println("Je suis connecté")
		connection = conn
		go listener(conn)
		isConnected = true

	} else {
		return gameStarted
	}
	return false
}

func listener(conn net.Conn) {
	scanner := bufio.NewScanner(conn)
	log.Println("démarrage du listener")
	var datas serverData
	for scanner.Scan() {
		fmt.Println("Message du serveur :", scanner.Text())
		err := json.Unmarshal([]byte(scanner.Text()), &datas) // formatage du JSON recu par le serveur
		if err != nil {
			fmt.Println("Erreur lors du décodage:", err)
			return
		}
		fmt.Println(datas.GameState)
		switch datas.GameState {
		case 1:
			id = datas.Id
			gameStarted = true
		case 2:
			colp2 = datas.ColorRow
			linep2 = datas.ColorLine
			serverHasChosenColor = datas.HasChosenColor
			oppColor = datas.OppColor
		case 3:
			choosenLine = datas.ChosenLine
			oppHasPlayed = datas.HasPlayed
		}
	}
	if err := scanner.Err(); err != nil {
		log.Println("Erreur de lecture :", err)
	}
	conn.Close()
}

func (g *game) sendData() {
	toSend := serverData{Id: id, GameState: g.gameState, OppColor: g.p1Color, ColorRow: col, ColorLine: line, HasChosenColor: hasChosenColor, ChosenLine: g.tokenPosition, HasPlayed: currentPlayerPlayed}
	jsonData, err := json.Marshal(toSend)
	if err != nil {
		fmt.Println("Erreur d'encodage des données", err)
		return
	}

	writer := bufio.NewWriter(connection)

	_, err = writer.WriteString(string(jsonData) + "\n")
	if err != nil {
		log.Println("erreur d'ecriture vers le serveur:", err)
		return
	}
	writer.Flush()
	fmt.Println("data send succesfully")

}
