<template>
  <div class="login" v-bind:data="win">
    <h3>连接信息</h3>
    <ul>
      <li><label>打开：</label><select v-model="curHost" @change="selectHost" class="in login_open">
            <option :value="0">新建</option>
            <option v-for="(value, key) in conf" :value="key">{{ value.name }}</option>
        </select>
      </li>
      <li><label>名称：</label><input class="in" type="text" v-model="name" /></li>
      <li><label>主机：</label><input class="in" type="text" v-model="host" /></li>
      <li><label>端口：</label><input class="in" type="text" v-model="port" /></li>
      <li><label>密码：</label><input class="in" type="text" v-model="auth" /></li>
      <li class="login_save"><input type="checkbox" v-model="save" name="save" :value="1" /> 保存信息</li>
      <li class="login_btn"><a href="javascript:void(0);" @click.prevent="connect">确认</a></li>
      <li class="login_msg">{{ msg }}</li>
    </ul>
  </div>
</template>

<script>
import $ from 'jquery'
export default {
  name: 'login',
  data () {
    return {
      curHost: 0,
      name: '',
      host: '',
      port: 6379,
      auth: '',
      msg: '',
      win: 1,
      save: 0,
      conf: window.Conf
    }
  },
  methods: {
    connect: function () {
      $.ajax({
        type: 'post',
        dataType: 'json',
        data: {action: 'connect', host: this.host, port: this.port, name: this.name, auth: this.auth, save: this.save ? 1 : 0},
        url: 'http://' + document.location.host
      }).done((resp) => {
        if (resp.err === '1') {
          this.msg = resp.msg
        } else {
          this.msg = resp.msg
          this.win = 2
          this.$emit('login-return', this.win)
        }
      }).fail((resp) => {
        this.msg = '获取信息失败：服务器错误'
      })
    },
    selectHost: function () {
      if (this.curHost === 0) {
        this.host = ''
        this.port = 6379
        this.name = ''
        this.auth = ''
        return
      }
      var config = this.conf[this.curHost]
      console.log(this.curHost, config)
      this.host = config['host']
      this.port = config['port']
      this.name = config['name']
      this.auth = config['auth']
    }
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
.login{
  float: left;
  width: 20%;
  margin-left: 40%;
  marign:0 auto;
  text-align: center;
}
.login h3{
  text-align: center;
  color: #888;
}
.login ul li{
  display: block;
  text-align: left;
}
.login .login_open{
  height: 40px;
}
.login .in{
  border: 1px solid #DADADA;
  color: #888;
  height: 30px;
  margin-bottom: 16px;
  margin-right: 6px;
  margin-top: 2px;
  outline: 0 none;
  padding: 3px 3px 3px 5px;
  width: 100%;
  font-size: 12px;
  line-height:15px;
  box-shadow: inset 0px 1px 4px #ECECEC;
  -moz-box-shadow: inset 0px 1px 4px #ECECEC;
  -webkit-box-shadow: inset 0px 1px 4px #ECECEC;
}
.login .login_save {
  height: 30px;
  line-height:30px;
  vertical-align: middle;
  margin-bottom: 16px;
  color: #888;
  text-align: center;
}
.login .login_save input{
  border: 1px solid #DADADA;
  color: #888;
  box-shadow: inset 0px 1px 4px #ECECEC;
  -moz-box-shadow: inset 0px 1px 4px #ECECEC;
  -webkit-box-shadow: inset 0px 1px 4px #ECECEC;
}
.login ul li label {
  font-size: 11px;
  color: #888;
}
.login .login_btn{
  text-align: center;
}
.login .login_btn a{
  text-decoration: none;
  background: #E27575;
  border: none;
  padding: 10px 25px 10px 25px;
  color: #FFF;
  box-shadow: 1px 1px 5px #B6B6B6;
  border-radius: 3px;
  text-shadow: 1px 1px 1px #9E3F3F;
  cursor: pointer;
  width: 80%;
}

.login .login_btn a:hover {
  background: #CF7A7A;
}

.login .login_msg{
  height: 30px;
  color: #F00;
  text-align: center;
  line-height: 30px;
  vertical-align: middle;
}
</style>
