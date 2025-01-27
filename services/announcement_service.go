package services

import (
	"NilCTF/dto"
	"NilCTF/error_code"
	managers_interface "NilCTF/managers/interface"
	"NilCTF/models"
	"NilCTF/repositories/interface"
)

type AnnouncementService struct {
    AM managers_interface.AnnouncementManagerInterface
    UM managers_interface.UserManagerInterface
}

// NewAnnouncementService 返回一个新的 AnnouncementService 实例
func NewAnnouncementService(AM repositories_interface.AnnouncementRepositoryInterface, UM managers_interface.UserManagerInterface) *AnnouncementService {
    return &AnnouncementService{AM: AM, UM: UM}
}

func (s *AnnouncementService) Create(announcementDTO *dto.AnnouncementDTO) error {
    if announcementDTO.Title == "" || announcementDTO.Content == "" {
        return error_code.ErrInvalidInput
    }

    announcement := &models.Announcement{
        Title:       announcementDTO.Title,
        Content:     announcementDTO.Content,
        AuthorID:    announcementDTO.AuthorID,
        PublishedAt: announcementDTO.PublishedAt,
        Status:      announcementDTO.Status,
        Priority:    announcementDTO.Priority,
    }

    return s.AM.Create(announcement)
}

func (s *AnnouncementService) Get(ID uint) (*dto.AnnouncementDTO, error) {
    if ID == 0 {
        return nil, error_code.ErrInvalidID
    }

    announcement, err := s.AM.Get(ID)
    if err != nil {
        return nil, err
    }

    return &dto.AnnouncementDTO{
        ID:          announcement.ID,
        Title:       announcement.Title,
        Content:     announcement.Content,
        AuthorID:    announcement.AuthorID,
        PublishedAt: announcement.PublishedAt,
        Status:      announcement.Status,
        Priority:    announcement.Priority,
    }, nil
}

func (s *AnnouncementService) Update(announcementDTO *dto.AnnouncementDTO) error {
    if announcementDTO.ID == 0 {
        return error_code.ErrInvalidID
    }
    if announcementDTO.Title == "" || announcementDTO.Content == "" {
        return error_code.ErrInvalidInput
    }

	var announcement models.Announcement
	announcement.ID = announcementDTO.ID
	announcement.Title = announcementDTO.Title
	announcement.Content = announcementDTO.Content
	announcement.AuthorID = announcementDTO.AuthorID
	announcement.PublishedAt = announcementDTO.PublishedAt
	announcement.Status = announcementDTO.Status
	announcement.Priority = announcementDTO.Priority

    return s.AM.Update(&announcement)
}

func (s *AnnouncementService) Delete(ID uint) error {
    if ID == 0 {
        return error_code.ErrInvalidID
    }

    announcement, err := s.AM.Get(ID)
    if err != nil {
        return err
    }

    return s.AM.Delete(announcement)
}

func (s *AnnouncementService) List(filters map[string]interface{}, limit, offset int, isFuzzy bool) ([]dto.AnnouncementDTO, error) {
    announcements, err := s.AM.List(filters, limit, offset, isFuzzy)
    if err != nil {
        return nil, err
    }

    var announcementDTOs []dto.AnnouncementDTO
    for _, announcement := range announcements {
        user, err := s.UM.Get(announcement.AuthorID, "", "")
        if err != nil {
            return nil, err
        }
        announcementDTOs = append(announcementDTOs, dto.AnnouncementDTO{
            ID:          announcement.ID,
            Title:       announcement.Title,
            Content:     announcement.Content,
            AuthorID:    announcement.AuthorID,
            AuthorName:  user.Username,
            PublishedAt: announcement.PublishedAt,
            Status:      announcement.Status,
            Priority:    announcement.Priority,
        })
    }
    return announcementDTOs, nil
}