package auth

import (
	"lutra"
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
