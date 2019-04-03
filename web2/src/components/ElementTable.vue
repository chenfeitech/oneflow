<template>
  	<div class="main">
 	 <h1 class="page-header"><span id="page_title">流程实例列表</span> <div class="pull-right"></div></h1>
       <div class="table-responsive">
       <form class="navbar-form navbar-left" role="form" method="post">
          <div class="form-group">
              <input type="text" class="input-date form-control" id="begin_date" value="2019-03-29" name="begin_date">
              <span>至</span>
              <input type="text" class="input-date form-control" id="end_date" value="2019-03-31" name="end_date">
              <input type="hidden" class="form-control" name="env" value="data_flow">
              <select name="set_nu" id="set_nu" class="form-control">
                 <Option v-for="item in set_nu" :value="item.value" :key="item.value" name="set_type">
                      {{ item.label }}
                    </Option>
              </select>
              <select name="process_type" id="process_type" class="form-control">
                <option value="all">所有流程</option>
                <option value="NEW_FLOW">新流程</option>
                <option value="NEW_FLOW2">新流程</option>
                </select>
                <Select name="process_state_type" id="process_state_type" class="form-control">
                    <Option v-for="item in process_state_type" :value="item.value" :key="item.value" name="state_type">
                      {{ item.label }}
                    </Option>
                </Select>
              <select class="chosen-select" style="width: 200px; display: none;" tabindex="-1" id="pid">
                  <option value=""></option>
                  <option selected="selected" value="all">所有业务</option><option value="222">222(223)</option>              
                  </select>

                <select  class="form-control" mame="pids" id="pids">
                  <option value="all">所有</option>
                  <Option v-for="item in pids" :value="item.value" :key="item.value" name="state_type">
                      {{ item.label }}
                    </Option>
                </select>
              <button type="button" class="btn btn-success" onclick="show_flow_list()">查询</button>
          </div>
        </form>
        <div id="processes_detail"></div>
        <div id="flow_insts_list">
            {{info}}
        </div>
       </div>
    </div>
</template>

<style>
.page-header {
    padding-bottom: 9px;
    margin: -10px 0 20px; 
    border-bottom: 1px solid #eee;
    text-align: left;
}
</style>

<script>
import axios from 'axios';
export const test = 'http://oneflow.com';
export default {
  data () {
    return {
       process_state_type: [
          {
              value: '-1', label: '未成功'
          },
          {
              value: 'all', label: '所有状态'
          },
          {
              value: '2', label: '成功'
          },
          {
              value: '3', label: '失败'
          },
       ],
       set_nu: [
          {
              value: 'all', label: '所有集群'
          },
          {
              value: '0', label: 'set0'
          },
          {
              value: '1', label: 'set1'
          },
       ],
       pids: [
          {
              value: '222', label: '222'
          },

       ],
      info: null
    }
  },
  mounted(){
    this.show_flow_list()
  },
    methods: {
    show_flow_list: function () {
        let params = { project_id: 1};
        let headers = {"Content-Type": "application/json"};
        axios.get(`${test}/oneflow/getInst`, { params: params, headers:headers}).then(
            res => (this.info = res.data));
    }
}
}
</script>
