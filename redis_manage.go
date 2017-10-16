package main

import (
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"io"
	//"net"
	"github.com/gorilla/websocket"
	"net/http"
	"text/template"
	"strconv"
)

var (
	upgrader    = websocket.Upgrader{}
	redisServer = &server{}
)

func main() {
	http.Handle("/favicon.ico", http.FileServer(http.Dir("page/dist")))
	http.Handle("/static/", http.FileServer(http.Dir("page/dist")))
	http.HandleFunc("/", handle)
	http.HandleFunc("/cmd", cmdHandle)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Error(err.Error())
	} else {
		log.Info("service start.")
	}
}

//处理请求
func handle(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	data := make(map[string]string)
	if _, ok := r.Form["action"]; ok {
		action := r.Form["action"][0]
		switch action {
		case "connect":
			connect(w, r)
			return
		case "logout":
			logout(w, r)
			return
		case "getkey":
			getKey(w, r)
			return
		case "getval":
			getVal(w, r)
			return
		case "save":
			save(w, r)
			return
		case "del":
			delete(w, r)
			return
		}
	}
	conf := getConfig()
	data["Conf"] = conf
	data["Win"] = "1"
	data["Mode"] = "1"
	if redisServer.Host != "" {
		data["Win"] = "2"
	}
	html, err := template.ParseFiles("page/dist/index.html")
	if err != nil {
		log.Error(err.Error())
		return
	}
	html.Execute(w, data)
}

//连接
func connect(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	name := r.Form["name"][0]
	host := r.Form["host"][0]
	port := r.Form["port"][0]
	auth := r.Form["auth"][0]
	isSave := r.Form["save"][0]
	redisServer = &server{name, host, port, auth}
	if isSave == "1" {
		saveConfig(redisServer)
	}
	client, err := conn(redisServer)
	defer client.Close()
	ret := make(map[string]string)
	if err != nil {
		log.Error(err.Error())
		ret["msg"] = "connect failed. error:" + err.Error()
		ret["err"] = "1"
		data, _ := json.Marshal(ret)
		io.WriteString(w, string(data))
		return
	}
	ret["msg"] = "connect success."
	ret["err"] = "0"
	dats, _ := json.Marshal(ret)
	io.WriteString(w, string(dats))
	return
}

func logout(w http.ResponseWriter, r *http.Request) {
	redisServer.Host = ""
	return
}

//命令行处理
func cmdHandle(w http.ResponseWriter, r *http.Request) {
	redisClient, err := conn(redisServer)
	defer redisClient.Close()
	//处理websocket
	ws, errs := upgrader.Upgrade(w, r, nil)
	if errs != nil {
		log.Error("upgrade: " + errs.Error())
		return
	}
	defer ws.Close()
	//处理客户端输入
	msgType, reader, errn := ws.NextReader()
	if errn != nil {
		log.Error("next reader: " + errn.Error())
		return
	}
	bufs := make([]byte, 102400)
	n, _ := reader.Read(bufs)
	cmd := string(bufs[:n])
	fmt.Println(cmd)
	//执行命令
	out, erro := exec(cmd, redisClient)
	if erro != nil {
		log.Error("exec command: " + erro.Error())
		return
	}
	message, errm := json.Marshal(out)
	if errm != nil {
		log.Error("message json code: " + errm.Error())
		return
	}
	//处理输出
	err = ws.WriteMessage(msgType, message)
	if err != nil {
		log.Error("ws write: " + err.Error())
	}
	fmt.Println(string(message))
	return
}

/*
	获取keys
*/
func getKey(w http.ResponseWriter, r *http.Request) {
	redisClient, err := conn(redisServer)
	if err != nil {
		log.Error(err.Error())
	}
	defer redisClient.Close()
	keys, errk := getKeys("*", redisClient)
	if errk != nil {
		log.Error("keys *:" + errk.Error())
	}
	keysByte, errm := json.Marshal(keys)
	if errm != nil {
		log.Error("keys json code:" + errm.Error())
	}
	io.WriteString(w, string(keysByte))
}

/*
	获取值
*/
func getVal(w http.ResponseWriter, r *http.Request) {
	redisClient, err := conn(redisServer)
	defer redisClient.Close()
	if err != nil {
		log.Error(err.Error())
	}
	r.ParseForm()
	key := r.Form["key"][0]
	keyT, errt := getType(key, redisClient)
	if errt != nil {
		log.Error("get type:" + errt.Error())
	}
	keyV, errv := getValue(key, redisClient)
	if errv != nil {
		log.Error("get value:" + errv.Error())
	}
	value, erre := json.Marshal(keyV)
	if erre != nil {
		log.Error("value json code:" + erre.Error())
	}
	vals := string(value)
	data := make(map[string]string)
	data["type"] = keyT
	data["value"] = vals
	dataJson, errd := json.Marshal(data)
	if errd != nil {
		log.Error("json code:" + errd.Error())
	}
	io.WriteString(w, string(dataJson))
}

