package main

import (
	//"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// Mise à jour de l'état du jeu en fonction des entrées au clavier.
func (g *game) Update() error {

	g.stateFrame++

	switch g.gameState {
	case titleState:
		if g.titleUpdate() {
			g.gameState++

		}
	case waitingState:
		if g.waitingUpdate() {
			g.gameState++
			if id == 1 { //initialisation des tours des joueurs en fonction de l'id dans waiting state car il est forcement connecté au serveur à ce moment la
				g.turn = p1Turn
			} else {
				g.turn = p2Turn
			}

		}

	case colorSelectState:
		if g.colorSelectUpdate() {
			g.gameState++
		}
	case playState:
		g.tokenPosUpdate()
		var lastXPositionPlayed int
		var lastYPositionPlayed int
		if g.turn == p1Turn {
			lastXPositionPlayed, lastYPositionPlayed = g.p1Update()
		} else {
			lastXPositionPlayed, lastYPositionPlayed = g.p2Update()
		}
		if lastXPositionPlayed >= 0 {
			finished, result := g.checkGameEnd(lastXPositionPlayed, lastYPositionPlayed)
			if finished {
				g.result = result
				g.gameState++
			}
		}
	case resultState:
		if g.resultUpdate() {
			g.reset()
			g.gameState = playState
		}
	}

	return nil
}

// Mise à jour de l'état du jeu à l'écran titre.
func (g *game) titleUpdate() bool {
	g.stateFrame = g.stateFrame % globalBlinkDuration
	return inpututil.IsKeyJustPressed(ebiten.KeyEnter)
}
func (g *game) waitingUpdate() bool {
	//tant que le client n'arrive pas a se connecter au serveur, l'état ne change pas
	return connect()
}

// Mise à jour de l'état du jeu lors de la sélection des couleurs.
func (g *game) colorSelectUpdate() bool {

	col = g.p1Color % globalNumColorCol
	line = g.p1Color / globalNumColorLine
	if !lockKeyboard { // si l'on a choisi une couleur, on "bloque le clavier"
		if inpututil.IsKeyJustPressed(ebiten.KeyRight) {
			if (col+1)%globalNumColorCol != colp2 || line != linep2 {
				col = (col + 1) % globalNumColorCol
				g.sendData()
			}

		}

		if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
			if (col-1+globalNumColorCol)%globalNumColorCol != colp2 || line != linep2 {
				col = (col - 1 + globalNumColorCol) % globalNumColorCol
				g.sendData()
			}

		}

		if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
			if (line+1)%globalNumColorLine != linep2 || col != colp2 {
				line = (line + 1) % globalNumColorLine
				g.sendData()
			}

		}
		if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
			if (line-1+globalNumColorLine)%globalNumColorLine != linep2 || col != colp2 {
				line = (line - 1 + globalNumColorLine) % globalNumColorLine
				g.sendData()

			}
		}
	}

	g.p1Color = line*globalNumColorLine + col

	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {

		if g.p2Color == g.p1Color {
			g.p2Color = (g.p2Color + 1) % globalNumColor
		}
		//permet ici de déselectionner la couleur
		if hasChosenColor {
			lockKeyboard = false
			hasChosenColor = false
			g.sendData()
		} else {
			lockKeyboard = true
			hasChosenColor = true
			g.sendData()
		}
	}
	if serverHasChosenColor && hasChosenColor {
		lockKeyboard = false
		g.p2Color = oppColor
		return true
	}

	return false
}

// Gestion de la position du prochain pion à jouer par le joueur 1.
func (g *game) tokenPosUpdate() {
	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
		g.tokenPosition = (g.tokenPosition - 1 + globalNumTilesX) % globalNumTilesX
		g.sendData()
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyRight) {
		g.tokenPosition = (g.tokenPosition + 1) % globalNumTilesX
		g.sendData()
	}
}

