
var NewUI = function(workerStorage) {

  var sourceColors = ["#b5382f", "#5cb56b", "#b1b55c", "#4544b5", "#7031cd", "#00c9cd", "#bc8d60", "#90e600", "#e66900", "#4176e6", "#00ff94", "#ff0003", "#d9aaff", "#b8d0ff", "#ffc100", "#c6ffc3", "#077400", "#000674", "#700074", "#74001c"]
  var colors = sourceColors.slice(0)
  var graph = {}
  var stats = {}
  var workers = []

  

  var wIndex = function(wid) {
    for (var w in workers) {
      if (workers[w].id == wid) {
        return w
      }
    }
  }

  /* width of data points (seconds) */
  graph.n = 60
  /* line data info */
  graph.random = d3.random.normal(0, 0);
  graph.workerStorage = workerStorage;

    /* graph dimensions */
  graph.margin = {top: 0, right: 0, bottom: 20, left: 0};
  graph.width = $(window).width() - graph.margin.left - graph.margin.right;
  graph.height = 500 - graph.margin.top - graph.margin.bottom;

  /* starting y value */
  graph.currentMax = 5;
    /* definition of x and y scales */
  graph.x = d3.scale.linear().domain([1, graph.n - 2]).range([0, graph.width]);

  graph.xAxis = d3.scale.linear().domain([2 - graph.n, 1]).range([0, graph.width]);
  graph.y = d3.scale.linear().domain([0, graph.currentMax]).range([graph.height, graph.margin.bottom * 1.5]);

   graph.line = d3.svg.line()
      .interpolate("basis")
      .x(function(d, i) { return graph.x(i); })
      .y(function(d, i) { return graph.y(d); })

  stats.n = 10;
  stats.workerHeight = 50;
  stats.width = $(window).width();
  stats.height = (stats.n / 2 * stats.workerHeight)

  function worker(wid) {

    var rgb = [];
    for(var i = 0; i < 3; i++)
        rgb.push(Math.floor(Math.random() * 255));
    var colorI = Math.floor(Math.random() * colors.length);
    // var color = colors.splice(colorI, 1)[0]
    var color = 'rgb('+ rgb.join(',') +')';

    this.id = wid
    this.q_val = 0
    this.avg = 0
    this.active = true
    this.data = d3.range(graph.n).map(graph.random)
    this.line = d3.svg.line()
      .interpolate("basis")
      .x(function(d, i) { return graph.x(i); })
      .y(function(d, i) { return graph.y(d); })

    this.path = graph.svg.append("g")
                .attr("clip-path", "url(#clip)")
                .append("path")
                .datum(this.data)
                .attr("class", "line")
                .attr("id", "line-" + this.id)
                .attr("d", this.line),

    this.ui    = {

      "color" : color,

      "selectors" : {
        outer: ".worker-stats-"+wid,
        color: ".worker-stats-"+wid+" .color-key",
        label: ".worker-stats-"+wid+" .worker-label"
      },
      "elements" : {
        outer: $("<div class='worker-stat worker-stats-"+wid+"'/>"),
        color: $("<div class=color-key style='background:"+color+";'/>"),
        label: $("<p class='worker-label'>Worker "+wid+"</p>")
      }
    }
  }

  newWorker = function(wid) {

    var w = new worker(wid)
    workers.push(w)

  }


  var build = function() {

    graph.svg = d3.select(".line-graph").append("svg")
      .attr("width", graph.width + graph.margin.left + graph.margin.right)
      .attr("height", graph.height + graph.margin.top + graph.margin.bottom)
      .append("g")
        .attr("transform", "translate(" + graph.margin.left + "," + graph.margin.top + ")");

    graph.svg.append("defs").append("clipPath")
        .attr("id", "clip")
      .append("rect")
        .attr("class", "box")
        .attr("width", graph.width)
        .attr("height", graph.height);

    graph.svg.append("g")
        .attr("class", "x axis")
        .attr("transform", "translate(0," + graph.y(0) + ")")
        .call(d3.svg.axis().scale(graph.xAxis).orient("bottom").outerTickSize(0).tickFormat(formatTick))

    graph.svg.append("g")
        .attr("class", "y axis")
        .call(d3.svg.axis().scale(graph.y).orient("right"))

    stats.canvas = d3.select(".worker-stats")


  }

  var update = function() {

    graph.workerStorage.mapWorkers(function(wid, worker) {
      if(typeof(workers[wIndex(wid)]) != "undefined") {
        workers[wIndex(wid)]["avg"] = worker.avgWeighted
        workers[wIndex(wid)]["data"].push(worker.avgWeighted)
        workers[wIndex(wid)]["active"] = !worker.kill
        workers[wIndex(wid)]["q_val"] = worker.avgWeighted
        if(worker.kill) {
          workers[wIndex(wid)].path.remove()
          workers.splice(wIndex(wid), 1)
        }
        var max = graph.workerStorage.maxQueueAverage();
        if(max > graph.currentMax) {
          graph.currentMax = max
          graph.y = d3.scale.linear().domain([0, graph.currentMax * 1.5]).range([graph.height, graph.margin.bottom * 1.5])
          graph.svg.selectAll("g.y.axis").call(d3.svg.axis().scale(graph.y).orient("right"))
        }
      } else {
        newWorker(wid)
        tick(wid)
      }
    })

    $("input[name='worker-count']").val(workers.length)
    updateStats();
    updateStdDev();

  }

  var tick = function(wid) {

      if(typeof(workers[wIndex(wid)]) == "undefined" || graph.stop) {
        d3.select("#line-"+wid)
        .transition()
        .duration(500)
        .ease("linear")
        .style("opacity", 0)
        .each("end", function() {
          d3.select(this).remove()
        })
        return;
      }


      workers[wIndex(wid)].path.attr("d", workers[wIndex(wid)].line)
        .attr("transform", null)
        .transition()
        .duration(1000)
        .ease("linear")
        .attr("transform", "translate(" + graph.x(0) + ",0)")
        .style("stroke", function(d) {
          return workers[wIndex(wid)].ui.color
        })
        .each("end", function() {
          tick(wid);
        });

    // pop the old data point off the front
      workers[wIndex(wid)].data.shift();

  }

  var updateStats = function() {

    var cells = stats.canvas.selectAll(".worker-stat")
      .data(workers, function(d) {
        return d["id"]
      });

    var cellEnter = cells
      .enter()
      .append("div")
      .attr("class", function(d, i) {
        return "worker-stat worker-data-"+d["id"]
      })

    cells
      .exit()
      .remove()

    cellEnter.on("click", function(d) {
      d3.select("#line-"+d.id)
        .transition()
        .duration(500)
        .ease("linear")
        .style("opacity", 0)
        .each("end", function() {
          d3.select(this).remove()
          graph.workerStorage.killWorker(d["id"])
        })
      colors.push(d.ui.color)
    });

    cellEnter.append("div")
      .attr("class", "color-key")
      .style("background", function(d) {
        return d["ui"]["color"]
      })

    var dataEl = cellEnter.append("div")
      .attr("class", "data")

    dataEl.append("p")
      .attr("class", "worker-id")
      .text(function(d) {
        return "Worker " + d["id"]
      });

    dataEl.append("p")
      .attr("class", "q-val")
      .text(function(d) {
        return "Queue: " + d["q_val"]
      });


    cells.selectAll(".q-val").text(function(d, i) {
        return "Queue " + d["q_val"].toFixed(2)
    })

  }

  var start = function() {
    axisTick(2 - graph.n, 1);
    graph.stop = false;
  }

  var stop = function() {
    graph.axisReset = true;
    graph.stop = true;
    graph.currentMax = 5;
    graph.y = d3.scale.linear().domain([0, graph.currentMax]).range([graph.height, graph.margin.bottom * 1.5]);
    graph.svg.selectAll("g.y.axis").call(d3.svg.axis().scale(graph.y).orient("right"))
    colors = sourceColors.slice(0)

    workers = []

    updateStats()

    axisTick(2 - graph.n, 1);
  }

  var formatTick = function(d, i) {
    var neg = ""
    if(d < 0) {
      neg = "-"
    }
    if(d % 60 == 0) {
      return neg + d / 60 + ":00"
    } else {
      return neg + ":" + d % 60
    }
  }

  var axisTick = function(prevMin, prevMax) {
    if(graph.axisReset) {
      graph.xAxis = d3.scale.linear().domain([2 - graph.n, 1]).range([0, graph.width]);
      graph.axisReset = false;
      graph.svg.selectAll("g.x.axis")
        .transition()
        .duration(1000)
        .ease("linear")
        .call(d3.svg.axis()
        .scale(graph.xAxis).orient("bottom").outerTickSize(0).tickFormat(formatTick))
    } else {
      graph.xAxis = d3.scale.linear().domain([prevMin + 1, prevMax + 1]).range([0, graph.width])
      graph.svg.selectAll("g.x.axis")
        .transition()
        .duration(1000)
        .ease("linear")
        .call(d3.svg.axis()
        .scale(graph.xAxis).orient("bottom").outerTickSize(0).tickFormat(formatTick))
        .each("end", function() {
          axisTick(prevMin + 1, prevMax +1)
        })
    }
    

  }

  var updateStdDev = function() {

    if(graph.stop) {
      return 0;
    }

    var sum = 0
    for (var i in workers) {
      sum += workers[i].q_val
    }

    if(sum == 0) {
      $(".std-dev h2").html(0)
      return;
    }

    var avg = sum / workers.length

    console.log(sum)

    average = sum / workers.length
    
    var varSum = 0;

    for (var i in workers) {
      varSum +=  Math.pow((workers[i].q_val - avg),2)
    }
  
    $(".std-dev h2").html((Math.sqrt(varSum / workers.length)).toFixed(4))

  }

  return {
    build: build,
    start: start,
    stop: stop,
    update: update
  }

}



