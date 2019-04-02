<?php $this->load->view($header);?>


<link rel="stylesheet" href="/thirdlib/static/CodeMirror/lib/codemirror.css">
<script src="/thirdlib/static/jquery-sortable.js"></script>
<script src="/thirdlib/static/CodeMirror/lib/codemirror.js"></script>
<script src="/thirdlib/static/CodeMirror/mode/lua/lua.js"></script>
<link rel="stylesheet" href="/thirdlib/static/CodeMirror/theme/mdn-like.css">

<style type="text/css">
body.dragging, body.dragging * {
  cursor: move !important;
}

.dragged {
  position: absolute;
  top: 0;
  opacity: .5;
  z-index: 2000;
}

ol.vertical {
  margin: 0 0 9px 0;
  min-height: 30px;
}
  ol.vertical li {
    display: block;
    margin: 5px;
    padding: 5px;
    border: 1px solid #cccccc;
    color: #0088cc;
    background: #eeeeee;
    min-height: 60px;
  }
  ol.vertical li.placeholder {
    position: relative;
    margin: 0;
    padding: 0;
    border: none; }
    ol.vertical li.placeholder:before {
      position: absolute;
      content: "";
      width: 0;
      height: 0;
      margin-top: -5px;
      left: -5px;
      top: -4px;
      border: 5px solid transparent;
      border-left-color: red;
      border-right: none; }

ol.vertical li:hover
{
  background-color:#cccccc;
}
ol i.move {
  cursor: pointer;
}
ol span.remove {
  cursor: pointer;
}

li.selected {
  background:#aaaaaa !important;
}

.CodeMirror {
  border: 1px solid #eee;
  min-height: 100px;
  height: 100%;
}
</style>
<script src="/thirdlib/bootstrap/js/bootstrap.min.js"></script>

<div class="container-fluid">
<div class="row">
  <?php $this->load->view($list);?>
  <div class="col-sm-9 col-sm-offset-3 col-md-10 col-md-offset-2 main">
    <h1 class="page-header"><span id="page_title"><?php echo ($flow->name) ?: '新流程' ?></span>
<div class="pull-right">
<button class="btn btn-primary" role="button" onclick="add_task(event);">
添加任务
</button>

<button id="save_btn" class="btn btn-primary" role="button" data-loading-text="Saving..." onclick="save(event);">
保存
</button>
</div>
</h1>
  <div class="col-md-4 flow">
    <ol class='example vertical'>
    <?php foreach($flow->tasks as $i => $task):?>
      <li onclick="on_task_click(event);" task_id="<?php echo $task->id?>" new_task_id="<?php echo $task->id?>" description="<?php echo $task->description?>" max_retries="<?php echo $task->max_retries?>">
        <i class="glyphicon glyphicon-move move" aria-hidden="true"></i> <span class="name"><?php echo $task->name?></span>
        <div class="pull-right"><span onclick="on_task_remove(event);" class="glyphicon glyphicon glyphicon-trash remove" aria-hidden="true"></span></div>
        <div style="display:none"><textarea class="script"><?php echo $task->script?></textarea></div>
      </li>
    <?php endforeach; ?>
    </ol>
    <ol class='template' style="display:none">
      <li onclick="on_task_click(event);" task_id="" new_task_id="" description="" max_retries="0">
      <i class="glyphicon glyphicon-move move" aria-hidden="true"></i> <span class="name"></span>
      <div class="pull-right"><span onclick="on_task_remove(event);" class="glyphicon glyphicon glyphicon-trash remove" aria-hidden="true"></span></div>
      <div style="display:none"><textarea class="script"></textarea></div>
      </li>
    </ol>
  </div>
  <div class="col-md-8">
    <div id="control_panel" class="panel panel-default">
      <div class="panel-body">
        <p class="title">流程属性</p>
        <form class="form-horizontal">
          <div class="form-group">
            <label for="inputID" class="col-sm-2 control-label">ID</label>
            <div class="col-sm-10">
              <input type="text" class="form-control" id="inputID" onkeyup="on_id_change(this, event);" onchange="on_id_change(this, event);">
            </div>
          </div>
          <div class="form-group">
            <label for="inputName" class="col-sm-2 control-label">名称</label>
            <div class="col-sm-10">
              <input type="text" class="form-control" id="inputName" onkeyup="on_name_change(this, event);" onchange="on_name_change(this, event);">
            </div>
          </div>
          <div class="form-group">
            <label for="inputDescription" class="col-sm-2 control-label">描述</label>
            <div class="col-sm-10">
              <textarea class="form-control" id="inputDescription" rows="5" onchange="on_description_change(this, event);"></textarea>
            </div>
          </div>
          <div class="form-group" id="formMaxRetries" style="display: none;">
            <label for="inputMaxRetries" class="col-sm-2 control-label">失败重试次数</label>
            <div class="col-sm-10">
              <select class="form-control" id="inputMaxRetries" onchange="on_max_retries_change(this, event);">
                <option>0</option>
                <option>1</option>
                <option>2</option>
                <option>3</option>
                <option>4</option>
                <option>5</option>
              </select>
            </div>
          </div>
          <div class="form-group">
            <label for="inputScript" class="col-sm-2 control-label">脚本</label>
            <div class="col-sm-10">
              <textarea class="form-control" id="inputScript" rows="5"></textarea>
            </div>
          </div>
        </form>
      </div>
    </div>
  </div>
