package postgres

import (
	"pg_permissions/database"
)

/*
 Permissions data structure

        "user" => {

                "databases": [ // ObjectMap
                        "db1": {"insert": "true", "read": "false"}, // ObjectDesc
                        "db2": ...
                ]
                "functions": [
                ]
        }
*/

type PermissionMap map[string]bool

type ObjectDesc map[string]PermissionMap

type ObjectMap struct {
	FailedToCheck map[string]string     `json:"failed_to_check,omitempty"`
	ObjectDesc    map[string]ObjectDesc `json:"objects,omitempty"`
	Permissions   map[string][]string   `json:"-"`
}

type UserPermissionMap map[string]ObjectMap

func PrivilegeDescriptorToPermissionMap(privilegeDescriptor *database.PrivilegeDescriptor, permissions []string) PermissionMap {
	permissionMap := make(PermissionMap)

	for idx, _ := range permissions {
		status := privilegeDescriptor.HasPrivilegeEnabled(&permissions[idx])
		permissionMap[permissions[idx]] = status
	}

	return permissionMap
}

/*
 * Role object must have different logic applied.
 */
func PrivilegeDescriptorToRolePermissionMap(roleDescriptor *database.Role, permissions []string) PermissionMap {
	permissionMap := make(PermissionMap)

	permissionMap["CREATEDB"] = false
	permissionMap["SUPERUSER"] = false
	permissionMap["REPLICATION"] = false
	permissionMap["CREATEROLE"] = false
	permissionMap["CANLOGIN"] = false

	if roleDescriptor.CreateDBPrivilege == "t" || roleDescriptor.CreateDBPrivilege == "true" {
		permissionMap["CREATEDB"] = true
	}

	if roleDescriptor.SuperuserPrivilege == "t" || roleDescriptor.SuperuserPrivilege == "true" {
		permissionMap["SUPERUSER"] = true
	}

	if roleDescriptor.ReplicationPrivilege == "t" || roleDescriptor.ReplicationPrivilege == "true" {
		permissionMap["REPLICATION"] = true
	}

	if roleDescriptor.CreateRolePrivilege == "t" || roleDescriptor.CreateRolePrivilege == "true" {
		permissionMap["CREATEROLE"] = true
	}

	if roleDescriptor.CanLogin == "t" || roleDescriptor.CanLogin == "true" {
		permissionMap["CANLOGIN"] = true
	}

	return permissionMap
}
