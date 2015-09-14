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
						$('#toggle_consoles').css("border", "none")
						$('#toggle_games').css("border", "none")
						$('#toggle_'+filter).css("border", "3px solid")
					}
				});
			} else {
				$('#toggle_container').html('');
			}
			$('#content_div').html(html);
			$('#header-collection').css("border", "none")
			$('#header-consoles').css("border", "none")
			$('#header-games').css("border", "none")
			$('#header-'+thingtoshow).css("border", "3px solid")
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
function toggle_owned(id)
{
	var data="action=toggle";
	$.ajax({type: "POST", url: '/thing/'+id, data:data, success:function(html)
	{
		$('#div_thing_'+id).css("background-color", html);
	}})
}
