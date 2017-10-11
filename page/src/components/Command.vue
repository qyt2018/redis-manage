<template>
  <div class="command" v-bind:data="mode">
    <div class="show_command">
      <p v-for="item in msg">
        <label v-if="item.type==='content'">
          <template v-if="item.content instanceof Array">
            <em v-for="(cont,index) in item.content">{{ cont }}<br v-if="(index <= item.content.length-1) && cont!=''" /></em>
          </template>
          <template v-else>
            {{ item.content }}<br />
          </template>
        </label>
        <label v-if="item.type==='command'" class="command_text">${{ item.content }}</label>
      </p>
    </div>
    <div class="input_command"><span>$</span><input type="text" v-model="command" @keyup.enter="sendCommand" @keyup.up="historyUp" @keyup.down="historyDown" /></div>
  </div>
</template>
<script>
  import $ from 'jquery'
  export default {
    name: 'command',
    data () {
      return {
        msg: [{type: 'content', content: 'Connection success.'}],
        command: '',
        conn: null,
        historyCommand: [],
        history: 0,
        mode: 2
      }
    },
    mounted: function () {
      this.$nextTick(function () {
        $('.input_command input').focus()
      })
    },
    methods: {
      sendCommand: function () {
        var data = {}
        if (this.command === '') {
          data.content = ''
          data.type = 'command'
          this.msg.push(data)
          _this.scrollTop()
          return
        }
        if (this.command === 'exit') {
          this.win = 1
          return
        }
        var _this = this
        data = {}
        data.content = this.command
        data.type = 'command'
        this.msg.push(data)
        this.historyCommand.push(this.command)
        this.history = this.historyCommand.length
        this.conn = new WebSocket('ws://' + document.location.host + '/cmd')
        this.conn.onopen = function () {
          _this.conn.send(_this.command)
          _this.command = ''
        }
        this.conn.onmessage = function (event) {
          data = {}
          data.content = event.data.split('\n')
          data.type = 'content'
          _this.msg.push(data)
          _this.scrollTop()
        }
        _this.scrollTop()
      },
      scrollTop: function () {
        this.$nextTick(() => {
          var container = document.body
          container.scrollTop = container.scrollHeight
        })
      },
      historyUp: function () {
        this.history--
        if (this.history <= 0) {
          this.history = this.historyCommand.length - 1
        }
        this.command = this.historyCommand[this.history]
      },
      historyDown: function () {
        this.history++
        if (this.history >= this.historyCommand.length - 1) {
          this.history = 0
        }
        this.command = this.historyCommand[this.history]
      }

    }
  }
</script>
<style scoped>
  .command{
    width: 100%;
    height: 100%;
    background: #000000;
    padding-left: 0px;
  }
  .command .show_command{
    background: #000000;
    width: 100%;
    color: #fff;
  }
  .command .input_command{
    width: 100%;
    background: #000;
    color: #fff;
    padding-left: 0;
    padding-right: 0;
    padding-bottom: 10px;
  }
  .command .input_command input{
    width: 80%;
    border: none;
    background: #000;
    color:#fff;
    padding-left: 5px;
    outline:none;
  }
  .command .show_command .command_text{
    color: #ccc;
  }
  p {
    margin:0px;
    display: block;
    -webkit-margin-before: 0px;
    -webkit-margin-after: 0px;
    -webkit-margin-start: 0px;
    -webkit-margin-end: 0px;
  }
</style>
