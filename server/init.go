package server

import (
	"fmt"
	"ws-gateway/config"
	"ws-gateway/server/router"
)

func StartGinServer() error {
	r := router.SetupRouter()
	if err := r.Run(fmt.Sprintf("0.0.0.0:%s", config.Conf.Http.Port)); err != nil {
		return err
	} else {
		return nil
	}
}
