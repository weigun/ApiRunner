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
		"token" : "{%^&*",
		"sid" :1
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

		}
	]
}

#--用例配置完毕

if __name__ == '__main__':
	try:
		json.dump(caseset,open("{}.conf".format(caseset["name"]),"w"))
	except:
		traceback.print_exc()
		exit(-1)

