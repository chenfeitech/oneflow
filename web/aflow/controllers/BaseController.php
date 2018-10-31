<?php
date_default_timezone_set('Asia/Chongqing');
class BaseController extends CI_Controller {

	public $data = array();
	public $menu = array();

	function __construct(){
		parent::__construct();

		session_start();

		$this->data['userEnName'] = "helight";
		$this->data['header'] = 'comm/header';
		$this->data['footer'] = 'comm/footer';
		$this->data['list'] = 'comm/list';
		$this->username = "helight";
		setcookie('t_uid', $this->username);
	}

	function gbk2utf($para){
		return mb_convert_encoding($para,'utf8','gbk');
	}

	function checkipaddres($ipaddres) {
		$preg="/\A((([0-9]?[0-9])|(1[0-9]{2})|(2[0-4][0-9])|(25[0-5]))\.){3}(([0-9]?[0-9])|(1[0-9]{2})|(2[0-4][0-9])|(25[0-5]))\Z/";
		if(preg_match($preg,$ipaddres))return true;
		return false;
	}

	function check_auth($group_name){
		return true;
	}

	function flow_operation_log(){

	}
}
