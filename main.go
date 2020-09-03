// Copyright 2017 Zack Guo <zack.y.guo@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT license that can
// be found in the LICENSE file.

package main

import (
  "log"
  "time"
  "strings"
  "strconv"

  ui "github.com/gizak/termui"
  "github.com/gizak/termui/widgets"
)

func main() {
  ctrls, max_len := list_ctrls()
  lines := []string{}
  for _, ctrl := range ctrls {
    data := gen_line(ctrl, max_len)
    lines = append(lines, data)
  }

  if err := ui.Init(); err != nil {
    log.Fatalf("failed to initialize termui: %v", err)
  }
  defer ui.Close()

  p := widgets.NewParagraph()
  p.Title = "Text Box"
  p.Text = "PRESS q TO QUIT DEMO"
  p.SetRect(0, 0, 50, 5)
  p.TextStyle.Fg = ui.ColorWhite
  p.BorderStyle.Fg = ui.ColorCyan

  updateParagraph := func(count int) {
    if count%2 == 0 {
      p.TextStyle.Fg = ui.ColorRed
    } else {
      p.TextStyle.Fg = ui.ColorWhite
    }
  }

  list := widgets.NewList()
  list.Title = "Controls"
  list.Rows = lines
  list.SetRect(0, 5, 200, 20)
  list.TextStyle.Fg = ui.ColorYellow

  draw := func() {
    ui.Render(p, list)
  }

  tickerCount := 1
  draw()
  tickerCount++
  uiEvents := ui.PollEvents()
  ticker := time.NewTicker(time.Second).C
  for {
    select {
    case e := <-uiEvents:
      switch e.ID {
      case "q", "<C-c>":
        return
      case "j":
        list.SelectedRow++
        if list.SelectedRow >= len(list.Rows) {
          list.SelectedRow--
        }
      case "k":
        list.SelectedRow--
        if list.SelectedRow < 0 {
          list.SelectedRow++
        }
      }
      draw()
    case <-ticker:
      updateParagraph(tickerCount)
      draw()
      tickerCount++
    }
  }
}

func gen_line(ctrl *control, max_len int) string {
  line := ""
  line = add_col(line, ctrl.name,              max_len)
  line = add_col(line, strconv.Itoa(ctrl.min), 6)
  line = add_col(line, strconv.Itoa(ctrl.val), 6)
  line = add_col(line, strconv.Itoa(ctrl.max), 6)
  line += "["
  percent := float64(ctrl.val - ctrl.min) / float64(ctrl.max - ctrl.min)
  num_chars := int(percent * 50)
  line += strings.Repeat("-", num_chars)
  line += strings.Repeat(" ", 50 - num_chars)
  line += "] "
  return line
}

func add_col(line string, val string, max_len int) string {
  return line + val + strings.Repeat(" ", max_len - len(val) + 1)
}
