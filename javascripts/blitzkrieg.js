"use strict";

var Average = 0

var Blitzkrieg = function() {

    var W = {}

    var blitzkrieg = {

        beginTest: function() {
            var test1Type = $("input[name='req-1-per-sec']").val()
            var test2Type = $("input[name='req-2-per-sec']").val()
            var test3Type = $("input[name='req-3-per-sec']").val()
            console.log(test1Type)
            $.ajax({
                url: "/blitz",
                type: "post",
                dataType: "json",
                data: {
                    "type_1": test1Type,
                    "type_2": test2Type,
                    "type_3": test3Type,
                    "alg_type" : "rr",
                    "worker_count" : 10
                },
                success: function(data) {
                    blitzkrieg.updateDisplay();
                },
                error: function(header, status, error) {
                    console.log(header);
                }
            })
        },

        endTest: function() {
            $.ajax({
                url: "/stop",
                type: "post",
                dataType: "json",
                success: function(data) {
                    console.log(data)
                }
            })
        },

        handleResponse: function(workers, data) {
            workers.pushResponse(data)
        },
    }

    return blitzkrieg;

}

$(function() {

    var b = Blitzkrieg();

    var workers = NewWorkerStorage()

    var UI = NewUI(workers)
    UI.build()

    if (window["WebSocket"]) {
        var conn = new WebSocket("ws://localhost:9010/ws");
        conn.onclose = function(evt) {
            console.log("connection closed")
        }
        conn.onmessage = function(evt) {
            var data = JSON.parse(evt["data"])
            b.handleResponse(workers, data)
        }

    } else {
        console.log($("<div><b>Your browser does not support WebSockets.</b></div>"))
    }

    $(".blitzkrieg-button").on("click", function() {
        b.beginTest();
        UI.start()
        refreshUI();
    })

    $(".stop-button").on("click", function() {
        b.endTest();
        UI.stop();
    })



    var refreshUI = function() {
        UI.update(workers)
        workers.mapWorkers(function(wid, worker) {
            if(worker.kill) {
                workers.removeWorker(wid)
            }
        })
        window.setTimeout(function() {
            refreshUI();
        }, 1000)
    }

})