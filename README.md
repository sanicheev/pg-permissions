# pg-permissions

## Description
The purpose of this tool is to gather and display permissions for postgres objects for all users.
It tries to solve a problem when someone can grant granular permission to an object(e.g: column) which is not then visible/detected until you explicitly execute the function(e.g: has_column_privilege).
Taking into account postgres has lots of objects in its possession, humans cannot possibly check all of its permissions from time to time and make audit decisions.
Thus such routing must be automated and a report must be generated.

## How it works
For each user pg-permissions will gather permissions(read-only or read-write) information about following objects:
| Object Type | Read Permissions | Write Permissions| Check function |
|-------------|------------------|------------------|----------------|
| Column      | SELECT           | INSERT, UPDATE, REFERENCES | has_column_privilege(user, table, column, privilege) |
| Database    | CONNECT          | CREATE, TEMPORARY | has_database_privilege(user, database, privilege) |
| FDW | USAGE | N/A | has_foreign_data_wrapper_privilege(user, fdw, privilege) |
| Function | EXECUTE | N/A | has_function_privilege(user, function, privilege) |
| Language | USAGE | N/A | has_language_privilege(user, language, privilege) |
| Role | CANLOGIN, REPLICATION | SUPERUSER, CREATEROLE, CREATEDB | pg_has_role(user, role, usage), pg_has_role(user, role, member) |
| Schema | USAGE | CREATE | has_schema_privilege(usage, schema, privilege) |
| Sequence | USAGE, SELECT | UPDATE | has_sequence_privilege(user, sequence, privilege) |
| Server | USAGE | N/A | has_server_privilege(user, server, privilege) |
| Table | SELECT | INSERT, UPDATE, DELETE, REFERENCES, TRIGGER, TRUNCATE | has_table_privilege(user, table, privilege) |
| Tablespace | N/A | CREATE | has_tablespace_privilege(user, tablespace, privilege) |
| Type | USAGE | N/A | has_type_privilege(user, type, privilege) |

For each of the objects it will execute a function in a column: 'Check function'.
It will try its best to parallelize the work in several goroutines to make it faster.
In the end it will generate a report in json format.
Additional flags can be set in the config file and it will also generate reports in html format.

## Build & Usage Instructions
In order to use this tool you need to compile it first by running:
```
go build
```

Afterwards you can execute it by running:
```
./pg_permission postgres
```
Please note that in order to get the full picture it's best to run this tool as a user with superuser privilege.

Please note that you must create the correct configuration beforehand. Example:
```
database:
  host: 127.0.0.1
  name: mydatabase
  user: myusername
  password: 12345
  permissions: W # if you do not specify the W flag it will fetch all possible permissions even for objects which have only Read permissions(e.g: function) thus increasing its runtime dramatically.
```

Runtime for the objects which has write permissions when 'W' is set in configuration file: ~5-8 minutes
Runtime for ALL objects which have read and write permissions when 'W' is omitted in configuration file: ~20 minutes.
In the latter case runtime increases because it will also try to process functions and types for each user. And by default there are a lot of these objects.

## HTML report example
![rdsadmin](/images/rdsadmin.png?raw=true "")
![rdsrepladmin](/images/rdsrepladmin.png?raw=true "")

# Authors & Contributors:
Serghei Anicheev - serghei.anicheev@gmail.com
