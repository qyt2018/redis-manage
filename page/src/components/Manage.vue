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
    min-height: 30px;
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
  .clear{
    clear:both;
  }
  .layout-menu-left{
    height: 800px;
    overflow-y: auto;
  }
</style>
<template>
  <div class="manage clear" v-bind:data="mode">
    <Row>
      <Col span="5" class="layout-menu-left">
      <Menu active-name="1-2" width="auto" :open-names="['1']">
        <Submenu name="1">
          <template slot="title">
            <Icon type="ios-navigate"></Icon>
            KEYS
          </template>
          <MenuItem v-for="(item, index) in keys" :name="'1-'+index" @click.native="getValue(item)" :title="item">{{item.length > 45 ? item.substring(0,40)+'...' : item}}</MenuItem>
        </Submenu>
      </Menu>
      </Col>
      <Col span="19">
        <div v-show="show_detail==1">
            <div class="layout-content">
                <div class="layout-content-main height15">
                    <h5>{{'类型：'+type}}</h5>
                </div>
            </div>
            <hash :value="value" :ikey="ikey" v-if="type=='hash'" v-on:listenHash="del_ikey"></hash>
            <list :value="value" :ikey="ikey" v-if="type=='list'" v-on:listenList="del_ikey"></list>
            <set :value="value" :ikey="ikey" v-if="type=='set'" v-on:listenSet="del_ikey"></set>
            <string :value="value" :ikey="ikey" v-if="type=='string'" v-on:listenString="del_ikey"></string>
        </div>
      </Col>
    </Row>
  </div>
</template>
<script>
  import $Msg from 'iview/src/components/message'
  import $ from 'jquery'
  import Hash from './Manage/Hash'
  import List from './Manage/List'
  import Set from './Manage/Set'
  import String from './Manage/String'
  export default {
    name: 'manage',
    components: {
      Hash,
      List,
      Set,
      String
    },
    data () {
      return {
        mode: 1,
        keys: [],
        type: null,
        value: null,
        show_detail: 0,
        ikey: null
      }
    },
    methods: {
      getValue: function (item) {
        $.ajax({
          type: 'post',
          dataType: 'json',
          data: {action: 'getval', key: item},
          url: 'http://' + document.location.host
        }).done((resp) => {
          if (resp) {
            this.show_detail = 1
            this.type = resp.type
            let values = JSON.parse(resp.value)
            if (this.type === 'hash' || this.type === 'string') {
              this.value = values[0]
            } else {
              this.value = values
            }
            this.ikey = item
          } else {
            $Msg.warning('没有获取到任何信息')
            this.show_detail = 0
          }
        }).fail((resp) => {
          $Msg.error('获取信息失败：服务器错误')
          this.show_detail = 0
        })
      },
      del_ikey: function (data) {
        let j
        for (var i in this.keys) {
          if (this.keys[i] === data) {
            j = i
            break
          }
        }
        if (!j) {
          return false
        }
        this.keys.splice(j, 1)
        this.type = null
        this.ikey = null
        this.value = null
        return true
      }
    },
    beforeMount () {
      $.ajax({
        type: 'post',
        dataType: 'json',
        data: {action: 'getkey'},
        url: 'http://' + document.location.host
      }).done((resp) => {
        if (resp) {
          this.keys = resp
        } else {
          $Msg.warning('没有获取到任何KEYS')
        }
      }).fail((resp) => {
        $Msg.error('获取信息失败：服务器错误')
      })
    }
  }
</script>
