"use strict";

var Blitzkrieg = function() {

    var blitzkrieg = {

        sendRequest: function() {
            
            $.ajax({
                url: "localhost:9000/api/work/1",
                type: "GET",
                success: function(data) {
                    console.log(data);
                },
                failure: function(data) {
                    console.log(data);
                }
            })

        }

    }

    return blitzkrieg;

}

$(function() {
    var b = Blitzkrieg();
    b.sendRequest();
})