<?php $this->load->view($header);?>

<style type="text/css">

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
          <h1 class="page-header"><span id="page_title">List</span> <div class="pull-right">
      <a class="btn btn-primary" href="/flow/add_flow" role="button">Add Flow</a>
      </div>
      </h1>

      <table class="table table-striped">
        <thead>
          <tr>
            <th>#</th>
            <th>Name</th>
            <th>Description</th>
            <th>Creator</th>
            <th>Create At</th>
          </tr>
        </thead>
        <tbody>

      <?php foreach($flows as $i => $flow):?>
          <tr data-href="/flow/show_flow/<?php echo $flow->id?>">
            <td><?php echo $i+1?></td>
            <td><?php echo $flow->name?></td>
            <td><?php echo $flow->description?></td>
            <td><?php echo $flow->creator?></td>
            <td><?php echo $flow->create_time?></td>
          </tr>
      <?php endforeach; ?>

        </tbody>
      </table>

        </div>
  </div>
</div>


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
            }
        );
    });

  //$(window).scroll(bindScroll);
});

function loadMore()
{
   console.log("More loaded");
   $("body").append("<div>");
   $(window).bind('scroll', bindScroll);
 }


function bindScroll(){
   if($(window).scrollTop() + $(window).height() > $(document).height() - 100) {
       $(window).unbind('scroll');
       loadMore();
   }
}


</script>
