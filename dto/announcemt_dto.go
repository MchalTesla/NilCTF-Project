package dto

import "time"

type AnnouncementDTO struct {
    ID          uint      `json:"id"`
    Title       string    `json:"title"`
    Content     string    `json:"content"`
    AuthorID    uint      `json:"author_id"`
    PublishedAt time.Time `json:"published_at"`
    Status      string    `json:"status"`
    Priority    string    `json:"priority"`
}