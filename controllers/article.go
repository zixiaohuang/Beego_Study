package controllers

import (
	"beego_Study/models"
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"math"
	"path"
	"time"
)

type ArticleController struct {
	beego.Controller
}

// 显示文章列表页
func (c* ArticleController) ShowArticleList(){
	// 1.查到有多少条数据
	// 2.共几页
	// 3.首页和末页
	// 4.上一页和下一页
	o := orm.NewOrm()
	//创建文章表查询器，但不查询
	qs := o.QueryTable("Article")
	var articles []models.Article  //qs.All(&articles) //select * from article

	//获取类型数据
	var articletypes []models.ArticleType
	o.QueryTable("ArticleType").All(&articletypes)
	var buffer bytes.Buffer
	enc := gob.NewEncoder(&buffer)
	enc.Encode(articletypes)


	//count, err := qs.Count() // 返回数据条目数
	//定义每页大小，即本次请求的条数
	pageSize := 6

	//获取本次查询的页码
	pageIndex, err := c.GetInt("pageIndex")
	if err != nil {
		//若未获取到页码，设置默认页码1
		pageIndex = 1
	}
	searchStart := pageSize * (pageIndex - 1)

	//使用文章查询器，简单获得记录总数
	count, err := qs.RelatedSel("ArticleType").Count()
	if err != nil {
		fmt.Println("获取记录数错误：", err)
		return
	}
	//根据查询头和查询量，开始查询数据
	//参数1：限制获取的条数，参数2，偏移量，即开始位置
	//qs.Limit(pageSize, searchStart).RelatedSel("ArticleType").All(&articles)

	// 加入文章类型筛选，默认全部，选择类型后，再次筛选
	selectedtype := c.GetString("select")
	if selectedtype == "" || selectedtype == "全部类型" {
		fmt.Println("本次GET请求全部，未加入select参数，默认全部")
	} else {
		count, err = qs.RelatedSel("ArticleType").Filter("ArticleType__TypeName", selectedtype).Count()
		if err != nil {
			fmt.Println("获取记录数错误: ", err)
			return
		}
		fmt.Printf("请求类型：%v count: %v",selectedtype, count)
		// 分页取
		// 好处：1.防止一次性读取太多数据到内存，导致卡顿，2. 减少页面条目数，提高网页浏览速度
		// Filter过滤器函数 相当于sql语句中的where，指定查询条件 参数1：查询字段"表名__表字段" 参数2：要匹配的值
		// 多表查询的时候，orm使用的是惰性查询，需要使用RelatedSel关联表
		_, err = qs.Limit(pageSize, searchStart).RelatedSel("ArticleType").Filter("ArticleType__TypeName", selectedtype).All(&articles)// 1.pageSize 一页显示多少 2.start 起始位置
		//_, err = qs.All(&articles) // select * from Article
		if err != nil {
			logs.Error("查询文章信息出错")
			return
		}
	}

	// 得出总页
	pageCount := int(math.Ceil(float64(count) / float64(pageSize))) // 向上取整
	if err != nil {
		logs.Error("查询所有文章条数出错")
		return
	}
	logs.Info(articles)
	logs.Info(count)

	// 定义页码按钮启动状态
	enablelast, enablenext := true, true
	if pageIndex == 1 {
		enablelast = false
	}
	if pageIndex == pageCount {
		enablenext = false
	}

	fmt.Println(articletypes)
	c.Data["typename"] = selectedtype
	c.Data["articletypes"] = articletypes
	c.Data["articles"] = articles
	c.Data["count"] = count
	c.Data["EnableNext"] = enablenext
	c.Data["EnableLast"] = enablelast
	c.Data["pageCount"] = pageCount
	c.Data["pageIndex"] = pageIndex

	c.TplName = "index.html"
}

func (c* ArticleController) HandleTypeSelected() {
	selectedtype := c.GetString("select")
	articles := []models.Article{}
	o := orm.NewOrm()
	fmt.Println("Post select type", selectedtype)
	o.QueryTable("article").RelatedSel("ArticleType").Filter("ArticleType__TypeName", selectedtype).All(&articles)
	c.Data["articles"] = articles

	// 文章类型下拉
	articletypes := []models.ArticleType{}
	o.QueryTable("article_type").All(&articletypes)
	c.Data["articletypes"] = articletypes
	c.Data["username"] = c.GetSession("username")
	c.TplName = "index.html"
}


// 显示添加文章界面
func (c* ArticleController) ShowAddArticle() {
	o := orm.NewOrm()
	articletypes := []models.ArticleType{}
	cnt, err := o.QueryTable("article_type").All(&articletypes)
	if err != nil {
		logs.Error("查询文章类型失败")
	}
	if cnt > 0 {
		//logs.Info("查询文章类型少于零")
		c.Data["articletypes"] = articletypes
	}
	//c.Data["username"] = c.GetSession("username")
	c.TplName = "add.html"
}

// 处理上传的图片
func (c* ArticleController) HandleUpatePic() (filename string){
	f, h , err:= c.GetFile("uploadname")// 文件流（注意关闭，不让会内存泄漏） 文件头
	if err != nil {
		logs.Error("上传文件失败")
		return
	}else {
		/*保存之前先做校验处理:
		1.校验文件类型
		2.校验文件大小
		3.防止重名，重新命名
		*/
		defer f.Close()
		ext := path.Ext(h.Filename)
		fmt.Println(ext)
		if ext != ".jpg" && ext != ".png" && ext != ".jpeg" {
			logs.Error("上传文件类型错误")
			return
		}

		if h.Size > 5000000 {
			logs.Error("上传文件超出大小 ")
			return
		}
		filename = time.Now().Format("20060102150405") + ext // 名字一样可能会覆盖 format要用go语言诞生的时间，否则转换不正确

		c.SaveToFile("uploadname", "./static/img/" + filename) // 文件存储
		if err != nil {
			logs.Error("文件保存失败：", err)
			return
		}
	}
	return
}


