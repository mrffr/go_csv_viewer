package main

import (
  "github.com/jroimartin/gocui"
  "strconv"
  "fmt"
  "sort"
)


func run_ui(){
  g, err := gocui.NewGui(gocui.OutputNormal)
  if err != nil { panic(err) }

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
    if err != nil { panic(err) }
    v.Clear()
    for j := 0; j < rec_len; j++ {
      fmt.Fprintln(v, mv.records[j][i])
    }
  }
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
  //col_w := maxX / fields_n //TODO var size set in read func
  mx_fl := float64(maxX)
  lx := 0
  helper_h := 3
  for i := 0; i < mv.fields_n; i++ {
//    if v, err := g.SetView(strconv.Itoa(i), col_w*i, 0, col_w*(i+1), maxY-1); err != nil {
    col_w := int(mx_fl * mv.width_ratios[i])
    if v, err := g.SetView(strconv.Itoa(i), lx, 0, lx+col_w, maxY-1-helper_h); err != nil {
      if err != gocui.ErrUnknownView { return err }

      v.Frame = false //no border
      v.Editable = false

      if mv.has_header {
        v.Title = mv.header[i]
        v.Frame = true
      }
    }
    lx += col_w
  }

  fill_cols(g)

  //setup view on first run
  if g.CurrentView() == nil {
    if _, err := g.SetCurrentView(strconv.Itoa(0)); err != nil { return err }
  }

  //helper height
  if v, err := g.SetView("helper", 0, maxY-helper_h, maxX-1, maxY-1); err != nil {
    if err != gocui.ErrUnknownView { return err }
    v.Editable = false

    fmt.Fprintln(v, "Ctrl-C: quit |", "Ctrl-S: sort |")
  }

	return nil
}

func keybinds(g *gocui.Gui) {
  err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit)
  if err != nil { panic(err) }

  err = g.SetKeybinding("", gocui.KeyArrowDown, gocui.ModNone,
  func(g *gocui.Gui, v *gocui.View) error {
    return nextLine(g, v, 1)
  })
  if err != nil { panic(err) }

  err = g.SetKeybinding("", gocui.KeyArrowUp, gocui.ModNone,
  func(g *gocui.Gui, v *gocui.View) error {
    return nextLine(g, v, -1)
  })
  if err != nil { panic(err) }

  // left right
  err = g.SetKeybinding("", gocui.KeyArrowRight, gocui.ModNone,
  func(g *gocui.Gui, v *gocui.View) error {
    return nextView(g, v, 1)
  })
  if err != nil { panic(err) }

  err = g.SetKeybinding("", gocui.KeyArrowLeft, gocui.ModNone,
  func(g *gocui.Gui, v *gocui.View) error {
    return nextView(g, v, -1)
  })
  if err != nil { panic(err) }

  //paging
  err = g.SetKeybinding("", gocui.KeyPgdn, gocui.ModNone,
  func(g *gocui.Gui, v *gocui.View) error {
    _, h := v.Size()
    return nextLine(g, v, h)
  })
  if err != nil { panic(err) }

  err = g.SetKeybinding("", gocui.KeyPgup, gocui.ModNone,
  func(g *gocui.Gui, v *gocui.View) error {
    _, h := v.Size()
    return nextLine(g, v, -h)
  })
  if err != nil { panic(err) }

  //sort
  err = g.SetKeybinding("", gocui.KeyCtrlS, gocui.ModNone, sortCol)
  if err != nil { panic(err) }

}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

//Move to next line dir -1,1 for u,d does not wrap around
func nextLine (g *gocui.Gui, v *gocui.View, dir int) error {
  x, y := v.Cursor()

  //moving lines
  y += dir
  _, pageH := v.Size()
  _, oy := v.Origin()

  if y >= pageH { //scroll down
    scrollViews(g, oy + pageH)

    //if top of page is now blank we scrolled too far
    //so revert scroll
    if ! notEmptyLine(v, 0) {
      scrollViews(g, oy)
    }

    if notEmptyLine(v, y - pageH) {
      v.SetCursor(x, y - pageH)
    }else{
      //set cursor on last occupied line
      v.SetCursor(x, getLastLine(v))
    }
  }else if y < 0{ //scroll up

    //we are already at the top
    if oy == 0 { return nil }
    //scroll up page
    ny := oy - pageH
    //make sure we don't overshoot top
    if ny < 0 { ny = 0 }
    scrollViews(g, ny)
    v.SetCursor(x, y + pageH)
  }else{
    //normal movement
    if notEmptyLine(v, y) { //TODO only need to check on scroll down
      v.SetCursor(x, y)
    }
  }

  return nil
}

func notEmptyLine(v *gocui.View, y int) bool{
  ln, _ := v.Line(y)
  return (len(ln) > 0)
}

func getLastLine(v *gocui.View) int {
  _, h := v.Size()
  for i := h; i >= 0; i-- {
    if notEmptyLine(v, i){ return i }
  }
  panic(-1)
  return -1
}

// scrolls all columns
func scrollViews(g *gocui.Gui, ny int){
  for i := 0; i < mv.fields_n; i++ {
    v, err := g.View(strconv.Itoa(i))
    if err != nil { panic(err) }
    v.SetOrigin(0, ny)
  }
}

// Move to next view dir is -1, 1 for l,r wraps around.
func nextView (g *gocui.Gui, v *gocui.View, dir int) error {
  //store cursor position so we are on correct line
  _, y := v.Cursor()

  //moving columns
  n, err := strconv.Atoi(v.Name())
  if err != nil { return err }
  n += dir
  if n < 0 { n = mv.fields_n - 1 }
  n %= mv.fields_n
  new_v, err := g.SetCurrentView(strconv.Itoa(n))
  if err != nil { return err }

  //restore "correct" position
  new_v.SetCursor(0, y)
  return nil
}

// Sort column in desc, asc cycle
func sortCol(g *gocui.Gui, v *gocui.View) error{
  ind, err := strconv.Atoi(v.Name())
  if err != nil { return err }

  sortbyCol := func(i, j int) bool {
    //TODO terrible code need better solution
    if v1, err := strconv.Atoi(mv.records[i][ind]); err == nil {
      v2, err := strconv.Atoi(mv.records[j][ind])
      if err != nil { panic(err) }
      return v1 < v2
    }
    return mv.records[i][ind] < mv.records[j][ind]
  }

  //switch order on subseq calls
  if sort.SliceIsSorted(mv.records, sortbyCol) { //TODO wasteful keep a variable and switch instead of this
    sort.Slice(mv.records, func(i, j int) bool {
      return sortbyCol(j, i)
    })
  }else{
    sort.Slice(mv.records, sortbyCol)
  }

  fill_cols(g)
  return nil
}
