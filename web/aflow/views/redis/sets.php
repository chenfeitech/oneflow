<?php $this->load->view($header);?>
<link href="/static/dashboard.css" rel="stylesheet">

<script src="/static/util.js"></script>
<script src="/static/commonalert.js"></script>
<style type="text/css">
.page{display:block;width:200px;margin: 0px auto;}
form label.field {float: left;width: 20%;text-align: right;padding: 1px;margin-right: 10px;}
input.text {border: 1px solid #B2D0EA;padding: 4px;font-size: 12px;line-height: 100%;width: 300px;}
.input_select {border: 1px solid #B2D0EA;font-size: 12px;line-height: 25px;width: 300px;height: 25px;}
form div {margin: 5px;}
.edit {padding: 4px 16px;height: 30px;font-size: 14px;border-radius:7px;}
.btn-save{
    padding: 6px 25px;
    height: 36px;
    font-size: 14px;
    border-radius: 3px;
    color: #fff;
    background-color: #0c89d4;
}
.btn-reset{
    padding: 6px 25px;
    height: 36px;
    font-size: 14px;
    border-radius: 3px;
    color: #fff;
    background-color: #ffa940;
}
th {
    text-align: center;
}
</style>
<script type="text/javascript">
var currenturl='<?php echo $currenturl ?>';
$(document).ready(function(){
  setTable();
});
function editredisserver(RedisId,RedisName,state,Host,User,Passwd,Capacity,port){
    var dialogbox=commonalert.cover();
    var $testDiv=$($('#faddForm').get(0).outerHTML);//测试的div
          var self=this;
            var testDiv=$testDiv.get(0).outerHTML;
             dialogbox.confirm(testDiv,{
                    title:'添加新的集群',
                    confirmTitle:'确认',
                    cancelTitle:'关闭',
                    onConfirm:function(content){
                       saveredisserver(content);
                    },
                    onSuccess:function($cnt){
                        $cnt.find('#faddForm').css({
                            'display':'block',
                            'position':'relative',
                            'left':'0px',
                            'top':'0px',
                            'margin':'0px'
                        });
                        $cnt.find("#RedisId").val(RedisId);
            $cnt.find("#Redisname").val(RedisName);
            $cnt.find("#Redisstatus").val(state);
            $cnt.find("#host").val(Host);
                        $cnt.find("#Port").val(port);
            $cnt.find("#user").val(User);
            $cnt.find("#password").val(Passwd);
            $cnt.find("#Capacity").val(Capacity);
            $cnt.find("#RedisId").attr('disabled',true);
          }
            });
}
function setTable(){
    //设置表格宽度
    $("table th:eq(0)").width("5%");
    $("table th:eq(1)").width("20%");
    $("table th:eq(2)").width("10%");
    $("table th:eq(3)").width("20%");
    $("table th:eq(4)").width("5%");
    $("table th:eq(5)").width("10%");
    $("table th:eq(6)").width("10%");
    $("table th:eq(7)").width("10%");
    $("table th:eq(7)").width("10%");
}
function saveredisserver(data){
    var dialogbox=commonalert.cover();
    var RedisId = data.find("#RedisId").val();
    var RedisName=data.find("#Redisname").val();
      var state=data.find("#Redisstatus").val();
        var Port=data.find("#Port").val();
      var Host = data.find("#host").val();
      var User = data.find("#user").val();
      var Passwd = data.find("#password").val();
      var Capacity=data.find("#Capacity").val();
      if (RedisId==''){
          var url=currenturl+"/redis/addredisServer";
        }
      else{
          var url=currenturl+"/redis/editredisServer";
      }
      var data={
         'RedisId':RedisId,
         'RedisName':RedisName,
         'state':state,
         'Host':Host,
           'Port':Port,
         'User':User,
         'Passwd':Passwd,
         'Capacity':Capacity
      }
      $.ajax({
                  type:'POST',
                  url:url,
                  data:data,
                  dataType:'json',
                  success:function (result) {
                        dialogbox.alert(result.msg);
                      showPageInfo(result.data);
                  }
      });
  }
function addredisserver(){
    var dialogbox=commonalert.cover();
    var $testDiv=$($('#faddForm').get(0).outerHTML);//测试的div
          var self=this;
            var testDiv=$testDiv.get(0).outerHTML;
             dialogbox.confirm(testDiv,{
                    title:'添加新的集群',
                    confirmTitle:'确认',
                    cancelTitle:'关闭',
                    onConfirm:function(content){
                       saveredisserver(content);
                    },
                    onSuccess:function($cnt){
                        $cnt.find('#faddForm').css({
                            'display':'block',
                            'position':'relative',
                            'left':'0px',
                            'top':'0px',
                            'margin':'0px'
                        });
            }
                });
}
function searchdataserver(){
      var url = currenturl+"/redis/searchdataserver";
      var RedisId= $("#redis_id").val();
      $.post(url,
          {RedisId:RedisId},
          function(result) {
              showPageInfo(result);
      });
  }
function setColor(obj,act,flag){
    if(act==1){
        $(obj).css('background-color','#d9f0ff');
    }
    if(act==0&&flag==1)
    {
        $(obj).css('background-color','#FFF');
    }
    if(act==0&&flag==2)
    {
       $(obj).css('background-color','#F1F7FC');
    }
}
//}
function showPageInfo(result,formname,res){
    var tablediv = document.getElementById("tableinfo");
    var pagediv = document.getElementById("pageinfo");
    var dataObj= typeof(result) == "string"? eval("("+result+")") : result;
    var pid=dataObj.pid;
    renderPageInfo(result,tablediv,pagediv);
    setTable();
}
// 定制表格数据和pageinfo
function renderPageInfo(result,tablediv,pagediv){
    if(result=="fail"){
        alert("操作失败！");
    }
    var dataObj= typeof(result) == "string"? eval("("+result+")") : result;
    var table = dataObj.table;
    var recordCount =  dataObj.recordCount;
    var pagelink = dataObj.pagelink;
    var pageid = dataObj.pageid;
    var pid = dataObj.pid;
    $("#pageid").val(pageid);
    $("#pid").val(pid);

    tablediv.innerHTML = "<p>" + table + "</p>";
    pagediv.innerHTML = '<p style="color:#1F74C7;font-size:16px;"><span class="page"> 共&nbsp;'+recordCount+'&nbsp;条记录&nbsp;&nbsp;'+pagelink+'</span></p>';
}
</script>
<div class="container-fluid">
    <div class="row">
        <div class="col-sm-9 col-sm-offset-3 col-md-10 col-md-offset-2 main">
            <?php $this->load->view("redis/list");?>
            <h2 class="page-header">Redis集群管理</h2>
      <div>
    <input style="margin-left:10px;width:100px" type="button" class="btn btn-success" value="添加新集群" onclick="addredisserver();return false;">
    <input type="text" id="redis_id" style="margin-left:10px;width:100px;display:inline-block" class="form-control" name="redis_id" value="<?php echo $RedisId ?>" placeholder="redis_id">
    <input style="margin-left:10px;width:80px" type="button" class="btn btn-success" value="查 询" onclick="searchdataserver();return false;">
</div>
<div>
  <div id="tableinfo" style="word-wrap:break-word">
          <p> <?php echo $tableInfo;?> </p>
  </div>
      <!--pagination-->
  <div id="pageinfo">
          <p style="color:#1F74C7;font-size:16px;"> <span class="page">共<?php echo $recordCount;?>条记录<?php echo $page_links;?></span> </p>
  </div>
      <!--end pagination-->
</div>
<div id="shadow" style="display: none; "></div>
<div id="f_layer" style="display: none; "></div>
<!-- 添加表单 -->
<form id="faddForm" name="faddForm" action="" method="post" class="modal-box" validate="true" style="display:none" >
    <div style="display:none;">
        <label class="field">RedisId</label>
        <input type="text" id="RedisId" name="RedisId" class="text" size="25" rule="required" tip="请输入!" onblur=""></input>
        <span style="width:100px;display:inline-block;"><div class="checRedisId" style="display:none;color:red;">请输入Redis名称</div></span>
    </div>
    <div>
        <label class="field">Redis名称</label>
        <input type="text" id="Redisname" name="Redisname" class="text" size="25" rule="required" tip="请输入!" onblur=""></input>
        <span style="width:100px;display:inline-block;"><div class="checkRedisname" style="display:none;color:red;">请输入Redis名称</div></span>
    </div>
    <div>
        <label class="field">Redis状态</label>
        <select id="Redisstatus" name="Redisstatus" class="input_select">
            <option value="1">激活</option>
            <option value="0">未激活</option>
        </select>
        <span style="width:100px;display:inline-block;"><div class="checkRedisstatus" style="display:none;color:red;"></div></span>
    </div>
    <div>
        <label class="field">Host</label>
        <input type="text" id="host" name="host" class="text" size="25" rule="required" tip="请输入!" onblur=""></input>
        <span style="width:100px;display:inline-block;"><div class="checkhostname" style="display:none;color:red;">请输入host名称</div></span>
    </div>
     <div>
        <label class="field">Port</label>
        <input type="text" id="Port" name="Port" class="text" size="25" rule="required" tip="请输入!" onblur=""></input>
        <span style="width:100px;display:inline-block;"><div class="checkProtname" style="display:none;color:red;">请输入host名称</div></span>
    </div>
    <div>
        <label class="field">用户名</label>
        <input type="text" id="user" name="user" class="text" size="25" onblur=""></input>
        <span style="width:100px;display:inline-block;"><div class="checkusername" style="display:none;color:red;">请输入帐号</div></span>
    </div>
    <div>
        <label class="field">Redis容量</label>
        <input type="text" id="capacity" name="capacity" class="text" size="25" ></input>
        <span style="width:100px;display:inline-block;"><div class="" style="display:none;color:red;"></div></span>
    </div>
</form>

        </div>
    </div>
</div>
<?php $this->load->view($footer); ?>
