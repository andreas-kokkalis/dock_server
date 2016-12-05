$('#logout-button').click(function(){
    $.ajax({
        url: "/v0/admin/logout",
        contentType: "application/json",
        success: function(response) {
            window.location.replace("/ui/login.html");
        },
        error: function(response) {
            handleError(response)
        }
    });
});
