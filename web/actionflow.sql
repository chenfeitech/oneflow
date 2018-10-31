
CREATE TABLE `tbServer` (
      `host` varchar(30) NOT NULL,
      `port` int(11) NOT NULL DEFAULT '36000',
      `username` varchar(20) DEFAULT 'root',
      `password` varchar(20) DEFAULT NULL,
      `crypto_password` varchar(100) NOT NULL DEFAULT '',
      `supervisors` varchar(200) DEFAULT NULL,
      `cont_failures` int(11) NOT NULL DEFAULT '0',
      `failure_uuid` varchar(40) DEFAULT NULL,
      `tags` varchar(100) CHARACTER SET utf8 COLLATE utf8_bin NOT NULL DEFAULT '|',
      PRIMARY KEY (`host`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `tbProducts` (
  `PId` varchar(50) NOT NULL DEFAULT '',
  `State` int(11) DEFAULT '0',
  `Ptype` int(11) DEFAULT '0',
  `Name` varchar(50) DEFAULT '0',
  `DBHost` varchar(30) NOT NULL,
  `DBName` varchar(50) DEFAULT '0',
  `StarLevel` int(11) NOT NULL DEFAULT '0',
  PRIMARY KEY (`PId`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8

-- Create syntax for TABLE 'tbAccessFlowLog'
CREATE TABLE `tbAccessFlowLog` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `pid` varchar(30) NOT NULL,
  `flow_id` varchar(128) NOT NULL,
  `log_date` date NOT NULL COMMENT '????',
  `flow_name` varchar(255) NOT NULL,
  `step_name` varchar(255) NOT NULL DEFAULT '',
  `log_msg` varchar(512) NOT NULL DEFAULT '',
  `creator` varchar(128) NOT NULL,
  `log_time` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `flow_name` (`pid`,`flow_id`,`log_date`)
) ENGINE=InnoDB AUTO_INCREMENT=11295 DEFAULT CHARSET=utf8;

-- Create syntax for TABLE 'tbFlow'
CREATE TABLE `tbFlow` (
  `id` varchar(50) NOT NULL DEFAULT '',
  `name` varchar(255) NOT NULL DEFAULT '',
  `description` varchar(255) NOT NULL DEFAULT '',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `creator` varchar(30) NOT NULL DEFAULT '',
  `start_timer` varchar(30) NOT NULL DEFAULT '',
  `next_run_time` datetime DEFAULT NULL,
  `startup_script` text NOT NULL,
  `last_run_time` datetime DEFAULT NULL,
  `last_run_log` text NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- Create syntax for TABLE 'tbFlowConf'
CREATE TABLE `tbFlowConf` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `process_id` varchar(50) NOT NULL DEFAULT '' COMMENT '流程id，如itsUpdateDailyProcess',
  `pid` varchar(30) NOT NULL COMMENT 'id，如xxx',
  `isactive` int(11) NOT NULL DEFAULT '0' COMMENT '表示状态是否有效',
  `start_time` varchar(50) NOT NULL DEFAULT '00 07 * * *' COMMENT '开始时间mm hh * * *',
  `active_date` date NOT NULL COMMENT '生效日期yyyy-mm-dd',
  `data_delay` varchar(128) NOT NULL DEFAULT '1' COMMENT '运行时延时的天数',
  `create_time` datetime DEFAULT NULL COMMENT '创建日期yyyy-mm-dd hh:mm:ss',
  `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新日期timestamp',
  `last_run_time` datetime DEFAULT NULL,
  `next_run_time` datetime DEFAULT NULL,
  `last_result` int(11) NOT NULL DEFAULT '0',
  `last_error` text,
  `creator` varchar(128) DEFAULT NULL COMMENT '创建人',
  `watcher` varchar(512) DEFAULT NULL COMMENT '管理员',
  PRIMARY KEY (`id`),
  UNIQUE KEY `process_id` (`process_id`,`pid`)
) ENGINE=InnoDB AUTO_INCREMENT=1441 DEFAULT CHARSET=utf8;

-- Create syntax for TABLE 'tbFlowSchdRunLog'
CREATE TABLE `tbFlowSchdRunLog` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `job_id` int(11) unsigned NOT NULL,
  `output` text NOT NULL,
  `errors` text,
  `schedule_time` datetime NOT NULL,
  `begin_time` datetime NOT NULL,
  `end_time` datetime NOT NULL,
  `result` int(11) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=222353 DEFAULT CHARSET=utf8;

-- Create syntax for TABLE 'tbFlowInst'
CREATE TABLE `tbFlowInst` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `flow_id` varchar(50) NOT NULL DEFAULT '',
  `pid` varchar(30) NOT NULL,
  `key` varchar(255) NOT NULL DEFAULT '',
  `running_day` date NOT NULL,
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `creator` varchar(30) NOT NULL DEFAULT '',
  `last_task_id` varchar(50) NOT NULL DEFAULT '',
  `last_task_state` int(11) NOT NULL,
  `last_update_time` datetime DEFAULT NULL,
  `state` int(11) NOT NULL DEFAULT '0',
  `startup_script` text NOT NULL,
  `begin_task` varchar(50) NOT NULL DEFAULT '',
  `end_task` varchar(50) NOT NULL DEFAULT '',
  PRIMARY KEY (`id`),
  UNIQUE KEY `flow_id` (`flow_id`,`pid`,`key`,`running_day`)
) ENGINE=InnoDB AUTO_INCREMENT=7726990 DEFAULT CHARSET=utf8;

-- Create syntax for TABLE 'tbJobRunLog'
CREATE TABLE `tbJobRunLog` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `job_id` int(11) unsigned NOT NULL,
  `output` mediumtext NOT NULL,
  `errors` mediumtext,
  `schedule_time` datetime NOT NULL,
  `begin_time` datetime NOT NULL,
  `end_time` datetime NOT NULL,
  `result` int(11) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=53 DEFAULT CHARSET=utf8;

-- Create syntax for TABLE 'tbJobSchedule'
CREATE TABLE `tbJobSchedule` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `pid` varchar(30) NOT NULL,
  `job_name` varchar(500) NOT NULL DEFAULT '',
  `pattern` varchar(100) NOT NULL DEFAULT '',
  `script` text NOT NULL,
  `last_run_time` datetime DEFAULT NULL,
  `next_run_time` datetime DEFAULT NULL,
  `creator` varchar(50) NOT NULL DEFAULT '',
  `create_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `enabled` tinyint(1) NOT NULL DEFAULT '0',
  `last_result` int(11) NOT NULL DEFAULT '0',
  `last_error` text,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=9 DEFAULT CHARSET=utf8;

-- Create syntax for TABLE 'tbStateReportLog'
CREATE TABLE `tbStateReportLog` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `pid` varchar(30) NOT NULL DEFAULT '',
  `key` varchar(255) NOT NULL DEFAULT '',
  `flow_id` varchar(50) NOT NULL DEFAULT '',
  `task_id` varchar(50) NOT NULL DEFAULT '',
  `state` int(11) NOT NULL,
  `creator` varchar(30) NOT NULL DEFAULT '',
  `date` date NOT NULL,
  `report_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `extra_data` text NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=4122793 DEFAULT CHARSET=utf8;

-- Create syntax for TABLE 'tbTask'
CREATE TABLE `tbTask` (
  `id` varchar(50) NOT NULL,
  `flow_id` varchar(50) NOT NULL DEFAULT '',
  `name` varchar(255) NOT NULL DEFAULT '',
  `description` text NOT NULL,
  `order_id` int(11) unsigned NOT NULL,
  `parent_id` int(11) unsigned NOT NULL DEFAULT '0',
  `script` text,
  `max_retries` int(10) unsigned NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`,`flow_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- Create syntax for TABLE 'tbTaskInst'
CREATE TABLE `tbTaskInst` (
  `flow_inst_id` int(11) unsigned NOT NULL,
  `task_id` varchar(50) NOT NULL DEFAULT '',
  `state` int(11) NOT NULL,
  `ready_time` datetime DEFAULT NULL,
  `running_time` datetime DEFAULT NULL,
  `succeed_time` datetime DEFAULT NULL,
  `failed_time` datetime DEFAULT NULL,
  `last_update_time` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00' ON UPDATE CURRENT_TIMESTAMP,
  `script_output` mediumtext,
  `retries` int(10) unsigned NOT NULL DEFAULT '0',
  `remote_exec_host` varchar(30) NOT NULL DEFAULT '',
  `remote_exec_uuid` varchar(50) NOT NULL DEFAULT '',
  `status_timestamp` bigint(20) NOT NULL DEFAULT '0',
  PRIMARY KEY (`flow_inst_id`,`task_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- Create syntax for TABLE 'tbTaskInstAlarm'
CREATE TABLE `tbTaskInstAlarm` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `flow_inst_id` int(11) unsigned NOT NULL,
  `task_id` varchar(50) NOT NULL DEFAULT '',
  `content` text NOT NULL,
  `timestamp` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=43 DEFAULT CHARSET=utf8;

-- Create syntax for TABLE 'tbToolRunLog'
CREATE TABLE `tbToolRunLog` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `tool_name` varchar(500) NOT NULL DEFAULT '',
  `arguments` text,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
