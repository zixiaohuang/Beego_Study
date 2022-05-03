package main

import (
	_ "beego_Study/models" //运行的时候会先去掉该package的init函数
	_ "beego_Study/routers"
	"github.com/astaxie/beego"
)

func main() {
	// 要在Run前映射
	beego.AddFuncMap("ShowNextPage", getNextPage)
	beego.AddFuncMap("ShowLastPage", getLastPage)
	beego.Run()
}

// 1. 在html定义视图函数
// 2. 后台添加代码
// 3. 建立两个名字映射
func getNextPage(pageindex int) int {
	pageindex++
	return pageindex
}
func getLastPage(pageindex int) int {
	if pageindex--; pageindex < 0 {
		pageindex = 0
	}
	return pageindex
}

