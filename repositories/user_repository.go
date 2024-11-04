package repositories

import (
	"NilCTF/error_code"
	"NilCTF/models"

	"gorm.io/gorm"
)

type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepository(DB *gorm.DB) *UserRepository {
	return &UserRepository{DB: DB}
}

// Create 插入新的用户记录
func (r *UserRepository) Create(user *models.User) error {
	if user.ID != 0 {
		return error_code.ErrInvalidID
	}
	if err := r.DB.Create(&user).Error; err != nil {
		return error_code.ErrInternalServer
	}
	return nil
}
func (r *UserRepository) Update(user *models.User) error {
	if user.ID == 0 {
		return error_code.ErrInvalidID
	}
	if err := r.DB.Updates(&user).Error; err != nil {
		return error_code.ErrInternalServer
	}
	return nil
}

// Get 查找单个用户记录，通过 ID、邮箱或用户名进行条件查询
func (r *UserRepository) Get(ID uint, email, username string) (*models.User, error) {
	var user models.User
	query := r.DB.Model(&user)

	// 动态构建查询条件
	if ID != 0 {
		query = query.Where("id = ?", ID)
	} else if email != "" {
		query = query.Where("email = ?", email)
	} else if username != "" {
		query = query.Where("username = ?", username)
	} else {
		return nil, error_code.ErrInvalidInput
	}

	// 执行查询
	if err := query.First(&user).Error; err != nil {
		// 捕获系统错误
		if err == gorm.ErrRecordNotFound {
			return nil, error_code.ErrUserNotFound
		}
		return nil, error_code.ErrInternalServer
	}

	return &user, nil
}

// Delete 删除用户记录
func (r *UserRepository) Delete(user *models.User) error {
	if user.ID == 0 {
		return error_code.ErrInvalidID
	}
	if err := r.DB.Delete(user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return error_code.ErrUserNotFound
		}
		return error_code.ErrInternalServer
	}
	return nil
}

// List 列出所有符合条件的记录
func (r *UserRepository) List(filters map[string]interface{}, limit, offset int, isFuzzy bool) ([]models.User, error) {
	var users []models.User
	query := r.DB

	// 应用过滤条件
	for key, value := range filters {
		if isFuzzy { // 如果启用模糊搜索
			// 使用 LIKE 查询并在值两端添加 % 通配符
			query = query.Where(key+" LIKE ?", "%"+value.(string)+"%")
		} else {
			// 精确匹配
			query = query.Where(key+" = ?", value)
		}
	}

	// 设置分页
	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset >= 0 {
		query = query.Offset(offset)
	}

	// 执行查询
	if err := query.Find(&users).Error; err != nil {
		return []models.User{}, error_code.ErrInternalServer
	}
	return users, nil
}