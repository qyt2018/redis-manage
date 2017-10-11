<template>
  <div class="app" id="app">
    <login :data="win" v-if="win==1" v-on:login-return="getWin"></login>
    <home :data="win" v-if="win==2" v-on:logout-return="getWin"></home>
  </div>
</template>

<script>
  import Login from './components/Login'
  import Home from './components/Home'
  export default {
    name: 'app',
    components: {
      Login,
      Home
    },
    data () {
      return {
        win: window.Win
      }
    },
    methods: {
      getWin: function (data) {
        this.win = data
      }
    }
  }
</script>

<style>
  html,body{
    padding:0;
    margin:0;
    height: 100%;
  }
  .app{
    width: 100%;
    height: 100%;
  }
</style>
