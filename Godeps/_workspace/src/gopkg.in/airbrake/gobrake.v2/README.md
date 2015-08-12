# Airbrake Golang Notifier [![Build Status](https://circleci.com/gh/airbrake/gobrake.png?circle-token=4cbcbf1a58fa8275217247351a2db7250c1ef976)](https://circleci.com/gh/airbrake/gobrake)

<img src="http://f.cl.ly/items/3J3h1L05222X3o1w2l2L/golang.jpg" width=800px>

Example
---

```go
package main

import (
	"errors"

	"gopkg.in/airbrake/gobrake.v2"
)

var airbrake = gobrake.NewNotifier(1234567, "FIXME")

func init() {
	airbrake.AddFilter(func(notice *gobrake.Notice) *gobrake.Notice {
		notice.Context["environment"] = "production"
		return notice
	})
}

func main() {
	defer airbrake.Flush()

	airbrake.Notify(errors.New("qqq"), nil)
}
```
