package main

import (
  "github.com/jroimartin/gocui"
  "strconv"
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

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
  col_w := maxX / fields_n
  for i := 0; i < fields_n; i++ {
    if v, err := g.SetView(strconv.Itoa(i), col_w*i, 0, col_w*(i+1), maxY-1); err != nil {
      if err != gocui.ErrUnknownView { return err }

      //v.Frame = false //no border
      v.Editable = true

      if _, err := g.SetCurrentView(strconv.Itoa(i)); err != nil { return err }
    }
  }

	return nil
}

func keybinds(g *gocui.Gui) {
  err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit)
  if err != nil { panic(err) }
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}


