function show(thingtoshow,filter)
{
	$.ajax({
		type: "GET",
		url: '/list/'+thingtoshow+'?filter='+filter, 
		success:function(html){
			if(thingtoshow != 'collection') {
				$.ajax({
					type:"GET",
					url:'/toggle/'+thingtoshow,
					success:function(togglehtml){
						$('#toggle_container').html(togglehtml);
						$('#toggle_consoles').css("text-decoration", "none")
						$('#toggle_games').css("text-decoration", "none")
						$('#toggle_'+filter).css("text-decoration", "underline")
					}
				});
			} else {
				$('#toggle_container').html('');
			}
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
	$.ajax({
		type: "GET", url: '/list/games?filter=console&console_id='+console_id,
		success:function(html){
			$('#div_thing'+console_id).css("text-decoration","3px solid");
			$('#secondary_div').html(html)
		}
	})
}
function add_console(form)
{
        $('menu_status').innerHTML='Adding...';
        var newcon =form.add_console_text.value;
        $.ajax({type: "GET",url: '/console/new/'+newcon,success:function(html)
        {
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
function toggle_owned(id)
{
	var data="action=toggle";
	$.ajax({type: "POST", url: '/thing/'+id, data:data, success:function(html)
	{
		$('#div_thing_'+id).css("background-color", html);
	}})
}
