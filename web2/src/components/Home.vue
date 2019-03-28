<template>
  <div class="main">
    <div class="app_content clearfix">
        <div class="app_nav">
            <app-nav></app-nav>
        </div>
        <div class="app_wrap">
            <!-- 此处放置el-tabs代码 -->
            <div class="template-tabs">
                <el-tabs v-model="activeIndex" type="border-card" closable @tab-click="tabClick" v-if="options.length" @tab-remove="tabRemove">
                    <el-tab-pane v-for="(item, index) in options" :key="item.name" :label="item.name" :name="item.route">
                        <div class="page_content">
                            <router-view/>
                        </div>
                    </el-tab-pane>
                </el-tabs>
            </div>
        </div>
    </div>
  </div>
</template>

<script>
import appNav from './appNav'
export default {
  components:{
    appNav
  },
  watch: {
    '$route'(to) {
      let flag = false;//判断是否页面中是否已经存在该路由下的tab页
      //options记录当前页面中已存在的tab页
      for (let option of this.options) {
      //用名称匹配，如果存在即将对应的tab页设置为active显示桌面前端
          if (option.name === to.name) {
              flag = true;
              this.$store.commit('set_active_index', '/' + to.path.split('/')[1]);
              break;
          }
      }
      //如果不存在，则新增tab页，再将新增的tab页设置为active显示在桌面前端
      if (!flag) {
          this.$store.commit('add_tabs', { route: '/' + to.path.split('/')[1], name: to.name });
          this.$store.commit('set_active_index', '/' + to.path.split('/')[1]);
      }
    }
  },
  created() {
      this.tabClick()
  },
  methods: {
    // tab切换时，动态的切换路由
    tabClick(tab) {
        let path = this.activeIndex;
        // 用户详情页的时候，对应了二级路由，需要拼接添加第二级路由
        if (this.activeIndex === '/userInfo') {
            path = this.activeIndex + '/' + this.$store.state.userInfo.name;
        }
        this.$router.push({ path: path });//路由跳转
    },
    tabRemove(targetName) {
        // 首页不可删除
        if (targetName == '/') {
            return;
        }
        //将改tab从options里移除
        this.$store.commit('delete_tabs', targetName);
        
        //还同时需要处理一种情况当需要移除的页面为当前激活的页面时，将上一个tab页作为激活tab
        if (this.activeIndex === targetName) {
            // 设置当前激活的路由
            if (this.options && this.options.length >= 1) {
                this.$store.commit('set_active_index', this.options[this.options.length - 1].route);
                this.$router.push({ path: this.activeIndex });
            } else {
                this.$router.push({ path: '/' });
            }
        }
    }
  },
  computed: {
      options() {
          return this.$store.state.options;
      },
      //动态设置及获取当前激活的tab页
      activeIndex: {
          get() {
              return this.$store.state.activeIndex;
          },
          set(val) {
              this.$store.commit('set_active_index', val);
          }
      }
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
.app_nav{
    float:left;
    width: 15%;
}
.app_wrap {
    float:right;
    width: 85%;
}
.clearfix:before{
    content:'.';
    visibility: hidden;
    display: block;
    clear: both;
    width: 0;
    height: 0;
}
.clearfix:after{
    content:'.';
    visibility: hidden;
    display: block;
    clear: both;
    width: 0;
    height: 0;
}
.clearfix:before,.clearfix:after{
    *zoom: 1;
}
.page_content {
    min-height: 500px;
}
</style>
