package main

import (
	"fmt"
	"github.com/awesome-gocui/gocui"
	"sort"
	"strconv"
	//"log"
)

func run_ui() {
	g, err := gocui.NewGui(gocui.OutputNormal, true)
	if err != nil {
		panic(err)
	}
	defer g.Close()

	g.Highlight = true
	g.Cursor = true
	g.SelFgColor = gocui.ColorGreen
  //g.FgColor = gocui.ColorBlue
	//g.SelBgColor = gocui.ColorRed

	g.SetManagerFunc(layout)

	keybinds(g)

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		panic(err)
	}
}

func fill_cols(g *gocui.Gui) {
	for i := 0; i < mv.fields_n; i++ {
    // setup ith columns view
		v, err := g.View(strconv.Itoa(i))
		if err != nil {
			panic(err)
		}
		v.Clear()
    // print to ith column
		for j := 0; j < mv.records_len; j++ {
      if mv.records[j][i] == "" {
        fmt.Fprintln(v, " ") //handle issue with empty fields messing up columns
      }else{
        fmt.Fprintln(v, mv.records[j][i])
      }
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


// try to get columns to size in a useful way
func set_col_width(maxX int) []int{
	mx_fl := float64(maxX)
  col_widths := make([]int, mv.fields_n)
  copy(col_widths, mv.max_widths)

  s := sum(col_widths, mv.fields_n)
  //log.Println(maxX, s)
  //fields are too big for screen size
  if(s > maxX){
    //sort them to reduce from largest
    //add index to maintain orig order probably nicer way
    indexable_widths := make([][2]int, mv.fields_n)
    for i:=0; i<mv.fields_n; i++ {
      indexable_widths[i][0] = i
      indexable_widths[i][1] = mv.max_widths[i]
    }

    //sort widths so we can start reducing size with largest
    sort.Slice(indexable_widths, func(i, j int) bool {
      return indexable_widths[i][1] > indexable_widths[j][1]
    })

    //go through reducing size from largest until
    ind := 0
    for {
      next_largest_index := indexable_widths[ind][0]
      col_widths[next_largest_index] = int(mx_fl * mv.width_ratios[ind])
      //log.Println(maxX, s, next_largest_index, indexable_widths[ind][1], "->", col_widths[next_largest_index])
      ind ++
      s = sum(col_widths, mv.fields_n)
      //expand it out in case we undershoot
      if (s < maxX) {
        col_widths[next_largest_index] += (maxX - s)
        break
      }
    }
    //log.Println(ind, col_widths)
  } else {
    //TODO fix this it's not correct
    extra_w := float64(maxX - s)
    for i:=0;i<mv.fields_n;i++ {
      v := int(extra_w * mv.width_ratios[i])
      col_widths[i] += v
    }
  }

  return col_widths
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	lx := 0
	helper_h := 3

  col_widths := set_col_width(maxX-1)

	for i := 0; i < mv.fields_n; i++ {
		if v, err := g.SetView(strconv.Itoa(i), lx, 0, lx + col_widths[i], maxY-1-helper_h, 0); err != nil {
			if ! gocui.IsUnknownView(err) {
        //log.Println(lx, col_widths, mv.width_ratios, mv.max_widths, mv.fields_n, mv.records_len)
        panic(err)
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
	if v, err := g.SetView("helper", 0, maxY-helper_h, maxX-1, maxY-1, 0); err != nil {
    if ! gocui.IsUnknownView(err) {
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
		//TODO need better solution
    r1 := mv.records[i][ind]
    r2 := mv.records[j][ind]

    //handle if the fields are empty
    if r1 == "" {
      return false
    }
    if r2 == "" {
      return true
    }

    //try to sort as ints otherwise fall through
		if v1, err := strconv.Atoi(r1); err == nil {
			v2, err := strconv.Atoi(r2)
			if err != nil {
				//panic(err) //field has non-empty mixed type values
        //perhaps should just revert to string comparison here
        return r1 < r2
			}
			return v1 < v2
		}

    //revert to sort by string
		return r1 < r2
	}

	//switch order on subseq calls
	if sort.SliceIsSorted(mv.records, sortbyCol) { //TODO wasteful keep a variable and switch instead of this
		sort.Slice(mv.records, func(i, j int) bool {
			return sortbyCol(j, i)
		})
	} else {
		sort.SliceStable(mv.records, sortbyCol)
	}


	fill_cols(g)
	return nil
}
