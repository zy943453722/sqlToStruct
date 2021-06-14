# sqlToStruct
sqlToStruct tools
sql语句转struct工具

通过读取mysql的COLUMNS库，得到对应列后，再利用go原生的template填充成struct。其中命令行用到的是cobra包

用法：
go run main.go sql struct --username 用户名 --password 密码 --db=数据库名 --table 表名
