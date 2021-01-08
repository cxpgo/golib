package test

import (
	"fmt"
	"github.com/cxpgo/golib/lib"
	"testing"
	"time"
)

type Test2 struct {
	Id        int64     `json:"id" gorm:"primary_key"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}
var (
	createTableSQL2 = "CREATE TABLE `test2` (`id` int(12) unsigned NOT NULL AUTO_INCREMENT" +
		" COMMENT '自增id',`name` varchar(255) NOT NULL DEFAULT '' COMMENT '姓名'," +
		"`created_at` datetime NOT NULL,PRIMARY KEY (`id`)) ENGINE=InnoDB " +
		"DEFAULT CHARSET=utf8"
	dropTableSQL2 = "DROP TABLE `test2`"
)


func Test_GORM(t *testing.T) {
	testInitOnce()
	lib.InitGormPool(lib.GConfig.MySqlConfList)
	//获取链接池
	dbpool, err := lib.GetGormPool("default")
	if err != nil {
		t.Fatal(err)
	}
	db := dbpool.Begin()
	//traceCtx := lib.NewTrace()


	//设置trace信息
	//db = db.SetCtx(traceCtx)
	if err := db.Exec(createTableSQL2).Error; err != nil {
		db.Rollback()
		t.Fatal(err)
	}

	//插入数据
	t1 := &Test2{Name: "test_name1", CreatedAt: time.Now()}
	t2 := &Test2{Name: "test_name2", CreatedAt: time.Now()}
	t3 := &Test2{Name: "test_name3", CreatedAt: time.Now()}
	tList := []*Test2{t1,t2,t3}
	if err := db.Save(tList).Error; err != nil {
		db.Rollback()
		t.Fatal(err)
	}

	//查询数据
	list := []Test2{}
	if err := db.Where("name=?", "test_name1").Find(&list).Error; err != nil {
		db.Rollback()
		t.Fatal(err)
	}

	fmt.Printf("list==%+v\n",list)

	// 原生 SQL
	rows, err := db.Raw("select id, name, created_at from test2 ").Rows()
	if err != nil {
		db.Rollback()
		t.Fatal(err)
	}

	defer rows.Close()
	var myType Test2
	for rows.Next() {
		//rows.Scan(&name, &age, &email)
		rows.Scan(&myType.Id,&myType.Name,&myType.CreatedAt)
		fmt.Printf("mytype=%+v\n",myType)
		// 业务逻辑...
	}



	//删除表数据
	if err := db.Exec(dropTableSQL2).Error; err != nil {
		db.Rollback()
		t.Fatal(err)
	}
	db.Commit()

	//Close()
}
