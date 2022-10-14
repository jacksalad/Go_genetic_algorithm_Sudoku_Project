package main

import "fmt"

// 棋盘数组
type Board [9][9]int

// 显示状态
func (this Board) Show() {
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			fmt.Print(this[i][j])
			if j != 8 {
				fmt.Print(" ")
			}
		}
		fmt.Println()
	}
}
