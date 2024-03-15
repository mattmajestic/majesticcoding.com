package main

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    "os"

    "github.com/gin-gonic/gin"
    "github.com/joho/godotenv"
)

type Video struct {
    Id string `json:"videoId"`
}

type Item struct {
    Id Video `json:"id"`
    Statistics struct {
        ViewCount    string `json:"viewCount"`
        LikeCount    string `json:"likeCount"`
        CommentCount string `json:"commentCount"`
    } `json:"statistics"`
}

type Response struct {
    Items []Item `json:"items"`
}

func getYoutubeMetrics(c *gin.Context) {
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }

    API_KEY := os.Getenv("YT_API_KEY")
    if API_KEY == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "YouTube API key not found"})
        return
    }

    CHANNEL_ID := "UCjTavL86-CW6j58fsVIjTig"
    url := fmt.Sprintf("https://www.googleapis.com/youtube/v3/search?key=%s&channelId=%s&part=id&maxResults=100", API_KEY, CHANNEL_ID)
    resp, err := http.Get(url)
    if err != nil {
        log.Printf("Error fetching YouTube search: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch data from YouTube"})
        return
    }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        log.Printf("Error reading response body: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response from YouTube"})
        return
    }

    var searchData Response
    json.Unmarshal(body, &searchData)

    videoData := make([]map[string]string, 0)

    for _, item := range searchData.Items {
        videoID := item.Id.Id
        statsUrl := fmt.Sprintf("https://www.googleapis.com/youtube/v3/videos?key=%s&id=%s&part=statistics", API_KEY, videoID)
        statsResp, err := http.Get(statsUrl)
        if err != nil {
            log.Printf("Error fetching YouTube video statistics: %v", err)
            continue
        }
        defer statsResp.Body.Close()

        statsBody, err := ioutil.ReadAll(statsResp.Body)
        if err != nil {
            log.Printf("Error reading statistics response body: %v", err)
            continue
        }

        var statsData Response
        json.Unmarshal(statsBody, &statsData)

        for _, statsItem := range statsData.Items {
            videoMetrics := map[string]string{
                "Video ID":  videoID,
                "Views":     statsItem.Statistics.ViewCount,
                "Likes":     statsItem.Statistics.LikeCount,
                "Comments":  statsItem.Statistics.CommentCount,
            }

            videoData = append(videoData, videoMetrics)
        }
    }

    c.JSON(http.StatusOK, videoData)
}

func main() {
    router := gin.Default()
    router.GET("/youtube-metrics", getYoutubeMetrics)
    router.Run(":8080")
}