</div>
<div class="modal fade" id="myMsgBox" tabindex="-1" role="dialog" aria-labelledby="myMsgBoxTitle" aria-hidden="true">
  <div class="modal-dialog">
    <div class="modal-content">
      <div class="modal-header">
        <button type="button" class="close" data-dismiss="modal" aria-label="Close"><span aria-hidden="true">&times;</span></button>
        <h4 class="modal-title" id="myMsgBoxTitle">Modal title</h4>
      </div>
      <div class="modal-body">
      <div id="myMsgBoxBody"></div>
      </div>
      <div class="modal-footer">
        <button type="button" class="btn btn-default" data-dismiss="modal">Close</button>
      </div>
    </div>
  </div>
</div>
</div>

<script>
  var flow_id = "<?php echo $flow->id?>";
  var new_flow_id = "<?php echo ($flow->id) ?: 'NEW_FLOW' ?>";
  var li_template = null;
  var selected_item = null;
  var task_num = 1;
  var flow_name = "<?php echo ($flow->name) ?: '新流程' ?>";
  var flow_description = <?php echo json_encode($flow->description); ?>;
  var flow = <?php echo json_encode($flow); ?>;
  var delete_item = [];

  function change_select(sel) {
    if (selected_item == sel) {
      return;
    } else if (selected_item != null && sel != null && selected_item.get(0) == sel.get(0)) {
      return;
    }
    save_values();
    if (selected_item != null) {
      selected_item.removeClass("selected");
    }
    selected_item = sel;
    if (selected_item != null) {
      selected_item.addClass("selected");
    };

    var control_panel = $("#control_panel");
    tmp_control_panel = control_panel.clone(true, true).removeAttr("id");
    tmp_control_panel.insertBefore(control_panel);

    if (sel == null) {
      control_panel.find(".title").text("流程属性");
      control_panel.find("#inputID").val(new_flow_id).attr("readonly", flow_id=="" ? null : "");
      control_panel.find("#inputName").val(flow_name);
      control_panel.find("#inputDescription").val(flow_description);
      control_panel.find("#formMaxRetries").hide();
      editor.setValue(flow.startup_script ?flow.startup_script: "");
    } else {
      control_panel.find(".title").text("任务属性");
      control_panel.find("#inputID").val(selected_item.attr("new_task_id")).attr("readonly", null);
      control_panel.find("#inputName").val(sel.find(".name").text());
      control_panel.find("#inputDescription").val(selected_item.attr("description"));
      control_panel.find("#formMaxRetries").show();
      control_panel.find("#inputMaxRetries").val(selected_item.attr("max_retries"));
      editor.setValue(sel.find(".script").val());
    }
    control_panel.finish().hide().slideDown(200, function(){control_panel.show();});
    tmp_control_panel.finish().show().slideUp(200, function(){tmp_control_panel.remove();});
  }

  function save_values() {
    if (selected_item != null) {
      selected_item.find(".script").val(editor.getValue());
    } else {
      flow.startup_script = editor.getValue();
    }
  }

  function add_task(event) {
    event.stopPropagation();

    var lis = $("ol.example li");
    var new_task_ids = {};
    for (var i = 0; i < lis.length; i++) {
      var ele = $(lis[i]);
      new_task_ids[ele.attr("new_task_id")] = true;
    }

    for (var id = lis.length+1; new_task_ids["TASK_" + id]; ++id) {
    }

    new_li = li_template.clone().appendTo($("ol.example"));
    new_li.attr("new_task_id", "TASK_" + id);
    new_li.find(".name").text("任务" + id);
    new_li.click(function (e)
      {
        change_select($(e.target));
      });
    change_select(new_li);
    ++task_num;
  }


  function on_id_change(obj, event) {
    obj = $(obj);
    if (selected_item == null) {
      new_flow_id = obj.val();
    } else {
      selected_item.attr("new_task_id", obj.val());
    }
  }

  function on_name_change(obj, event) {
    obj = $(obj);
    if (selected_item == null) {
      flow_name = obj.val();
      $("#page_title").text(flow_name);
      document.title = flow_name;
    } else {
      selected_item.find(".name").text(obj.val());
    }
  }

  function on_description_change(obj, event) {
    obj = $(obj);
    if (selected_item == null) {
      flow_description = obj.val();
    } else {
      selected_item.attr("description", obj.val());
    }
  }

  function on_max_retries_change(obj, event) {
    obj = $(obj);
    console.log(obj.val());
    if (selected_item == null) {
    } else {
      selected_item.attr("max_retries", obj.val());
    }
  }

  function on_task_click(e)
  {
    if (e.target.tagName == "LI")
    {
      change_select($(e.target));
    }
  }


  function on_task_remove(e)
  {
    var ele = $(e.target).parents("li");
    var task_id = ele.attr("task_id");
    if (task_id != "") {
      delete_item.push(task_id);
    }
    ele.remove();
  }


  function msg_box(title, message) {
    $("#myMsgBoxTitle").text(title);
    $("#myMsgBoxBody").html(message);
    $('#myMsgBox').modal();
  }

  function save(event) {
    event.stopPropagation();

    if (new_flow_id.trim() == "") {
        change_select(null);
        msg_box("错误", "流程ID不能为空！");
        return;
    }
    save_values();

    var url = "/oneflow/API";

    var lis = $("ol.example li");
    var tasks = new Array(lis.length);
    var new_task_ids = {};
    for (var i = 0; i < lis.length; i++) {
      var ele = $(lis[i]);
      var task_id = ele.attr("task_id");
      var new_task_id = ele.attr("new_task_id").trim();
      if (new_task_id == "") {
        change_select(ele);
        msg_box("错误", "任务ID不能为空！");
        return;
      }
      if (new_task_ids[new_task_id]) {
        msg_box("错误", "任务ID重复！");
        change_select(ele);
        return;
      }
      new_task_ids[new_task_id] = true;

      if (task_id != new_task_id && task_id != "") {
        delete_item.push(task_id);
      }

      tasks[i] = {
        "id": new_task_id,
        "name": ele.find(".name").text(),
        "description": ele.attr("description"),
        "max_retries": ele.attr("max_retries"),
        "script": ele.find(".script").val()
      };
    }

    var request_params = {
      "id": flow_id == "" ? new_flow_id : flow_id,
      "name": flow_name,
      "creator": "<?php echo $userEnName ?>",
      "description": flow_description,
      "startupScript": (flow.startup_script?flow.startup_script:""),
      "tasks": tasks,
      "deleteTaskIds": delete_item
    }

    var request = {};
    request.id = 2;
    if (flow_id != "") {
      // url += "UpdateFlow"
      request.method = "FlowService.UpdateFlow";
    } else {
      // url += "AddFlow"
      request.method = "FlowService.AddFlow";
    }
    request.params = [request_params];


    console.log(JSON.stringify(request));


    var $btn = $("#save_btn").button('loading');
    $.ajax({
      url: url,
      data: JSON.stringify(request),
      type: "POST",
      contentType: "application/json",
      success: function(rpcRes) {
        $btn.button('reset');
        if (rpcRes.error != null && rpcRes.error != "") {
          $("#myMsgBoxTitle").text("错误");
          $("#myMsgBoxBody").text("保存失败：" + rpcRes.error);
          $('#myMsgBox').modal();
        } else {
          document.location.href = "/flow/show_flow/" + new_flow_id;
        }
      },
      error: function(err, status, thrown) {
        $btn.button('reset');
        $("#myMsgBoxTitle").text("请求失败");
        $("#myMsgBoxBody").html(err.statusText + "<br/>" + err.responseText);
        $('#myMsgBox').modal();
      }
    });
  }

  var editor = CodeMirror.fromTextArea(document.getElementById("inputScript"), {
      lineNumbers: true,
      styleActiveLine: true,
      matchBrackets: true,
        theme: "mdn-like"
      });

  $(function  () {
    var control_panel = $("#control_panel");

    control_panel.find(".title").text("流程属性");
    control_panel.find("#inputID").val(new_flow_id).attr("readonly", flow_id=="" ? null : "");
    control_panel.find("#inputName").val(flow_name);
    control_panel.find("#inputDescription").val(flow_description);

    editor.setValue((flow.startup_script?flow.startup_script:""));

    li_template = $("ol.template li");

    $("ol.example").sortable({
      group: 'no-drop',
      handle: 'i.glyphicon'
    });

  $(document).click(function (e)
  {
      var container = $("ol.example");
      var control_panel = $("#control_panel");
      var msgBox = $("#myMsgBox");
      if (!container.is(e.target) // if the target of the click isn't the container...
          && container.has(e.target).length === 0 // ... nor a descendant of the container
          && !control_panel.is(e.target)
          && control_panel.has(e.target).length === 0
          && !msgBox.is(e.target)
          && msgBox.has(e.target).length === 0
          )
      {
          change_select(null);
      }
  });
  })
</script>
</div>
</script>
