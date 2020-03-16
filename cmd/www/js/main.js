
     $("[data-hide]").on("click", function(){
           $("." + $(this).attr("data-hide")).hide();
      
     });		
	
	 $(document).on("click","[data-url]",function(e){
		 
		 e.preventDefault();
		 
		  typ = $(this).attr("data-method");
		  
		  console.log(this.href);
		  
		  if(this.href.includes("/deviceremove")){
			  
			  if(!confirm("Are you sure you want to remove this device?")){
				  
				  return;
			  }
			  
		  }
		  
		  
		  $.ajax({
			  url:this.href, 
			  cache: false,
			  type: (typ=="post" ? "post" : "get"),
			  success: function(resp){
   
  			   $(".alert").fadeIn(function(){
  			   	$("#message").html(resp.message);
  			   });
  			   getStatus();

  		  }});
		 
	 });




function getStatus(){
		
		 $.get("/status",function(resp){
	

		 
		 
		 ht="<tr><th>Devices</th><th></th></tr>";
		 
		

		 d = Object.keys(resp.payload.devices).length;
		 
		 nodes = Object.keys(resp.payload.nodes).length;
		 
		 for (var key in resp.payload.devices){ 
			 
			 ht+="<tr>";
			 ht+="<td>"+key+" ("+resp.payload.devices[key]["label"]+")</td>"; 
		     ht+="<td><a data-url='true' data-method='post' href='/deviceremove/"+key+"' class='btn btn-danger' role='button'  ><span class='oi oi-circle-x' title='icon name' aria-hidden='true'></span> Remove</a></td>";
				 
			 
			 
			 ht+="</tr>";
						 
			 
			 
			 
		 }
		 
		  $("#dev").html("<table class='table'>"+ht+"</table>");
		 
		  ht="<tr><th>Nodes</th><th>Status</th></tr>";
		 
		 for (var key in resp.payload.nodes){ 

			 ht+="<tr>";
			 ht+="<td>"+key+" ("+resp.payload.nodes[key]["addr"]+")</td>";
			 ht+="<td>"+(resp.payload.nodes[key]["status"]=="master" ? '<span class="btn btn-success" >Master</span>' : '<span class="btn btn-info" >Slave</span>' )+"</td>";
			 ht+="</tr>";

		 }
		 

		 
		 $("#device_ct").html(d);
		 
		
		 
		 $("#nodes").html("<table class='table'>"+ht+"</table>");
		 
		 
		 
		 
	 },"json");	
		
	}

