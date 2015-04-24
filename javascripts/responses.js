function NewWorkerStorage() {

    var storage = {}
    var maxQAvg = 0

    function worker(num_responses, responses_array) {
        this.respCount    = num_responses;
        this.resps        = responses_array;
        this.active       = true;
        this.push  = function(r) {
            this.resps.push(r);
            this.respCount += 1;
            this.lastPush = Math.floor(Date.now() / 1000);
            if(typeof(this.avgWeighted) == "undefined") {
                this.avgWeighted = 0
            } else {
                this.avg(r)
            }
        },
        this.avg   = function(new_response) {
            this.avgWeighted = this.avgWeighted * 0.95 + new_response["q_val"] * 0.05
            if(this.avgWeighted > maxQAvg) {
                maxQAvg = this.avgWeighted
            }
        }
    }

    var workerStorage = {
        pushResponse: function(response) {
            var wid = response["id"];
            var w   = storage[wid];
            if(typeof(w) != "undefined") {
                w.push(response)
            } else {
                storage[wid] = new worker(0, []);
                storage[wid].push(response)
            }
        },

        mapWorkers: function(mapFn) {
            for (var wid in storage) {
                if(storage.hasOwnProperty(wid)) {
                    storage[wid].active = Math.floor(Date.now() / 1000) - storage[wid].lastPush < 10
                    if(!storage[wid].active) 
                        storage[wid].avgWeighted = 0
                    mapFn(wid, storage[wid])
                }
            }
        },

        maxQueueAverage: function() {
            return maxQAvg
        },

        killWorker: function(wid) {
            storage[wid].kill = true
        },

        removeWorker: function(wid) {
            delete storage[wid]
        },

        killAll: function() {
            storage = {}
            maxQAvg = 0;
        }
    }

    return workerStorage;

}