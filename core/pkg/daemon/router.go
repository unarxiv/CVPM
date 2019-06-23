// Copyright 2019 The CVPM Authors. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package daemon

import (
	"github.com/cvpm-contrib/auth"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/unarxiv/cvpm/pkg/config"
	"github.com/unarxiv/cvpm/pkg/contrib"
	"github.com/unarxiv/cvpm/pkg/runtime"
	api "gopkg.in/appleboy/gin-status-api.v1"
	"net/http"
	"path/filepath"
)

func getRouter() *gin.Engine {
	config := config.Read()
	webuiFolder := filepath.Join(config.Local.LocalFolder, "webui")
	r := gin.Default()
	r.Use(BeforeResponse())
	r.Use(gin.Recovery())
	watchLogs(socketServer)
	r.Use(static.Serve("/webui", static.LocalFile(webuiFolder, false)))
	r.Use(auth.InspectorStats())
	r.Use(gin.Logger())
	// Reverse Proxy for solvers
	r.POST("/engine/solvers/:vendor/:name/:solver", ReverseProxy)
	// Status Related Handlers
	r.GET("/status", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"daemon": "running",
		})
	})
	// System Related Handlers
	r.GET("/system", GetSystemHandler)
	// Repo Related Routes
	r.GET("/repo/meta/:vendor/:name", GetRepoMetaHandler)
	r.GET("/repo/env/:vendor/:name", runtime.QueryRepoEnvironments)
	r.POST("/repo/running", PostRunningRepoHandler)
	r.POST("/repo/envs/:vendor/:name", runtime.AddRepoEnvironments)
	r.GET("/repos", GetReposHandler)
	r.GET("/repos/running", GetRunningReposHandler)
	r.POST("/repos", PostReposHandler)
	// Solver Related Routers
	r.GET("/solvers/running", GetRunningSolversHandler)
	r.GET("/solvers/running/:vendor/:package", GetRunningSolversByPackageHandler)
	r.DELETE("/solvers/running/:vendor/:name/:solver", StopInferProcess)
	// Socket Related Routes
	r.GET("/socket.io/", socketHandler)
	r.POST("/socket.io/", socketHandler)
	// Contrib Related Routes
	// Datasets
	r.GET("/contrib/datasets", contrib.GetAllDatasets)
	r.POST("/contrib/datasets/registries", contrib.AddNewRegistry)
	// Files
	r.POST("/contrib/files/upload/:type", contrib.UploadFile)
	r.GET("/contrib/files/list", contrib.QueryFilesList)
	// Camera
	// Plugin Related Routes
	r.GET("/_inspector", func(c *gin.Context) {
		c.JSON(200, auth.GetRequests())
	})
	r.GET("/_api/status", api.StatusHandler)
	// Socket Related Routes
	r.Handle("WS", "/socket.io/", socketHandler)
	r.Handle("WSS", "/socket.io/", socketHandler)
	return r
}
