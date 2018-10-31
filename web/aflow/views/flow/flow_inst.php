<?php $this->load->view($header);?>

<style type="text/css">
  pre.log{
    display: block;
    padding: ;
    margin: 0 0 10px;
    font-size: 13px;
    line-height: 1;
    color: #333;
    word-break: break-all;
    word-wrap: break-word;
    background-color: rgba(0, 0, 0, 0);
    border: 0;
    border-radius: 0;
  }

  .state_btn {
    text-align: left;
  }
  ol.vertical {
    margin: 0 0 9px 0;
    min-height: 30px;
  }
  ol.vertical button {
    margin: 5px;
    padding: 5px;
    border: 1px solid #cccccc;
    background: #eeeeee;
    min-height: 60px;
    width: 100%;
  }

  .state-heading span
  {
    color: #FFFFFF;
  }
  .header- span
  {
    color: #000000 !important;
  }
  .Ready .panel-heading
  {
    background-color:#5bc0de !important;
    border-color: #5bc0de !important;
  }
  .Running  .panel-heading
  {
    background-color:#6ABE05 !important;
    border-color:#6ABE05 !important;
  }

  .Succeed  .panel-heading
  {
    background-color: #337ab7 !important;
    border-color: #337ab7 !important;
  }

  .Failed  .panel-heading
  {
    background-color:#FF5050 !important;
    border-color:#FF5050 !important;
  }

  .Unknown  .panel-heading
  {
    background-color: #f0ad4e !important;
    border-color:#f0ad4e !important;
  }

  .title-link
  {
    color: #FFFFFF;
  }

  .link-
  {
    color: #000000 !important;
  }

  .bs-callout-info {
    border-left-color: #1b809e !important;
  }
  .bs-callout {
    padding: 5px 20px 5px 20px;
    margin: 20px 0;
    border: 1px solid #eee;
    border-left-width: 5px;
    border-radius: 3px;
  }
  .fa-stack {
    position: relative;
    display: inline-block;
    width: 2em;
    height: 2em;
    line-height: 2em;
    vertical-align: middle;
  }
  .fa-2x {
    font-size: 2em;
  }
  .pull-left {
    float: left !important;
  }
  .clear {
    display: block;
    overflow: hidden;
  }
  .m-t-xs {
    margin-top: 5px;
    margin-bottom: 0px;
  }
  .block {
    display: block;
  }
  .b-light {
    border-color: #e4e4e4;
  }
  .b-r {
    border-right: 1px solid #cfcfcf;
  }
</style>
<script src="/thirdlib/bootstrap/js/bootbox.min.js"></script>
<script src="/thirdlib/bootstrap/js/bootstrap.min.js"></script>

