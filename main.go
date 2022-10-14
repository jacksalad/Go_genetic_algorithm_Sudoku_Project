package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"time"
)

func readData(path string) Board {
	var board Board
	file, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
	}
	reader := bufio.NewReader(file)
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			c, err := reader.ReadByte()
			if err == io.EOF {
				break
			}
			for !('0' <= c && c <= '9') {
				c, err = reader.ReadByte()
				if err == io.EOF {
					break
				}
			}
			board[i][j] = int(c - '0')
		}
	}
	err = file.Close()
	if err != nil {
		fmt.Println(err)
	}
	return board
}

func main() {
	board := readData("data.txt")
	ga := NewGAModel()
	ga.ModelInit(board)
	ui := newUI(&ga)
	go func() {
		time.Sleep(2 * time.Second)
		ga.Train()
		ui.flush()
	}()
	go func() {
		for {
			ui.flush()
			time.Sleep(time.Second)
		}
	}()
	ui.myWindow.ShowAndRun()
}
