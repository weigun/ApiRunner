# -*- coding: utf-8 -*-

import json
import os
import sys
import traceback
import socket
import time
 

#--下面是用例配置部分
#这是公共用例列表，主要是避免编写大量的重复用例，需要用到的，在对应的地方引用就可以
common = {
	#--login的名字不能改
	"login": {
		"api": "/api/users/login",
		"method": "POST",
		"param": {
			"head":{},
			"body":{
				"username":"{{randUser}}", #--这里是随机用户名的用法
				# "username":"22222222223",
				"verifyCode" : "1111",
				"loginType": "0"
			}
		}
	}
}

caseset = {
	"name" : "demo",
	"host" : "https://www.ixbow.com/",
	"headers": {
		"Content-Type": "application/json",
		"Authorization": ""
	},
	"globalVars":{
	# 全局变量与用例导出的变量共用同一个命名空间，所以相同的变量名会出现覆盖的情况
	# 只针对header和params字段 有效
		"name" : "weigun",
		"sid" :"1",
		"lucky": "{{randRange 10 100}}"
	},
	"cases":[
		{
			"name":"login",
			"api": common["login"]["api"],
			"method":"POST",
			# "params":{
			# 	"username":"{{randUser}}",
			# 	"loginType":"0"
			# },
			"params":common["login"]["param"],
			"export":{
				#该用例导出的变量
				#如果是json格式，可以用{{body.data.token}}方式导出
				#如果是文本形式，则可以用正则表达式导出
				#如需引用变量，只要在变量名前加上.即可，如{{$token}}
				# 全局变量与用例导出的变量共用同一个命名空间，所以相同的变量名会出现覆盖的情况
				"token" : "{{.body.data.token}}",
			},
			"validate":[
				{
					"op" : "eq",
					"source": "{{.body.code}}",
					"verified":"200",
				},
				{
					"op" : "gt",
					"source": "{{.body.data.firstTime}}",
					"verified":"0",
				},
				{
					"op" : "ne",
					"source": "{{.body.data.token}}",
					"verified":"",
				},
				{
					"op" : "regx",
					"source": "{{.body.code}}}",
					"verified":"\d+",
				},
			],

		},
		{
			"name":"info",
			"api": "/api/users/info",
			"method":"GET",
			"headers":{
				"Authorization": "{{$token}}"
			},
			"validate":[
				{
					"op" : "eq",
					"source": "{{.body.code}}",
					"verified":"200",
				},
				{
					"op" : "gt",
					"source": "{{.body.data.firstTime}}",
					"verified":"10",
				},
				{
					"op" : "ne",
					"source": "{{.body.data.token}}",
					"verified":"",
				},
				{
					"op" : "regx",
					"source": "{{.body.code}}}",
					"verified":"\d+",
				},
			],

		}
	]
}

#--用例配置完毕

if __name__ == '__main__':
	try:
		json.dump(caseset,open("{}.conf".format(os.path.join(os.getcwd(),"conf",caseset["name"])),"w"),indent=4)
	except:
		traceback.print_exc()
		exit(-1)

