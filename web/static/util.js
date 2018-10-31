function check_ip(host_ip)
{
    if (host_ip == "") {
        return false;
    }
    var part = host_ip.split(".");
    if (part.length != 4) {
        return false;
    }
    var i;
    for (i = 0; i < part.length; i++) {
        if (isNaN(part[i]) || parseInt(part[i]) > 255 || parseInt(part[i]) < 0) {
            return false;
        }
    }
    return true;
}

function draw_line(divid, mtitle, ytitle, xtitle, data)
{
    $('#'+divid).highcharts({
        title: {
            text: mtitle,
            x: -20 //center
        },
        xAxis: {
            categories: xtitle,
            labels:{
                step:3,
            }
        },
        yAxis: {
            title: {
                text: ytitle
            },
            plotLines: [{
                value: 0,
                width: 1,
                color: '#808080'
            }]
        },
        legend: {
            layout: 'vertical',
            align: 'right',
            verticalAlign: 'middle',
            borderWidth: 0
        },
        series: data
    });
}

/**
 *   è½¬ä¸ºæ¥ææ ¼å¼
 *
 *   @method parseDate
 *   @param {string} s æ¶é´å­ç¬¦ä¸²
 *   @return {date} æ¥æ
 */
function parseDate(s) {
    if(typeof s == 'object') {
        return s;
    }

    //å¦ææ¯"HH:mm:ss"æ ¼å¼ æè'HH:mm'
    if (new RegExp('^([0-9]{2,2}\:)').test(s)) {
        var tempDate = new Date();
        var dateString = 'yyyy-MM-dd';
        dateString = dateString.replace('yyyy', tempDate.getFullYear().toString())
            .replace('MM', (tempDate.getMonth() < 9 ? '0' : '') + (tempDate.getMonth() + 1).toString())
            .replace('dd', (tempDate.getDate() < 10 ? '0' : '') + tempDate.getDate().toString());

        s = dateString + " " + s;
    }

    var ar = (s + ",0,0,0").match(/\d+/g);
    return ar[5] ? (new Date(ar[0], ar[1] - 1, ar[2], ar[3], ar[4], ar[5])) : (new Date(s));
}
/**
 * * æ ¼å¼åæ¶é´
 * */
function formatDate(date, format) {
    date = date || new Date();
    format = format || 'yyyy-MM-dd HH:mm:ss';
    var result = format.replace('yyyy', date.getFullYear().toString())
        .replace('MM', (date.getMonth()< 9?'0':'') + (date.getMonth() + 1).toString())
        .replace('dd', (date.getDate()< 10?'0':'')+date.getDate().toString())
        .replace('HH', (date.getHours() < 10 ? '0' : '') + date.getHours().toString())
        .replace('mm', (date.getMinutes() < 10 ? '0' : '') + date.getMinutes().toString())
        .replace('ss', (date.getSeconds() < 10 ? '0' : '') + date.getSeconds().toString());

    return result;
}

// var obj = "2012-11-21 16:39:44";
// FormatDateTime(obj);
function hichatFormatDateTime(obj)
{
    var val = parseDate(obj);
    return formatDate(val, "MM-dd HH:mm");
}

