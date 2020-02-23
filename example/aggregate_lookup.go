package example

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dev-means/undao"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//----------------------------------------------------------------------------------------------------------------------

// 1. 连接数据库，创建集合句柄

var DB = undao.NewStorageDatabase(false, "127.0.0.1:28000", "demo", "", "", "")

var (
	coll_User = DB.Collection("users")
	coll_Post = DB.Collection("posts")
)

//----------------------------------------------------------------------------------------------------------------------

// 2. 构建用于聚合查询的演示数据模型

type User struct {
	UserId primitive.ObjectID `bson:"_id" json:"userId"`
	Name   string             `bson:"name" json:"name"`
}

// 关联查询容器
type PostLookup struct {
	User []User `bson:"user" json:"user"`
}

// 关联查询引用字段
type PostAs struct {
	UserId primitive.ObjectID `bson:"userId" json:"userId"`
}

type Post struct {
	PostId primitive.ObjectID `bson:"_id" json:"postId"`
	Lookup PostLookup         `bson:"lookup" json:"lookup"`
	As     PostAs             `bson:"as" json:"as"`
	Title  string             `bson:"title" json:"title"`
	Num    int                `bson:"num" json:"num"`
}

//----------------------------------------------------------------------------------------------------------------------

// 3. 添加测试数据

func insertTestSet() {
	for i := 0; i < 10; i++ {
		it := fmt.Sprintf("%d", i+1)

		u := User{
			UserId: primitive.NewObjectID(),
			Name:   it,
		}

		p := Post{
			PostId: primitive.NewObjectID(),
			As:     PostAs{UserId: u.UserId},
			Title:  it,
			Num:    i + 1,
		}

		if e := undao.Add(coll_User, context.Background(), &u); e != nil {
			panic(e)
		}
		if e := undao.Add(coll_Post, context.Background(), &p); e != nil {
			panic(e)
		}
	}
}

//----------------------------------------------------------------------------------------------------------------------

// 4. 使用 undao.GetList() 作聚合管道查询

func find() {
	var lookup = undao.AggregateLookup{
		From:         "users",
		LocalField:   "as.userId",
		ForeignField: "_id",
		As:           "lookup.user",
		Project:      bson.M{"name": 1},
	}
	ctx, cur, total, e := undao.GetList(coll_Post, 3, 1,
		bson.M{},
		bson.D{
			{Key: "_id", Value: -1},
		},
		nil,
		lookup,
	)
	if e != nil {
		panic(e)
	} else {
		defer cur.Close(ctx)
	}
	var list = make([]Post, 0)
	if e = cur.All(ctx, &list); e != nil {
		panic(e)
	}
	js, err := json.MarshalIndent(list, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("total:%d\nlistLen:%d\n%s\n", total, len(list), js)
}

//----------------------------------------------------------------------------------------------------------------------

// 5. 执行程序

func main() {
	insertTestSet()
	find()
}

//----------------------------------------------------------------------------------------------------------------------

// 6. 输出结果

/*
total:10
listLen:3
{
  "postId": "5e4e35b15488a7e05c568c40",
  "lookup": {
    "user": [
      {
        "userId": "5e4e35b15488a7e05c568c3f",
        "name": "10"
      }
    ]
  },
  "as": {
    "userId": "5e4e35b15488a7e05c568c3f"
  },
  "title": "10",
  "num": 10
}
*/
