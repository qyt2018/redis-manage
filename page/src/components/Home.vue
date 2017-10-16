<style scoped>
  .layout{
    border: 1px solid #d7dde4;
    background: #f5f7f9;
  }
  .layout-logo{
    width: 100px;
    height: 30px;
    background: #5b6270;
    border-radius: 3px;
    float: left;
    position: relative;
    top: 15px;
    left: 20px;
    line-height:30px;
    vertical-align: middle;
    color:#fff;
    padding-left:20px;
  }
  .layout-nav{
    width: 420px;
    margin: 0 auto;
  }
  .layout-assistant{
    width: 300px;
    margin: 0 auto;
    height: inherit;
  }
  .layout-breadcrumb{
    padding: 10px 15px 0;
  }
  .layout-content{
    min-height: 200px;
    margin: 15px;
    overflow: hidden;
    background: #fff;
    border-radius: 4px;
  }
  .layout-content-main{
    padding: 10px;
  }
  .layout-copy{
    text-align: center;
    padding: 10px 0 20px;
    color: #9ea7b4;
  }
  .height15{
    height:15px;
  }
</style>
<template>
  <div class="layout">
    <Menu mode="horizontal" theme="dark" active-name="1">
      <div class="layout-logo"><a href="javascript:void(0);" @click.prevent="logout">退出登录</a></div>
      <div class="layout-nav">
        <MenuItem @click.native="getMode(1)">
          <Icon type="ios-eye"></Icon>
          视图模式
        </MenuItem>
        <MenuItem @click.native="getMode(2)">
          <Icon type="code-working"></Icon>
          命令模式
        </MenuItem>
      </div>
    </Menu>
    <div class="layout-breadcrumb">
      <Breadcrumb>
        <BreadcrumbItem>{{mode===1 ? '视图模式' : '命令模式'}}</BreadcrumbItem>
      </Breadcrumb>
    </div>
    <div class="layout-content">
      <manage :data="mode" v-if="mode==1"></manage>
      <command :data="mode" v-if="mode==2"></command>
    </div>
  </div>
</template>
<script>
  import Manage from './Manage'
  import Command from './Command'
  import $Msg from 'iview/src/components/message'
  import $ from 'jquery'

  export default {
    name: 'home',
    components: {
      Manage,
      Command
    },
    data () {
      return {
        mode: window.Mode,
        win: 2
      }
    },
    methods: {
      getMode: function (data) {
        this.mode = data
      },
      logout: function () {
        $.ajax({
          type: 'post',
          dataType: 'json',
          data: {action: 'logout'},
          url: 'http://' + document.location.host
        }).done((resp) => {
          this.win = 1
          this.$emit('logout-return', this.win)
        }).fail((resp) => {
          $Msg.error('操作失败：服务器错误')
        })
      }
    }
  }
</script>
