package main

import (
	"fmt"
	"github.com/jroimartin/gocui"
	"sort"
	"strconv"
//	"log"
)

func run_ui() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		panic(err)
	}

	defer g.Close()

	g.Highlight = true
	g.Cursor = true
	g.SelFgColor = gocui.ColorGreen
	//g.SelBgColor = gocui.ColorRed

	g.SetManagerFunc(layout)

	keybinds(g)

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		panic(err)
	}
}

func fill_cols(g *gocui.Gui) {
	rec_len := len(mv.records)
	for i := 0; i < mv.fields_n; i++ {
		v, err := g.View(strconv.Itoa(i))
		if err != nil {
			panic(err)
		}
		v.Clear()
		for j := 0; j < rec_len; j++ {
			fmt.Fprintln(v, mv.records[j][i])
		}
	}
}

func sum(arr []int, n int) int{
  s:=0
  for i:=0;i<n;i++{
    s += arr[i]
  }
  return s
}

func get_col_width(maxX int) []int{
	mx_fl := float64(maxX)
  col_widths := make([]int, mv.fields_n)
  copy(col_widths, mv.max_widths)

  s := sum(col_widths, mv.fields_n)
  if(s > maxX){
    ind := 0
    for (sum(col_widths, mv.fields_n) > maxX){
      //TODO reduce starting with largest
      col_widths[ind] = int(mx_fl * mv.width_ratios[ind])
      ind ++
    }
  } else {
    for i:=0;i<mv.fields_n;i++ {
      col_widths[i] = int(mx_fl * mv.width_ratios[i])
    }
  }

  return col_widths
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	lx := 0
	helper_h := 3

  col_widths := get_col_width(maxX)

	for i := 0; i < mv.fields_n; i++ {
		if v, err := g.SetView(strconv.Itoa(i), lx, 0, lx + col_widths[i], maxY-1-helper_h); err != nil {
			if err != gocui.ErrUnknownView {
				return err
			}

			v.Frame = false //no border
			v.Editable = false

			if mv.has_header {
				v.Title = mv.header[i]
				v.Frame = true
			}
		}
		lx += col_widths[i]
	}

	fill_cols(g)

	//setup view on first run
	if g.CurrentView() == nil {
		if _, err := g.SetCurrentView(strconv.Itoa(0)); err != nil {
			return err
		}
	}

	//helper height
	if v, err := g.SetView("helper", 0, maxY-helper_h, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Editable = false

		fmt.Fprintln(v, "Ctrl-C: quit |", "Ctrl-S: sort | Ctrl-F/B: scroll horiz | <arrow keys/PgUp/PgDn>: move |")
	}

	return nil
}

func call_move_control(foo func(*gocui.Gui, *gocui.View, int) error, dir int) func(*gocui.Gui, *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		return foo(g, v, dir)
	}
}

func keybinds(g *gocui.Gui) {
	err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit)
	if err != nil {
		panic(err)
	}

	err = g.SetKeybinding("", gocui.KeyArrowDown, gocui.ModNone, call_move_control(nextLine, 1))
	if err != nil {
		panic(err)
	}

	err = g.SetKeybinding("", gocui.KeyArrowUp, gocui.ModNone, call_move_control(nextLine, -1))
	if err != nil {
		panic(err)
	}

	// left right columns
	err = g.SetKeybinding("", gocui.KeyArrowRight, gocui.ModNone, call_move_control(nextView, 1))
	if err != nil {
		panic(err)
	}

	err = g.SetKeybinding("", gocui.KeyArrowLeft, gocui.ModNone, call_move_control(nextView, -1))
	if err != nil {
		panic(err)
	}

	// horiz scroll
	err = g.SetKeybinding("", gocui.KeyCtrlF, gocui.ModNone, call_move_control(scrollHoriz, 1))
	if err != nil {
		panic(err)
	}

	err = g.SetKeybinding("", gocui.KeyCtrlB, gocui.ModNone, call_move_control(scrollHoriz, -1))
	if err != nil {
		panic(err)
	}

	//paging
	err = g.SetKeybinding("", gocui.KeyPgdn, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			_, h := v.Size()
			return nextLine(g, v, h)
		})
	if err != nil {
		panic(err)
	}

	err = g.SetKeybinding("", gocui.KeyPgup, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			_, h := v.Size()
			return nextLine(g, v, -h)
		})
	if err != nil {
		panic(err)
	}

	//sort
	err = g.SetKeybinding("", gocui.KeyCtrlS, gocui.ModNone, sortCol)
	if err != nil {
		panic(err)
	}

}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