/*
	保存值
*/
func save(w http.ResponseWriter, r *http.Request) {
	redisClient, err := conn(redisServer)
	defer redisClient.Close()
	if err != nil {
		log.Error(err.Error())
	}
	data := make(map[string]string)
	r.ParseForm()
	data_type := r.Form["type"][0]
	ikey := r.Form["ikey"][0]
	switch data_type {
	case "hash":
		key := r.Form["key"][0]
		val := r.Form["val"][0]
		_, err := hSet(ikey, key, val, redisClient)
		if err != nil {
			data["err"] = "1"
			data["msg"] = err.Error()
		} else {
			data["err"] = "0"
		}
		jsonData, errj := json.Marshal(data)
		if errj != nil {
			log.Error("json code:" + errj.Error())
		}
		io.WriteString(w, string(jsonData))
	case "list":
		index := r.Form["index"][0]
		val := r.Form["val"][0]
		var err error
		if index == "" || index == "null" {
			_, err = lPush(ikey, val, redisClient)
		} else {
		    idx, _ := strconv.Atoi(index)
			_, err = lSet(ikey, idx, val, redisClient)
		}
		if err != nil {
			data["err"] = "1"
			data["msg"] = err.Error()
		} else {
			data["err"] = "0"
		}
		jsonData, errj := json.Marshal(data)
		if errj != nil {
			log.Error("json code:" + errj.Error())
		}
		io.WriteString(w, string(jsonData))
	case "set":
		index := r.Form["index"][0]
		val := r.Form["val"][0]
		var err error
		if index == "" || index == "null" {
		    _, err = sAdd(ikey, val, redisClient)
		} else {
		    old_val := r.Form["old_val"][0]
		    _, err = sMod(ikey, old_val, val, redisClient)
		}
		if err != nil {
		    data["err"] = "1"
		    data["msg"] = err.Error()
		} else {
		    data["err"] = "0"
		}
		jsonData, errj := json.Marshal(data)
		if errj != nil {
		    log.Error("json code:" + errj.Error())
		}
		io.WriteString(w, string(jsonData))
    case "string":
        val := r.Form["val"][0]
        _, err := set(ikey, val, redisClient)
        if err != nil {
            data["err"] = "1"
            data["msg"] = err.Error()
        } else {
            data["err"] = "0"
        }
        jsonData, errj := json.Marshal(data)
        if errj != nil {
            log.Error("json code:" + errj.Error())
        }
        io.WriteString(w, string(jsonData))
	}
}

/*
	删除值
*/
func delete(w http.ResponseWriter, r *http.Request) {
	redisClient, err := conn(redisServer)
	defer redisClient.Close()
	if err != nil {
		log.Error(err.Error())
	}
	data := make(map[string]string)
	r.ParseForm()
	data_type := r.Form["type"][0]
	ikey := r.Form["ikey"][0]
	switch data_type {
	case "hash":
		key := r.Form["key"][0]
		var err error
		if key == "" || key == "null" {
			_, err = del(ikey, redisClient)
		} else {
			_, err = hDel(ikey, key, redisClient)
		}
		if err != nil {
			data["err"] = "1"
			data["msg"] = err.Error()
		} else {
			data["err"] = "0"
		}
		jsonData, errj := json.Marshal(data)
		if errj != nil {
			log.Error("json code:" + errj.Error())
		}
		io.WriteString(w, string(jsonData))
	case "list":
		index := r.Form["index"][0]
		var err error
		if index == "" || index == "null" {
			_, err = del(ikey, redisClient)
		} else {
		    idx, _ := strconv.Atoi(index)
			_, err = lDel(ikey, idx, redisClient)
		}
		if err != nil {
			data["err"] = "1"
			data["msg"] = err.Error()
		} else {
			data["err"] = "0"
		}
		jsonData, errj := json.Marshal(data)
		if errj != nil {
			log.Error("json code:" + errj.Error())
		}
		io.WriteString(w, string(jsonData))
    case "set":
        val := r.Form["val"][0]
        var err error
        if val == "" || val == "null" {
            _, err = del(ikey, redisClient)
        } else {
            _, err = sRem(ikey, val, redisClient)
        }
        if err != nil {
            data["err"] = "1"
            data["msg"] = err.Error()
        } else {
            data["err"] = "0"
        }
        jsonData, errj := json.Marshal(data)
        if errj != nil {
            log.Error("json code:" + errj.Error())
        }
        io.WriteString(w, string(jsonData))
    case "string":
        _, err := del(ikey, redisClient)
        if err != nil {
            data["err"] = "1"
            data["msg"] = err.Error()
        } else {
            data["err"] = "0"
        }
        jsonData, errj := json.Marshal(data)
        if errj != nil {
            log.Error("json code:" + errj.Error())
        }
        io.WriteString(w, string(jsonData))

	}
}
