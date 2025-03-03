# projet R305-Connect4_multiplayer

Le but de ce projet était de prendre un jeu déjà existant puis de le rendre jouable en multijoueurs via un serveur et un client fonctionnel. Vous pouvez trouver le sujet de ce projet [ici](./0-presentation.pdf) pour une explication plus précise

Espace gitlab du projet de la ressource R305 - Programmation système

membres du groupe : 
- Noan LISSILLOUR
- Abel GERAULT

## Getting started

Pour lancer le projet :
- Lancer le serveur
    1. Aller dans le dossier `server/`
    2. Lancer le serveur via la commande `go run main.go`

- Lancer le jeu
    1. Aller dans `client/`
    2. fait un `go build` afin de construire le projet
    3. Lancer le projet via la commande `./puissancequatre -ip 127.0.0.1`

## fichier client

### format des données stockées
```go
type serverData struct {
	Id             int         `json:"id"`
	GameState      int         `json:"gameState"`
	OppColor       color.Color `json:"oppColor"`
	ColorRow       int         `json:"colorRow"`
	ColorLine      int         `json:"colorLine"`
	HasChosenColor bool        `json:"hasChosenColor"`
	ChosenLine     int         `json:"chosenLine"`
	HasPlayed      bool        `json:"hasPlayed"`
}
```
## fichier server

### format des données reçus et envoyées

```json
{
  "id": int,
  "gameState": int,
  "oppColor": color.Color,
  "colorRow": int,
  "colorLine": int,
  "choosenLine": int,
  "hasPlayed": bool
}
```
