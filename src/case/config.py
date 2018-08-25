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
	"host" : host,
	"headers" : {
		"auth" : "asdasda456789",
		"type" : "json",
	},
	"globalVars":{
		"token" : "{{getToken}}",
		"sid" :1
	},
	"cases":[
		{
			"name":"login",
			"api": "/api/user/userinfo",
			"method":"GET",
			"headers": {
				"cache": True,
			},
			"params":{
				"time":time.time()
				"cid":"{{getRand 50 55}}"
			},
			"validate":[
				{
					"op" : "eq",
					"source": "{{.body.code}}",
					"verified":200,
				},
				{
					"op" : "gt",
					"source": "{{.body.data.num}}",
					"verified":0,
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

