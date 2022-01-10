$(document).ready(function() {
    $('#create_button').click(function(){
        var account = document.getElementById("username").value
        var password = document.getElementById("passwd").value
        var check_password = document.getElementById("ckpasswd").value
        if(account !=""){
            let result = account.includes("'")
            if(result){
                alert("不可有特殊字元")
            }
        }
        else if(password !=""){
            let result = password.includes("'")
            if(result){
                alert("不可有特殊字元")
            }
        }
        else if(check_password !=""){
            let result = check_password.includes("'")
            if(result){
                alert("不可有特殊字元")
            }
        }
        if(password == check_password){
           send_to_create(account,password)
           window.location.href='/login';
        }else{
            alert("兩次密碼不相同")
        }
    });
});
function send_to_create(username,password){
    $. ajax ({
        type: "POST",
        url: "/create",
        dataType: 'json',
        contentType: "application/json",
        data: JSON.stringify({
            "username":username,
            "password":password
        })
    })
}