package com

import (
  _ "github.com/go-sql-driver/mysql"
  xormCore "github.com/go-xorm/core"
  "github.com/go-xorm/xorm"
  "fmt"
  "os"
)

func ModelConnetMysql(groupName string) *xorm.Engine {
  var err error
  if groupName != "" {
    groupName = groupName + "::"
  }
  dbHost := Conf.String(groupName + "host")
  dbDataBase := Conf.String(groupName + "database")
  dbName := Conf.String(groupName + "user")
  dbPwd := Conf.String(groupName + "pwd")
  engine, err := xorm.NewEngine("mysql", dbName + ":" + dbPwd + "@tcp(" + dbHost + ")/" + dbDataBase + "?charset=utf8mb4")
  if err != nil {
    Logger.Error("%s", err.Error())
    os.Exit(1)
  }

  if Conf.String("runmode") != "pro" {
    engine.ShowSQL(true)
    engine.Logger().SetLevel(xormCore.LOG_DEBUG)
  }
  return engine
}

//CHARACTER
//COLLATE
//ALTER DATABASE `ecbox_zuqiuzhandui` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci
//
//SELECT CONCAT('ALTER TABLE ', table_name, ' CONVERT TO CHARACTER SET  utf8 COLLATE utf8_unicode_ci;')
//FROM information_schema.TABLES
//WHERE TABLE_SCHEMA = 'ecbox_zuqiuzhandui'
//
//SELECT CONCAT('ALTER TABLE `', table_name, '` MODIFY `', column_name, '` ', DATA_TYPE, '(', CHARACTER_MAXIMUM_LENGTH, ') CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci', (CASE WHEN IS_NULLABLE = 'NO' THEN ' NOT NULL' ELSE '' END), ';')
//FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'ecbox_zuqiuzhandui';

func AlterDbCharacter(engine *xorm.Engine, dbName, character, collate string) {

  _, err := engine.Exec(fmt.Sprintf("ALTER DATABASE `%s` DEFAULT CHARACTER SET %s COLLATE %s", dbName, character, collate))
  if err == nil {
    //get table
    queryList, err := engine.Query(fmt.Sprintf(`
    SELECT CONCAT('ALTER TABLE ', table_name, ' CONVERT TO CHARACTER SET  %s COLLATE %s;')
FROM information_schema.TABLES
WHERE TABLE_SCHEMA = '%s'
    `, character, collate, dbName))

    if err != nil {
      fmt.Println(err)
      return
    }

    for _, v := range queryList {
      for _, vv := range v {
        _, err := engine.Exec(string(vv))
        if err != nil {
          fmt.Println(err)
        }
      }
    }

    queryList, err = engine.Query(fmt.Sprintf("" +
      "SELECT CONCAT('ALTER TABLE `', table_name, '` MODIFY `', column_name, '` ', DATA_TYPE, '(', CHARACTER_MAXIMUM_LENGTH, ') " +
      "CHARACTER SET %s COLLATE %s', (CASE WHEN IS_NULLABLE = 'NO' THEN ' NOT NULL' ELSE '' END), ';')" +
      "FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = '%s';", character, collate, dbName))

    if err != nil {
      fmt.Println(err)
      return
    }

    for _, v := range queryList {
      for _, vv := range v {
        _, err := engine.Exec(string(vv))
        if err != nil {
          fmt.Println(err)
        }
      }
    }

  } else {
    fmt.Println(err)
  }
}