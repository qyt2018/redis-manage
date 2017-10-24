package main

import(
	"fmt"
	"net"
	"strings"
)

const (
	SEND_PRE = "*"
	SEND_NUM = "$"
	SEND_END = "\r\n"
	REPLY_STATUS = "+"
	REPLY_ERROR = "-"
	REPLY_INTEGER = ":"
	REPLY_BULK = "$"
	REPLY_MULTI = "*"
)

type server struct {
	Name string `json:"name"`
	Host string `json:"host"`
	Port string `json:"port"`
	Auth string `json:"auth"`
}

/*
	连接redis服务器
*/
func conn(s *server) (net.Conn, error){
	ip := fmt.Sprintf("%s:%s", s.Host, s.Port)
	c, err := net.Dial("tcp", ip)
	return c, err
}

/*
	执行命令
*/
func exec(cmd string, client net.Conn) ([]string, error){
	var data []string
	cmdF := formatCmd(cmd)
	cmdBuf := []byte(cmdF)
	n, err := client.Write(cmdBuf)
	if err != nil {
		return data, err
	}
	buf := make([]byte, 204800000)
	n, err = client.Read(buf)
	if err != nil {
		return data, err
	}
	data = formatOut(string(buf[:n]))
	//fmt.Println("cmd:"+cmd+",output:"+string(buf[:n]))
	return data, nil
}

/*
	格式化命令
*/
func formatCmd(sendCmd string) string{
	var cmd string = ""
	cmds := strings.Split(sendCmd, " ")
	cmd += fmt.Sprintf("%s%d%s",SEND_PRE, len(cmds), SEND_END)
	for _,cmdUnit := range cmds {
		cmd += fmt.Sprintf("%s%d%s%s%s", SEND_NUM, len(cmdUnit), SEND_END, cmdUnit, SEND_END)
	}
	return cmd
}

/*
	格式化命令
*/
func formatOut(reply string) []string {
	var data []string
	if(strings.HasPrefix(reply, REPLY_STATUS)){
		data = append(data, strings.Trim(strings.TrimLeft(reply, REPLY_STATUS), SEND_END))
		return data
	}
	if(strings.HasPrefix(reply, REPLY_ERROR)){
		data = append(data, strings.Trim(strings.TrimLeft(reply, REPLY_STATUS), SEND_END))
		return data
	}
	if(strings.HasPrefix(reply, REPLY_INTEGER)){
		data = append(data, strings.Trim(strings.TrimLeft(reply, REPLY_INTEGER), SEND_END))
		return data
	}
	replies := strings.Split(reply, SEND_END)
	if(strings.HasPrefix(reply, REPLY_BULK)){
		for _, bRep := range replies {
			if bRep == "" || bRep == " "{
				continue
			}
			if !strings.HasPrefix(bRep, REPLY_BULK) {
				data = append(data, bRep)
			}
		}
		return data
	}
	if(strings.HasPrefix(reply, REPLY_MULTI)){
		for _, mRep := range replies {
			if mRep == "" || mRep == " "{
				continue
			}
			if !strings.HasPrefix(mRep, REPLY_MULTI) && !strings.HasPrefix(mRep, REPLY_BULK) {
				data = append(data, mRep)
			}
		}
		return data
	}
	return data
}