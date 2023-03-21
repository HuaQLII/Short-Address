package main

import (
	"github/im-lauson/Short-Address/conf"
)

func main() {
	a := App{}
	a.Initialize(conf.GetEnv())
	a.Run(":9090")

}