function create_paging(paging_id,total_rows,per_page)
{
    var show_pages = 10;
    var total_pages = Math.floor((total_rows-1)/per_page)+1;
    if (total_pages <=1){$('.'+paging_id).html(""); return;}
    var show_max_page_num = Math.min(show_pages,total_pages);
    var ret = "<ul class='pagination' style='margin-top:0px; margin-bottom:0px'>"
        + "<li style='display:table-cell'><a href='#' tyle='cursor:pointer'>首页</a></li>"
        + "<li style='display:none'><a href='#' tyle='cursor:pointer'>&laquo;</a></li>"
        + "<li style='display:table-cell' class='active'><a href='#' tyle='cursor:pointer'>1</a></li>"
    for (var i=2; i<=total_pages; i++){
        if (i <= show_max_page_num){
            ret = ret + "<li style='display:table-cell'><a href='#' tyle='cursor:pointer'>"+i+"</a></li>";
        }else{
            ret = ret + "<li style='display:none'><a href='#' tyle='cursor:pointer'>"+i+"</a></li>";
        }
    }
    if (total_pages > show_pages){
        ret = ret + "<li style='display:table-cell'><a href='#' tyle='cursor:pointer'>&raquo;</a></li>";
    }else{
        ret = ret + "<li style='display:none'><a href='#' tyle='cursor:pointer'>&raquo;</a></li>";
    }
    ret = ret + "<li style='display:table-cell'><a href='#' tyle='cursor:pointer'>尾页</a></li>";
        + "</ul>";
    $('#'+paging_id).html(ret);
}
function get_click_paging_index(obj){
    var show_pages = 10;
    var click_paging_index = $(obj).index();
    var max_index = $(obj).parent().children("li").length - 1;
    var max_page_index = max_index - 2;
    var left_index = 1;
    var right_index = max_index - 1;
    var active_index = $(obj).parent().find("li.active").index();

    if (click_paging_index == 0){
        click_paging_index = 2;
    }
    else if (click_paging_index == max_index){
        click_paging_index = max_page_index;
    }
    else if (click_paging_index == left_index){
        click_paging_index = (active_index - active_index % show_pages) - show_pages + 2;
    }
    else if (click_paging_index == right_index){
        click_paging_index = (active_index - active_index % show_pages) + show_pages + 2;
    }
    return click_paging_index;
}
function click_paging(obj)
{
    var show_pages = 10;
    var paging_id = $(obj).parent().parent().attr('id');
    var max_index = $(obj).parent().children("li").length - 1;
    var max_page_index = max_index - 2;
    var left_index = 1;
    var right_index = max_index - 1;
    var active_index = $(obj).parent().find("li.active").index();

    var click_paging_index = get_click_paging_index(obj);

    var show_page_index_begin = (click_paging_index - 2) - (click_paging_index - 2) % show_pages + 2;
    var show_page_index_end = Math.min(((click_paging_index - 2) - (click_paging_index - 2) % show_pages + show_pages + 1),max_page_index);

    var show_left = 1;
    var show_right = 1;
    if (show_page_index_begin <= (show_pages + 1)){show_left = 0;}
    if (show_page_index_end == max_page_index){show_right = 0;}

    $('#'+paging_id+' ul li').eq(show_page_index_begin).prevUntil($('#'+paging_id+' ul li').eq(1)).hide();
    $('#'+paging_id+' ul li').eq(show_page_index_begin).nextUntil($('#'+paging_id+' ul li').eq(show_page_index_end+1)).andSelf().css("display","table-cell");
    $('#'+paging_id+' ul li').eq(show_page_index_end).nextUntil($('#'+paging_id+' ul li').eq(max_page_index+1)).hide();

    if (show_left == 1){$('#'+paging_id+' ul li').eq(1).css("display","table-cell");}else{$('#'+paging_id+' ul li').eq(1).hide();}
    if (show_right == 1){$('#'+paging_id+' ul li').eq(max_index -1).css("display","table-cell");}else{$('#'+paging_id+' ul li').eq(max_index -1).hide();}

    $('#'+paging_id+' ul li.active').removeClass('active');
    $('#'+paging_id+' ul li').eq(click_paging_index).addClass('active');
}

function on_click_table_tab(table_id,tab_id,per_page){
    $('table#'+table_id+' tbody tr').hide();
    $('table#'+table_id+' tbody tr.'+tab_id+':lt(10)').each(function(i, ele){var obj=$(ele); obj.find('td:first-child').text(i+1); obj.show();});
    $('#'+table_id+'_tabs button.selected').removeClass('selected');
    $('#'+tab_id).addClass('selected');
    $('#'+tab_id+"_paging").show();
    click_paging($('#'+tab_id+"_paging ul li").eq(0));
    $('#'+tab_id+"_paging").prevAll().hide();
    $('#'+tab_id+"_paging").nextAll().hide();
}

function on_click_table_tab_id_paging(table_id,tab_id,page,per_page){
    var per_page = 10;
    var page_begin = (page - 1) * per_page;
    var page_end = (page) * per_page;
    $('table#'+table_id+' tbody tr').hide();
    $('table#'+table_id+' tbody tr.'+tab_id+':eq('+page_begin+')').nextUntil($('table#'+table_id+' tbody tr.'+tab_id+':eq('+page_end+')')).andSelf().each(
        function(i, ele){
            var obj=$(ele); obj.find('td:first-child').text((page-1)*per_page+i+1); obj.show();
    });
}
function on_click_table_paging(table_id,page,per_page){
    var page_begin = (page - 1) * per_page;
    var page_end = (page) * per_page;
    $('table#'+table_id+' tbody tr').hide();
    $('table#'+table_id+' tbody tr:eq('+page_begin+')').nextUntil($('table#'+table_id+' tbody tr:eq('+page_end+')')).andSelf().each(
        function(i, ele){
            var obj=$(ele); obj.find('td:first-child').text((page-1)*per_page+i+1); obj.show();
    });
}
