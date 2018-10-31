<!DOCTYPE html>
<html lang="zh">
  <head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <meta name="description" content="">
    <meta name="author" content="">
    <link rel="shortcut icon" href="/static/imgs/favicon.png">

    <title>Aflow</title>

    <!-- Bootstrap core CSS -->
    <link href="/thirdlib/bootstrap/css/bootstrap.min.css" rel="stylesheet">

    <!-- Custom styles for this template -->
    <link href="/static/dashboard.css" rel="stylesheet">
    <script src="/thirdlib/static/jquery-1.11.0.min.js"></script>
    <script src="/thirdlib/highcharts/js/highcharts.js"></script>
    <script src="/static/util.js"></script>
  </head>

  <body>

    <div class="navbar navbar-inverse navbar-fixed-top" role="navigation">
      <div class="container-fluid">
        <div class="navbar-header">
          <a class="navbar-brand" href="/">数据Flow运营系统</a>
        </div>
        <div class="navbar-collapse collapse">
        <ul class="nav navbar-nav navbar-right">
        <li <?php if($menu == 10) echo 'class="active"' ?>><a href="/flow/flow_insts">运行管理</a></li>
        <li <?php if($menu == 30) echo 'class="active"' ?>><a href="/flow/flows">流程配置</a></li>
        <li <?php if($menu == 40) echo 'class="active"' ?>><a href="/flow/config_products">Time配置</a></li>
        <li><a><span class="badge pull-right">helight</span></a></li>
        </ul>
        </div>
      </div>
    </div>
