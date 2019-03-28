import Vue from 'vue'
import Router from 'vue-router'
import Home from '@/components/Home'
import ElementTable from '@/components/ElementTable'
import DetailInfo from '@/components/DetailInfo'
import Template from '@/components/Template'
import Monthform from '@/components/Monthform'
import Theform from '@/components/Theform'
import AgreeInfo from '@/components/AgreeInfo'

Vue.use(Router)

export default new Router({
  routes: [
    {
      path: '/',
      name: '首页',
      component: Home,
      children: [
        {
          path: '/user',
          name: '用户管理',
          component: ElementTable,
        },
        {
          path: '/userInfo/:id',
          name: '用户详情页',
          component: DetailInfo
        },
        {
          path: '/feedback',
          name: '意见反馈',
          component: AgreeInfo
        },
        {
          path: '/perform',
          name: '绩效考核',
          component: Template,
          children: [
            {
              path: '/month',
              name: '月度绩效',
              component: Monthform
            },
            {
              path: '/year',
              name: '年度绩效',
              component: Theform
            }
          ]
        }
      ]
    },
    {
      path: '*',
      redirect: '/user'
    }
  ]
})