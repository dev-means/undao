package undao

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AggregateLookup struct {
	From         string
	LocalField   string
	ForeignField string
	As           string
	Project      bson.M
}

func CheckOID(coll *mongo.Collection, id string) (oid primitive.ObjectID, c int, m string, b bool) {
	/* example
	oid, c, m, b := mdao.CheckOID(coll, "ObjectID")
	if !b {
		// return c, m
	} else {
		// call: Get()
	}
	*/
	var e error
	oid, e = primitive.ObjectIDFromHex(id)
	if e != nil {
		c, m, b = 1103, "ID格式非法", false
	} else {
		b = true
		if total, _ := coll.CountDocuments(context.Background(), bson.M{"_id": oid}); total != 1 {
			c, m, b = 1104, "ID不存在", false
		}
	}
	return
}

func Update(coll *mongo.Collection, ctx context.Context, match, update interface{}) (e error) {
	/* example
	if e := mdao.Update(coll, ctx, bson.M{<match>}, bson.M{<update>}}); e != nil {
		panic(e)
	}
	*/
	if match != nil {
		_, e = coll.UpdateMany(ctx, match, update)
	}
	return
}

func Add(coll *mongo.Collection, ctx context.Context, data interface{}) (e error) {
	/* example
	if e := mdao.Add(coll, ctx, data); e != nil {
		panic(e)
	}
	*/
	if data != nil {
		_, e = coll.InsertOne(ctx, data)
	}
	return
}

func Del(coll *mongo.Collection, ctx context.Context, match interface{}) (e error) {
	/* example
	if e := mdao.Del(coll, ctx, bson.M{<match>}); e != nil {
		panic(e)
	}
	*/
	if match != nil {
		_, e = coll.DeleteMany(ctx, match)
	}
	return
}

func Len(coll *mongo.Collection, match interface{}) (total int64) {
	/* example
	total := mdao.Len(coll, bson.M{<match>})
	*/
	if match == nil {
		match = bson.M{}
	}
	total, _ = coll.CountDocuments(context.Background(), match)
	return
}

func Get(coll *mongo.Collection, match interface{}, opts ...*options.FindOneOptions) (single *mongo.SingleResult) {
	/* example
	var op <Your struct>
	if e := mdao.Get(coll, bson.M{"_id": ObjectID}).Decode(&op); e != nil {
		panic(e)
	}
	*/
	return coll.FindOne(context.Background(), match, opts...)
}

func GetList(coll *mongo.Collection, limit, page int,
	match, sort, project interface{}, lookups ...AggregateLookup) (ctx context.Context, cur *mongo.Cursor, total int64, e error) {
	/* example
	ctx, cur, total, e := mdao.GetList(coll, <Limit>, <Page>,
		bson.M{<match>},
		bson.D{
				{Key: "_id", Value: -1},
			},
		bson.M{<project>},
	)
	if e != nil {
		panic(e)
	} else {
		defer cur.Close(ctx)
	}
	var list = make([]<Your struct>, 0)
	if e = cur.All(ctx, &list); e != nil {
		panic(e)
	}
	*/
	if limit == -1 {
		limit = 1<<63 - 1
	}
	if match == nil {
		match = bson.M{}
	}
	if sort == nil {
		sort = bson.M{"_id": -1}
	}
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: match}},
		{{Key: "$sort", Value: sort}},
		{{Key: "$skip", Value: limit * (page - 1)}},
		{{Key: "$limit", Value: limit}},
	}
	if project != nil {
		pipeline = append(pipeline, mongo.Pipeline{
			{{Key: "$project", Value: project}},
		}...)
	}
	for _, al := range lookups {
		pipeline = append(pipeline, mongo.Pipeline{
			{{Key: "$lookup", Value: bson.M{
				"from": al.From,
				"let":  bson.M{"id": "$" + al.LocalField},
				"pipeline": mongo.Pipeline{
					{{Key: "$match", Value: bson.M{"$expr": bson.M{"$eq": []string{"$" + al.ForeignField, "$$id"}}}}},
					{{Key: "$project", Value: al.Project}},
				},
				"as": al.As,
			}}},
		}...)
	}
	cur, e = coll.Aggregate(ctx, pipeline)
	if e != nil {
		return
	}
	total, e = coll.CountDocuments(ctx, match)
	if e != nil {
		defer cur.Close(ctx)
		return
	}
	return
}
