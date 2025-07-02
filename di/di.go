package di

import "hunter-backend/di/server"

func InitApplication() error {
	err := server.InitApiServer()
	if err != nil {
		return err
	}
	return nil
}
