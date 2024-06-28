package main

import (
	"fmt"
	"log"
	"math/rand/v2"
	"os"
	"os/exec"
	"time"

	"github.com/eiannone/keyboard"
	"golang.org/x/crypto/ssh/terminal"
)

type Pos struct {
	x    int
	y    int
	next *Pos
	prev *Pos
}

type Board struct {
	width     int
	height    int
	apple_pos Pos
	running   bool
}

type Snake struct {
	head   *Pos
	tail   *Pos
	grow   bool
	dir    int //
	length int
}

var board Board
var snake Snake

func clear() {
	cmd := exec.Command("clear")

	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func is_snake_pos(row int, col int) bool {
	t := snake.head.next

	for {
		if t == nil {
			break
		}

		if row == t.y && col == t.x {
			return true
		}

		t = t.next
	}

	return false
}

func printfBoard() {
	for i := 0; i < board.height-1; i++ {
		for j := 0; j < board.width; j++ {

			if i == board.apple_pos.y && j == board.apple_pos.x {
				fmt.Print("A")
				continue
			}

			if i == snake.head.y && j == snake.head.x {
				fmt.Print("S")
				continue
			}

			if is_snake_pos(i, j) {
				fmt.Print("s")
				continue
			}

			if (i == 0 && j == 0) || (i == 0 && j == board.width-1) || (i == board.height-2 && j == 0) || (i == board.height-2 && j == board.width-1) {
				fmt.Print("+")
				continue
			}

			if j == 0 || j == board.width-1 {
				fmt.Print("|")
				continue
			}

			if i == 0 || i == board.height-2 {
				fmt.Print("-")
				continue
			}

			fmt.Print(" ")
		}

		if i != board.height-2 {
			fmt.Print("\n")

		}
	}

	fmt.Printf("Snake length: %d", snake.length)
	t := snake.head

	for {
		if t == nil {
			break
		}

		fmt.Printf("(%d, %d) ->", t.x, t.y)
		t = t.next
	}
}

func gen_apple_pos() {
	board.apple_pos.x = rand.IntN(board.width-2) + 1
	board.apple_pos.y = rand.IntN(board.height-3) + 1
}

func keys() {
	keysEvent, err := keyboard.GetKeys(10)

	if err != nil {
		log.Fatal(err)
		return
	}

	defer func() {
		_ = keyboard.Close()
	}()

	for board.running {
		event := <-keysEvent
		if event.Err != nil {
			log.Fatal(event.Err)
			return
		}

		if event.Key == keyboard.KeyEsc {
			board.running = false
			break
		}

		if event.Rune == 'w' {
			snake.dir = 1
			// break
		}

		if event.Rune == 'd' {
			snake.dir = 2
			// break
		}

		if event.Rune == 's' {
			snake.dir = 3
			// break
		}

		if event.Rune == 'a' {
			snake.dir = 4
			// break
		}
	}
}

func move_snake() {

	if !snake.grow {
		t := snake.tail
		for {
			if t == snake.head {
				break
			}

			t.x = t.prev.x
			t.y = t.prev.y

			t = t.prev
		}
	} else {

		snake.grow = false
		nHead := Pos{
			x:    snake.head.x,
			y:    snake.head.y,
			next: snake.head,
			prev: nil,
		}

		snake.head.prev = &nHead
		snake.head = &nHead
	}

	switch snake.dir {
	case 1: // UP
		snake.head.y--
		break
	case 2: // RIGHT
		snake.head.x++
		break
	case 3: // DOWN
		snake.head.y++
		break
	case 4: // LEFT
		snake.head.x--
		break
	default:
		snake.dir = 1
	}
}

func check_apple() {
	if snake.head.x == board.apple_pos.x && snake.head.y == board.apple_pos.y {
		snake.length++
		snake.grow = true
		gen_apple_pos()
	}
}

func check_collision() bool {
	t := snake.head.next

	// touching self
	for {
		if t == nil {
			break
		}

		if t.x == snake.head.x && t.y == snake.head.y {
			return true
		}

		t = t.next
	}

	// touching side walls
	if snake.head.x == 0 || snake.head.x == board.width-1 {
		return true
	}

	// touching top or bottom
	if snake.head.y == 0 || snake.head.y == board.height-2 {
		return true
	}

	return false
}

func main() {
	// var err error
	board.width, board.height, _ = terminal.GetSize(0)

	board.running = true

	snake.head = &Pos{
		x:    10,
		y:    10,
		next: nil,
		prev: nil,
	}

	snake.tail = snake.head

	snake.dir = 2
	snake.length = 1
	snake.grow = false
	go keys()
	gen_apple_pos()

	for board.running {
		clear()
		move_snake()
		if check_collision() {
			board.running = false
			break
		}
		check_apple()
		printfBoard()
		time.Sleep(400 * time.Millisecond)
	}

	fmt.Println("Game over\n")

}
