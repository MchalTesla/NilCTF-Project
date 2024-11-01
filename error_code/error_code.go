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
	ErrInvalidTeamname			= errors.New("ERR_INVALID_TEAMNAME")		   // 队伍名无效
	
	// 用户相关错误
	ErrUserAlreadyInTeam        = errors.New("ERR_USER_ALREADY_IN_TEAM")       // 用户已在队伍中
	ErrUserNotInTeam            = errors.New("ERR_USER_NOT_IN_TEAM")           // 用户不在队伍中
	ErrInvalidEmail               = errors.New("ERR_EMAIL_TAKEN")                // 邮箱无效
	ErrEmailExists              = errors.New("ERR_EMAIL_EXISTS")                // 邮箱已存在
	ErrInvalidUsername           = errors.New("ERR_INVALID_USERNAME")            // 用户名无效
	ErrUsernameExists           = errors.New("ERR_USERNAME_EXISTS")            // 用户名已存在
	ErrUserNotFound             = errors.New("ERR_USER_NOT_FOUND")             // 未找到用户
	ErrInvalidCredentials 		= errors.New("ERR_INVALID_CREDENTIALS") 	   // 用户名或密码错误
	
	// 通用错误
	ErrInvalidInput             = errors.New("ERR_INVALID_INPUT")              // 无效输入
	ErrInternalServer           = errors.New("ERR_INTERNAL_SERVER")            // 内部服务器错误
	ErrPermissionDenied			= errors.New("ERR_PERMISSION_DENIED")		   // 权限不足
	ErrNotFound					= errors.New("ERR_NOT_FOUND")				   // 未找到
	ErrInvalidDescription		= errors.New("ERR_INVALID_DESCRIPTION")		   // 描述无效
	ErrInvalidID				= errors.New("ERR_INVALID_ID")				   // ID无效

	// 中间件错误
	ErrFailedToParseForm		= errors.New("ERR_FAILED_TO_PARSE_FORM")		// 解析表单错误
	ErrTooManyParameters		= errors.New("ERR_TOO_MANY_PARAMETERS")			// 太多参数
	ErrKeyTooLong				= errors.New("ERR_KEY_TOO_LONG")				// 键值太长
	ErrInputTooLong				= errors.New("ERR_INPUT_TOO_LONG")				// 字段太长
	ErrTooManyRequests			= errors.New("ERR_TOO_MANY_REQUESTS")			// 太多请求
)