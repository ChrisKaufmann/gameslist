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
function show_owned(console_id)
{
	$('#secondary_div').html('Loading...');
	$.ajax({
		type: "GET", url: '/list/games?filter=consoleowned&console_id='+console_id,
		success:function(html){
			$('#content_div').html(html);
			$('#console_div_'+console_id).css("text-decoration", "underline");
			$('#floating_sidebar').css("visibility", "visible");
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
    window.location.replace("/search/?query="+ss)
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
function save_game_has(id,elem)
{
	var data="action=have_not&id="+id;
	if (document.getElementById(elem).checked)
	{
		data="action=have&id="+id;
	}
	$.ajax({type: "POST", url: '/set/game/', data:data, success:function(html)
	{
	}})
}
function save_console_has(id,elem)
{
    console.log("id:"+id+" elem:"+elem)
	var data="action=have_not&name="+id;
	if (document.getElementById(elem).checked)
	{
        $.ajax({type: "GET", url: "/set/console/?action=have&value=true&name="+id, success:function(html){}})
	} else {
        $.ajax({type: "GET", url: "/set/console/?action=have_not&value=false&name="+id, success:function(html){}})

	}
}
function save_game_manual(id,elem)
{
	var data="action=hasnot_manual&id="+id;
	if (document.getElementById(elem).checked)
	{
		data="action=has_manual&id="+id;
	}
	$.ajax({type: "POST", url: '/set/game/'+id, data:data, success:function(html)
	{
	}})
}
function save_console_manual(id,elem)
{
	var data="action=hasnot_manual&name="+id;
	if (document.getElementById(elem).checked)
	{
		data="action=has_manual&name="+id;
	}
	$.ajax({type: "POST", url: '/set/console/'+id, data:data, success:function(html)
	{
	}})
}
function save_game_box(id,elem)
{
	var data="action=hasnot_box&id="+id;
	if (document.getElementById(elem).checked)
	{
		data="action=has_box&id="+id;
	}
	$.ajax({type: "POST", url: '/set/game/'+id, data:data, success:function(html)
	{
	}})
}
function save_console_box(id,elem)
{
	var data="action=hasnot_box&name="+id;
	if (document.getElementById(elem).checked)
	{
		data="action=has_box&name="+id;
	}
	$.ajax({type: "POST", url: '/set/console/'+id, data:data, success:function(html)
	{
	}})
}
function set_game_rating(id, rating)
{
	console.log("set_game_rating:"+id+", rating:"+rating)
	var data="id="+id+"&action=setrating&rating="+rating
	$.ajax({type: "POST", url: '/set/game/', data:data, success:function(html)
	{
		$('#star_container_'+id).html(html)
	}})
}
function set_console_rating(id, name, rating)
{
	var data="name="+name+"&action=setrating&rating="+rating
	$.ajax({type: "POST", url: '/set/console/', data:data, success:function(html)
	{
		$('#star_container_'+id).html(html)
	}})
}
function save_game_review(id)
{
	var rev = document.getElementById("review_text_"+id).value;
	var data="action=set_review&id="+id+"&review="+rev;
	$.ajax({type: "POST", url: '/set/game/', data:data, success:function(html){
		document.getElementById("review_text_"+id).value=html;
	}})
}
function save_console_review(id,name)
{
	var rev = document.getElementById("review_text_"+id).value;
	var data="action=set_review&name="+name+"&review="+rev;
	$.ajax({type: "POST", url: '/set/console/', data:data, success:function(html){
		document.getElementById("review_text_"+id).value=html;
	}})
}
