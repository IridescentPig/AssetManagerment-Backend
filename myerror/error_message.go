package myerror

const (
	INVALID_BODY_INFO             = "Invalid request body"
	USER_NOT_FOUND_INFO           = "User not found"
	PERMISSION_DENIED_INFO        = "Permission denied"
	DUPLICATED_NAME_INFO          = "Dupilicated name"
	INVALID_PARAM_INFO            = "Invalid param in router"
	TOKEN_EMPTY_INFO              = "Cannot find token in request header"
	TOKEN_EXPIRED_INFO            = "Token has expired"
	TOKEN_INVALID_INFO            = "Invaild token"
	ENTITY_NOT_FOUND_INFO         = "Entity not found"
	USER_HAS_EXISTED_INFO         = "User has existed"
	USER_NOT_IN_ENTITY_INFO       = "User does not exist in this entity"
	NAME_CANNOT_EMPTY_INFO        = "Name cannot be empty"
	DEPARTMENT_NOT_FOUND_INFO     = "department not found"
	DEPARTMENT_NOT_IN_ENTITY_INFO = "Department not in entity"
	USER_NOT_IN_DEPARTMENT_INFO   = "User not in department"
	ENTITY_HAS_USERS_INFO         = "Entity has users, cannot be deleted"
	DELETE_USER_SELF_INFO         = "You cannot delete yourself"
	DEPARTMENT_HAS_USERS_INFO     = "Department has users, cannot be deleted"
	ASSET_CLASS_NOT_FOUND_INFO    = "Asset class not found"
)
