<?php
require_once "BaseController.php";

class Flow extends BaseController
{
    public function __construct()
    {
        parent::__construct();
        $this->config->set_item('enable_query_strings', FALSE);
        $this->config->set_item('global_xss_filtering', TRUE);

        $this->load->model('flow_manager');
        $this->load->helper('url');
        $this->active_ip = "";

        $this->data['header'] = 'comm/header';
        $this->data['footer'] = 'comm/footer';
        $this->data['list'] = 'flow/list.php';

        error_reporting(E_ALL);
        ini_set('display_errors', 'on');

        $this->output->set_header("Access-Control-Allow-Origin: *");
        $this->output->set_header("Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept");
    }


    public function index() {
        $this->data["menu"] = 10;
        $this->flow_insts();
    }

    public function flows()
    {
        $this->data["menu"] = 30;
        if (!$this->check_auth('group_admin')) {
            $this->load->view('comm/noright',$this->data);
            return;
        }
        $flows = $this->flow_manager->find_flows();
        // var_dump($flows);
        $this->data['flows'] = $flows;
        $this->load->view('flow/flows', $this->data);
    }

    function show_flow_list()
    {
        $begin_date = $_GET['begin_date'];
        $end_date = $_GET['end_date'];
        $set_nu = $_GET['set_nu'];
        $process_type = $_GET['process_type'];
        $process_state_type = $_GET['process_state_type'];
        $pid = $_GET['pid'];

        $flow_insts = $this->flow_manager->find_inst_by_all($begin_date, $end_date,
                                                            $set_nu, $process_type, $process_state_type, $pid);
        // var_dump($flow_insts);
        foreach ($flow_insts as $inst) {
            $inst->last_task_state_text = $this->get_task_state_text($inst->last_task_state);
            $inst->task_insts = $this->flow_manager->find_task_inst_by_flow($inst->id, $inst->flow_id);
            foreach ($inst->task_insts as $fi) {
                $fi->state_text = $this->get_task_state_text($fi->state);
            }
        }
        $this->data['flow_insts'] = $flow_insts;
        $this->load->view('flow/flow_insts_list', $this->data);
    }
    //all - all
    public function flow_insts()
    {
        $this->data["menu"] = 10;
        if (!$this->check_auth('group_admin')) {
            $this->load->view('comm/noright',$this->data);
            return;
        }
        $this->load->helper('form');
        $this->load->library('form_validation');

        $options_of_set = array('all' => "所有集群");
        $options_of_set['0'] = "set0";
        $options_of_set['1'] = "set1";
        $options_of_set['2'] = "set2";
        $options_of_set['3'] = "set3";
        $this->data['set_list'] = $options_of_set;

        $flows = $this->flow_manager->get_flow_list();
        // var_dump($flows);
        $options_of_process = array('all' => "所有流程");
        foreach ($flows as $key => $value) {
            $options_of_process[$value->id] = $value->name;
        }
        $this->data['process_list'] = $options_of_process;

        $products = $this->flow_manager->get_products();
        // var_dump($products);
        $options_of_pid = array('all' => "所有");
        foreach ($products as $key => $value) {
            $options_of_pid[$value->PId] = $value->PId."(".$value->Name.")";
        }
        $this->data['pids'] = $options_of_pid;
        $options_of_process_state = array('-1' => "未成功");
        $options_of_process_state['all'] = "所有状态";
        $options_of_process_state['2'] = "成功";
        $options_of_process_state['3'] = "失败";
        $this->data['process_state_list'] = $options_of_process_state;

        $env="data_flow";
        $this->data['env'] = $env;
        $this->load->view('flow/flow_insts', $this->data);
    }

