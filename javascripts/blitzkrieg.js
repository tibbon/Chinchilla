"use strict";

var Average = 0

var Blitzkrieg = function() {

    var W = {}

    var blitzkrieg = {

        beginTest: function() {
            var test1Type = $("input[name='req-1-per-sec']").val()
            var test2Type = $("input[name='req-2-per-sec']").val()
            var test3Type = $("input[name='req-3-per-sec']").val()
            var workerCount = $("input[name='worker-count']").val()
            var distType;
            if($(".selected").hasClass("shortest-queue-button")) {
                console.log("here")
                distType = "sq"
            } else if ($(".selected").hasClass("roundrobin-button")) {
                distType = "rr"
            }
            console.log(test1Type)
            $.ajax({
                url: "/start_test",
                type: "post",
                dataType: "json",
                data: {
                    "type_1": test1Type,
                    "type_2": test2Type,
                    "type_3": test3Type,
                    "alg_type" : distType,
                    "worker_count" : workerCount
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
                url: "/stop_test",
                type: "post",
                dataType: "json"
            })
        },

        handleResponse: function(workers, data, stop) {
            if(!stop) {
                workers.pushResponse(data)
            }
        },
    }

    return blitzkrieg;

}

$(function() {

    var b = Blitzkrieg();

    var workers = NewWorkerStorage()

    var stop = false;

    var UI = NewUI(workers)
    UI.build()

    if (window["WebSocket"]) {
        var conn = new WebSocket("ws://localhost:9010/ws");
        conn.onclose = function(evt) {
            console.log("connection closed")
        }
        conn.onmessage = function(evt) {
            var data = JSON.parse(evt["data"])
            b.handleResponse(workers, data, stop)
        }

    } else {
        console.log($("<div><b>Your browser does not support WebSockets.</b></div>"))
    }

    $(".blitzkrieg-button").on("click", function() {
        stop = false
        b.beginTest();
        UI.start()
        refreshUI();
    })

    $(".stop-button").on("click", function() {
        stop = true
        b.endTest();
        workers.killAll();
        UI.stop();

    })

    $(".add-worker").on("click", function() {
        $.ajax({
            url: "/add",
            type: "post",
            dataType: "json"
        })
    })

    $(".option-button").on("click", function(e) {
        console.log('clicked')
        if(!$(e.currentTarget).hasClass("selected")) {
            $(".option-button").removeClass("selected")
            $(e.currentTarget).addClass("selected");
        }
    })


    var refreshUI = function() {
        
        if(stop) {
            return
        }

        UI.update(workers)
        workers.mapWorkers(function(wid, worker) {
            if(worker.kill) {
                workers.removeWorker(wid)
                $.ajax({
                    url: "/kill/" + wid,
                    type: "post",
                    dataType: "json"
                })
            }
        })
        window.setTimeout(function() {
            refreshUI();
        }, 1000)
    }

})