package main

import (
	"mall/cache"
	"mall/conf"
	"mall/routes"
)

func main() {
	// Ek1+Ep1==Ek2+Ep2
	conf.Init()
	cache.Init()
	r := routes.NewRouter()
	_ = r.Run(conf.HttpPort)
}