    public function insert_access_flow_log($pid, $flow_id, $flow_name, $step_name, $get_params){
        $log_msg='';
        foreach ($get_params as $key=>$value)
        {
            $log_msg=$log_msg.$key.':'.$value.'&';
        }
        $this->flow_manager->insert_access_flow_log($pid, $flow_id, date("Y-m-d"),
                                                    $flow_name, $step_name, $log_msg,
                                                    $this->data['userEnName'], date('y-m-d h:i:s',time()));
    }

//config_
    public function config_products()
    {
        $this->data["menu"] = 40;
        if (!$this->check_auth('group_admin')) {
            $this->load->view('comm/noright',$this->data);
            return;
        }
        $this->load->helper('form');
        $this->load->library('form_validation');
        $config_products = $this->flow_manager->find_config_products('all', '');

        $flows = $this->flow_manager->get_flow_list();
        $options_of_process = array('all' => "所有流程");
        foreach ($flows as $key => $value) {
            $options_of_process[$value->id] = $value->name;
        }
        $this->data['process_list'] = $options_of_process;
        $this->form_validation->set_rules('process_type', 'Process_Type', 'trim');
        $this->form_validation->set_rules('pid', 'Pid', 'trim');

        if ($this->form_validation->run() == TRUE)
        {
	      $config_products = $this->flow_manager->find_config_products($this->input->post('process_type'),
                                                                       $this->input->post('pid'));
        }

        $env="data_flow";
        $this->data['env'] = $env;
        $this->data['config_products'] = $config_products;
        $this->load->view('flow/config_products', $this->data);
    }

    // edit_config_products
    public function edit_config_products() {
        if (!$this->check_auth('group_admin')) {
            $this->load->view('comm/noright',$this->data);
            return;
        }

        $this->load->helper('form');
        $this->load->library('form_validation');
        $this->data["menu"] = 0;
        $config_products = $this->flow_manager->find_config_products('all', '');

    	$options_of_actions['0'] = "插入";
    	$options_of_actions['1'] = "删除";
    	$options_of_actions['2'] = "修改";
    	$this->data['action_list'] = $options_of_actions;

        $flows = $this->flow_manager->get_flow_list();
        $options_of_process = array();
        foreach ($flows as $key => $value) {
            $options_of_process[$value->id] = $value->name;
        }

        $this->data['process_list'] = $options_of_process;

    	$options_of_isactive['0'] = "无效";
    	$options_of_isactive['1'] = "有效";
    	$this->data['isactive_list'] = $options_of_isactive;

    	$this->form_validation->set_rules('action_type', 'Action_Type', 'trim');
        $this->form_validation->set_rules('process_type', 'Process_Type', 'trim');
        $this->form_validation->set_rules('pid', 'Pid', 'trim');
        $this->form_validation->set_rules('isactive', 'Isactive', 'trim');
        $this->form_validation->set_rules('start_time', 'Start_Time', 'trim');
        $this->form_validation->set_rules('active_date', 'Active_Date', 'trim');
        $this->form_validation->set_rules('data_delay', 'Data_Delay', 'trim');
        if ($this->form_validation->run() == TRUE)
        {
            switch($this->input->post('action_type')) {
            case '0':
                if ($this->input->post('pid') == '' ||
                    $this->input->post('start_time') == '' ||
                    $this->input->post('data_delay') == '') {
                    echo "input error: pid: ".$this->input->post('pid')
                        ." start_time: ".$this->input->post('start_time')
                        ." data_delay: ".$this->input->post('data_delay');
                    return;
                }

                $this->flow_manager->insert_config_products($this->input->post('process_type'),
                                                            $this->input->post('pid'),
                                                            $this->input->post('isactive'),
                                                            $this->input->post('start_time'),
                                                            $this->input->post('active_date'),
                                                            $this->input->post('data_delay'),
                                                            date('Y-m-d H:i:s'),
                                                            $this->data['userEnName']);
                break;
            case '1':
                $this->flow_manager->delete_config_products($this->input->post('process_type'),
                                                            $this->input->post('pid'));
                break;
            case '2':
                $this->flow_manager->update_config_products($this->input->post('process_type'),
                                                            $this->input->post('pid'),
                                                            $this->input->post('isactive'),
                                                            $this->input->post('start_time'),
                                                            $this->input->post('active_date'),
                                                            $this->input->post('data_delay'));
                break;
            default:

            }

            $config_products = $this->flow_manager->find_config_products('all', '');
        }

        $env="data_flow";
        $this->data['env'] = $env;
        $this->data['config_products'] = $config_products;
        $this->load->view('flow/edit_config_products', $this->data);
    }

    public function flow_inst($id)
    {
        $this->data["menu"] = 0;
        if (!$this->check_auth('group_user')) {
            $this->load->view('comm/noright',$this->data);
            return;
        }
        $query_result = $this->flow_manager->find_inst_by_id($id);

        if (count($query_result) == 0) {
            show_404();
            return;
        }

        $flow_inst = $query_result[0];

        $flow_inst->task_insts = $this->flow_manager->find_task_inst_by_flow($flow_inst->id, $flow_inst->flow_id);
        foreach ($flow_inst->task_insts as $fi) {
            $fi->state_text = $this->get_task_state_text($fi->state);
        }

        $this->data['flow_inst'] = $flow_inst;
        $this->load->view('flow/flow_inst', $this->data);
    }

