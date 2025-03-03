package main

import (
	"bufio"
	"log"
	"net"
	"sync"
)

var wg sync.WaitGroup

func main() {
	// Définition du listener
	listener, err := net.Listen("tcp", ":1563")
	if err != nil {
		log.Println("listen error:", err)
		return
	}
	defer listener.Close()

	// Gestion de la première connexion
	connect1, err := listener.Accept()
	if err != nil {
		log.Println("accept error for client 1:", err)
		return
	}
	log.Println("Client 1 connecté")
	writer1 := bufio.NewWriter(connect1)
	_, err = writer1.WriteString("{\"id\": 1,\"gameState\": 1,\"oppColor\": -1,\"colorRow\": -1,\"colorLine\":-1, \"choosenLine\": -1,\"hasPlayed\": false}\n")
	if err != nil {
		log.Println("error writing for client 1:", err)
		return
	}

	// Gestion de la deuxième connexion
	connect2, err := listener.Accept()
	if err != nil {
		log.Println("accept error for client 2:", err)
		return
	}
	log.Println("Client 2 connecté")
	writer2 := bufio.NewWriter(connect2)
	_, err = writer2.WriteString("{\"id\": 2,\"gameState\": 1,\"oppColor\": -1,\"colorRow\": -1,\"colorLine\":-1, \"choosenLine\": -1,\"hasPlayed\": false}\n")
	if err != nil {
		log.Println("error writing for client 2:", err)
		return
	}
	writer2.Flush()
	writer1.Flush() // Envoie les données une fois le deuxième joueur connecté

	// Lancer les goroutines avec les connexions
	wg.Add(2)
	go listenToClient(1, connect1, connect2)
	go listenToClient(2, connect2, connect1)

	wg.Wait()
}

func listenToClient(id int, conn net.Conn, targetConn net.Conn) {
	defer conn.Close() // Fermer la connexion source à la fin
	defer wg.Done()    // Signaler la fin de la goroutine au WaitGroup

	reader := bufio.NewReader(conn)

	for {
		// Lire le message depuis la connexion
		message, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("Erreur de lecture depuis la connexion (Client %d): %v\n", id, err)
			return
		}

		// Rediriger vers la connexion cible
		writer := bufio.NewWriter(targetConn)
		_, err = writer.WriteString(message)
		if err != nil {
			log.Printf("Erreur lors de la redirection des données du client %d : %v\n", id, err)
			return
		}
		err = writer.Flush()
		if err != nil {
			log.Printf("Erreur lors de l'envoi des données au client cible : %v\n", err)
			return
		}

		log.Printf("Message du client %d redirigé : %s", id, message)
	}
}
