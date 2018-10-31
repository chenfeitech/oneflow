commonalert={
	cover:function(){
		var openTemplate='<div class=\"ui-modal\" style=\"min-width:400px;background-color:white\"><div class=\"header\" style=\"width:auto;\"><a href=\"#\">\u5173\u95ed</a><span>open-Title</span></div><div class=\"content\"><div style=\"max-height:800px;overflow-y:auto;overflow-x:hidden;width:auto;\">提示语。。。</div><div class=\"footer text-right\" style=\"margin-top:20px;text-align:right;width:auto\"><a href=\"javascript:void(0)\" class=\"btn btn-reset\" style=\"margin-right:30px;\">open-Cancel</a><a href=\"javascript:void(0)\" class=\"btn btn-save\">open-OK</a></div></div></div>';
	  	var openCover='<div id="open-cover" style="z-index:500;display:none"></div>';
		var currentZIndex=100000; 
		//判断是否为HTML DOM
	 	var isDOM = ( typeof HTMLElement === 'object' ) ?
	                function(obj){
	                    return obj instanceof HTMLElement;
	                } :
	                function(obj){
	                    return obj && typeof obj === 'object' && obj.nodeType === 1 && typeof obj.nodeName === 'string';
	                } 

		var _this=this;
		var open={};
		function openObject(){
		    this.setFrame=function()
		      {		
		            var _frame=$(openTemplate).attr('id','openFrame_'+(new Date()).getTime());
		            var _cover=$(openCover).attr('id','openCover_'+(new Date()).getTime());
		            
		              $(document.getElementsByTagName('body')[0]).append(_frame);
		            //_frame.css
		            _frame.css({
		              'font-size':'14px',
		              'font-family':'Microsoft YaHei, Arial, sans-serif'
		            });
		              
		                //cover
		            $(document.getElementsByTagName('body')[0]).append(_cover);
		            _cover.ready(function(){
		            _cover.css(
		              {
		              'position':'fixed',
		              'background-color':'#777',
		              'opacity':0.6,
		              'filter':'alpha(opacity=60)',
		              'left':'0px',
		              'top':'0px',
		              'border':'none'
		              }
		            );
		            });
		            
		            //z-index
		            _cover.css('z-index',(currentZIndex+=2));
		            _frame.css('z-index',(currentZIndex+=2));
		            
		            _this.openShow(_frame,{frame:_frame,cover:_cover});
		            _this.openHide(_frame,{frame:_frame,cover:_cover,notRemove:true});//notRemove:不remove框架
		              return {frame:_frame,cover:_cover};
		      };      
		    this.alert=function()
		      {// 格式1：(msg,onConfirm,onClose,onSuccess)  格式2：(msg,{...})
		         var self=this;
		           var _opt={};
		           var args=['','onConfirm','onClose','onSuccess'];
		           if(!arguments.length){
		             throw 'alert 必须至少带一个参数';
		             return;
		           }
		           if(typeof arguments[0] !== 'number' && typeof arguments[0]!=='string' && ((typeof arguments[0]=='object') && !arguments[0].css) && !isDOM(arguments[0])){
		             throw 'alert 第一个参数必须为string,number,jquery object,DOM之一';
		             return;
		           }
		           _opt.msg=arguments[0];
		           if(arguments.length>=2){
		             if(typeof arguments[1] == 'function'){
		               for(var i=1;i<arguments.length;i++){
		                 _opt[args[i]]=arguments[i];
		               }
		             }else{
		               _opt=$.extend(_opt,arguments[1]);
		             }
		           }
		         $.extend(true,_opt,self.setFrame());
		           _this.openInit(_opt.frame,_opt,'alert');
		      };
		    this.confirm=function()
		      {// 格式1：(msg,onConfirm,onCancel,onClose,onSuccess)  格式2：(msg,{...})
		        var self=this;
		          var _opt={};
		           var args=['','onConfirm','onCancel','onClose','onSuccess'];
		           if(!arguments.length){
		             throw 'confirm 必须至少带一个参数';
		             return;
		           }
		           if(typeof arguments[0] !== 'number' && typeof arguments[0]!=='string' && ((typeof arguments[0]=='object') && !arguments[0].css) && !arguments[0].css && !isDOM(arguments[0])){
		             throw 'confirm 第一个参数必须为string,number,jquery object,DOM之一';
		             return;
		           }
		           _opt.msg=arguments[0];
		           if(arguments.length>=2){
		             if(typeof arguments[1] == 'function'){
		               for(var i=1;i<arguments.length;i++){
		                 _opt[args[i]]=arguments[i];
		               }
		             }else{
		               _opt=$.extend(_opt,arguments[1]);
		             }
		           }
		           if(!_opt.onCancel){
		             _opt.onCancel=function(){};
		           }
		         $.extend(true,_opt,self.setFrame());
		          _this.openInit(_opt.frame,_opt,'confirm');
		      }
	  	} 
		  //接口API
		open={
		     alert:function()
		     {// 格式1：(msg,onConfirm,onClose,onSuccess)  格式2：(msg,{...})
		       var newObj=new openObject();
		      	newObj.alert.apply(newObj,arguments);
		     },
		    confirm:function(){// 格式1：(msg,onConfirm,onCancel,onClose,onSuccess)  格式2：(msg,{...})
		      var newObj=new openObject();
		      newObj.confirm.apply(newObj,arguments);
		    	}
		  };
		 return open;
	},
	  //传值初始化[$_obj:dialog的jquery对象   opt：初始化参数：{title,confirmTitle,cancelTitle,onSuccess,onConfirm,onClose}}
	openInit:function($_obj,_opt,alerttype){
			var _this=this;
			var defaultOption={
		    title:'\u7cfb\u7edf\u6d88\u606f',//系统消息
		    confirmTitle:'\u786e\u8ba4',//确认
		    cancelTitle:'\u53d6\u6d88',//取消
		    msg:'hello,world',
		    coverColor:'#777',
		    dragable:true,
		    width:'auto',//初始化内容区宽度
		    height:'auto',//初始化内容区高度
		    onSuccess:function(selector){
		    },
		    onConfirm:function(){},
		   	onCancel:function(){},
		    onClose:function(){}
		  }
	      //从defaultOption中深拷贝
	          var opt=$.extend(true,{},defaultOption);
	            //将调用者的选项整合到默认option
	          opt=$.extend(true,opt,_opt);
	          if(!$_obj.css){
	            $_obj=$($_obj);
	          }
	          var divTitle,divClose,divContent,divConfirm,divCancel;
	          divTitle=$_obj.find('.header>span');
	          divClose=$_obj.find('.header>a');
	          divContent=$_obj.find('.content>div:first');
	          divCancel=$_obj.find('.footer>a:eq(0)');
	          divConfirm=$_obj.find('.footer>a:eq(1)');
	          
	          //设置内容
	            divTitle.html(opt.title);
	            if(typeof opt.msg =='string' || typeof opt.msg=='number'){
	            divContent.html(opt.msg);
	            divConfirm.css('visibility','hidden');
	            }else{
	            opt.msg=$(opt.msg);
	            divContent.html('');
	            divContent.append(opt.msg);
	            }
	            divConfirm.html(opt.confirmTitle);
	            divCancel.html(opt.cancelTitle);
	          //遮罩层的颜色
	            $('#open-cover').css('background-color',opt.coverColor);
	          // 清除已有事件
	            divClose.unbind('click');
	            divConfirm.unbind('click');
	            divCancel.unbind('click');
	            divContent.unbind('load');
	          //加载新事件
	           
	            //confirm
	            if(opt.onConfirm && typeof opt.onConfirm=='function'){
	            divConfirm.css('visibility','visible');
	            divConfirm.click(function(){
	              if(opt.onConfirm(divContent)!==false){//return false 不关闭
	            _this.openHide($_obj,opt);
	          }
	            });
	            }else{
	            divConfirm.css('visibility','hidden');
	            }
	            //cancel
	            if(opt.onCancel && typeof opt.onCancel=='function'){
	            divCancel.css('visibility','visible');
	            divCancel.click(function(){
	                if(opt.onCancel(divContent)!==false){//return false 不关闭
	             _this.openHide($_obj,opt); 
	            }
	            });
	            }else{
	            divCancel.css('visibility','hidden');
	            }
	            //close
	            if(opt.onClose && typeof opt.onClose=='function'){
	              divClose.click(function(){ 
	                if(opt.onClose(divContent)!==false){//return false 不关闭
	                		
	            _this.openHide($_obj,opt);  
	            }
	              });
	            }else{
	              divClose.click(function(){
	              	$_obj.fadeIn();
	               // _this.openHide($_obj,opt);
	              });
	            }
	            
	             //success
	            if(opt.onSuccess && typeof opt.onSuccess=='function'){
	          
	            opt.onSuccess(divContent);
	            window.setTimeout(function(){
	            _this.openShow($_obj,opt);
	                  defaultOption.onSuccess(divContent);
	            },200);
	                

	            }else{

	            window.setTimeout(function(){
	            _this.openShow($_obj,opt);
	                  defaultOption.onSuccess(divContent);
	            },200);

	            }
	            if(alerttype=='alert'){
	            	 divConfirm.css('visibility','hidden');
	            }
	            //dragable
	           /* if(opt.dragable){
	              $_obj.dragable();
	            }else{
	              $_obj.dragable('disable');
	            }*/       
	},
		  //显示dialog
	openShow:function($_frame,_opt){
		    if(!$_frame.css){
		      $_frame=$($_frame);
		    }
		    var global_body=document.compatMode=='BackCompat'?document.body:document.documentElement;
		    var gWidth=global_body.clientWidth;
		    var gHeight=global_body.clientHeight;
		    $_frame.css('display','block');
		      //刷新divContent的maxHeight和width:auto
		    var cntWidth=(_opt && _opt.width)?_opt.width:'auto';
		      var maxHeight=(gHeight-150)<50?50:(gHeight-150);
		    if(_opt.height!='auto'){
		      $_frame.find('.content>div:first').css({
		        'height':_opt.height,
		        'overflow-y':'auto'
		      });
		    }else{
		      $_frame.find('.content>div:first').css({
		       'max-height':maxHeight+'px', 
		       'overflow-y':'auto'
		      });
		    }
		      $_frame.find('.content>div:first').css('width',cntWidth);
		    var fWidth=$_frame.get(0).clientWidth;
		    var fHeight=$_frame.get(0).clientHeight;
		  
		   //刷新divContent的maxHeight和width:auto
		    var cntWidth=(_opt && _opt.width)?_opt.width:'auto';
		      var maxHeight=(gHeight-150)<50?50:(gHeight-150-(gHeight-fHeight)/2);
		    if(_opt.height!='auto'){
		      $_frame.find('.content>div:first').css({
		        'height':_opt.height,
		        'overflow-y':'auto'
		      });
		    }else{
		      $_frame.find('.content>div:first').css({
		       'max-height':maxHeight+'px', 
		       'overflow-y':'auto'
		      });
		    }
		      $_frame.find('.content>div:first').css('width',cntWidth);
		    
		    //console.log(gWidth,gHeight,fWidth,fHeight);
		      //刷新遮层的尺寸
		      _opt.cover.width(gWidth);
		      _opt.cover.height(gHeight);
		      _opt.cover.css('display','block');
		      
		      
		      //居中
		      $_frame.css('position','fixed');
		      $_frame.css('left',(gWidth-fWidth)/2+'px');
		      $_frame.css('top',(gHeight-fHeight)/2+'px');
	},
	   //隐藏dialog
	openHide:function($_frame,_opt){
	      if($_frame.css){
		      	$_frame.css('display','none');
		        if(!_opt['notRemove']){
		        $_frame.remove();
		        }
	        }
	        else if($_frame.addEventListener){
	            $($_frame).css('display','none');
	        	if(!_opt['notRemove']){
	          	$($_frame).remove();
	        	}
	        }
	        //遮层
	        console.log(!_opt['notRemove']);
	        if(!_opt['notRemove']){
	         _opt.cover.animate({width:'0px',height:'0px'},500,function(){if(!_opt['notRemove']){
	          _opt.cover.remove();
	        }});
	        }
	         else{
	        	 _opt.cover.css('display','none');
	         
	        }
	}

}