    public function add_flow() {
        if (!$this->check_auth('group_user')) {
            $this->load->view('comm/noright',$this->data);
            return;
        }

        $this->data["menu"] = 1;
        $flow = (object) array('id' => null, 'name' => null, 'description' => null, 'tasks' => array());

        $this->data['flow'] = $flow;
        $this->load->view('flow/edit_flow', $this->data);
    }


    public function edit_flow($id) {
        if (!$this->check_auth('group_user')) {
            $this->load->view('comm/noright',$this->data);
            return;
        }

        $this->data["menu"] = 1;

        $query_result = $this->flow_manager->find_flow_by_id($id);

        if (count($query_result) == 0) {
            show_404();
            return;
        }
        $flow = $query_result[0];
        $flow->tasks = $this->flow_manager->find_tasks_by_flow($id);
        $this->data['flow'] = $flow;
        $this->load->view('flow/edit_flow', $this->data);
    }

    public function show_flow($id) {
        if (!$this->check_auth('group_user')) {
            $this->load->view('comm/noright',$this->data);
            return;
        }
        $this->data["menu"] = 1;
        $query_result = $this->flow_manager->find_flow_by_id($id);
        $products = $this->flow_manager->get_products();       

        if (count($query_result) == 0) {
            show_404();
            return;
        }
        $flow = $query_result[0];
        $type = 0;

        $flow->tasks = $this->flow_manager->find_tasks_by_flow($id);
        $this->data['flow'] = $flow;
        $this->data['products'] = $products;
        $this->load->view('flow/show_flow', $this->data);
    }
/*
    public function api() {
        if ($this->input->server('REQUEST_METHOD')!="POST") {
	        var_dump($this->input->server('REQUEST_METHOD'));
            return;
        }
        if (!$this->check_auth('group_user')) {
            $this->load->view('comm/noright',$this->data);
            return;
        }

        $url = 'http://localhost/data_flow/api';
        $data = file_get_contents('php://input');

        $headers = array();
        foreach (getallheaders() as $name => $value) {
            array_push($headers, "$name: $value");
        }

        $ch = curl_init($url);
        curl_setopt($ch, CURLOPT_CUSTOMREQUEST, "POST");
        curl_setopt($ch, CURLOPT_RETURNTRANSFER, true);
        curl_setopt($ch, CURLOPT_HTTPHEADER, $headers);
        curl_setopt($ch, CURLOPT_POSTFIELDS,$data);
        curl_setopt($ch, CURLOPT_HEADER, 1);
        curl_setopt($ch, CURLOPT_FOLLOWLOCATION, 1);
        $httpcode = curl_getinfo($ch, CURLINFO_HTTP_CODE);
        $response = curl_exec($ch);
        $header_size = curl_getinfo($ch, CURLINFO_HEADER_SIZE);
        $headers = substr($response, 0, $header_size);
        $body = substr($response, $header_size);

        if(curl_errno($ch)) {
            show_error(curl_error($ch));
        } else {
            foreach (explode("\r\n", $headers) as $header) {
                $header = trim($header);
                if ($header) header($header);
            }
            echo $body;
        }
        curl_close($ch);
        return;
    }
*/
    public function get_task_state_text($id)
    {
        $text = "";
        switch ($id) {
            case "0":
                $text = "Ready";
                break;
            case "1":
                $text = "Running";
                break;
            case "2":
                $text = "Succeed";
                break;
            case "3":
                $text = "Failed";
                break;

            default:
                $text = "";
                break;
        }
        return $text;
    }
}

if (!function_exists('getallheaders'))
{
    function getallheaders()
    {
        $headers = array();
       foreach ($_SERVER as $name => $value)
       {
       $iname = str_replace(' ', '-', ucwords(strtolower(str_replace('_', ' ', substr($name, 5)))));
	   // var_dump($iname);
           if (substr($name, 0, 5) == 'HTTP_')
           {
               $headers[$iname] = $value;
           }
       }
       return $headers;
    }
}

if (!function_exists('json_exception'))
{
    function json_exception($code, $message) {
        throw new Exception($message, $code);
    }
}
