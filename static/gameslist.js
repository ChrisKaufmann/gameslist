var searchtimeout;
$(document).on('input propertychange paste', '#search_text',  function() {
	clearTimeout(searchtimeout);
	searchtimeout = setTimeout(function() {
		var sc = document.getElementById('search_text');
		var sv = search_text.value;
		if(sv.length >=2){
			searchthings(sv);
		}
	}, 1000);
});
$(document).keyup(function (e) {
	if(e.keyCode == 13) {
		var sc = document.getElementById('search_text');
		var sv = search_text.value;
		searchthings(sv);
	}
});
$("#search_text").keyup(function(e){
	if(e.keyCode == 13) {
		console.log("enter");
		return false;
	}
});
function show(thingtoshow,filter)
{
	$.ajax({
		type: "GET",
		url: '/list/'+thingtoshow+'?filter='+filter, 
		success:function(html){
			$('#content_div').html(html);
			$('#header-collection').css("text-decoration", "none")
			$('#header-consoles').css("text-decoration", "none")
			$('#header-games').css("text-decoration", "none")
			$('#header-'+thingtoshow).css("text-decoration", "underline")
    		$('#secondary_div').html("")
			$('#floating_sidebar').css("visibility", "hidden");
			}
		})
}
function show_games(console_id)
{
	$('#secondary_div').html('Loading...');
	$.ajax({
		type: "GET", url: '/list/games?filter=console&console_id='+console_id,
		success:function(html){
			$('#content_div').html(html);
			$('#console_div_'+console_id).css("text-decoration", "underline");
			$('#floating_sidebar').css("visibility", "visible");
		}
	})
}
function searchthings(ss)
{
	$.ajax({type: "GET",url: '/search/',data:"query="+ss, success:function(html){
		$('#content_div').html(html);
		$('#floating_sidebar').css("visibility", "hidden");
	}})
}
function add_console(form)
{
	$('menu_status').innerHTML='Adding...';
	var newcon =form.add_console_text.value;
	$.ajax({type: "GET",url: '/console/new/'+newcon,success:function(html){
		$('#menu_status').html(html);
		form.add_console_text.value="";
	}})
}
function add_game(form)
{
	var newgame = form.add_game_text.value;
	var index=form.add_game_select.selectedIndex;
	var selvalue=form.add_game_select.options[index].value;
	var data="console_id="+selvalue+"&game_name="+newgame
	$.ajax({type: "POST", url: '/console/newgame', data:data, success:function(html){
		$('#menu_status').html(html);
		form.add_game_text.value="";
	}})
}
function save_change(id,elem)
{
	var data="action=have_not";
	if (document.getElementById(elem).checked)
	{
		data="action=have";
	}
	$.ajax({type: "POST", url: '/thing/'+id, data:data, success:function(html)
	{
	}})
}
function toggle_owned(id)
{
	var data="action=toggle";
	$.ajax({type: "POST", url: '/thing/'+id, data:data, success:function(html)
	{
		$('#div_thing_'+id).css("background-color", html);
	}})
}
function setrating(id, rating)
{
	console.log("setrating:"+id+", rating:"+rating)
	var data="action=setrating&rating="+rating
	$.ajax({type: "POST", url: '/thing/'+id, data:data, success:function(html)
	{
		$('#star_container_'+id).html(html)
	}})
}
function enable_review(id)
{
	console.log("enable_review: "+id)
	var data="action=get_review_html"
	$.ajax({type: "GET", url: '/thing/'+id, data:data, success:function(html){
		$('#div_review_'+id).html(html)
	}})
}
