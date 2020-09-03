package main

import (
  "os/exec"
  "strings"
  "strconv"
)

type control struct {
  name string
  min int
  max int
  val int
  default_val int
}

type string_writer struct {
  data []byte
}

func (s *string_writer) Write(data []byte) (int, error) {
  s.data = append(s.data, data...)
  return len(data), nil
}

func list_ctrls() ([]*control, int) {
  cmd := exec.Command("v4l2-ctl", "--list-ctrls")
  w := &string_writer{}
  cmd.Stdout = w
  cmd.Run()
  lines := strings.Split(string(w.data), "\n")
  max_len := 0
  controls := []*control{}
  for _, line := range lines {
    sections := strings.Split(line, " ")
    i := 0
    control := control{}
    for _, section := range sections {
      if len(section) == 0 {
        continue
      }
      if i == 0 {
        control.name = section
        if len(section) > max_len {
          max_len = len(section)
        }
      } else if i == 2 {
        if section == "(bool)" {
          control.min = 0
          control.max = 1
        }
      }
      parts := strings.Split(section, "=")
      if len(parts) == 2 {
        name := parts[0]
        num, err := strconv.ParseInt(parts[1], 10, 64)
        if err == nil {
          switch name {
          case "min":
            control.min = int(num)
          case "max":
            control.max = int(num)
          case "value":
            control.val = int(num)
          case "default":
            control.default_val = int(num)
          }
        }
      }
      i++
    }
    if control.name != "" {
      controls = append(controls, &control)
    }
  }
  return controls, max_len
}
