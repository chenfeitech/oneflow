<?php
class Flow_manager extends CI_Model {
	public function __construct()
	{
		parent::__construct();
		$this->db = $this->load->database('flow', TRUE);
	}

	function getField($tableName, $column, $bykey, $order)
	{
		$this->db->select("$column");
		$this->db->from($tableName);
		$this->db->order_by($bykey, $order);
		$result = $this->db->get()->result();

		return $result;
	}

	function find_flows()
	{
		$query = $this->db->query(
			'SELECT `id`, `name`, `description`, `create_time`, `creator`'
			.', `start_timer`, `next_run_time`, `startup_script`, `last_run_time`'
			.', `last_run_log` FROM `tbFlow`'
		);
		return $query->result();
	}

	function find_inst_by_date($date)
	{
		$query = $this->db->query(
			'SELECT fi.`id`, fi.`flow_id`, f.`name`, f.`description`'
			.', fi.`pid`,  fi.`key`, fi.`running_day`, fi.`create_time` '
			.', ifnull(t.name, fi.`last_task_id`) last_task_id, fi.`last_task_state`'
			.', fi.last_update_time from `tbFlowInst` fi'
			.' inner join `tbFlow` f on fi.flow_id = f.id left join `tbTask` t on fi.last_task_id=t.id'
			.'  and t.`flow_id` = fi.`flow_id` left join `tbProducts` tbg on tbg.`PId` = fi.`pid` WHERE fi.`running_day`=? ORDER BY last_update_time DESC'
			, array($date));
		return $query->result();
	}

	function find_inst_by_all($begin_date, $end_date, $set_nu, $flow_id, $process_state_type, $pid)
	{
		$sql = 'SELECT fi.`id`, fi.`flow_id`, f.`name`, f.`description`'
			.', fi.`pid`,  fi.`key`, fi.`running_day`, fi.`create_time`'
			.', ifnull(t.name, fi.`last_task_id`) last_task_id, fi.`last_task_state`'
			.', fi.last_update_time from `tbFlowInst` fi'
			.' inner join `tbFlow` f on fi.flow_id = f.id left join `tbTask` t on fi.last_task_id=t.id'
			.' and t.`flow_id` = fi.`flow_id` left join `tbProducts` tbg on tbg.`PId` = fi.`pid` WHERE (fi.`running_day` between \''.$begin_date.' 00:00:00\' and \''.$end_date.' 23:59:59\')'
			.' and (\''.$set_nu.'\' = \'all\')'
			.' and (\''.$flow_id.'\' = \'all\' or f.`id` = \''.$flow_id.'\')'
			.' and (\''.$process_state_type.'\' = \'all\' or fi.`state` = \''.$process_state_type.'\' or (\''.$process_state_type.'\' = -1 and fi.`state` in (0,1,3)))'
			.' and (\''.$pid.'\' = \'all\' or fi.`pid` = \''.$pid.'\')'
			.' ORDER BY fi.`running_day` desc, fi.`flow_id`, last_update_time DESC';
		// var_dump($sql);
		$query = $this->db->query($sql);
		return $query->result();
	}

	//config_products
	function find_config_products($flow_id, $pid)
	{
		$sql = 'SELECT  fc.`process_id`, f.`name`, fc.`pid`, fc.`pid`, fc.`isactive`, fc.`start_time`, fc.`active_date`, fc.`data_delay`, fc.`create_time`, fc.`update_time`'
			.' from `tbFlowConf` fc left join `tbFlow` f on fc.`process_id` = f.`id`'
			.' where (\''.$pid.'\' = \'\' or fc.`pid` = \''.$pid.'\')'
			.' and (\''.$flow_id.'\' = \'all\' or f.`id` = \''.$flow_id.'\')';
		$query = $this->db->query($sql);
		return $query->result();
	}
	//insert_config_products
	function insert_config_products($process_type, $pid, $isactive, $start_time, $active_date, $data_delay, $create_time,$userEnName)
	{
		$sql = 'insert into `tbFlowConf`(`process_id`, `pid`, `isactive`, `start_time`, `active_date`, `data_delay`, `create_time`, creator)'
			.' values'
			.' (\''.$process_type.'\', \''.$pid.'\', '.$isactive.', \''.$start_time.'\', \''.$active_date.'\', \''.$data_delay.'\', \''.$create_time.'\',\''.$userEnName.'\')';
		$query = $this->db->query($sql);
	}

	//delete_config_products
	function delete_config_products($process_type, $pid)
	{
		$sql = 'delete from `tbFlowConf` where `process_id` = \''.$process_type.'\' and `pid` = \''.$pid.'\'';
		$query = $this->db->query($sql);
	}

