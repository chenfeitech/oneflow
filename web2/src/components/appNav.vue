<template>
    <div class="app-nav-wrap">
    <!-- default-active设置当前激活的菜单。设置router使用 vue-router 的模式，启用该模式会在激活导航时以 index 作为 path 进行路由跳转 -->
    <el-menu :default-active="$route.path" class="el-menu-vertical-demo" router>
    <!-- index设置的值为path -->
        <el-menu-item v-for="menu in menus" :index="menu.route" :key="menu.route" v-if="!menu.children">
            <i class="el-icon-menu"></i>{{menu.name}}
        </el-menu-item>
        <el-submenu v-for="menu in menus" :index="menu.route" :key="menu.route" v-if="menu.children">
            <template slot="title"><i class="el-icon-menu"></i>{{menu.name}}</template>
            <el-menu-item :index="item.route" v-for="item in menu.children" :key="item.route"> <i class="el-icon-location"></i>{{item.name}}</el-menu-item>
        </el-submenu>
    </el-menu>
</div>
</template>

<script>
export default {
    data(){
        return {
            menus: [
                { route: '/', name: '首页' , children:false},
                { route: '/user', name: '用户管理' , children:false},
                // { route: '/psd', name: '密码管理' , children:false},
                // { route: '/salary', name: '工资管理' , children:false},
                // { route: '/attendence', name: '考勤管理' , children:false},
                // { route: '/admin', name: '系统管理' , children:false},
                { route: '/perform', name: '绩效考核', children: [{ route: '/month', name: '月度绩效' }, { route: '/year', name: '年度绩效' }] },
                { route: '/feedback', name: '意见反馈' , children:false}
            ]
        }
    },
    
}
</script>


<style scoped>

</style>