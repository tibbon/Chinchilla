"use strict";

var Blitzkrieg = function() {

    var blitzkrieg = {

        sendRequest: function() {
            
            $.ajax({
                url: "http://localhost:9000/api/work/1",
                type: "GET",
                dataType: "json",
                success: function(data) {
                    console.log(data);
                },
                error: function(data) {
                    console.log(data);
                }
            })

        }

    }

    return blitzkrieg;

}

$(function() {

    var b = Blitzkrieg();
    $(".blitzkrieg-button").on("click", function() {
        b.sendRequest();    
    })

})