	//update_config_products
	function update_config_products($process_type, $pid, $isactive, $start_time, $active_date, $data_delay)
	{
		$sql = 'update `tbFlowConf` set `isactive` = '.$isactive.' where `process_id` = \''.$process_type.'\' and `pid` = \''.$pid.'\'';
		$query = $this->db->query($sql);

		$sql = 'update `tbFlowConf` set `start_time` = \''.$start_time.'\', next_run_time=null where `process_id` = \''.$process_type.'\' and `pid` = \''.$pid.'\' and \''.$start_time.'\' != \'\'';
		$query = $this->db->query($sql);

		$sql = 'update `tbFlowConf` set `active_date` = \''.$active_date.'\' where `process_id` = \''.$process_type.'\' and `pid` = \''.$pid.'\'';
		$query = $this->db->query($sql);

		$sql = 'update `tbFlowConf` set `data_delay` = \''.$data_delay.'\' where `process_id` = \''.$process_type.'\' and `pid` = \''.$pid.'\' and \''.$data_delay.'\' != \'\'';
		$query = $this->db->query($sql);
	}

	//
	function find_inst_by_date_set($date, $setnu)
	{
		$query = $this->db->query(
			'SELECT fi.`id`, fi.`flow_id`, f.`name`, f.`description`'
			.', fi.`pid`,  fi.`key`, fi.`running_day`, fi.`create_time`'
			.', ifnull(t.name, fi.`last_task_id`) last_task_id, fi.`last_task_state`'
			.', fi.last_update_time from `tbFlowInst` fi'
			.' inner join `tbFlow` f on fi.flow_id = f.id left join `tbTask` t on fi.last_task_id=t.id'
			.' and t.`flow_id` = fi.`flow_id` left join `tbProducts` tbg on tbg.`PId` = fi.`pid` WHERE fi.`running_day`=?'
			.' ORDER BY last_update_time DESC'
			, array($date, $setnu));
		return $query->result();
	}

	function find_inst_by_date_id($date)
	{
		$query = $this->db->query(
			'SELECT fi.`id`, fi.`flow_id`, f.`name`, f.`description`'
			.', fi.`pid`,  fi.`key`, fi.`running_day`, fi.`create_time`'
			.', ifnull(t.name, fi.`last_task_id`) last_task_id, fi.`last_task_state`'
			.', fi.last_update_time from `tbFlowInst` fi'
			.' inner join `tbFlow` f on fi.flow_id = f.id left join `tbTask` t on fi.last_task_id=t.id'
			.'  and t.`flow_id` = fi.`flow_id` left join `tbProducts` tbg on tbg.`PId` = fi.`pid` WHERE fi.`running_day`=? ORDER BY last_update_time DESC'
			, array($date));
		return $query->result();
	}

	//
	function find_inst_by_date_id_set($date, $setnu)
	{
		$query = $this->db->query(
			'SELECT fi.`id`, fi.`flow_id`, f.`name`, f.`description`'
			.', fi.`pid`,  fi.`key`, fi.`running_day`, fi.`create_time`'
			.', ifnull(t.name, fi.`last_task_id`) last_task_id, fi.`last_task_state`'
			.', fi.last_update_time from `tbFlowInst` fi'
			.' inner join `tbFlow` f on fi.flow_id = f.id left join `tbTask` t on fi.last_task_id=t.id'
			.'  and t.`flow_id` = fi.`flow_id` left join `tbProducts` tbg on tbg.`PId` = fi.`pid`'
			.' WHERE fi.`running_day`=? and tbg.`StorageSet` = ? ORDER BY last_update_time DESC'
			, array($date, $setnu));
		return $query->result();
	}

	//to be modified
	function find_inst_by_id($id)
	{
		$query = $this->db->query(
			'SELECT fi.`id`, fi.`flow_id`, f.`name`, f.`description`, fi.`creator` as `last_operator`'
			.', fi.`pid`,  fi.`key`, fi.`running_day`,  fi.`state`, fi.`create_time`'
			.', ifnull(t.name, fi.`last_task_id`) last_task_id, fi.`last_task_state`'
			.', fi.last_update_time from `tbFlowInst` fi'
			.' inner join `tbFlow` f on fi.flow_id = f.id left join `tbTask` t on fi.last_task_id=t.id'
			.'  and t.`flow_id` = fi.`flow_id` WHERE fi.`id`=? ORDER BY last_update_time DESC'
			, array($id));
		return $query->result();
	}

