# projet R305-Connect4_multiplayer

Vous pouvez trouver le sujet de ce projet [ici](./0-presentation.pdf)

Espace gitlab du projet de la ressource R305 - Programmation système

membres du groupe : 
- Noan LISSILLOUR
- Abel GERAULT

## Getting started



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
