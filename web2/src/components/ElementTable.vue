<template>
  	<div class="main">
 	 <h1 class="page-header"><span id="page_title">流程实例列表</span> <div class="pull-right"></div></h1>
       <div class="table-responsive">
       <form class="navbar-form navbar-left" role="form" method="post">
          <div class="form-group">
              <input type="text" class="input-date form-control" id="begin_date" value="" name="begin_date"/>
              <span>至</span>
              <input type="text" class="input-date form-control" id="end_date" value="" name="end_date"/>
              <input type="hidden" class="form-control" name='env' value="<?php echo $env?>"/>
              <select class="chosen-select" style="width:200px;" tabindex="2" id="pid">
                  <option value=""></option>  
              </select>
              <button type="button" class="btn btn-success" onclick="show_flow_list()">查询</button>
          </div>
      </form>
        <div id="processes_detail"></div>
        <div id="flow_insts_list"></div>
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
export default {
     methods: {
        show_flow_list: function () {
            var begin_date = document.getElementById("begin_date").value;
    var end_date = document.getElementById("end_date").value;
    var set_nu = document.getElementById("set_nu").value;
    var process_type = document.getElementById("process_type").value;
    var process_state_type = document.getElementById("process_state_type").value;
    var pid = document.getElementById("pid").value;
    var ajax_url = "/flow/show_flow_list?begin_date="+begin_date+"&end_date="+end_date+"&set_nu="+set_nu
    	+"&process_type="+process_type+"&process_state_type="+process_state_type+"&pid="+pid;
    var self = $(this);
    var data = $.ajax({
        type: 'GET',
        url: ajax_url,
        success: function(result) {
            var container_id = "#flow_insts_list";
            $(container_id).html(result);
        }
    });
    }
}
}
</script>
