package main

import (
	"github.com/Fajar-Islami/simple_manage_products/internal/infrastructure/container"
	"github.com/Fajar-Islami/simple_manage_products/internal/infrastructure/mysql"

	rest "github.com/Fajar-Islami/simple_manage_products/internal/server/http"
)

func main() {
	containerConf := container.InitContainer()
	defer mysql.CloseDatabaseConnection(containerConf.Mysqldb)

	rest.HTTPRouteInit(containerConf, containerConf)

}
