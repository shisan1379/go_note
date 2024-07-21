package main

import "gin_gorm_redis_demo/router"

func main() {

	//g := gen.NewGenerator(gen.Config{
	//	//  设置输出路径
	//	OutPath: "./models",
	//	Mode:    gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface, // 选择生成模式
	//})
	////  建立数据库连接
	//gormdb, _ := gorm.Open(mysql.Open(config.Mysqldb))
	//g.UseDB(gormdb) // 选择数据库连接
	//
	//// 为结构模型生成基本类型安全的 DAO API。用户的以下约定
	//
	//g.ApplyBasic(
	//	//基于表`users`生成结构`User``
	//	g.GenerateModel("user"),
	//
	//	// 基于表`users`生成结构`Employee``
	//	g.GenerateModelAs("user", "Employee"),
	//
	//	// 基于表`users'生成结构`User`并生成选项
	//	g.GenerateModel("user", gen.FieldIgnore("address"), gen.FieldType("id", "int64")),
	//)
	//g.ApplyBasic(
	//	// 从当前数据库的所有表生成结构
	//	g.GenerateAllTable()...,
	//)
	//// 生成代码
	//g.Execute()

	r := router.Router()

	r.Run(":9999")
}
