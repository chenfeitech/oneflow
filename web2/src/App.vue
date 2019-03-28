<template>
  <div id="app">
    <div class="app-header">
        <div class="title">OneFlow流程管理系统</div>
    </div>
    <router-view/>
  </div>
</template>

<script>
export default {
  name: 'App'
}
</script>

<style>
#app {
  font-family: 'Avenir', Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
}
*{
  margin: 0;
  padding: 0;
}
#app .app-header{color: #fff; background:#324057; line-height:60px;  text-align: left; padding-left: 15px;}
#app .app-header .title{font-size: 24px;}

.el-submenu .el-menu-item{
  padding-left:0!important;
  padding-right:0!important;
  min-width: 150px;
}

</style>
