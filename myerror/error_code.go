package myerror

const (
	INVALID_BODY                 = -1
	SUCCESS                      = 0
	USER_NOT_FOUND               = 1
	PERMISSION_DENIED            = 2
	DUPLICATED_NAME              = 3
	TOKEN_EMPTY                  = 5
	TOKEN_INVALID                = 6
	TOKEN_EXPIRED                = 7
	INVALID_PARAM                = 8
	ENTITY_NOT_FOUND             = 9
	USER_HAS_EXISTED             = 10
	USER_NOT_IN_ENTITY           = 11
	NAME_CANNOT_EMPTY            = 12
	DEPARTMENT_NOT_FOUND         = 13
	DEPARTMENT_NOT_IN_ENTITY     = 14
	USER_NOT_IN_DEPARTMENT       = 15
	ENTITY_HAS_USERS             = 16
	DELETE_USER_SELF             = 17
	DEPARTMENT_HAS_USERS         = 18
	ASSET_CLASS_NOT_FOUND        = 19
	PARENT_ASSET_CLASS_NOT_FOUND = 20
	INVALID_TYPE_OF_CLASS        = 21
	PARENT_CANNOOT_BE_SUCCESSOR  = 22
	CLASS_HAS_ASSET              = 23
	ASSET_NOT_FOUND              = 24
	ASSET_NOT_IN_DEPARTMENT      = 25
	PARENT_ASSET_NOT_FOUND       = 26
)
