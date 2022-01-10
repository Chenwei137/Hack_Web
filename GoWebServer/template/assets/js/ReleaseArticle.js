$(document).ready(function() {
    $('#release').click(function(e){
        var TitleStr = $('#article_title').val();
        var TypeResult = $('#type-result').val();
        var now = new Date();        
        console.log(now)
        if (TitleStr == ""){
            TitleStr = "None"
        }
        var TextStr = $('#ArticleTextArea').val();
        //var TextStr = $('#ArticleTextArea').val().replace(/\n/g,"<br>").replace(/ /g,"&nbsp");
        if (TextStr == ""){
            $('#texterror').html("內容不可為空!");
            document.getElementById('texterror').setAttribute("style", "color:red;");
        }
        else{
            $('#texterror').html("");
            send_article(TitleStr,TextStr,TypeResult);
        }
    })
});
function send_article(title,text,category){
    $. ajax ({
        type: "POST",
        url: "/acceptarticle",
        dataType: 'json',
        contentType: "application/json",
        data: JSON.stringify({
            "title": title,
            "category": category,
            "text": text
        })
    })
}