	//to be modified
	function find_task_inst_by_flow($inst_id, $flow_id)
	{
		$query = $this->db->query(
			'SELECT t.`id`, t.`flow_id`, t.`name`, t.`description`, t.`order_id`'
			.', t.`parent_id`, ti.`state`, ti.ready_time, ti.running_time, ti.succeed_time'
			.', ti.failed_time, ti.last_update_time, ti.script_output FROM `tbTask` t '
			.' inner join `tbFlowInst` fi ON t.`flow_id`=fi.flow_id'
			.' AND fi.id=? left join `tbTaskInst` ti ON t.id = ti.`task_id`'
			.' AND fi.id = ti.flow_inst_id'
			.' WHERE t.`flow_id`=? ORDER BY t.order_id'
			, array($inst_id, $flow_id));
		return $query->result();
	}

	//to be modified
	function find_flow_by_id($flow_id)
	{
		$query = $this->db->query(
			'SELECT `id`, `name`, `description`, `create_time`, `creator`'
			.', `start_timer`, `next_run_time`, `startup_script`, `last_run_time`'
			.', `last_run_log` FROM `tbFlow`'
			.' WHERE `id`=?'
			, array($flow_id));
		return $query->result();
	}

	function find_tasks_by_flow($flow_id)
	{
		$query = $this->db->query(
			'SELECT `id`, `flow_id`, `name`, `description`, `order_id`'
			.', `parent_id`, `script`, `max_retries` FROM `tbTask`'
			.' WHERE `flow_id`=? ORDER BY order_id'
			, array($flow_id));
		return $query->result();
	}

	function get_products_list()
	{
		$sql='select PId,concat(PId,\'(\',Name,\')\') as Name from tbProducts order by PId';
		$query = $this->db->query($sql);
		return $query->result();
	}

	function get_productsinfo_by_pid($pid)
	{
		$query = $this->db->query('SELECT * from tbProducts WHERE `PId`=?', array($pid));
		return $query->result();
	}

	function get_products()
	{
		$sql='select PId, State, Name, StarLevel from tbProducts order by PId';
		$query = $this->db->query($sql);
		return $query->result();
	}

	function get_products_type()
	{
		$sql='select PId, State as Ptype from tbProducts';
		$query = $this->db->query($sql);
		return $query->result();
	}

	function get_products_config($env,$pid)
	{
		$sql = 'select *'
			.'  from '.$env.'.tbProducts where \'all\' = \''.$pid.'\' or PId = \''.$pid.'\'';
		$query = $this->db->query($sql);
		return $query->result();
	}

	function update_products_resource($pid,$col_name,$col_value)
	{
		$sql = 'update tbProducts set '.$col_name.'=\''.$col_value.'\' where PId=\''.pid.'\'';
		$this->db->query($sql);
	}

	function get_production_id($pid)
	{
		$sql='SELECT `ProductionId` FROM tbFromFront_ProductionName WHERE (`TargetId`="'.$pid.'") limit 1';
		$query = $this->db->query($sql);
		return $query->result();
	}

	function get_tbFlowConf($pid, $process_id)
	{
		$query = $this->db->query('select a.*,SUBSTRING_INDEX(SUBSTRING_INDEX(a.start_time, \' \', 2),\' \',-1) as start_hour,SUBSTRING_INDEX(a.start_time, \' \', 1) as start_minute,b.Name,c.name from tbFlowConf a join tbProducts b on a.pid=b.PId join tbFlow c on a.process_id=c.id where a.pid=? and a.process_id=?', array($pid, $process_id));
		return $query->result();
	}

	function update_tbFlowConf($pid,$process_id,$isactive,$start_time,$active_date,$data_delay,$watcher,$userEnName)
	{
		$query = $this->db->query('select pid from tbFlowConf where pid = ? and process_id=?',array($pid,$process_id));
		$result = $query -> result();
		if (count($result) > 0) {
			$query = $this->db->query('update tbFlowConf set isactive=?, start_time=?, active_date=?, data_delay=?, watcher=?, update_time=now(), next_run_time=null where pid=? and process_id=?'
				, array($isactive,$start_time,$active_date,$data_delay,$watcher,$pid,$process_id));
		}else{
			$query = $this->db->query('insert into tbFlowConf (pid,process_id,isactive,start_time,active_date,data_delay,creator,watcher,create_time,update_time) values (?,?,?,?,?,?,?,?,now(),now())'
				,array($pid,$process_id,$isactive,$start_time,$active_date,$data_delay,$userEnName,$watcher));
		}
	}

	function update_tbFlowConf_isactive($pid,$process_id,$isactive)
	{
		$query = $this->db->query('select pid from tbFlowConf where pid = ? and process_id=?',
			array($pid, $process_id));
		$result = $query -> result();
		if (count($result) > 0) {
			$query = $this->db->query('update tbFlowConf set isactive=?, update_time=now() where pid=? and process_id=?'
				, array($isactive,$pid,$process_id));
		}
	}
	function update_tbProducts_CanLevelStat($pid,$CanLevelStat)
	{
		$sql = 'update tbProducts set CanLevelStat=? where PId=?';
		$query = $this->db->query($sql, array($CanLevelStat,$pid));
	}

