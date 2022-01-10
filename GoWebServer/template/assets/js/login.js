$(document).ready(function() {
    $('#login_button').click(function(){
        var account = $("input[name=account]").val()
        var password = $("input[name=password]").val()
        send_to_api(account,password)
        event.preventDefault()
    });
    $('#create_button').click(function(){
        window.location.href='/create';
    });
    $('#wordchange').click(function(){
        const password = document.getElementById("password");
        if (this.checked){
            password.type = "text";
        }
        else{
            password.type = "password";
        }
    });
});
function send_to_api(username,password){
    $. ajax ({
        type: "POST",
        url: "/login",
        dataType: 'json',
        contentType: "application/json",
        data: JSON.stringify({
            "username":username,
            "password":password
        })
    })
}