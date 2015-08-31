function showConsoles()
{
	$.ajax({
		type: "GET",
		url: '/list/consoles', 
		success:function(html){
			$('#all_consoles_list').html(html);
			}
		})
}