	function get_flow_list()
	{
		$sql = 'select * from tbFlow';
		$query = $this->db->query($sql);
		return $query->result();
	}

	function get_flows_by_pid($pid)
	{
		$sql = 'SELECT a.*,b.name FROM tbFlowConf a join tbFlow b on a.process_id=b.id WHERE a.pid = ?';
		$query = $this->db->query($sql,array($pid));
		return $query->result();
	}

	function update_product_state($pid, $state)
	{
		$sql = 'update tbProducts set state=? where PId=?';
		$query = $this->db->query($sql, array($state, $pid));
	}

	function get_access_flow_log($pid, $flow_id, $begin_date, $end_date)
	{
		$sql = 'SELECT * FROM tbAccessFlowLog WHERE ("all" = ? or pid = ?) and flow_id = ? and log_date between ? and ? order by id desc';
		$query = $this->db->query($sql,array($pid, $pid, $flow_id, $begin_date, $end_date));
		return $query->result();
	}

	function insert_access_flow_log($pid, $flow_id, $log_date, $flow_name, $step_name, $log_msg, $creator, $log_time)
	{
		$sql = 'insert into tbAccessFlowLog(pid, flow_id, log_date, flow_name, step_name, log_msg, creator, log_time) values (?,?,?,?,?,?,?,?)';
		$query = $this->db->query($sql,array($pid, $flow_id, $log_date, $flow_name, $step_name, $log_msg, $creator, $log_time));
	}

	function get_server_info($host)
	{
		$sql = 'SELECT * FROM tbServer WHERE host = ?';
		$query = $this->db->query($sql,array($host));
		return $query->result();
	}

    function find_flow_instance_by_flow_day_products($flow_id, $day, $pids)
    {
        if (!is_array($pids) || count($pids) == 0) {
            return array();
        }

        $in = join(',', array_fill(0, count($pids), '?'));

        $sql = <<<SQL
    SELECT fi.`id`, ? as `flow_id`, f.`name`, f.`description`
    , tbg.`PId` as pid,  fi.`key`, date(?) as `running_day`, fi.`create_time`
    , fi.`state`
    , fi.last_update_time from `tbFlowInst` fi
     inner join `tbFlow` f on fi.flow_id = f.id
      and fi.flow_id = ? AND fi.`running_day`=date(?)
      right join `tbProducts` tbg on tbg.`PId` = fi.`pid` where tbg.`PId` in ($in)
SQL;

        $result = $this->db->query($sql, array_merge(array($flow_id, $day, $flow_id, $day), $pids))->result();
        foreach ($result as $instance) {
            if (!is_null($instance->id)) {
                $instance->tasks = $this->find_task_inst_by_flow($instance->id, $instance->flow_id);
                foreach ($instance->tasks as $task) {
                    $task->script_output = null;
                }
            }
        }
        return $result;
    }

    function find_flow_instance_by_id($flow_inst_id)
    {
        $sql = <<<SQL
    SELECT fi.`id`, f.`id` as `flow_id`, f.`name`, f.`description`
    , tbg.`PId` as pid,  fi.`key`, `running_day` as `running_day`, fi.`create_time`
    , fi.`state`
    , fi.last_update_time from `tbFlowInst` fi
     inner join `tbFlow` f on fi.flow_id = f.id
      and fi.id = ?
      inner join `tbProducts` tbg on tbg.`PId` = fi.`pid`
SQL;

        $result = $this->db->query($sql, array($flow_inst_id))->result();
        foreach ($result as $instance) {
            if (!is_null($instance->id)) {
                $instance->tasks = $this->find_task_inst_by_flow($instance->id, $instance->flow_id);
            }
        }
        return $result;
    }


    function find_flow_instance($flow_id, $pid, $running_day, $key)
    {
        $sql = <<<SQL
    SELECT fi.`id`, f.`id` as `flow_id`, f.`name`, f.`description`
    , tbg.`PId` as pid,  fi.`key`, `running_day` as `running_day`, fi.`create_time`
    , fi.`state`
    , fi.last_update_time from `tbFlowInst` fi
     inner join `tbFlow` f on fi.flow_id = f.id
      and fi.flow_id = ? and fi.pid= ? and fi.running_day = ? and fi.key = ?
      inner join `tbProducts` tbg on tbg.`PId` = fi.`pid`
SQL;

        $result = $this->db->query($sql, array($flow_id, $pid, $running_day, $key))->result();
        foreach ($result as $instance) {
            if (!is_null($instance->id)) {
                $instance->tasks = $this->find_task_inst_by_flow($instance->id, $instance->flow_id);
            }
        }
        return $result;
    }
}
