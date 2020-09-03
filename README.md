This is a simple wrapper for `v4l2-ctrl`. It useses `v4l2-ctrl --list-ctrls` to generate the startgin list. Every time you change one of the sliders, it writes that using `v4l2 --set-ctrl name=val`.
