package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Komplementariteten/lutra/auth"
	"github.com/Komplementariteten/lutra/db"
	"github.com/Komplementariteten/lutra/util"
)

type portalHandle struct {
	db   *db.Db
	mail *util.SmtpClient
}

func (a *portalHandle) ShouldBeExecuted(n int) bool {
	targetTime := 1
	return n%targetTime == 0
}
func (a *portalHandle) Cleanup(ctx context.Context) {
	// os.RemoveAll(FileCachePath)
	// os.Mkdir(FileCachePath, 0777)
	fmt.Println("portalHandle Cleanup called")
}

func (*portalHandle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	claims, err := auth.CheckRefresh(w, r)
	if err != nil {
		w.Write([]byte("HTTP Error"))
	} else {
		w.Write([]byte(fmt.Sprintf("Hallo %s", claims.Username)))
	}
}
