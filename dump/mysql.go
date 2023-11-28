package dump

import (
	"database/sql"
	"fmt"

	"github.com/aliakseiz/go-mysqldump"
	"github.com/go-sql-driver/mysql"
	"github.com/pluswing/datasync/data"
	"github.com/spf13/cobra"
	"github.com/tanimutomo/sqlfile"
)

func Dump(dumpDir string, cfg data.TargetMysqlType) {

	config := mysql.NewConfig()
	config.User = cfg.User
	config.Passwd = cfg.Password
	config.DBName = cfg.Database
	config.Net = "tcp"
	config.Addr = fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)

	dumpFilenameFormat := fmt.Sprintf("%s-%s", "mysql", cfg.Database)

	db, err := sql.Open("mysql", config.FormatDSN())
	cobra.CheckErr(err)

	dumper, err := mysqldump.Register(db, dumpDir, dumpFilenameFormat, config.DBName)
	cobra.CheckErr(err)

	err = dumper.Dump()
	cobra.CheckErr(err)

	dumper.Close()
}

func Import(dumpFile string, cfg data.TargetMysqlType) {
	config := mysql.NewConfig()
	config.User = cfg.User
	config.Passwd = cfg.Password
	config.DBName = cfg.Database
	config.Net = "tcp"
	config.Addr = fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)

	db, err := sql.Open("mysql", config.FormatDSN())
	cobra.CheckErr(err)

	s := sqlfile.New()
	err = s.File(dumpFile)
	cobra.CheckErr(err)

	_, err = s.Exec(db)
	cobra.CheckErr(err)
}
