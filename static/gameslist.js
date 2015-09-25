function show(thingtoshow,filter)
{
console.log("show("+thingtoshow+","+filter+")")
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
			}
		})
}
function show_games(console_id)
{
	console.log("show_games("+console_id+")")
	$('#secondary_div').html('Loading...');
	$.ajax({
		type: "GET", url: '/list/games?filter=console&console_id='+console_id,
		success:function(html){
			$('#content_div').html(html)
			$('#console_div_'+console_id).css("text-decoration", "underline")
		}
	})
}
function searchthings(form)
{
	var ss = form.search.value;
	console.log("search("+ss+")");
	$.ajax({type: "GET",url: '/search/',data:"query="+ss, success:function(html){
		$('#content_div').html(html);
	}})
}
function add_console(form)
{
	console.log("add_console("+form+")")
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
	console.log("add_game("+form+")")
	var data="console_id="+selvalue+"&game_name="+newgame
	$.ajax({type: "POST", url: '/console/newgame', data:data, success:function(html){
		$('#menu_status').html(html);
		form.add_game_text.value="";
	}})
}
function save_change(id,elem)
{
	console.log("id="+id+"elem="+elem);
	var data="action=have_not";
	if (document.getElementById(elem).checked)
	{
		data="action=have";
	}
	$.ajax({type: "POST", url: '/thing/'+id, data:data, success:function(html)
	{
		console.log(html);	
	}})
}
function toggle_owned(id)
{
	console.log("toggle_owned("+id+")")
	var data="action=toggle";
	$.ajax({type: "POST", url: '/thing/'+id, data:data, success:function(html)
	{
		$('#div_thing_'+id).css("background-color", html);
	}})
}
