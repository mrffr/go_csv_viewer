package main

import (
  "github.com/jroimartin/gocui"
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
  if v, err := g.SetView("main", 0, 0, maxX-1, maxY-1); err != nil {
    if err != gocui.ErrUnknownView { return err }

    //v.Frame = false //no border
    v.Editable = true

    if _, err := g.SetCurrentView("main"); err != nil { return err }
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


