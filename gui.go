package main

import (
	"image/color"

	fyne "fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
)

const N int = 9

// 主界面
type UI struct {
	model    *GA
	myApp    fyne.App
	myWindow fyne.Window
	canGrid  [][]fyne.CanvasObject
	gridSize float32
}

// 添加网格标签
func genCanGrid(vis Board) [][]fyne.CanvasObject {
	ans := make([][]fyne.CanvasObject, N+2)
	for i := 0; i < N; i++ {
		ans[i] = make([]fyne.CanvasObject, N+2)
	}
	for i := 0; i < N; i++ {
		for j := 0; j < N; j++ {
			num := string(vis[i][j] + '0')
			if num == "0" {
				ans[i][j] = canvas.NewText(num, color.Black)
			} else {
				ans[i][j] = canvas.NewText(num, color.RGBA{0, 0, 255, 255})
			}
			ans[i][j].(*canvas.Text).Alignment = 1
		}
	}
	return ans
}

// 生成
func newUI(model *GA) UI {
	gridSize := float32(10)
	myApp := app.New()
	myWindow := myApp.NewWindow("Sudoku")
	myWindow.CenterOnScreen()
	CanGrid := genCanGrid(model.vis)
	ctnSlice := make([]fyne.CanvasObject, N*N)
	k := 0
	for i := 0; i < N; i++ {
		for j := 0; j < N; j++ {
			ctnSlice[k] = CanGrid[i][j]
			k++
		}
	}
	grid := container.NewAdaptiveGrid(N, ctnSlice...)
	myWindow.SetContent(grid)
	myWindow.Resize(fyne.NewSize(340, 340))
	return UI{model, myApp, myWindow, CanGrid, gridSize}
}

// 刷新
func (this UI) flush() {
	grid := this.model.bestIdv.grid
	for i := 0; i < N; i++ {
		for j := 0; j < N; j++ {
			this.canGrid[i][j].(*canvas.Text).Text = string(grid[i][j] + '0')
			this.canGrid[i][j].Refresh()
		}
	}
}
