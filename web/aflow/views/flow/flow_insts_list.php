<link href="/thirdlib/bootstrap/css/datepicker3.css" rel="stylesheet">
<script src="/thirdlib/bootstrap/js/bootstrap-datepicker.js" charset="UTF-8"></script>
<script src="/thirdlib/bootstrap/js/locales/bootstrap-datepicker.zh-CN.js" charset="UTF-8"></script>
<script src="/static/oneflow.js"></script>
<script src="/static/util.js"></script>
<script src="/thirdlib/bootstrap/js/bootstrap.min.js"></script>
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
      <section class="panel panel-default">
        <div class="table-responsive">
          <table class="table table-striped">
            <thead>
              <tr>
                <th>序号</th>
                <th>流程</th>
                <th>业务</th>
                <th>Key</th>
                <th>数据日期</th>
                <th>开始时间</th>
                <th>更新时间</th>
              </tr>
            </thead>
            <tbody>
      <?php foreach($flow_insts as $i => $inst):?>
              <tr data-href="/flow/flow_inst/<?php echo $inst->id?>" id="row_<?php echo $inst->id?>">
                <td rowspan="1"><?php echo $i?></td>
                <td><a href="/flow/show_flow/<?php echo $inst->flow_id?>"><?php echo $inst->name?></a></td>
                <td><?php echo $inst->pid?></td>
                <td><?php echo $inst->key?></td>
                <td><?php echo date("m-d", strtotime($inst->running_day))?></td>
                <td><?php echo date("m-d G:i:s", strtotime($inst->create_time))?></td>
                <td><?php echo date("m-d G:i:s", strtotime($inst->last_update_time))?></td>
              <td data-href="/flow/flow_inst/<?php echo $inst->id?>" id="row_<?php echo $inst->id?>_flowchart">
                <td colspan="8">
                <?php foreach($inst->task_insts as $ii => $ti):?>
                <?php if ($ii > 0) {?>
                  <span class="glyphicon glyphicon-arrow-right" aria-hidden="true"></span>
                <?php }?>
                  <span class="label label-<?php echo $ti->state?>" data-toggle="tooltip" data-placement="bottom" title="<?php echo $ti->state_text?>"><?php echo $ti->name?></span>
                <?php endforeach; ?>
                </td>
              </td>
              </tr>
      <?php endforeach; ?>
            </tbody>
          </table>
        </div>
      </section>
<script type="text/javascript">
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
    $(function(){
      $('.table tr[data-href]').each(function(){
        $(this).css('cursor','pointer').hover(
          function(){
            $(this).addClass('active');
          },
          function(){
            $(this).removeClass('active');
          }).click( function(){
            var strWindowFeatures = "location=yes,height=1000,width=900,scrollbars=yes,status=yes";
            window.open($(this).attr('data-href'), "_blank", strWindowFeatures);
            //document.location = $(this).attr('data-href');
          });
        });
    });
    function on_page_click(){
         $("#job_state_list button.btn-lg").removeClass('btn-lg');
         $(this).addClass('btn-lg');
    }

    $(document).ready(function(){
        $("#job_state_list button").click(on_page_click);
    });
var list_data = <?php echo json_encode($flow_insts)?>;
</script>
