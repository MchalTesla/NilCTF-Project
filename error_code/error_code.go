package error_code

import "errors"

// 定义所有错误
var (
	// 比赛相关错误
	ErrCompetitionAlreadyExists = errors.New("ERR_COMPETITION_ALREADY_EXISTS") // 比赛已存在
	ErrCompetitionNotFound      = errors.New("ERR_COMPETITION_NOT_FOUND")      // 未找到比赛
	
	// 队伍相关错误
	ErrTeamAlreadyExists        = errors.New("ERR_TEAM_ALREADY_EXISTS")        // 队伍已存在
	ErrTeamNotFound             = errors.New("ERR_TEAM_NOT_FOUND")             // 未找到队伍
	ErrTeamAlreadyInCompetition	= errors.New("ERR_TEAM_ALREADY_IN_COMPETITION")// 队伍已在比赛中
	ErrTeamNotInCompetition		= errors.New("ERR_TEAM_NOT_IN_COMPETITION")    // 队伍不在比赛中
	
	// 用户相关错误
	ErrUserAlreadyInTeam        = errors.New("ERR_USER_ALREADY_IN_TEAM")       // 用户已在队伍中
	ErrUserNotInTeam            = errors.New("ERR_USER_NOT_IN_TEAM")           // 用户不在队伍中
	ErrEmailTaken               = errors.New("ERR_EMAIL_TAKEN")                // 邮箱已被占用
	ErrUsernameExists           = errors.New("ERR_USERNAME_EXISTS")            // 用户名已存在
	ErrUserNotFound             = errors.New("ERR_USER_NOT_FOUND")             // 未找到用户
	
	// 通用错误
	ErrInvalidInput             = errors.New("ERR_INVALID_INPUT")              // 无效输入
	ErrInternalServer           = errors.New("ERR_INTERNAL_SERVER")            // 内部服务器错误
)