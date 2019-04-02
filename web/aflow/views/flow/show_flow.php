<?php $this->load->view($header);?>

<link href="/thirdlib/bootstrap/css/datepicker3.css" rel="stylesheet">
<script src="/thirdlib/bootstrap/js/bootstrap-datepicker.js" charset="UTF-8"></script>
<script src="/thirdlib/bootstrap/js/locales/bootstrap-datepicker.zh-CN.js" charset="UTF-8"></script>
<script src="/thirdlib/bootstrap/js/bootstrap.min.js"></script>

<link rel="stylesheet" href="/thirdlib/static/CodeMirror/lib/codemirror.css">
<script src="/thirdlib/static/CodeMirror/lib/codemirror.js"></script>
<script src="/thirdlib/static/CodeMirror/mode/lua/lua.js"></script>
<link rel="stylesheet" href="/thirdlib/static/CodeMirror/theme/mdn-like.css">
<style type="text/css">
.CodeMirror {
  border: 1px solid #eee;
  min-height: 100px;
  height: 100%;
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

.bs-callout-info {
  border-left-color: #1b809e !important;
}
.bs-callout {
  padding: 20px;
  margin: 20px 0;
  border: 1px solid #eee;
  border-left-width: 5px;
  border-radius: 3px;
}
</style>


<div class="container-fluid">
  <div class="row">
    <?php $this->load->view($list);?>

    <div class="col-sm-9 col-sm-offset-3 col-md-10 col-md-offset-2 main">
          <h1 class="page-header"><span id="page_title"><?php echo ($flow->name)?></span>
<div class="pull-right">
<a class="btn btn-primary" role="button" href="/flow/edit_flow/<?php echo ($flow->id)?>">
编辑
</a>
</div>
</h1>

<div class="bs-callout bs-callout-info">
    <h4>ID: <?php echo ($flow->id)?></h4>
    <p></p>
    <p>PID：
          <select id="pid" name="pid" class="input">
          <?php foreach($products as $product):?>
            <option value="<?php echo $product->PId ?>"><?php echo $product->PId ?> - <?php echo $product->Name ?></option>
          <?php endforeach; ?>
          </select></p>
    <p>Day：
          <input type="text" class="input-date" size="20" value="<?php echo date('Y-m-d',time()-1*24*60*60);?>" id="day" name="day"/>
          </p>
     <div>
<textarea id="inputScript" rows="3"><?php echo ($flow->startup_script)?></textarea>
</div>
     <p><input class="btn btn-primary" id="start_flow_btn" type="button" value="Start flow" onclick="start_flow(this)"></p>
</div>

<div class="row">
  <div class="col-md-8 flow">
    <ol class="example vertical">
      <?php foreach($flow->tasks as $i => $task):?>
      <li>
        <span class="name"><b><?php echo $task->name?></b></span>
        <div class="pull-right">
          <span class="name">ID: <?php echo $task->id?></span>
        </div>
        <div style="margin: 5px;"></div>
      </li>
      <?php endforeach; ?>

    </ol>
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
<script type="text/javascript">
  var editor = CodeMirror.fromTextArea(document.getElementById("inputScript"), {
      lineNumbers: true,
      styleActiveLine: true,
      matchBrackets: true,
        theme: "mdn-like"
      });

  // var url = "/flow/api";
  var url = "/oneflow/StartFlow";
  function start_flow() {
    var request_params = {
      "id": "<?php echo ($flow->id)?>",
      "creator": "<?php echo $userEnName; ?>",
      "pid": $("#pid").val(),
      "date": $("#day").val(),
      "startup_script": editor.getValue()
    }

    var request = {};
    request.id = <?php echo rand()?>;
    request.method = "FlowService.StartFlow";
    request.params = [request_params];

    console.log(JSON.stringify(request));

    var $btn = $("#start_flow_btn").button('loading');
    $.ajax({
      url: url,
      data: JSON.stringify(request_params),
      type: "POST",
      contentType: "application/json",
      success: function(rpcRes) {
        $btn.button('reset');
        if (rpcRes.code != 0) {
          $("#myMsgBoxTitle").text("错误");
          $("#myMsgBoxBody").text("启动失败：" + rpcRes.message);
          $('#myMsgBox').modal();
        } else {
          document.location.href = "/flow/flow_inst/" + rpcRes.flow_inst_id;

          // $btn.button('reset');
          // $("#myMsgBoxTitle").text("成功");
          // $("#myMsgBoxBody").text("流程启动成功！");
          // $('#myMsgBox').modal();
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


  $(function(){
      $('.input-date').datepicker({
        format: "yyyy-mm-dd",
        todayBtn: "linked",
        language: "zh-CN",
        autoclose: true,
        todayHighlight: true
      });
    });
</script>
