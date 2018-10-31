<link href="/thirdlib/static/chosen/chosen.css" rel="stylesheet">
<script src="/thirdlib/static/chosen/chosen.jquery.min.js"></script>
  <div class="table-responsive">
      <form class="navbar-form navbar-left" role="form" method="post">
          <div class="form-group">
              <input type="text" class="input-date form-control" id="begin_date" 
                value="<?php 
                  if (set_value('begin_date') != null) {
                    echo set_value('begin_date');
                  }elseif (isset($begin_date)) {
                    echo $begin_date;
                  }else{
                    echo date('Y-m-d',time()-3*24*60*60);
                }?>" name="begin_date"/>
              <span>至</span>
              <input type="text" class="input-date form-control" id="end_date" 
                value="<?php 
                  if (set_value('end_date') != null) {
                    echo set_value('end_date');
                  }elseif (isset($end_date)) {
                    echo $end_date;
                  }else{
                    echo date('Y-m-d',time()-24*60*60);
                }?>" name="end_date"/>
              <input type="hidden" class="form-control" name='env' value="<?php echo $env?>"/>
              <?php echo form_dropdown('set_nu', $set_list, set_value('set_nu'), 'id="set_nu" class="form-control"'); ?>
              <?php echo form_dropdown('process_type', $process_list, set_value('process_type'), 'id="process_type" class="form-control"'); ?>
              <?php 
                if (isset($process_state_type)) {
                  $default_type='all';
                }else{
                  $default_type=set_value('process_state_type');
                }
                echo form_dropdown('process_state_type', $process_state_list, $default_type, 'id="process_state_type" class="form-control"'); 
              ?>
              <select class="chosen-select" style="width:200px;" tabindex="2" id="pid">
                  <option value=""></option>
                  <?php
                    foreach ($pids as $key => $value) {
                      if ($key == "all" or count($pids) == 1) {
                        echo "<option selected = \"selected\" value=\"".$key."\">".$value."</option>";
                      }
                      else{
                        echo "<option value=\"".$key."\">".$value."</option>";
                      }
                    }
                  ?>
              </select>
              <button type="button" class="btn btn-success" onclick="show_flow_list()">查询</button>
          </div>
      </form>
      <br><br><br>
      <div id="processes_detail"></div>
      <div id="flow_insts_list"></div>
  </div>
<script type="text/javascript">
var config = {
   '.chosen-select'           : {},
   '.chosen-select-deselect'  : {allow_single_deselect:true},
   '.chosen-select-no-single' : {disable_search_threshold:10},
   '.chosen-select-no-results': {no_results_text:'Oops, nothing found!'},
   '.chosen-select-width'     : {width:"95%"}
}
for (var selector in config) {
   $(selector).chosen(config[selector]);
}
function show_flow_list()
{
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
show_flow_list();
</script>
