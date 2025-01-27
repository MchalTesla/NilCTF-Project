package controllers

import (
    "NilCTF/dto"
    "NilCTF/error_code"
    "NilCTF/services/interface"
    "net/http"

    "github.com/gin-gonic/gin"
)

type AdminAnnouncementController struct {
    AS services_interface.AnnouncementServiceInterface
}

// NewAdminAnnouncementController 创建一个新的 AdminAnnouncementController 实例
func NewAdminAnnouncementController(AS services_interface.AnnouncementServiceInterface) *AdminAnnouncementController {
    return &AdminAnnouncementController{AS: AS}
}

// HandleAnnouncement 处理公告相关操作
func (aac *AdminAnnouncementController) HandleAnnouncement(c *gin.Context) {
    var request struct {
        Action      string           `json:"action"`
        Announcement dto.AnnouncementDTO `json:"announcement"`
    }

    if err := c.ShouldBindJSON(&request); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": error_code.ErrInvalidInput.Error()})
        return
    }

    switch request.Action {
    case "create":
        aac.createAnnouncement(c, &request.Announcement)
    case "update":
        aac.updateAnnouncement(c, &request.Announcement)
    case "delete":
        aac.deleteAnnouncement(c, &request.Announcement)
    default:
        c.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "未知操作"})
    }
}

func (aac *AdminAnnouncementController) createAnnouncement(c *gin.Context, announcementDTO *dto.AnnouncementDTO) {
	
	announcementDTO.AuthorID = c.GetUint("userID")
	announcementDTO.PublishedAt = c.GetTime("currentTime")

    if err := aac.AS.Create(announcementDTO); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (aac *AdminAnnouncementController) updateAnnouncement(c *gin.Context, announcementDTO *dto.AnnouncementDTO) {
    if err := aac.AS.Update(announcementDTO); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (aac *AdminAnnouncementController) deleteAnnouncement(c *gin.Context, announcementDTO *dto.AnnouncementDTO) {
    if err := aac.AS.Delete(announcementDTO.ID); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"status": "ok"})
}