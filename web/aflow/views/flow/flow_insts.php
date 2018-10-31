<?php $this->load->view($header);?>
<div class="container-fluid">
  <div class="row">
  	<div class="col-sm-9 col-sm-offset-3 col-md-10 col-md-offset-2 main">
 	 <h1 class="page-header"><span id="page_title">流程实例列表</span> <div class="pull-right"></div></h1>
     <?php $this->load->view($list);?>
     <?php $this->load->view('flow/flow_insts_content', $this->data);?>
    </div>
  </div>
</div>
<?php $this->load->view($footer); ?>
