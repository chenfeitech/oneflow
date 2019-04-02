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
                <option value="all">所有集群</option>
                <option value="0">set0</option>
                <option value="1">set1</option>
                <option value="2">set2</option>
                <option value="3">set3</option>
              </select>
              <select name="process_type" id="process_type" class="form-control">
                <option value="all">所有流程</option>
                <option value="NEW_FLOW">新流程</option>
                <option value="NEW_FLOW2">新流程</option>
                </select>
              <select name="process_state_type" id="process_state_type" class="form-control">
                <option value="-1">未成功</option>
                <option value="all">所有状态</option>
                <option value="2">成功</option>
                <option value="3">失败</option>
                </select>
              <select class="chosen-select" style="width: 200px; display: none;" tabindex="-1" id="pid">
                  <option value=""></option>
                  <option selected="selected" value="all">所有</option><option value="222">222(223)</option>              </select><div class="chosen-container chosen-container-single" style="width: 200px;" title="" id="pid_chosen"><a class="chosen-single" tabindex="-1" data-original-title="" title=""><span>所有</span><div><b></b></div></a><div class="chosen-drop"><div class="chosen-search"><input type="text" autocomplete="off" tabindex="2"></div><ul class="chosen-results"><li class="active-result result-selected" data-option-array-index="1" style="">所有</li><li class="active-result" data-option-array-index="2" style="">222(223)</li></ul></div></div>
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
