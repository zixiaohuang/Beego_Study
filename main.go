package main

import (
	_ "beego_Study/routers"
	_ "beego_Study/models" //运行的时候会先去掉该package的init函数
	"github.com/astaxie/beego"
)

func main() {
	beego.Run()
}

