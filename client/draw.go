package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// Affichage des graphismes à l'écran selon l'état actuel du jeu.
func (g *game) Draw(screen *ebiten.Image) {

	screen.Fill(globalBackgroundColor)

	switch g.gameState {
	case titleState:
		g.titleDraw(screen)
	case waitingState:
		g.waitingDraw(screen)
	case colorSelectState:
		g.colorSelectDraw(screen)
	case playState:
		g.playDraw(screen)
	case resultState:
		g.resultDraw(screen)
	}
}

// Affichage des graphismes de l'écran titre.
func (g game) titleDraw(screen *ebiten.Image) {
	text.Draw(screen, "Puissance 4 en réseau", largeFont, 90, 150, globalTextColor)
	text.Draw(screen, "Projet de programmation système", smallFont, 105, 190, globalTextColor)
	text.Draw(screen, "Année 2023-2024", smallFont, 210, 230, globalTextColor)

	if g.stateFrame >= globalBlinkDuration/3 {
		text.Draw(screen, "Appuyez sur entrée", smallFont, 210, 500, globalTextColor)
	}
}
func (g game) waitingDraw(screen *ebiten.Image) {
	text.Draw(screen, "Waiting room", largeFont, 90, 150, globalTextColor)
	text.Draw(screen, "attente d'un joueur", smallFont, 105, 190, globalTextColor)

}

// Affichage des graphismes de l'écran de sélection des couleurs des joueurs.
func (g game) colorSelectDraw(screen *ebiten.Image) {
	text.Draw(screen, "Quelle couleur pour vos pions ?", smallFont, 110, 80, globalTextColor)

	line := 0
	col := 0

	for numColor := 0; numColor < globalNumColor; numColor++ {

		xPos := (globalNumTilesX-globalNumColorCol)/2 + col
		yPos := (globalNumTilesY-globalNumColorLine)/2 + line

		x2Pos := (globalNumTilesX-globalNumColorCol)/2 + colp2
		y2Pos := (globalNumTilesY-globalNumColorLine)/2 + linep2

		if numColor == g.p2Color {
			vector.DrawFilledCircle(screen, float32(globalTileSize/2+x2Pos*globalTileSize), float32(globalTileSize+globalTileSize/2+y2Pos*globalTileSize), globalTileSize/2, p2globalSelectColor, true)
		}

		if numColor == g.p1Color {
			vector.DrawFilledCircle(screen, float32(globalTileSize/2+xPos*globalTileSize), float32(globalTileSize+globalTileSize/2+yPos*globalTileSize), globalTileSize/2, globalSelectColor, true)
		}

		vector.DrawFilledCircle(screen, float32(globalTileSize/2+xPos*globalTileSize), float32(globalTileSize+globalTileSize/2+yPos*globalTileSize), globalTileSize/2-globalCircleMargin, globalTokenColors[numColor], true)

		col++
		if col >= globalNumColorCol {
			col = 0
			line++
		}
	}

	//gestion de l'affichage des status des joueurs
	j1Color := color.NRGBA{R: 239, G: 0, B: 0, A: 255}
	j1status := "Pas pret"
	if hasChosenColor {
		j1status = "Pret"
		j1Color = color.NRGBA{R: 0, G: 239, B: 0, A: 255}
	}
	j2Color := color.NRGBA{R: 239, G: 0, B: 0, A: 255}
	j2status := "Pas pret"
	if serverHasChosenColor {
		j2status = "Pret"
		j2Color = color.NRGBA{R: 0, G: 239, B: 0, A: 255}
	}
	text.Draw(screen, "Joueur 1 : "+j1status, smallFont, 30, globalHeight-100, j1Color)
	text.Draw(screen, "Joueur 2 : "+j2status, smallFont, 30, globalHeight-50, j2Color)
}

// Affichage des graphismes durant le jeu.
func (g game) playDraw(screen *ebiten.Image) {
	g.drawGrid(screen)

	vector.DrawFilledCircle(screen, float32(globalTileSize/2+g.tokenPosition*globalTileSize), float32(globalTileSize/2), globalTileSize/2-globalCircleMargin, globalTokenColors[g.p1Color], true)
	//affichage du pion du joueur 2
	vector.DrawFilledCircle(screen, float32(globalTileSize/2+choosenLine*globalTileSize), float32(globalTileSize/2), globalTileSize/2-globalCircleMargin, globalTokenColors[g.p2Color], true)

	//quand le joueur 1 et 2 ont la meme position de pion on genere un pion avec la couleur des 2 joueurs
	if g.tokenPosition == choosenLine {
		vector.DrawFilledCircle(screen, float32(globalTileSize/2+choosenLine*globalTileSize), float32(globalTileSize/2), globalTileSize/2, globalTokenColors[g.p2Color], true)
		vector.DrawFilledCircle(screen, float32(globalTileSize/2+choosenLine*globalTileSize), float32(globalTileSize/2), globalTileSize/2-globalCircleMargin*2, globalTokenColors[g.p1Color], true)
	}

}

// Affichage des graphismes à l'écran des résultats.
func (g game) resultDraw(screen *ebiten.Image) {
	g.drawGrid(offScreenImage)

	options := &ebiten.DrawImageOptions{}
	options.ColorScale.ScaleAlpha(0.2)
	screen.DrawImage(offScreenImage, options)

	message := "Égalité"
	if g.result == p1wins {
		message = "Gagné !"
	} else if g.result == p2wins {
		message = "Perdu…"
	}
	text.Draw(screen, message, smallFont, 300, 400, globalTextColor)
}

// Affichage de la grille de puissance 4, incluant les pions déjà joués.
func (g game) drawGrid(screen *ebiten.Image) {
	vector.DrawFilledRect(screen, 0, globalTileSize, globalTileSize*globalNumTilesX, globalTileSize*globalNumTilesY, globalGridColor, true)

	for x := 0; x < globalNumTilesX; x++ {
		for y := 0; y < globalNumTilesY; y++ {

			var tileColor color.Color
			switch g.grid[x][y] {
			case p1Token:
				tileColor = globalTokenColors[g.p1Color]
			case p2Token:
				tileColor = globalTokenColors[g.p2Color]
			default:
				tileColor = globalBackgroundColor
			}

			vector.DrawFilledCircle(screen, float32(globalTileSize/2+x*globalTileSize), float32(globalTileSize+globalTileSize/2+y*globalTileSize), globalTileSize/2-globalCircleMargin, tileColor, true)
		}
	}
}
