function handleError(response) {
    // undefined when ???
    if (typeof response.responseJSON !== "undefined") {
        console.log(response.responseJSON.errors)
        $.each(response.responseJSON.errors, function(k, v) {
            $("#error-panel").append(
                '<div class="alert alert-danger alert-dismissible" role="alert">' +
                    '<button type="button" class="close" data-dismiss="alert" aria-label="Close"><span aria-hidden="true">&times;</span></button>' +
                '<span class="glyphicon glyphicon-exclamation-sign" aria-hidden="true"></span>' +
                '<span class="sr-only">Error:</span> ' + v + '</div>');
        })
    } else {
        console.log(response)
        if (response.status === 401) {
            window.location.replace("/ui/login.html");
        }
    }
}
