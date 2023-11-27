package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"time"
)

const (
	width       = 12
	height      = 8
	borderChar  = '+'
	playerChar  = '☺'
	enemyChar   = '☠'
	attackChar  = '*'
	invalidChar = '!'
)

type Character struct {
	x, y, health int
}

type Movable interface {
	move(dx, dy int, game *Game) bool
}

type Attacker interface {
	attack(target Movable)
}

type Player struct {
	Character
}

type Enemy struct {
	Character
}

type Game struct {
	player  *Player
	enemies []*Enemy
}

func NewPlayer(x, y int) *Player {
	return &Player{Character: Character{x: x, y: y, health: 100}}
}

func NewEnemy(x, y int) *Enemy {
	return &Enemy{Character: Character{x: x, y: y, health: 50}}
}

func (c *Character) move(dx, dy int, game *Game) bool {
	newX, newY := c.x+dx, c.y+dy

	if newX >= 0 && newX < width && newY >= 0 && newY < height {
		c.x, c.y = newX, newY
		return true
	}

	return false
}

func (e *Enemy) moveRandomly() {
	newX, newY := e.x+rand.Intn(3)-1, e.y+rand.Intn(3)-1
	if newX >= 0 && newX < width && newY >= 0 && newY < height {
		e.x, e.y = newX, newY
	}
}

func (p *Player) attack(target Movable) {
	enemy, ok := target.(*Enemy)
	if !ok {
		return // Only attack enemies
	}

	damage := rand.Intn(10) + 1
	fmt.Printf("Player attacks, dealing %d damage!\n", damage)
	enemy.health -= damage

	if enemy.health <= 0 {
		fmt.Println("Enemy defeated!")
		enemy.x = -1
	}
}

func (e *Enemy) attack(target Movable) {
	player, ok := target.(*Player)
	if !ok {
		return // Only attack the player
	}

	damage := rand.Intn(10) + 1
	fmt.Printf("Enemy attacks player, dealing %d damage!\n", damage)
	player.health -= damage

	if player.health <= 0 {
		fmt.Println("Player defeated! Game Over!")
		os.Exit(0)
	}
}

func (g *Game) DrawGame() {
	clearScreen()

	// Draw the top border
	for i := 0; i < width+2; i++ {
		fmt.Printf("\033[1;%df%c", i, borderChar)
	}

	for _, enemy := range g.enemies {
		enemy.draw(enemyChar)
	}

	// Draw the player
	g.player.draw(playerChar)

	for i := 2; i < height+2; i++ {
		// Draw the left border
		fmt.Printf("\033[%d;%df%c", i, 1, borderChar)
		// Draw the right border
		fmt.Printf("\033[%d;%df%c", i, width+2, borderChar)
	}

	// Draw the bottom border
	for i := 0; i < width+2; i++ {
		fmt.Printf("\033[%d;%df%c", height+2, i, borderChar)
	}
}

func clearScreen() {
	fmt.Print("\033[H\033[2J")
}

func (c *Character) draw(ch rune) {
	fmt.Printf("\033[%d;%df%c", c.y+2, c.x+1, ch)
}

func (g *Game) MovePlayer(move rune) {
	switch move {
	case 'w':
		g.player.move(0, -1, g)
	case 'a':
		g.player.move(-1, 0, g)
	case 's':
		g.player.move(0, 1, g)
	case 'd':
		g.player.move(1, 0, g)
	case ' ':
		g.player.attackNearbyEnemies(g)
	case 'q':
		fmt.Println("Game Over!")
		os.Exit(0)
	default:
		fmt.Println("Invalid move! Use w/a/s/d to move or space to attack.")
		time.Sleep(1 * time.Second)
	}
}

func (g *Game) MoveEnemies() {
	for _, enemy := range g.enemies {
		enemy.moveRandomly()
		enemy.attack(g.player)
	}
}

func (p *Player) attackNearbyEnemies(game *Game) {
	for _, enemy := range game.enemies {
		if abs(enemy.x-p.x) <= 1 && abs(enemy.y-p.y) <= 1 {
			p.attack(enemy)
		}
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func main() {
	rand.Seed(time.Now().UnixNano())

	player := NewPlayer(width/2, height/2)
	enemies := []*Enemy{
		NewEnemy(rand.Intn(width), rand.Intn(height)),
		NewEnemy(rand.Intn(width), rand.Intn(height)),
	}

	game := Game{
		player:  player,
		enemies: enemies,
	}

	gameLoop(&game)
}

func gameLoop(game *Game) {
	reader := bufio.NewReader(os.Stdin)

	for {
		game.DrawGame()

		fmt.Print("\033[K") // Clear the current line
		fmt.Print("\nEnter a move (w/a/s/d), space to attack, q to quit: ")
		move, _, err := reader.ReadRune()
		if err != nil {
			fmt.Println("Error reading input:", err)
			break
		}

		game.MovePlayer(move)
		game.MoveEnemies()

		time.Sleep(500 * time.Millisecond)
	}
}
