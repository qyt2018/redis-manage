<style scoped>
    .mt20 {
        margin-top:20px;
    }
    .h400 {
        max-height: 400px;
    }
    .h360 {
        max-height: 360px;
    }
    .top,.bottom{
        text-align: center;
    }
    .center{
        width: 300px;
        margin: 10px auto;
        overflow: hidden;
    }
    .center-left{
        float: left;
    }
    .center-right{
        float: right;
    }
</style>
<template>
    <div class="layout-content">
        <div class="layout-content-main">
            <template>
                <Input v-model="value" type="textarea" :rows="6" placeholder="请输入..." style="width: 100%; margin-top:10px;"><span slot="prepend">值</span></Input>
            </template>
        </div>
        <div class="layout-content-main mt20">
            <template>
                <Button type="primary" @click.prevent="save">保存</Button>
                <Tooltip content="将删除整个KEY" placement="right">
                <Button type="ghost" @click.prevent="del" style="margin-left: 8px">删除</Button>
                </Tooltip>
            </template>
        </div>
    </div>
</template>
<script>
import $Msg from 'iview/src/components/message'
import $Modal from 'iview/src/components/modal'
import $ from 'jquery'
export default {
  name: 'string',
  data () {
    return {

    }
  },
  props: ['value', 'ikey'],
  methods: {
    save: function () {
      if (!this.value) {
        $Msg.warning('请输入您想保存的数据~')
        return false
      }
      $.ajax({
        type: 'post',
        dataType: 'json',
        data: {action: 'save', type: 'string', ikey: this.ikey, val: this.value},
        url: 'http://' + document.location.host
      }).done((resp) => {
        if (resp.err === '0') {
          $Msg.success('保存成功')
        } else {
          $Msg.warning(resp.msg)
        }
      }).fail((resp) => {
        $Msg.error('获取信息失败：服务器错误')
      })
    },
    del: function () {
      let msg = '您确认要删除' + this.ikey + '的所有值吗？'
      $Modal.confirm({
        title: '是否确认删除',
        content: msg,
        onOk: () => {
          this.delItem()
        }
      })
    },
    delItem: function () {
      $.ajax({
        type: 'post',
        dataType: 'json',
        data: {action: 'del', type: 'string', ikey: this.ikey},
        url: 'http://' + document.location.host
      }).done((resp) => {
        if (resp.err === '0') {
          this.$emit('listenString', this.ikey)
          $Msg.success('删除成功')
        } else {
          $Msg.warning(resp.msg)
        }
      }).fail((resp) => {
        $Msg.error('获取信息失败：服务器错误')
      })
    }
  }
}
</script>