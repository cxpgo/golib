package lib

import (
	"fmt"
	"testing"
)

var (
	createTableSQL = "CREATE TABLE `test1` (`id` int(12) unsigned NOT NULL AUTO_INCREMENT" +
		" COMMENT '自增id',`name` varchar(255) NOT NULL DEFAULT '' COMMENT '姓名'," +
		"`created_at` datetime NOT NULL,PRIMARY KEY (`id`)) ENGINE=InnoDB " +
		"DEFAULT CHARSET=utf8"
	insertSQL    = "INSERT INTO `test1` (`id`, `name`, `created_at`) VALUES (NULL, '111', '2018-08-29 11:01:43');"
	dropTableSQL = "DROP TABLE `test1`"
	beginSQL     = "start transaction;"
	commitSQL    = "commit;"
	rollbackSQL  = "rollback;"
)

func Test_DBPool(t *testing.T) {
	//os.Chdir("..")
	//fmt.Println(os.Getwd())
	testInitOnce()
	InitDBPool(gConfig.MySqlConfList)
	//获取链接池
	dbpool, err := GetDBPool("default")
	if err != nil {
		t.Fatal(err)
	}
	//开始事务
	trace := NewTrace()
	if _, err := DBPoolLogQuery(trace, dbpool, beginSQL); err != nil {
		t.Fatal(err)
	}

	//创建表
	if _, err := DBPoolLogQuery(trace, dbpool, createTableSQL); err != nil {
		Log.Printf("[INFO] %s\n", " 创建表")
		DBPoolLogQuery(trace, dbpool, rollbackSQL)
		t.Fatal(err)
	}

	//插入数据
	if _, err := DBPoolLogQuery(trace, dbpool, insertSQL); err != nil {
		DBPoolLogQuery(trace, dbpool, rollbackSQL)
		t.Fatal(err)
	}

	//循环查询数据
	current_id := 0
	table_name := "test1"
	fmt.Println("begin read table ", table_name, "")
	fmt.Println("------------------------------------------------------------------------")
	fmt.Printf("%6s | %6s\n", "id", "created_at")
	for {
		rows, err := DBPoolLogQuery(trace, dbpool, "SELECT id,created_at FROM test1 WHERE id>? order by id asc", current_id)
		defer rows.Close()
		row_len := 0
		if err != nil {
			DBPoolLogQuery(trace, dbpool, "rollback;")
			t.Fatal(err)
		}
		for rows.Next() {
			var create_time string
			if err := rows.Scan(&current_id, &create_time); err != nil {
				DBPoolLogQuery(trace, dbpool, "rollback;")
				t.Fatal(err)
			}
			fmt.Printf("%6d | %6s\n", current_id, create_time)
			row_len++
		}
		if row_len == 0 {
			break
		}
	}
	fmt.Println("------------------------------------------------------------------------")
	fmt.Println("finish read table ", table_name, "")

	//删除表
	if _, err := DBPoolLogQuery(trace, dbpool, dropTableSQL); err != nil {
		DBPoolLogQuery(trace, dbpool, rollbackSQL)
		t.Fatal(err)
	}

	//提交事务
	DBPoolLogQuery(trace, dbpool, commitSQL)

	//Close()

}