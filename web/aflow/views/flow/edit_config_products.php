<?php $this->load->view($header);?>
<link href="/thirdlib/bootstrap/css/datepicker3.css" rel="stylesheet">
<script src="/thirdlib/bootstrap/js/bootstrap-datepicker.js" charset="UTF-8"></script>
<script src="/thirdlib/bootstrap/js/locales/bootstrap-datepicker.zh-CN.js" charset="UTF-8"></script>
<style type="text/css">
  .label {
    display: inline;
    padding: .2em .6em .3em;
    font-size: 75%;
    font-weight: 700;
    line-height: 1;
    color: #fff;
    text-align: center;
    white-space: nowrap;
    vertical-align: baseline;
    border-radius: .25em;
    background-color: #f0ad4e;
  }
  .label- {
    background-color: #aaaaaa;
  }
  .label--1 {
    background-color: #aaaaaa;
  }
  .label-0 {
    background-color: #5bc0de;
  }
  .label-1 {
    background-color: #6ABE05;
  }
  .label-2 {
    background-color: #337ab7;
  }
  .label-3 {
    background-color: #d9534f;
  }
  .wrapper {
    padding: 15px;
  }
  .m-b-xs {
    margin-bottom: 5px;
  }
  select.input-sm {
    height: 30px;
    line-height: 30px;
  }
  .v-middle {
    vertical-align: middle !important;
  }
  .inline {
    display: inline-block !important;
  }
  .input-s-sm {
    width: 120px;
  }
</style>
<div class="container-fluid">
  <div class="row">
  <?php $this->load->view($list);?>



    <div class="col-sm-9 col-sm-offset-3 col-md-10 col-md-offset-2 main">
      <h1 class="page-header"><span id="page_title">业务配置列表</span> <div class="pull-right"></div></h1>

	<div class="pull-right">
	<a class="btn btn-primary" href="/flow/config_products" role="button">返回</a>
	</div>

      <div class="table-responsive">
          <form class="navbar-form navbar-left" role="form" action = "/flow/edit_config_products" method="post">
              <div class="form-group">
                  <input type="hidden" class="form-control" name='env' value="<?php echo $env?>"/>
                  <?php echo form_dropdown('process_type', $process_list, set_value('process_type'), 'id="process_type" class="form-control"'); ?>
                  <input type="text" style="width:100px;" class="form-control" name='pid' value="<?php echo set_value('pid')?>" placeholder="pid"/>
                  <?php echo form_dropdown('isactive', $isactive_list, set_value('isactive'), 'id="isactive" class="form-control"'); ?>
		  开始时间
                  <input type="text" style="width:95px;" class="form-control" name='start_time' value="<?php echo set_value('start_time')?>" placeholder="mm hh * * *"/>
		  生效日期
                  <input type="text" class="input-date form-control" id="active_date" value="<?php if(set_value('active_date') == null)
                      {echo date('Y-m-d',time()-24*60*60);} else
                      {echo set_value('active_date');}?>" name="active_date"/>
		  延迟天数
                  <input type="text" style="width:50px;" class="form-control" name='data_delay' value="<?php echo set_value('data_delay')?>" placeholder="ie,1"/>
		  创建时间
                  <?php echo form_dropdown('action_type', $action_list, set_value('action_type'), 'id="action_type" class="form-control"'); ?>
                  <button type="submit" class="btn btn-success">go!</button>
              </div>
          </form>
          <br><br><br>
          <div id="processes_detail"></div>
      </div>


      <section class="panel panel-default">
        <div class="table-responsive">
          <table class="table table-striped">
            <thead>
              <tr>
		<th>序号</th>
                <th>流程名称</th>
                <th>业务ID</th>
		<th>业务名称</th>
                <th>是否有效</th>
                <th>开始时间</th>
                <th>生效时间</th>
                <th>延迟天数</th>
		<th>创建时间</th>
		<th>更新时间</th>
              </tr>
            </thead>
            <tbody>
      <?php foreach($config_products as $i => $inst):?>
		<tr>
                <td rowspan="1"><?php echo $i?></td>

                <td><a href="/flow/show_flow/<?php echo $inst->process_id?>"><?php echo $inst->name?></a></td>
                <td><?php echo $inst->pid?></td>
		<th><?php echo $inst->pid?></th>
		<td><?php echo $inst->isactive?></td>
		<td><?php echo $inst->start_time?></td>
		<td><?php echo $inst->active_date?></td>
		<td><?php echo $inst->data_delay?></td>
                <td><?php if($inst->create_time == NULL) {echo "not set yet";} else {echo date("m-d G:i:s", strtotime($inst->create_time));}?></td>
                <td><?php echo date("m-d G:i:s", strtotime($inst->update_time))?></td>

		</td>
		</tr>
      <?php endforeach; ?>
            </tbody>
          </table>
        </div>
      </section>
    </div>
  </div>
</div>
</div>


<script src="/static/oneflow.js"></script>
<script src="/static/util.js"></script>
<?php $this->load->view($footer); ?>
<script type="text/javascript">
  $(function(){
    $('.table tr[data-href]').each(function(){
      $(this).css('cursor','pointer').hover(
        function(){
          $(this).addClass('active');
        },
        function(){
          $(this).removeClass('active');
        }).click( function(){
          document.location = $(this).attr('data-href');
        });
      });
  });
    $(function(){
      $('a').popover({trigger : 'click'});
    });
    $('[data-toggle="tooltip"]').tooltip();

    $(function(){
      $('.input-date').datepicker({
        format: "yyyy-mm-dd",
        todayBtn: "linked",
        language: "zh-CN",
        autoclose: true,
        todayHighlight: true
      });
    });

    function on_page_click(){
         $("#job_state_list button.btn-lg").removeClass('btn-lg');
         $(this).addClass('btn-lg');
    }

    $(document).ready(function(){
        $("#job_state_list button").click(on_page_click);
    });
var list_data = <?php echo json_encode($config_products)?>;
</script>