// 处理添加文章界面数据
func (c* ArticleController) HandleAddArticle() {
	// 1.拿到数据
	articleName := c.GetString("articleName")
	articleContent := c.GetString("content")

	//logs.Info(articleName, articleContent)
	filename := c.HandleUpatePic()
	// 2.判断数据是否合法
	if articleContent == "" || articleName == "" {
		logs.Info("添加文章数据错误")
		return
	}
	// 3.插入数据
	o := orm.NewOrm()
	//取得文章类型
	selectedtype := c.GetString("select") // 获取传过来的数据
	//利用此类型获取完整对象
	articletype := models.ArticleType{TypeName: selectedtype}
	o.Read(&articletype, "TypeName")

	logs.Info("articleType id", articletype.Id)
	article := models.Article{Title: articleName, Content: articleContent, ArticleType: &articletype}
	if filename != "" {
		article.Img = "/static/img/" + filename // 存的时候可以不用点
	}

	_, err := o.Insert(&article)
	if err != nil {
		logs.Error("插入数据库错误")
		return
	}
	// 4.返回文章界面
	c.Redirect("/Article/ShowArticle", 302)
}


// 内容详情页
func (c* ArticleController) ShowContent() {
	// 1.获取文章ID
	id, err := c.GetInt("id")
	if err != nil {
		logs.Error("获取文章Id错误")
		return
	}
	logs.Info("id = ", id)
	// 2.查询数据库获取数据
	o := orm.NewOrm()
	article := models.Article{Id:id}
	err = o.Read(&article)
	if err != nil {
		logs.Error("查询文章失败")
		return
	}
	//阅读量+1并写回数据库
	article.Count++
	o.Update(&article)

	/*处理最近浏览,
	1. 首先需确定当前浏览者登录状态,获取浏览者信息
	2. 将浏览者信息插入数据表
	3. 将历史浏览者信息从表中读出，去重，显示*/

	// 3.传递数据给视图
	c.Data["content"] = article

	c.TplName = "content.html"
}

// 显示编辑界面
func (c* ArticleController) ShowUpdate() {
	/*思路
	1. 获取数据，填充数据
	2. 更新数据，更新数据库，返回列表页
	*/
	c.TplName = "update.html"
	id, err := c.GetInt("id")
	if err != nil {
		logs.Error("id获取失败", err)
		return
	}
	article := models.Article{Id: id}
	o := orm.NewOrm()
	err = o.ReadForUpdate(&article)
	if err != nil {
		logs.Error("数据库读取失败", err)
		return
	}
	c.Data["article"] = article
	//c.Data[""]
}

// 处理文章更新业务
func (c* ArticleController) HandleUpdate() {
	c.TplName = "update.html"
	//取得post数据，使用getfile取得文件，注意设置enctype
	name := c.GetString("articleName")
	content := c.GetString("content")
	oldimagepath := c.GetString("oldimagepath")

	id, err := c.GetInt("id")
	if err != nil {
		logs.Error("id获取失败", err)
		return
	}
	if name == "" || content == "" {
		logs.Info("数据类型错误, 更新失败")
		return
	}

	article := models.Article{Id: id, Title: name, Content: content, Img: oldimagepath}
	filename := c.HandleUpatePic()
	//若上传了新文件，则使用新文件路径，否则使用旧路径不变
	if filename != "" {
		article.Img = "/static/img/" + filename
	}

	// 更新数据库
	o := orm.NewOrm()
	_, err = o.Update(&article, "title", "content", "img", "create_time", "update_time")
	if err != nil {
		logs.Error("数据库更新失败", err)
		c.Data["errmsg"] = "更新失败"
		return
	}
	c.Redirect("/Article/ShowArticle", 302)
}

// 删除操作
func (c* ArticleController) HandleDelete() {
	/*思路
	1.被点击的url传值
	2.执行对应的删除操作
	*/
	id, err := c.GetInt("id")
	if err != nil {
		logs.Error("获取id失败", err)
		return
	}
	o := orm.NewOrm()
	article := models.Article{Id: id}
	err = o.Read(&article)
	if err != nil {
		logs.Info("查询数据错误", err)
		return
	}
	o.Delete(&article)
	c.Redirect("/Article/ShowArticle", 302)
}

func (c* ArticleController) ShowAddType() {
	// 读取类型表，显示数据
	c.TplName = "addType.html"
	var types []models.ArticleType
	o := orm.NewOrm()
	_, err := o.QueryTable("article_type").All(&types)
	if err != nil {
		logs.Error("查询类型错误")
		return
	}
	c.Data["types"] = types

}


func (c *ArticleController) HandleAddType() {
	var articleType models.ArticleType
	// 获取数据, 判断数据
	if articleType.TypeName = c.GetString("typeName"); articleType.TypeName == "" {
		fmt.Println("类型不能为空")
		c.Redirect("/Article/AddArticleType", 302)
		return
	}
	fmt.Println("您输入的类型名为", articleType.Id, articleType.TypeName)
	o := orm.NewOrm()
	// 执行插入操作
	_, err := o.Insert(&articleType)
	if err != nil {
		fmt.Println("插入数据失败", err)
		return
	}
	// 展示视图
	c.Redirect("/Article/AddArticleType", 302)
	//插入数据库成功后，此处不更新缓存，否则需要再次请求所有类型，刷新页面时更新更合适。
}


