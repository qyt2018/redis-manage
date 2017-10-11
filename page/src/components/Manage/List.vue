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
            <div class="ivu-table-wrapper h400">
                <div class="ivu-table ivu-table-with-fixed-top">
                    <div class="ivu-table-header">
                        <table cellspacing="0" cellpadding="0" border="0" style="width: 100%;">
                            <thead>
                                <tr>
                                    <th class=""><div class="ivu-table-cell"><span>排序</span></div></th>
                                    <th class=""><div class="ivu-table-cell"><span>值</span></div></th>
                                </tr>
                            </thead>
                        </table>
                    </div>
                    <div class="ivu-table-body h360">
                        <table cellspacing="0" cellpadding="0" border="0" style="width: 1450px;">
                            <tbody class="ivu-table-tbody">
                                <tr class="ivu-table-row ivu-table-row-hover" v-for="(item, index) in value">
                                    <td class="" @click.prevent="show(index)"><div class="ivu-table-cell"><span>{{ index + 1 }}</span></div></td>
                                    <td class="" @click.prevent="show(index)"><div class="ivu-table-cell"><span>{{ item }}</span></div></td>
                                </tr>
                            </tbody>
                        </table>
                    </div>
                </div>
            </div>
        </div>
        <div class="layout-content-main mt20">
            <template>
                <Input v-model="key" placeholder="请输入..." style="width: 100%"><span slot="prepend">键</span></Input>
                <Input v-model="val" placeholder="请输入..." style="width: 100%; margin-top:10px;"><span slot="prepend">值</span></Input>
            </template>
        </div>
        <div class="layout-content-main mt20">
            <template>
                <Button type="primary" @click.prevent="save">保存</Button>
                <Tooltip content="若未选择删除某一条数据，将删除整个KEY值" placement="right">
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
  name: 'list',
  data () {
    return {
      val: null,
      index: null
    }
  },
  props: ['value', 'ikey'],
  methods: {
    show: function (index) {
      var data = this.value[index]
      this.val = data
      this.index = index
    },
    save: function () {
      if (!this.val) {
        $Msg.warning('请在数据展示区选择要保存的数据，或输入您想新增的数据~')
        return false
      }
      $.ajax({
        type: 'post',
        dataType: 'json',
        data: {action: 'save', type: 'list', ikey: this.ikey, index: this.index, val: this.val},
        url: 'http://' + document.location.host
      }).done((resp) => {
        if (resp.err === '0') {
          this.add()
          $Msg.success('保存成功')
        } else {
          $Msg.warning(resp.msg)
        }
      }).fail((resp) => {
        $Msg.error('获取信息失败：服务器错误')
      })
    },
    add: function () {
      if (this.index === null) {
        this.value.push(this.val)
        this.index = this.value.length - 1
      }
      return
    },
    del: function () {
      let msg = ''
      if (!this.key) {
        msg = '您确认要删除' + this.ikey + '的所有值吗？'
      } else {
        msg = '您确认要删除' + this.val + '吗？'
      }
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
        data: {action: 'del', type: 'list', ikey: this.ikey, index: this.index},
        url: 'http://' + document.location.host
      }).done((resp) => {
        if (resp.err === '0') {
          if (!this.key) {
            this.$emit('del_ikey', this.ikey)
          } else {
            this.value.splice(this.index, 1)
            this.val = null
            this.index = null
          }
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