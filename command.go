package main

import(
	"fmt"
	"errors"
	"net"
	"encoding/json"
	"strconv"
)

/*
	keys xxx
*/
func getKeys(k string, client net.Conn) ([]string, error){
	cmd := fmt.Sprintf("KEYS %s", k)
	return exec(cmd, client)
}

/*
	keys type
*/
func getType(k string, client net.Conn) (string, error){
	cmd := fmt.Sprintf("TYPE %s", k)
	data, err := exec(cmd, client)
	return data[0], err
}

/*
	get xxx
*/
func get(k string, client net.Conn) ([]string, error){
	cmd := fmt.Sprintf("GET %s", k)
	data, err := exec(cmd, client)
	return data, err
}

/*
	hget xxx
*/
func hGet(k string, client net.Conn) ([]string, error){
	cmd := fmt.Sprintf("HGETALL %s", k)
	data, err := exec(cmd, client)
	json_data := make(map[string]string)
	var key string
	for k, v := range data {
		if k % 2 == 0 {
			key = v
		}else{
			json_data[key] = v
		}
	}
	json_byte, _ := json.Marshal(json_data)
	json_str := string(json_byte)
	json_arr := []string{json_str}
	return json_arr, err
}

/*
	hset k field value
*/
func hSet(k string, field string, value string, client net.Conn) (int, error) {
	cmd := fmt.Sprintf("HSET %s %s %s", k, field, value)
	fmt.Println(cmd)
	r, err := exec(cmd, client)
	if err != nil {
		return -1, err
	}
	d := 0
	d, err = strconv.Atoi(r[0])
	if err != nil {
		return -1, err
	}
	return d, nil
}

/*
	HDEL key field [field ...]
*/
func hDel(k string, field string, client net.Conn) (int, error) {
	cmd := fmt.Sprintf("HDEL %s %s", k, field)
	r, err := exec(cmd, client)
	if err != nil {
		return -1, err
	}
	d := 0
	d, err = strconv.Atoi(r[0])
	if err != nil {
		return -1, err
	}
	return d, nil
}

/*
	DEL key [key ...]
*/
func del(k string, client net.Conn) (int, error) {
	cmd := fmt.Sprintf("DEL %s", k)
	r, err := exec(cmd, client)
	if err != nil {
		return -1, err
	}
	d := 0
	d, err = strconv.Atoi(r[0])
	if err != nil {
		return -1, err
	}
	return d, nil
}

/*
	LLEN key
*/
func lLen(k string, client net.Conn) (int, error) {
	cmd := fmt.Sprintf("LLEN %s", k)
	r, err := exec(cmd, client)
	if err != nil {
		return 0, err
	}
	d := 0
	d, err = strconv.Atoi(r[0])
	if err != nil {
		return d, err
	}
	return d, nil
}

/*
	LRANGE key start stop
*/
func lRange(k string, s int, e int, client net.Conn) ([]string, error){
	cmd := fmt.Sprintf("LRANGE %s %d %d", k, s, e)
	data, err := exec(cmd, client)
	return data, err
}

/*
	LPUSH key value [value ...]
*/
func lPush(k string, value string, client net.Conn){
	cmd := fmt.Sprintf("LPUSH %s %s", k, value)
	r, err := exec(cmd, client)
	if err != nil {
		return 0, err
	}
	d := 0
	d, err = strconv.Atoi(r[0])
	if err != nil {
		return d, err
	}
	return d, nil
}

/*
	LSET key index value
*/
func lSet(k string, index int, value string, client net.Conn) {
	
}

/*
	value
*/
func getValue(k string, client net.Conn) ([]string, error){
	var data []string
	kType, err := getType(k, client)
	if err != nil {
		return data, err
	}
	switch kType {
	case "none":
		err = errors.New(k + " is not exists.")
		return data, err
	case "hash":
		return hGet(k, client)
	case "list":
		var len int = 0
		len, err = lLen(k, client)
		if err != nil {
			return data, err
		}
		return lRange(k, 0, len, client)
	case "string","set","zset":
		return get(k, client)
	}
	return data, nil
}
