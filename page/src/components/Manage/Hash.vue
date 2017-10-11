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
                                    <th class=""><div class="ivu-table-cell"><span>键</span></div></th>
                                    <th class=""><div class="ivu-table-cell"><span>值</span></div></th>
                                </tr>
                            </thead>
                        </table>
                    </div>
                    <div class="ivu-table-body h360">
                        <table cellspacing="0" cellpadding="0" border="0" style="width: 100%;">
                            <tbody class="ivu-table-tbody">
                                <tr class="ivu-table-row ivu-table-row-hover" v-for="(item, index) in values">
                                    <td class="" @click.prevent="show(item.idx)"><div class="ivu-table-cell"><span>{{ item.idx }}</span></div></td>
                                    <td class="" @click.prevent="show(item.idx)"><div class="ivu-table-cell"><span>{{ item.key }}</span></div></td>
                                    <td class="" @click.prevent="show(item.idx)"><div class="ivu-table-cell"><span>{{ item.val }}</span></div></td>
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
  name: 'hash',
  data () {
    return {
      key: null,
      val: null,
      idx: 0
    }
  },
  props: ['value', 'ikey'],
  computed: {
    values: {
      get: function () {
        return this.format()
      },
      set: function (data) {
        this.values = data
      }
    }
  },
  methods: {
    show: function (idx) {
      var index = idx - 1
      var data = this.values[index]
      this.key = data.key
      this.val = data.val
      this.idx = data.idx
    },
    format: function () {
      var data = []
      let value = JSON.parse(this.value)
      let index = 0
      for (let key in value) {
        index++
        var obj = {}
        obj.idx = index
        obj.key = key
        obj.val = value[key]
        data.push(obj)
      }
      return data
    },
    save: function () {
      if (!this.key) {
        $Msg.warning('请在数据展示区选择要保存的数据，或输入您想新增的数据~')
        return false
      }
      $.ajax({
        type: 'post',
        dataType: 'json',
        data: {action: 'save', type: 'hash', ikey: this.ikey, key: this.key, val: this.val},
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
      let j
      for (let i in this.values) {
        if (this.values[i]['key'] === this.key) {
          j = i
          break
        }
      }
      if (!j) {
        let value = JSON.parse(this.value)
        value[this.key] = this.val
        this.value = JSON.stringify(value)
      }
      return
    },
    del: function () {
      let msg = ''
      if (!this.key) {
        msg = '您确认要删除' + this.ikey + '的所有值吗？'
      } else {
        msg = '您确认要删除' + this.key + '的值吗？'
      }
      $Modal.confirm({
        title: '是否确认删除',
        content: msg,
        onOk: () => {
          this.delKey()
        }
      })
    },
    delKey: function () {
      $.ajax({
        type: 'post',
        dataType: 'json',
        data: {action: 'del', type: 'hash', ikey: this.ikey, key: this.key},
        url: 'http://' + document.location.host
      }).done((resp) => {
        if (resp.err === '0') {
          if (!this.key) {
            this.$emit('listenHash', this.ikey)
          } else {
            let data = {}
            let value = JSON.parse(this.value)
            for (let key in value) {
              if (key !== this.key) {
                data[key] = value[key]
              }
            }
            this.value = JSON.stringify(data)
            this.key = null
            this.val = null
            this.idx = 0
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