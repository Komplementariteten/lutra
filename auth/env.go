package auth

import (
	"github.com/Komplementariteten/lutra"
)

type Env struct {
	User   string
	Config *lutra.LutraConfig
}

var Environment Env

func init() {
	Environment.User = ""
	Environment.Config = &lutra.LutraConfig{}
}
