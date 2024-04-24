package main

var QaddUser = "INSERT INTO users (tgid, username) VALUES ('{id}', '{username}');"
var QaddRecord = "INSERT INTO content (text, author) VALUES ('{text}', '{author}');"
var QgetAll = "SELECT * FROM {table}"
var QgetRandom = "SELECT * FROM {table} ORDER BY RANDOM() LIMIT 1;"
var Qcounter = "SELECT COUNT(*) FROM {table};"
var Qexist = "SELECT EXISTS(SELECT 1 FROM users WHERE tgid = {id});"
