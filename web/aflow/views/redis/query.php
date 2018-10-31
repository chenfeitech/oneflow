<?php $this->load->view($header);?>

<script type="text/javascript">

var formatJson = function(json, options) {
	var reg = null,
		formatted = '',
		pad = 0,
		PADDING = '    '; // one can also use '\t' or a different number of spaces
 
	// optional settings
	options = options || {};
	// remove newline where '{' or '[' follows ':'
	options.newlineAfterColonIfBeforeBraceOrBracket = (options.newlineAfterColonIfBeforeBraceOrBracket === true) ? true : false;
	// use a space after a colon
	options.spaceAfterColon = (options.spaceAfterColon === false) ? false : true;
 
	// begin formatting...
	if (typeof json !== 'string') {
		// make sure we start with the JSON as a string
		json = JSON.stringify(json);
	} else {
		// is already a string, so parse and re-stringify in order to remove extra whitespace
		json = JSON.parse(json);
		json = JSON.stringify(json);
	}
 
	// add newline before and after curly braces
	reg = /([\{\}])/g;
	json = json.replace(reg, '\r\n$1\r\n');
 
	// add newline before and after square brackets
	reg = /([\[\]])/g;
	json = json.replace(reg, '\r\n$1\r\n');
 
	// add newline after comma
	reg = /(\,)/g;
	json = json.replace(reg, '$1\r\n');
 
	// remove multiple newlines
	reg = /(\r\n\r\n)/g;
	json = json.replace(reg, '\r\n');
 
	// remove newlines before commas
	reg = /\r\n\,/g;
	json = json.replace(reg, ',');
 
	// optional formatting...
	if (!options.newlineAfterColonIfBeforeBraceOrBracket) {			
		reg = /\:\r\n\{/g;
		json = json.replace(reg, ':{');
		reg = /\:\r\n\[/g;
		json = json.replace(reg, ':[');
	}
	if (options.spaceAfterColon) {			
		reg = /\:/g;
		json = json.replace(reg, ':');
	}
 
	$.each(json.split('\r\n'), function(index, node) {
		var i = 0,
			indent = 0,
			padding = '';
 
		if (node.match(/\{$/) || node.match(/\[$/)) {
			indent = 1;
		} else if (node.match(/\}/) || node.match(/\]/)) {
			if (pad !== 0) {
				pad -= 1;
			}
		} else {
			indent = 0;
		}
 
		for (i = 0; i < pad; i++) {
			padding += PADDING;
		}
 
		formatted += padding + node + '\r\n';
		pad += indent;
	});
 
	return formatted;
};
  function do_query() {
  	var RedisId = $("#redis_id").val();
  	var QueryKey = $("#QueryKey").val();

  	if (RedisId == "all" || QueryKey == "") {
  		$("#queryret").html("请选择集群或者输入key");
  	}

  	var url = "/api/redis/get?key="+QueryKey+"&redis_id="+RedisId;
  	$.get(url, function(result) {     	
	  	var ret = "<pre>" + formatJson(result) + "</pre>";   
		$("#queryret").html(ret);
      });
  }
</script>

<div class="container-fluid">
    <div class="row">
        <div class="col-sm-9 col-sm-offset-3 col-md-10 col-md-offset-2 main">
            <?php $this->load->view("redis/list");?>
            <h2 class="page-header">Redis查询测试</h2>
			<div class="top-menu" style="overflow:hidden;padding-bottom:10px">
	<span><label style="width:100px;">集群ID : </label><?php echo form_dropdown('Redis_list', $Redis_list, set_value('redis_id'), 'style="width:240px;display:inline-block" id="redis_id" class="form-control"'); ?></span><p><p>
	<span><label style="width:100px;">查询Key: </label><input id="QueryKey" type="text" style="width:400px;display:inline-block" class="form-control" name="QueryKey"></span><p><p>
	<span><input style="margin-left:200px;width:80px" type="button" id="searchbtn" class="btn btn-success" value="查 询" onclick="do_query();return false;"></span>
</div>
<div>
	<div id="queryret" style="word-wrap:break-word">
	
	</div>
</div>
        </div>
    </div>
</div>
<?php $this->load->view($footer); ?>