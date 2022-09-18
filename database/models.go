package database

import (
	"pg_permissions/log"
	"strings"

	"github.com/jmoiron/sqlx"
)

type Controller struct {
	Database *sqlx.DB
}

func NewController(db *sqlx.DB) *Controller {
	log.Debugln("Initializing database controller")
	return &Controller{
		Database: db,
	}
}

type PrivilegeDescriptor struct {
	Usage      string `db:"usage"`
	Select     string `db:"select"`
	Insert     string `db:"insert"`
	Update     string `db:"update"`
	Delete     string `db:"delete"`
	References string `db:"references"`
	Trigger    string `db:"trigger"`
	Truncate   string `db:"truncate"`
	Connect    string `db:"connect"`
	Create     string `db:"create"`
	Temporary  string `db:"temporary"`
	Execute    string `db:"execute"`
	Member     string `db:"member"`
}

func (p *PrivilegeDescriptor) HasPrivilegeEnabled(privilege *string) bool {
	switch *privilege {
	case "USAGE":
		if p.Usage == "t" || p.Usage == "true" {
			return true
		}
	case "SELECT":
		if p.Select == "t" || p.Select == "true" {
			return true
		}
	case "INSERT":
		if p.Insert == "t" || p.Insert == "true" {
			return true
		}
	case "UPDATE":
		if p.Update == "t" || p.Update == "true" {
			return true
		}
	case "DELETE":
		if p.Delete == "t" || p.Delete == "true" {
			return true
		}
	case "REFERENCES":
		if p.References == "t" || p.References == "true" {
			return true
		}
	case "TRIGGER":
		if p.Trigger == "t" || p.Trigger == "true" {
			return true
		}
	case "TRUNCATE":
		if p.Truncate == "t" || p.Truncate == "true" {
			return true
		}
	case "CONNECT":
		if p.Connect == "t" || p.Connect == "true" {
			return true
		}
	case "CREATE":
		if p.Create == "t" || p.Create == "true" {
			return true
		}
	case "TEMPORARY":
		if p.Temporary == "t" || p.Temporary == "true" {
			return true
		}
	case "EXECUTE":
		if p.Execute == "t" || p.Execute == "true" {
			return true
		}
	case "MEMBER":
		if p.Member == "t" || p.Member == "true" {
			return true
		}
	}
	return false
}

type User struct {
	Name                         string `db:"name"`
	CreateDBPrivilege            string `db:"createdb_privilege"`
	SuperuserPrivilege           string `db:"superuser_privilege"`
	ReplicationPrivilege         string `db:"replication_privilege"`
	BypassRowLevelSecurityPolicy string `db:"bypass_row_level_security_policy"`
	PasswordValidUntil           string `db:"password_valid_until"`
}

type Column struct {
	Schema string `db:"schema"`
	Table  string `db:"table"`
	Name   string `db:"name"`
}

type Database struct {
	Name         string `db:"name"`
	IsTemplate   string `db:"is_template"`
	AllowConnect string `db:"allow_connect"`
}

type FDW struct {
	Name string `db:"name"`
}

type Function struct {
	Name      string `db:"name"`
	Schema    string `db:"schema"`
	Catalog   string `db:"catalog"`
	Arguments string `db:"arguments"`
}

/*
 * Dirty hack in order to extract function definition
 * Maybe it will be fixed some day
 */
func (f *Function) SanitizeArguments() {
	before, _, found := strings.Cut(f.Arguments, ", OUT")
	if found {
		f.Arguments = before
	}

	before, _, found = strings.Cut(f.Arguments, "OUT")
	if found {
		f.Arguments = before
	}

	endResult := []string{}

	// Masterpiece of software engineering
	args := strings.Split(f.Arguments, ",")
	for idx, _ := range args {
		if strings.HasSuffix(args[idx], "double precision") {
			endResult = append(endResult, "double precision")
			continue
		}
		if strings.HasSuffix(args[idx], "double precision[]") {
			endResult = append(endResult, "double precision[]")
			continue
		}
		if strings.HasSuffix(args[idx], "timestamp without time zone") {
			endResult = append(endResult, "timestamp without time zone")
			continue
		}
		if strings.HasSuffix(args[idx], "time with time zone") {
			endResult = append(endResult, "time with time zone")
			continue
		}

		if strings.HasSuffix(args[idx], "timestamp with time zone") {
			endResult = append(endResult, "timestamp with time zone")
			continue
		}
		if strings.HasSuffix(args[idx], "time without time zone") {
			endResult = append(endResult, "time without time zone")
			continue
		}
		if strings.HasSuffix(args[idx], "character varying") {
			endResult = append(endResult, "character varying")
			continue
		}

		if strings.HasSuffix(args[idx], "bit varying") {
			endResult = append(endResult, "bit varying")
			continue
		}

		param := strings.Split(args[idx], " ")
		endResult = append(endResult, param[len(param)-1])
	}
	f.Arguments = strings.Join(endResult, ", ")
}

type Language struct {
	Language string `db:"language"`
}

type Schema struct {
	Name  string `db:"name"`
	Owner string `db:"owner"`
}

type Sequence struct {
	Name   string `db:"name"`
	Schema string `db:"schema"`
}

type Server struct {
	Name string `db:"name"`
	FDW  string `db:"fdw"`
}

type Table struct {
	Name   string `db:"name"`
	Schema string `db:"schema"`
}

type Tablespace struct {
	Name string `db:"name"`
}

type Type struct {
	Name   string `db:"name"`
	Schema string `db:"schema"`
}

type Role struct {
	Name                 string `db:"name"`
	SuperuserPrivilege   string `db:"superuser"`
	CreateRolePrivilege  string `db:"createrole"`
	CreateDBPrivilege    string `db:"createdb"`
	CanLogin             string `db:"canlogin"`
	ReplicationPrivilege string `db:"replication"`
}
