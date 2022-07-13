设计 rom 考虑的主要矛盾 

如果使用 sql 语句

```sql
CREATE TABLE `User` (`Name` text, `Age` integer);
INSERT INTO `User` (`Name`, `Age`) VALUES ("fengwei", 18);
SELECT * FROM `User`;
```

使用 orm 之后，应该是 定义一个 User 之后，使用 orm 中的方法将 struct 映射成一个 table 
使用 save 方法实现 insert 使用 find 进行查找 

```go
type User struct {
    Name string
    Age  int
}

orm.CreateTable(&User{})
orm.Save(&User{"Tom", 18})
var users []User
orm.Find(&users)
```

orm 框架相当于对象和数据库中间的一个桥梁，借助 orm 框架可以避免繁琐的 sql 语言，只要通过操作具体的对象，就可以完成对关系型数据库的操作 

- mysql sqlite 等数据库的 sql 语句是存在区别的，orm 框架如何在开发者没有感知的情况下适配多种数据库 
- 如果对象的字段发生改变，数据表的结构能够自动更新，是否支持自动的 migrate (django 就不支持) 
- 数据库支持很多功能，例如 事务，orm 框架能够实现哪些

以 sqlite 举例 
创建一个表格 并插入测试数据 
```sql
sqlite> CREATE TABLE User(Name text, Age integer);
sqlite> .schema User
CREATE TABLE User(Name text, Age integer);
sqlite> .table
User
    
INSERT INTO User(Name, Age) VALUES ("Tom", 18), ("Jack", 25);
```

- log 是支持分级的 log 库 
- 核心结构 session 实现和数据库的交互 
- 
