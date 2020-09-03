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
  bar_width := 50
  for _, ctrl := range ctrls {
    data := gen_line(ctrl, max_len, bar_width)
    lines = append(lines, data)
  }

  if err := ui.Init(); err != nil {
    log.Fatalf("failed to initialize termui: %v", err)
  }
  defer ui.Close()

  list := widgets.NewList()
  list.Title = "Controls"
  list.Rows = lines
  list.SetRect(0, 0, 200, 20)
  list.TextStyle.Fg = ui.ColorBlue
  list.SelectedRowStyle.Fg = ui.ColorWhite

  draw := func() {
    ui.Render(list)
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
      case "j", "<Down>":
        list.SelectedRow++
        if list.SelectedRow >= len(list.Rows) {
          list.SelectedRow--
        }
      case "k", "<Up>":
        list.SelectedRow--
        if list.SelectedRow < 0 {
          list.SelectedRow++
        }
      case "h", "<Left>":
        ctrl := ctrls[list.SelectedRow]
        change_value(ctrl, -1, bar_width)
        list.Rows[list.SelectedRow] = gen_line(ctrl, max_len, bar_width)
      case "l", "<Right>":
        ctrl := ctrls[list.SelectedRow]
        change_value(ctrl, 1, bar_width)
        list.Rows[list.SelectedRow] = gen_line(ctrl, max_len, bar_width)
      }
      draw()
    case <-ticker:
      draw()
      tickerCount++
    }
  }
}

func change_value(ctrl *control, amount, bar_width int) {
  size := float64(ctrl.max - ctrl.min)
  step := int(size / float64(bar_width))
  if step < 1 {
    step = 1
  }
  ctrl.val += amount * step
  if ctrl.val < ctrl.min {
    ctrl.val = ctrl.min
  }
  if ctrl.val > ctrl.max {
    ctrl.val = ctrl.max
  }
  set_ctrl(ctrl.name, ctrl.val)
}

func gen_line(ctrl *control, max_len, bar_width int) string {
  line := ""
  line = add_col(line, ctrl.name,              max_len)
  line = add_col(line, strconv.Itoa(ctrl.min), 6)
  line = add_col(line, strconv.Itoa(ctrl.val), 6)
  line = add_col(line, strconv.Itoa(ctrl.max), 6)
  line += "["
  percent := float64(ctrl.val - ctrl.min) / float64(ctrl.max - ctrl.min)
  num_chars := int(percent * float64(bar_width))
  line += strings.Repeat("-", num_chars)
  line += strings.Repeat(" ", bar_width - num_chars)
  line += "] "
  return line
}

func add_col(line string, val string, max_len int) string {
  return line + val + strings.Repeat(" ", max_len - len(val) + 1)
}
