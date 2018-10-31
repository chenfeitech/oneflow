<!doctype html>
<html>
<head>
<title>Terminal</title>
<meta charset="utf-8"/>
<script src="/thirdlib/static/jquery-1.11.0.min.js"></script>
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
</style>

    <!-- Bootstrap core CSS -->
    <link href="/thirdlib/bootstrap/css/bootstrap.min.css" rel="stylesheet">
    <script src="/thirdlib/static/jquery-1.11.0.min.js"></script>

<script src="/thirdlib/bootstrap/js/bootbox.min.js"></script>
<script src="/thirdlib/bootstrap/js/bootstrap.min.js"></script>
</head>
<body>

<div>
<textarea id="inputScript" rows="5">print(time.Now())</textarea>
</div>
<div>
<input type="button" id="run" value="Run"></input>
</div>
<h3>Result:</h3>
<pre>
<div id="result">
</div>
</pre>
    <script>
      var editor = CodeMirror.fromTextArea(document.getElementById("inputScript"), {
      lineNumbers: true,
      styleActiveLine: true,
      matchBrackets: true,
        theme: "mdn-like"
      });

  var url = "/flow/api";
$("#run").click(function () {
  var request_params = {
    "script": editor.getValue()
  }

  var request = {};
  request.id = 3;
  request.method = "FlowService.RunScript";
  request.params = [request_params];

  $.ajax({
    url: url, 
    data: JSON.stringify(request), 
    type: "POST",
    contentType: "application/json", 
    success: function(rpcRes) {
      if (rpcRes.error != null && rpcRes.error != "") {
        $("#result").html("请求失败！<br/>" + rpcRes.error);
      } else {
        $("#result").text(rpcRes.result.Output);
        $("#result").each(function(idx, ele) {
            that=$(ele); 
            log = that.text().replace(/\[(([0-9]{1,3}\.){3}[0-9]{1,3})\] RUN:UUID\[(([0-9]|[a-f]|\-){36})\]/g, "<a class='process_log'  data-toggle='tooltip' data-placement='bottom' title='Show logs' href='javascript:show_log(\"" + that.attr("date") + "\",\"$1\",\"$3\")'>$&</a>");
            log = log.replace(/\[(([0-9]{1,3}\.){3}[0-9]{1,3})\] RUN:DATE\[([0-9]{8})\]:UUID\[(([0-9]|[a-f]|\-){36})\]/g, "<a class='process_log'  data-toggle='tooltip' data-placement='bottom' title='Show logs' href='javascript:show_log(\"$3\",\"$1\",\"$4\")'>$&</a>");
            that.html(log);
          });

          $('[data-toggle="tooltip"]').tooltip();
      }
    }, 
    error: function(err, status, thrown) {
      $("#result").html("请求失败！<br/>" + err.responseText);
    }
  }); 
});


        function show_log(date, ip, uuid) {

          bootbox.dialog({
            title: "运行日志",
            message: '<div class="log-modal-body">' +
            '<div class="progress progress-striped active" style="margin-bottom:0;"><div class="progress-bar" style="width: 100%"></div></div>' +
            '</div>'})


          var request_params = {
            "host": ip,
            "uuid": uuid, 
            "date": "<?php echo date('Y-m-d')?>"
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
                '      ' + rpcRes.result.Cmdline +
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
                '      <pre>' + rpcRes.result.Output +'</pre>' +
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
                '      <pre>' + rpcRes.result.Error +'</pre>' +
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
    </script>
</body>
</html>