// Gestion du moment où le prochain pion est joué par le joueur 1.
func (g *game) p1Update() (int, int) {
	lastXPositionPlayed := -1
	lastYPositionPlayed := -1
	if inpututil.IsKeyJustPressed(ebiten.KeyDown) || inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		if updated, yPos := g.updateGrid(p1Token, g.tokenPosition); updated {
			g.turn = p2Turn
			currentPlayerPlayed = true
			lastXPositionPlayed = g.tokenPosition
			lastYPositionPlayed = yPos
			g.sendData()
			currentPlayerPlayed = false
		}
	}
	return lastXPositionPlayed, lastYPositionPlayed
}

// Gestion de la position du prochain pion joué par le joueur 2 et
// du moment où ce pion est joué.
func (g *game) p2Update() (int, int) {
	position := choosenLine //valeur recuperé depuis le serveur
	lastXPositionPlayed := -1
	lastYPositionPlayed := -1
	if oppHasPlayed { //si l'autre client a joué
		_, yPos := g.updateGrid(p2Token, position)
		lastXPositionPlayed = position
		lastYPositionPlayed = yPos
		g.turn = p1Turn
		g.sendData()
	}
	return lastXPositionPlayed, lastYPositionPlayed
}

// Mise à jour de l'état du jeu à l'écran des résultats.
func (g game) resultUpdate() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeyEnter)
}

// Mise à jour de la grille de jeu lorsqu'un pion est inséré dans la
// colonne de coordonnée (x) position.
func (g *game) updateGrid(token, position int) (updated bool, yPos int) {
	for y := globalNumTilesY - 1; y >= 0; y-- {
		if g.grid[position][y] == noToken {
			updated = true
			yPos = y
			g.grid[position][y] = token
			return
		}
	}
	return
}

// Vérification de la fin du jeu : est-ce que le dernier joueur qui
// a placé un pion gagne ? est-ce que la grille est remplie sans gagnant
// (égalité) ? ou est-ce que le jeu doit continuer ?
func (g game) checkGameEnd(xPos, yPos int) (finished bool, result int) {

	tokenType := g.grid[xPos][yPos]

	// horizontal
	count := 0
	for x := xPos; x < globalNumTilesX && g.grid[x][yPos] == tokenType; x++ {
		count++
	}
	for x := xPos - 1; x >= 0 && g.grid[x][yPos] == tokenType; x-- {
		count++
	}

	if count >= 4 {
		if tokenType == p1Token {
			return true, p1wins
		}
		return true, p2wins
	}

	// vertical
	count = 0
	for y := yPos; y < globalNumTilesY && g.grid[xPos][y] == tokenType; y++ {
		count++
	}

	if count >= 4 {
		if tokenType == p1Token {
			return true, p1wins
		}
		return true, p2wins
	}

	// diag haut gauche/bas droit
	count = 0
	for x, y := xPos, yPos; x < globalNumTilesX && y < globalNumTilesY && g.grid[x][y] == tokenType; x, y = x+1, y+1 {
		count++
	}

	for x, y := xPos-1, yPos-1; x >= 0 && y >= 0 && g.grid[x][y] == tokenType; x, y = x-1, y-1 {
		count++
	}

	if count >= 4 {
		if tokenType == p1Token {
			return true, p1wins
		}
		return true, p2wins
	}

	// diag haut droit/bas gauche
	count = 0
	for x, y := xPos, yPos; x >= 0 && y < globalNumTilesY && g.grid[x][y] == tokenType; x, y = x-1, y+1 {
		count++
	}

	for x, y := xPos+1, yPos-1; x < globalNumTilesX && y >= 0 && g.grid[x][y] == tokenType; x, y = x+1, y-1 {
		count++
	}

	if count >= 4 {
		if tokenType == p1Token {
			return true, p1wins
		}
		return true, p2wins
	}

	// egalité ?
	if yPos == 0 {
		for x := 0; x < globalNumTilesX; x++ {
			if g.grid[x][0] == noToken {
				return
			}
		}
		return true, equality
	}

	return
}
