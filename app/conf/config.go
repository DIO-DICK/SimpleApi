package conf

import (
        "fmt"
        "github.com/BurntSushi/toml"
        "os"
        "path"
)

type mysql struct {
        DB string
        DB_ip string
        DB_port string
        DB_user string
        DB_password string
        DB_charset string
}

type TomlConfig struct {
        Log_file string
        Mysql mysql
}

var conf TomlConfig

func init() {
        dir, err := os.Getwd()
        if err != nil {
                fmt.Println(err)
        }
        dirname := path.Join(dir, "/conf/config.toml")

        if _,err :=toml.DecodeFile(dirname, &conf); err != nil {
                fmt.Println(err)
                return
        }
}

func Get() TomlConfig {
        return conf
}