#Chinchilla
##Todd Pollak
##Teddy Cleveland

##Overview
Chinchilla is a project we built for our COMP0112 (Networking) course this semester (spring 2015). The idea behind the project
was to test out different load balancing algorithms and visualize them in a meaningful way. We decided to use go to build
the scheduler-worker infrastructure, and then built our testing server in go as well. The front end, which allows us to 
visualize the different scheduling algorithms at work, was built using d3. Though we didn't have much time to put intense 
testing and error handling in place, it's a fun project that we really enjoyed building. If you'd like to run the front end
interface on your local machine, follow the steps below:

## Building the executables:
```
sh compile
```
or
```
./compile
```

##Running the webserver:
```
./webserver 9010 9000
```

##Info about running
The first port is the one that the webserver runs on, and must be 9010 because we did a bad thing and hard coded the websocket port... The second port is where the scheduler will run, and can be anything you want. The scheduler uses port 9020 for TCP connections with workers so I wouldn't touch that.

##Understanding the GUI:
It's pretty straightforward, and everything is labeled, but I'll give it a quick explanation. The graph itself plots the length of each worker queue in seconds as a weighted average. The X-axis is just the time that the simulation has been running. In terms of options, you can switch between two algorithms, Round Robin and Shortest Queue. To begin, set the number of each type of request that will be sent per second (the approximate times that each takes is noted below the inputs). Then, set the number of workers want (you can add and drop them during the test, too), and hit begin. To end the test, hit end. To add workers, click the add worker button and to remove them, just click their little boxes that pop up down below. If you want to send a blast of requests to the scheduler all at once, just put the type (1, 2, or 3) in the blast box, the number you want to send, and hit blast.

