<?php $this->load->view($header);?>
    <div class="container-fluid">
      <div class="row">
        <div class="col-sm-9 col-sm-offset-3 col-md-10 col-md-offset-2 main">
        <?php $this->load->view($list);?>
        <div align="center" style="padding:20px;">
    <p><img src="/thirdlib/static/imgs/403.png" /></p><br><br>
    <p >sorry,您没有权限,开通权限请联系:helight.</p>
		<a href="/">返回首页</a>
    </div>
<?php $this->load->view($footer);?>
