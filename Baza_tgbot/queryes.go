package main

var QaddUser = "INSERT INTO users (tgid, username) VALUES ('{id}', '{username}');"
var QaddRecord = "INSERT INTO content (text, author, tgid) VALUES ('{text}', '{author}', '{tgid}');"
var QgetAll = "SELECT * FROM {table}"
var QgetRandom = "SELECT * FROM {table} ORDER BY RANDOM() LIMIT 1;"
var Qcounter = "SELECT COUNT(*) FROM {table};"
var Qexist = "SELECT EXISTS(SELECT 1 FROM users WHERE tgid = {id});"
var QgetById = "SELECT * FROM content WHERE tgid = '{tgid}';"
var QdelById = "DELETE FROM content WHERE id = {recid} AND tgid = '{tgid}';"
