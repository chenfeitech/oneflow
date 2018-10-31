<table class="table table-hover table-condensed table-striped">
    <thead>
        <tr>
            <th>业务ID</th>
            <th>别名</th>
            <th>业务名称</th>
        </tr>
    </thead>
    <tbody>
    	<?php 
    		while($row = mysql_fetch_array($result))
			{
				echo "<tr>"
					."<td>".$row['Id']."</td>"
					."<td>".$row['AliasName']."</td>"
					."<td>".$row['Name']."</td>"
					.'<td><a href="#" onclick=\'javascript:{$("#tips_info_box").hide();$("#pid").val('.$row['Id'].');return false;}\'>选择</a></td>'
					."</tr>";
			}
  		?>
    </tbody>
</table>
