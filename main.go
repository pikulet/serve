package main

import (
    "fmt"
    "net/http"
    "github.com/gin-gonic/gin"
    "time"
    "crypto/md5"

    "github.com/pikulet/ghost-wp"
    "github.com/pikulet/kuute"
)

func main() {
    //gin.SetMode(gin.ReleaseMode)

    // routing
    r := gin.Default()
    r.GET("/", home)

    gpack.DownloadWords()
    gpack.Init()
    ghostRoutes := r.Group("/ghost")
    {
        ghostRoutes.GET("/", getGhostPairs)
    }

    kuute.Init()
    defer kuute.Shutdown()
    kuuteRoutes := r.Group("/kuute")
    {
        kuuteRoutes.GET("/", getKuuteIndex) 
        kuuteRoutes.GET("/:name", getKuuteBadge)
    }

    r.Run()
}

func home(c *gin.Context) {
    c.Redirect(http.StatusMovedPermanently, "https://github.com/pikulet")
    c.Abort()
}

func getGhostPairs(c *gin.Context) {
    tWord, fWord := gpack.GetRandomPair()
    c.JSON(http.StatusOK, gin.H {
        "town"  :   tWord,
        "fool"  :   fWord,
    })
}

func getKuuteIndex(c *gin.Context) {
    c.Redirect(http.StatusMovedPermanently, "https://github.com/pikulet/kuute")
    c.Abort()
}

func getKuuteBadge(c *gin.Context) {
    name := c.Param("name")
    badge := kuute.GetCounterBadge(name)

    // to block caching
    etag := fmt.Sprintf("%x", md5.Sum(badge))
    c.Header("Cache-Control", "no-cache")
    c.Header("Content-Type", "image/svg+xml;charset=utf-8")
    c.Header("ETag", etag)
    c.Header("Expires", time.Now().UTC().Format(http.TimeFormat))
    c.Header("Last-Modified", time.Now().UTC().Format(http.TimeFormat))
    c.String(http.StatusOK, string(badge))
}
