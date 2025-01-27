package controllers

import (
    "NilCTF/error_code"
    "NilCTF/services/interface"
    "net/http"

    "github.com/gin-gonic/gin"
)

type AnnouncementController struct {
    AS services_interface.AnnouncementServiceInterface
}

// NewAnnouncementController 创建一个新的 AnnouncementController 实例
func NewAnnouncementController(AS services_interface.AnnouncementServiceInterface) *AnnouncementController {
    return &AnnouncementController{AS: AS}
}

// ListAnnouncements 列出所有公告
func (ac *AnnouncementController) ListAnnouncements(c *gin.Context) {
    announcements, err := ac.AS.List(nil, 0, 0, false)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": error_code.ErrInternalServer.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"status": "ok", "announcements": announcements})
}