<div class="container-fluid">
  <div class="row">
    <?php $this->load->view($list);?>
    <div class="col-sm-9 col-sm-offset-3 col-md-10 col-md-offset-2 main">
      <h1 class="page-header"><span id="page_title"><?php echo $flow_inst->name?></span> 
        <div class="pull-right">
        </div>
      </h1>

      <div class="bs-callout bs-callout-info">
        <h4>流程: <a href="/flow/show_flow/<?php echo $flow_inst->flow_id?>"><?php echo $flow_inst->name?> (<?php echo $flow_inst->flow_id?>)</a></h4>

        <div class="row">
          <div class="col-md-2 flow">
            <h6>PID: <?php echo $flow_inst->pid?></h6>
          </div>
          <div class="col-md-3 flow">
            <h6>RunningDay: <?php echo $flow_inst->running_day?></h6>
          </div>
          <div class="col-md-2 flow">
            <h6>LastTask: <?php echo $flow_inst->last_task_id?></h6>
          </div>
          <div class="col-md-2 flow">
            <h6>Last Operator: <?php echo $flow_inst->last_operator?></h6>
          </div>
          <div class="col-md-2 flow">
	    <h6>Key: <?php echo $flow_inst->key?></h6>
          </div>
        </div>
        <div class="row">
          <div class="col-md-10 flow">
                <?php if ($flow_inst->state == 3) : ?>
                  <button class="btn btn-default success" type="button" onclick="set_flow_success();">Set success</button>
                <?php endif; ?>
          </div>
        </div>

      </div>

      <div class="row">
        <div class="panel-group" id="accordion" role="tablist" aria-multiselectable="true">
        <?php foreach($flow_inst->task_insts as $i => $ti):?>
          <div class="panel panel-default <?php echo $ti->state_text?>">
            <div class="panel-heading state-heading header-<?php echo $ti->state_text?>" role="tab" id="heading<?php echo $ti->id?>">
              <div data-toggle="collapse" data-parent="#accordion" href="#collapse<?php echo $ti->id?>" aria-expanded="true" aria-controls="collapse<?php echo $ti->id?>">
                <span class="name">
                  <a role="button" class="title-link link-<?php echo $ti->state_text?>" data-toggle="collapse" data-parent="#accordion_log" href="#collapse<?php echo $ti->id?>" aria-expanded="true" aria-controls="collapse<?php echo $ti->id?>">
                    <?php echo $ti->name?>
                  </a>
                </span>
                <div class="pull-right">
                  <span class="name">ID: <?php echo $ti->id?></span>
                </div>
                <div style="margin: 5px;"><span></span></div>
                <div style="margin: 5px;">
                  <span>State: <?php echo $ti->state_text?>&nbsp;</span>
                  <div class="pull-right"><em style="color:#f0f0f0"><?php echo $ti->last_update_time?></em></div>
                </div>
              </div>
            </div>
            <div id="collapse<?php echo $ti->id?>" class="panel-collapse collapse" role="tabpanel" aria-labelledby="heading<?php echo $ti->id?>">
              <div class="panel-body">
                <div class="row" style="margin-bottom:10px">
                  <div class="col-sm-6 col-md-3 b-r b-light">
                    <span class="fa-stack fa-2x pull-left m-r-sm text-center"> 
                      <i class="glyphicon glyphicon-dashboard" style="color:#5bc0de"></i>
                    </span>
                    <span class="h5 block m-t-xs">
                      <strong><?php echo $ti->ready_time?></strong>
                    </span>
                    <small class="text-muted text-uc">Ready</small> 
                  </div>
                  <div class="col-sm-6 col-md-3 b-r b-light">
                    <span class="fa-stack fa-2x pull-left m-r-sm text-center"> 
                      <i class="glyphicon glyphicon-dashboard" style="color:#6ABE05"></i>
                    </span>
                    <span class="h5 block m-t-xs">
                      <strong><?php echo $ti->running_time?></strong>
                    </span>
                    <small class="text-muted text-uc">Running</small> 
                  </div>
                  <div class="col-sm-6 col-md-3 b-r b-light">
                    <span class="fa-stack fa-2x pull-left m-r-sm text-center"> 
                      <i class="glyphicon glyphicon-dashboard" style="color:#337ab7"></i>
                    </span>
                    <span class="h5 block m-t-xs">
                      <strong><?php echo $ti->succeed_time?></strong>
                    </span>
                    <small class="text-muted text-uc">Succeed</small> 
                  </div>
                  <div class="col-sm-6 col-md-3 b-r b-light">
                    <span class="fa-stack fa-2x pull-left m-r-sm text-center"> 
                      <i class="glyphicon glyphicon-dashboard" style="color:#FF5050"></i>
                    </span>
                    <span class="h5 block m-t-xs">
                      <strong><?php echo $ti->failed_time?>&nbsp;</strong>
                    </span>
                    <small class="text-muted text-uc">Failed</small> 
                  </div>
                </div>
                <div class="row" style="margin:10px">
                <?php if ($ti->state == 1): ?>
                  <button class="btn btn-default kill" type="button" onclick="kill(<?php echo $i?>);">Kill</button>
                <?php endif; ?>
                <?php if ($ti->state == 3): ?>
                  <button class="btn btn-default success" type="button" onclick="set_success(<?php echo $i?>);">Set success</button>
                <?php endif; ?>
                  <button class="btn btn-default rerun" type="button" onclick="rerun(<?php echo $i?>);">Rerun this task</button>
                  <button class="btn btn-default rerun" type="button" onclick="rerun_seqs(<?php echo $i?>);">Rerun task sequences</button>
                  <button class="btn btn-danger force_rerun" type="button" onclick="rerun(<?php echo $i?>);">Force rerun this task</button>
                  <button class="btn btn-danger force_rerun" type="button" onclick="rerun_seqs(<?php echo $i?>);">Force rerun task sequences</button>
                </div>
              </div>
              <ul class="list-group">
                <li class="list-group-item">日志</li>
                <li class="list-group-item">
                  <pre class="log" date="<?php echo date("Y-m-d", strtotime($ti->ready_time));?>"><?php echo $ti->script_output?></pre>
                </li>
              </ul>
            </div>
          </div>
        <?php endforeach; ?>
        </div>
      </div>

      <script type="text/javascript">
        var task_list = <?php echo json_encode($flow_inst->task_insts)?>;

        var run_finished = true;
        for (var i=0; i<task_list.length; ++i) {
          if (task_list[i].state < 2) {
            run_finished = false;
            break;
          }
        }

        if (run_finished) {
          $(".force_rerun").hide();
        } else {
          $(".rerun").hide();
        }



        function kill(i) {
          bootbox.confirm("Are you sure kill task " + task_list[i].name +"?", function(result) {
            if (result) {
              do_kill(i);
            }
          }); 
        }

        function set_flow_success() {
          bootbox.confirm("Are you sure set flow success?", function(result) {
            if (result) {
              do_set_flow_success(i);
            }
          }); 
        }


        function set_success(i) {
          bootbox.confirm("Are you sure set task " + task_list[i].name +" success?", function(result) {
            if (result) {
              do_set_success(i);
            }
          }); 
        }



        function rerun(i) {
          bootbox.confirm("Are you sure rerun task " + task_list[i].name +"?", function(result) {
            if (result) {
              do_rerun(i, true);
            }
          }); 
        }
        function rerun_seqs(i) {
          bootbox.confirm("Are you sure rerun task sequences from " + task_list[i].name +"?", function(result) {
            if (result) {
              do_rerun(i, false);
            }
          }); 
        }

        var url = "/flow/api";
        function do_rerun(i, single) {
          bootbox.dialog({
            title: "Rerun",
            message: '<div class="rerun-modal-body">' +
            '<div class="progress progress-striped active" style="margin-bottom:0;"><div class="progress-bar" style="width: 100%"></div></div>' +
            '</div>'})

          var request_params = {
            "flowInstId": <?php echo $flow_inst->id?>,
            "taskId": task_list[i].id, 
            "SingleTask": single,
            "Creator": "<?php echo $userEnName?>"
          }

          var request = {};
          request.id = 3;
          request.method = "FlowService.RerunTask";
          request.params = [request_params];

          $.ajax({
            url: url, 
            data: JSON.stringify(request), 
            type: "POST",
            contentType: "application/json", 
            success: function(rpcRes) {
              if (rpcRes.error != null && rpcRes.error != "") {
                $(".rerun-modal-body").html("请求失败！<br/>" + rpcRes.error);
              } else {
                document.location.reload();
              }
            }, 
            error: function(err, status, thrown) {
              $(".rerun-modal-body").html("请求失败！<br/>" + err.responseText);
            }
          }); 
        }

        function do_kill(i) {
          bootbox.dialog({
            title: "Kill",
            message: '<div class="rerun-modal-body">' +
            '<div class="progress progress-striped active" style="margin-bottom:0;"><div class="progress-bar" style="width: 100%"></div></div>' +
            '</div>'})

          var request_params = {
            "flowInstId": <?php echo $flow_inst->id?>,
            "taskId": task_list[i].id,
          }

          var request = {};
          request.id = 3;
          request.method = "FlowService.KillTaskInstance";
          request.params = [request_params];

          $.ajax({
            url: url, 
            data: JSON.stringify(request), 
            type: "POST",
            contentType: "application/json", 
            success: function(rpcRes) {
              if (rpcRes.error != null && rpcRes.error != "") {
                $(".rerun-modal-body").html("请求失败！<br/>" + rpcRes.error);
              } else {
                document.location.reload();
              }
            }, 
            error: function(err, status, thrown) {
              $(".rerun-modal-body").html("请求失败！<br/>" + err.responseText);
            }
          }); 
        }

        function do_set_success(i) {
          bootbox.dialog({
            title: "Process",
            message: '<div class="rerun-modal-body">' +
            '<div class="progress progress-striped active" style="margin-bottom:0;"><div class="progress-bar" style="width: 100%"></div></div>' +
            '</div>'})

          var request_params = {
            "flowInstId": <?php echo $flow_inst->id?>,
            "taskId": task_list[i].id,
          }

          var request = {};
          request.id = 3;
          request.method = "FlowService.SetTaskInstanceSuccess";
          request.params = [request_params];

          $.ajax({
            url: url, 
            data: JSON.stringify(request), 
            type: "POST",
            contentType: "application/json", 
            success: function(rpcRes) {
              if (rpcRes.error != null && rpcRes.error != "") {
                $(".rerun-modal-body").html("请求失败！<br/>" + rpcRes.error);
              } else {
                document.location.reload();
              }
            }, 
            error: function(err, status, thrown) {
              $(".rerun-modal-body").html("请求失败！<br/>" + err.responseText);
            }
          }); 
        }

        function do_set_flow_success() {
          bootbox.dialog({
            title: "Process",
            message: '<div class="rerun-modal-body">' +
            '<div class="progress progress-striped active" style="margin-bottom:0;"><div class="progress-bar" style="width: 100%"></div></div>' +
            '</div>'})

          var request_params = {
            "flowInstId": <?php echo $flow_inst->id?>
          }

          var request = {};
          request.id = 3;
          request.method = "FlowService.SetFlowInstanceSuccess";
          request.params = [request_params];

          $.ajax({
            url: url, 
            data: JSON.stringify(request), 
            type: "POST",
            contentType: "application/json", 
            success: function(rpcRes) {
              if (rpcRes.error != null && rpcRes.error != "") {
                $(".rerun-modal-body").html("请求失败！<br/>" + rpcRes.error);
              } else {
                document.location.reload();
              }
            }, 
            error: function(err, status, thrown) {
              $(".rerun-modal-body").html("请求失败！<br/>" + err.responseText);
            }
          }); 
        }


        function show_log(date, ip, uuid) {

          bootbox.dialog({
            title: "运行日志",
            message: '<div class="log-modal-body">' +
            '<div class="progress progress-striped active" style="margin-bottom:0;"><div class="progress-bar" style="width: 100%"></div></div>' +
            '</div>'})


          var request_params = {
            "host": ip,
            "uuid": uuid, 
            "date": date
          }

          var request = {};
          request.id = 3;
          request.method = "FlowService.GetRemoteLog";
          request.params = [request_params];

          $.ajax({
            url: url, 
            data: JSON.stringify(request), 
            type: "POST",
            contentType: "application/json", 
            success: function(rpcRes) {
              if (rpcRes.error != null && rpcRes.error != "") {
                $(".log-modal-body").html("加载日志失败！<br/>" + rpcRes.error);
              } else {
                $(".log-modal-body").html('<div class="panel-group" id="accordion_log" role="tablist" aria-multiselectable="true">' +
                '<div class="panel panel-default">' +
                '  <div class="panel-heading" role="tab" id="headingCmdline">' +
                '    <h4 class="panel-title">' +
                '      <a role="button" data-toggle="collapse" data-parent="#accordion_log" href="#collapseCmdline" aria-expanded="true" aria-controls="collapseCmdline">' +
                '        Command Line' +
                '      </a>' +
                '    </h4>' +
                '  </div>' +
                '  <div id="collapseCmdline" class="panel-collapse collapse in" role="tabpanel" aria-labelledby="headingCmdline">' +
                '    <div class="panel-body">' +
                '      ' + htmlEncode(rpcRes.result.Cmdline) +
                '    </div>' +
                '  </div>' +
                '</div>' +
                '<div class="panel panel-default">' +
                '  <div class="panel-heading" role="tab" id="headingOutput">' +
                '    <h4 class="panel-title">' +
                '      <a class="collapsed" role="button" data-toggle="collapse" data-parent="#accordion_log" href="#collapseOutput" aria-expanded="true" aria-controls="collapseOutput">' +
                '        Standard Output' +
                '      </a>' +
                '    </h4>' +
                '  </div>' +
                '  <div id="collapseOutput" class="panel-collapse collapse" role="tabpanel" aria-labelledby="headingOutput">' +
                '    <div class="panel-body">' +
                '      <pre>' + htmlEncode(rpcRes.result.Output) +'</pre>' +
                '    </div>' +
                '  </div>' +
                '</div>' +
                '<div class="panel panel-default">' +
                '  <div class="panel-heading" role="tab" id="headingError">' +
                '    <h4 class="panel-title">' +
                '      <a class="collapsed" role="button" data-toggle="collapse" data-parent="#accordion_log" href="#collapseError" aria-expanded="false" aria-controls="collapseError">' +
                '        Error Output' +
                '      </a>' +
                '    </h4>' +
                '  </div>' +
                '  <div id="collapseError" class="panel-collapse collapse" role="tabpanel" aria-labelledby="headingError">' +
                '    <div class="panel-body">' +
                '      <pre>' + htmlEncode(rpcRes.result.Error) +'</pre>' +
                '    </div>' +
                '  </div>' +
                '</div>' +
                '</div>');
              }
            }, 
            error: function(err, status, thrown) {
              $(".log-modal-body").html("请求失败！<br/>" + err.responseText);
            }
          });
        }

        $(function(){
          $("pre.log").each(function(idx, ele) {
            that=$(ele); 
            log = that.text().replace(/\[(([0-9]{1,3}\.){3}[0-9]{1,3})\] RUN:UUID\[(([0-9]|[a-f]|\-){36})\]/g, "<a class='process_log'  data-toggle='tooltip' data-placement='bottom' title='Show logs' href='javascript:show_log(\"" + that.attr("date") + "\",\"$1\",\"$3\")'>$&</a>");
            log = log.replace(/\[(([0-9]{1,3}\.){3}[0-9]{1,3})\] RUN:DATE\[([0-9]{8})\]:UUID\[(([0-9]|[a-f]|\-){36})\]/g, "<a class='process_log'  data-toggle='tooltip' data-placement='bottom' title='Show logs' href='javascript:show_log(\"$3\",\"$1\",\"$4\")'>$&</a>");
            that.html(log);
          });

          $('[data-toggle="tooltip"]').tooltip();
        });
function htmlEncode(value){
  return $('<div/>').text(value).html();
}

      </script>

    </div>
  </div>
</div>
