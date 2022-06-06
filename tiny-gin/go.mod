module my-gin

go 1.18

replace koo => ./koo

//从 go 1.11 版本开始，引用相对路径的 package 需要使用上述方式。

require koo v0.0.0-00010101000000-000000000000