//Move to next line dir -1,1 for u,d does not wrap around
func nextLine(g *gocui.Gui, v *gocui.View, dir int) error {
	x, y := v.Cursor()

	//moving lines
	y += dir
	_, pageH := v.Size()
	_, oy := v.Origin()

	if y >= pageH { //scroll down
		scrollViews(g, oy+pageH)

		if !isEmptyLine(v, y-pageH) {
			v.SetCursor(x, y-pageH)
		} else {
			//set cursor on last occupied line
			v.SetCursor(x, getLastLine(v))
		}
	} else if y < 0 { //scroll up
		if oy == 0 {
			return nil
		} //can't scroll up at top
		ny := oy - pageH

		scrollViews(g, ny)
		v.SetCursor(x, y+pageH)
	} else {
		//normal movement
		if y <= getLastLine(v) { //don't scroll past last line
			v.SetCursor(x, y)
		}
	}

	return nil
}

func isEmptyLine(v *gocui.View, y int) bool {
	ln, _ := v.Line(y)
	return (len(ln) == 0)
}

func getLastLine(v *gocui.View) int {
	_, oy := v.Origin()
	return mv.records_len - oy
}

// scrolls all views up or down won't over or undershoot records
func scrollViews(g *gocui.Gui, ny int) {
	if ny > mv.records_len {
		return
	} // don't overshoot records
	for i := 0; i < mv.fields_n; i++ {
		v, err := g.View(strconv.Itoa(i))
		if err != nil {
			panic(err)
		}
		v.SetOrigin(0, ny)
	}
}

// scrolls view horizontally
func scrollHoriz(g *gocui.Gui, v *gocui.View, dir int) error {
	_, y := v.Cursor()
	ox, oy := v.Origin()
	sx, _ := v.Size()
	vln, _ := v.Line(y)
	line_width := len(vln)

	nx := ox + dir

	//out of bounds
	if nx+sx > line_width || nx < 0 {
		return nil
	}

	v.SetOrigin(nx, oy)
	return nil
}

// Move to next view dir is -1, 1 for l,r wraps around.
func nextView(g *gocui.Gui, v *gocui.View, dir int) error {
	//store cursor position so we are on correct line
	_, y := v.Cursor()

	//moving columns
	n, err := strconv.Atoi(v.Name())
	if err != nil {
		return err
	}
	n += dir
	if n < 0 {
		n = mv.fields_n - 1
	}
	n %= mv.fields_n
	new_v, err := g.SetCurrentView(strconv.Itoa(n))
	if err != nil {
		return err
	}

	//restore "correct" position
	new_v.SetCursor(0, y)
	return nil
}

// Sort column in desc, asc cycle
func sortCol(g *gocui.Gui, v *gocui.View) error {
	ind, err := strconv.Atoi(v.Name())
	if err != nil {
		return err
	}

	sortbyCol := func(i, j int) bool {
		//TODO terrible code need better solution
		if v1, err := strconv.Atoi(mv.records[i][ind]); err == nil {
			v2, err := strconv.Atoi(mv.records[j][ind])
			if err != nil {
				panic(err)
			}
			return v1 < v2
		}
		return mv.records[i][ind] < mv.records[j][ind]
	}

	//switch order on subseq calls
	if sort.SliceIsSorted(mv.records, sortbyCol) { //TODO wasteful keep a variable and switch instead of this
		sort.Slice(mv.records, func(i, j int) bool {
			return sortbyCol(j, i)
		})
	} else {
		sort.Slice(mv.records, sortbyCol)
	}

	fill_cols(g)
	return nil
}
