package mysql

import (
	"edance/util"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func init() {
	var err error
	db, err = gorm.Open(mysql.New(mysql.Config{
		DSN:                       "root:root@/db_name?charset=utf8&parseTime=True&loc=Local",
		DefaultStringSize:         256,  // default size for string fields
		DisableDatetimePrecision:  true, // disable datetime precision, witch not supported before MySQL 5.6
		DontSupportRenameIndex:    true, // drop & create when rename index, rename index not supported before MySQL 5.7, MariaDB
		DontSupportRenameColumn:   true, // `change` when rename column, rename column not supported before MuSQL 8, MariaDB
		SkipInitializeWithVersion: false,
	}), &gorm.Config{})
	if err != nil {
		util.Error(util.Wrap(err).Error())
		return
	}
}
