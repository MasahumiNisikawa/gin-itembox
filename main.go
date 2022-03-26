package main

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"

	_ "github.com/mattn/go-sqlite3"
)

type Item struct {
	gorm.Model
	Name string 
	Menge int 
	Einkauf string
	Preis int
}

func dbInit()  {
	db, err := gorm.Open("sqlite3", "itembox.db")
	if err != nil {
		panic("データベースひらけず！(dbInit)")
	}
	db.AutoMigrate(&Item{})
	defer db.Close()
}

func dbInsert(name string, menge int, einkauf string, preis int) {
	db, err := gorm.Open("sqlite3", "itembox.db")
	if err != nil {
        panic("データベース開けず！(dbInsert)")
	}
	db.Create(&Item{Name: name, Menge: menge, Einkauf: einkauf, Preis: preis})
	defer db.Close()
}

func dbUpdate(id int, name string, menge int, einkauf string, preis int)  {
	db, err := gorm.Open("sqlite3", "itembox.db")
	if err != nil {
        panic("データベース開けず！(dbUpdate)")
	}
	var item Item
	db.First(&item, id)
	item.Name = name
	item.Menge = menge
	item.Einkauf = einkauf
	item.Preis = preis
	db.Save(&item)
	db.Close()
}

func dbDelete(id int)  {
	db, err := gorm.Open("sqlite3", "itembox.db")
	if err != nil {
        panic("データベース開けず！(dbDelete)")
	}
	var item Item
	db.First(&item, id)
	db.Delete(&item)
	db.Close()
}

func dbGetAll() []Item {
	db, err := gorm.Open("sqlite3", "itembox.db")
	if err != nil {
        panic("データベース開けず！(dbGetAll)")
	}
	var items []Item
	db.Order("created_at desc").Find(&items)
	db.Close() 
	return items
}

func dbGetOne(id int) Item {
	db, err := gorm.Open("sqlite3", "itembox.db")
	if err != nil {
        panic("データベース開けず！(dbGetOne())")
	}
	var item Item
	db.First(&item, id)
	db.Close()
	return item
}

func main()  {
	router := gin.Default()
	router.Static("/assets", "./assets")
	router.LoadHTMLGlob("templates/*.html")
	// router.Static("/assets", "assets")

	dbInit()

	//Index
	router.GET("/", func(ctx *gin.Context) {
		items := dbGetAll()
		ctx.HTML(200, "index.html", gin.H{
			"items": items,
		})
	})

	//Create
	router.POST("/new", func(ctx *gin.Context)  {
		name := ctx.PostForm("name")
		menge,_ := strconv.Atoi(ctx.PostForm("menge"))
		einkauf := ctx.PostForm("einkauf")
		preis,_ := strconv.Atoi(ctx.PostForm("preis"))
		dbInsert(name, menge, einkauf, preis)
		ctx.Redirect(302, "/")
	})

	//Detail
	router.GET("/detail/:id", func(ctx *gin.Context) {
		n := ctx.Param("id")
		id, err := strconv.Atoi(n)
		if err != nil{
			panic(err)
		}
		item := dbGetOne(id)
		ctx.HTML(200, "detail.html", gin.H{"item": item})
	})

	//Update
	router.POST("/update/:id", func(ctx *gin.Context) {
		n := ctx.Param("id")
		id, err := strconv.Atoi(n)
		if err != nil {
			panic("ERROR")
		}
		name := ctx.PostForm("name")
		menge,_:= strconv.Atoi(ctx.PostForm("menge"))
		einkauf := ctx.PostForm("einkauf")
		preis,_ := strconv.Atoi(ctx.PostForm("preis"))
		dbUpdate(id, name, menge, einkauf,  preis)
		ctx.Redirect(302, "/")
	})

	//削除確認
	router.GET("/delete_check/:id", func(ctx *gin.Context) {
		n := ctx.Param("id")
		id, err := strconv.Atoi(n)
		if err != nil{
			panic("ERROR")
		}
		item := dbGetOne(id)
        ctx.HTML(200, "delete.html", gin.H{"item": item})
	})

	//Delete
	router.POST("/delete/:id", func(ctx *gin.Context) {
		n := ctx.Param("id")
		id, err := strconv.Atoi(n)
		if err != nil {
			panic("ERROR")
		}
		dbDelete(id)
		ctx.Redirect(302, "/")
	})

	router.Run()

}