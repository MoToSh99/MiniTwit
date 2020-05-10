# DevOps report

The following document was written for the DevOps course at ITU. It describes the work and architecture for a project surrounding supporting and maintaining an online messaging platform similar to Twitter.

## System's perspective

### Architecture
This section gives a high-level abstraction of how the whole system has been implemented. Later in the report, a more concrete look at how the MiniTwit component has been designed will be described.

The MiniTwit system and its features are running in [docker](https://www.docker.com/) containers, stored on remote servers at [DigitalOcean](https://www.digitalocean.com/). The systems are split up on separate servers, with the MiniTwit system itself, the database, the monitoring, and the logging running on their own servers on DigitalOcean. This is shown in figure 1 below.


![Architecture][Architecture]
<center><i> Figure 1: The system architecture.</i> </center>
<br />

Developers work on the system on their development machines. The code is pulled from and pushed to [GitHub](https://github.com/), which is then built with [GitHub Actions](https://github.com/features/actions). From there, a docker image is created and sent to [Docker Hub](https://hub.docker.com/), and the system is deployed to a server on DigitalOcean. The API documentation, [Swagger UI](https://swagger.io/tools/swagger-ui/), is stored on the same server.

The MiniTwit service sends logs directly to [Logstash](https://www.elastic.co/logstash), which is run on the Logging server. [Prometheus](https://prometheus.io/), running on the monitoring server, gets metrics from the MiniTwit system every 5 seconds, and the data is then read by [Grafana](https://grafana.com/) and visualized on a dashboard. Data is read from and written to a [Postgres](https://www.postgresql.org/) database cluster and shown on the front-end part of the MiniTwit service. Both front-end and API users of the system communicate directly with the server storing MiniTwit and Swagger UI.
### Design
This section describes how the MiniTwit system is designed and what parts it's comprised of. The system is written in [Golang](https://golang.org/) and is split into multiple components. The Golang language was chosen because it seemed easier to work with when doing web development, and when taking the size of the project into consideration, Golang was the simple choice. 

Figure 2 below shows the different components present in the system. 

![Component diagram](https://i.imgur.com/ZTtKLJw.png)
<center><i> Figure 2: Component diagram of the system.</i> </center>
<br />

The following components can be seen in the diagram above:
- A main component containing three other components:
    - the API, to handle the similator.
    - the back-end of the clientside system.
    - the routes that are used by the API and the back-end.
- A handler component which handles the requests made to the back-end.
- A helper component that contains all the utilization methods that both the API and the back-end needs.
- A struct component that contains declarations of objects/structs.
- A metrics component which collects the metrics that have been specified.
- A logger component that sends the logs to Logstash.

### Dependencies
The system dependencies can be split into two parts: the dependencies of the code and the dependencies of the service.

The code dependencies are shown in figure 3 below. 
![Dependency Diagram][dependencyDiagram]
<center><i> Figure 3: The dependencies of the code.</i> </center>
<br />

The dependency diagram shows that internally, classes only depend on few other classes in the system. With most of these dependencies simply being helper methods, or classes only used within the designated category, the overall coupling of the system is low. The external dependencies, however, are quite numerous, which means that these continuously have to be accounted for. If any of the external dependencies stop being supported, the system will break, which makes the code more difficult to maintain. 

The essential technologies and tools that the system depends on are listed below. These are frameworks and services that are needed to keep the complete system running. 
- [Gorilla/mux](https://github.com/gorilla/mux): An HTTP router and URL matcher for the web server.
- DigitalOcean: For the PostgreSQL hosting and Ubuntu servers running MiniTwit.
- Github Actions: CI/CD tool to automate the workflow.
- Prometheus: Creating metrics to use for monitoring.
- Grafana: Visual monitoring dashboard solution.
- Swagger UI: Documentation for the API.

### Important interactions of subsystems

The code is not split into any strict subsystems, and there are no interfaces to adhere to. We could have created two subsystems, one for the front-end and one for the back-end, but the structure of the system was simply copied from the original python project, so no split was done.

#### API
The interface that the API subsystem has to adhere to was given by the creators of the simulator. The API interacts with every part of the system except for the client-side.

#### Internal
With the MiniTwit system not being split up into subsystems internally, no real communication happens between them. The system is comprised of several components, which can also be seen in figure 2.

### The current state of the system
Now that the simulator has been shut down, it is time to evaluate the current state of our application in terms of maintainability, reliability, portability and modifiability.
##### Maintainability
Overall, we would argue that our code is not very maintainable. Due to the fact that testing has mostly been omitted, there is no way to guarantee that changing the code does not break something in the system. On the other hand, we do have a fairly clear structured code base, and it is a small project, therefore it is not too critical. [SonarCloud](https://sonarcloud.io/) suggests that the maintainability rating is "A", and only points out one method as being "confusing and difficult to maintain". It also claims that there is only a code duplication of 1.4%, which is rather good. 

Figure 4 below shows the static analysis from SonarCloud.

![Latest SonarCloud status][SonarCloud]
<center><i> Figure 4: Latest SonarCloud status.</i> </center>

##### Reliability
The reliability of the system is highly dependent on the cloud services that it is hosted on, and no fault tolerance has been implemented. SonarCloud's static analysis insists that there are no bugs, but if something breaks, the system is shut down and we are not notified. Using restart policies for our docker containers, we have implemented that they automatically restart when they crash, unless the container has been manually stopped.
##### Portability
The MiniTwit service is packaged in a Docker container targeted for Linux. This means that the service is usable on all Linux systems as well as macOS systems that are able to run docker images. Windows with Docker should be able to use the docker image as well, since Docker will run a Linux VM under the hood. All of this gives our service high portability.

##### Modifiability
At the moment, modifiability is fairly low. To improve modifiability, subsystem interfaces would have to be created. Components are fairly replaceable at the moment, but having strict interfaces to adhere to would improve modifiability. 

## Process' perspective
### Interactions as developers
Initially, interactions were done in person, with communication taking place before and after lectures at school. Any communication between these sessions happened through [Messenger](https://www.messenger.com/). 

Due to COVID-19, the approach was changed, and everyone was forced to work remotely. The same strategy was kept, but instead of being physically in the same room, a [Discord](https://discord.com/) server was used. This gave the opportunity to have discussions in the group while the lectures were being given.
Generally, if additional communication was required, someone in the team called for a Discord meeting at a time where everyone was available.

Most of the development happened independently, but in the case of more people working together, screen sharing was used through Discord.

Decisions made in the group were decided through casual communication. Mostly, an open discussion was initiated, and an agreement on how to proceed was reached. No strict rules were set in place for decisions.

### Team organization
#### Fixed workdays
Two weekly workdays were planned, where the work to be done and the work already done would be discussed. In between, the individual team members would have tasks/issues that they would work on. If a task was too difficult, or the team member needed input from the team, the task would be discussed and worked on during the weekly meetings. 
#### Self-organizing team
The team has been working as a self-organized team. Each week after the lecture, the team would distribute the tasks to be done between them. This was done by the entire team, and then the team members themselves would pick their own assignments. If anyone ran into problems, they would contact the rest of the team and either ask for help or hand the issue over to another team member. The responsibility of reaching out to the rest of the team lay with the individual team member.
### CI/CD pipeline
A continuous delivery pipeline was chosen for this project because it allows for manual deployment to production. This means that any version of the system is ready to be released at any time, but the actual deployment is executed manually. This way, new versions of the system can be tested and worked on before deployment.

Figure 5 below depicts the stages in the CI/CD pipeline.


![CI pipeline][Pipeline]
<center><i> Figure 5: The CI/CD pipeline.</i> </center>
<br />


The sequence diagram above shows the steps involved in a deployment to production. Before any chain is started, a developer continuously pulls and pushes code to a feature branch on GitHub. When a pull request to the master branch is created, the following chain starts:

1. Github runs checks on the pull request. These checks include:
    -  **Deepsource.io** for finding and fixing issues before they land in production. These can be issues such as bugs, anti-patterns, performance issues, and security flaws. 
    -  **CodeFactor.io** for linting.
    -  **Sonarcloud.io** to ensure the quality of the code.
2. If all checks pass, at least two developers must do code review and approve the pull request before it can be merged with the master branch.
3. When the pull request is approved, GitHub Actions runs automatically, resulting in a release draft that is ready to be released to production. The release draft is created using [Release Drafter GitHub Action](https://github.com/release-drafter/release-drafter).
4. A developer manually publishes the release from GitHub, which again triggers GitHub Actions. Here, a Docker Image is built from the release and pushed to Docker Hub. 
5. GitHub Actions logs in to a machine on DigitalOcean and runs the following commands: `docker-compose pull`, which pulls the latest image from Docker Hub, and `docker-compose up -d`, which builds, creates, starts the MiniTwit service in detached mode.
6. When all the steps above have been executed, the newest version of the system is up and running.

Because new features of the system are heavily prioritised, most testing of the system is done manually with each release. Therefore, automatic unit testing isn't part of the deployment pipeline. It is instead manually executed after the release draft is created, but before the release is published. 

### Organization of the repository
A Mono repository was chosen to organise the artifacts in the system. Having all the source code in one repository makes the process of giving access to the entire code base and editing it very simple. This was chosen over distributed repositories due to the system not having large independent submodules that need their own repository.

Internally, the repository is organised in a folder called *app* containing all source code. Inside the app folder, the code is organised in folders representing the functionality (API, logging etc.)

After session 4, the task was to create a CI/CD setup for the project. Originally, the plan was to use [Travis CI](https://travis-ci.org/) as in the example from the session. However, Travis CI only works on private repositories, which means that alternatives had to be considered. In this process, Github Action was chosen. It turned out that this service was not available on the github.itu.dk domain, which meant that the project got transferred to github.com.

During this transfer, only the code was copied, which meant that all the 
commit history up to that point was left behind in the old repository. 

### Applied branching strategy
The git-flow branching model was chosen due to it being tried and tested in the group. Everyone had positive experiences with this strategy from similar size projects. Considering the small team size and modular work delegation, everyone agreed to use this model to avoid over-complicating the development process. 

### Applied development process and tools supporting it
 
The team is using the agile development process. At an early stage of the project, the team deployed the system and started doing weekly releases and deployments. These regular deployments, and the fact that the team is self organizing, is well in line with agile principles. 

In order to be able to have these regular releases, the team is using Git and GitHub. GitHub Actions is used to automate the release process. The team uses this tool to create releases and deploy the system onto a virtual machine hosted at DigitalOcean. Before the system is deployed to DigitalOcean, the code runs through a series of checks described in the *CI/CD pipeline* section.

In terms of git the team is using a centralized workflow, where each developer has a version of the software and pushes directly to the repository. The team chose the workflow as it goes well in line with the small size of the team and project.

### Monitoring of the system
Two tools are used for monitoring the system: Prometheus and Grafana. Whenever an API call is made, Prometheus registers the event. The data collected can then be retrieved in a Grafana dashboard that visualizes the data and collects statistics. 

The dashboard is shown in figure 6 below.

![Grafana dashboard][Grafana]
<center><i> Figure 6: The Grafana dashboard.</i> </center>
<br />

 The following data is shown on the dashboard:

* Total uptime of the system
* Number of messages sent
* Users registered
* Number of requests for the last hour
* Users followed and unfollowed
* Response times for the different API calls (follow, register, and post)

Monitoring played a crucial role in improving the performance of the system; whenever a change was made, any performance changes could quickly be observed on the Grafana dashboard. Furthermore, it was possible to observe the timeslots that users were most active during the day, which, in case of an unexpected error, allowed for affecting as few users as possible when releasing a new update to production.

### Logging
Since the simulator that was running was only sending requests to the API and not the front-end part, the API calls were prioritised in terms of logging. To open up the ability to find all possible errors that might occur with the API in the logs, every API call was logged.
Because the analysis that can be made by monitoring the system is quite limited, logging was utilised to diagnose any specific problems occurring with the API.

An [ELK stack](https://www.elastic.co/what-is/elk-stack) was used to aggregate logs. Everything that was logged was sent directly to the Logstash, meaning that there wasn't a need for a local log file. ElasticSearch, which is part of the ELK stack, was used for searching in the logs themselves, and [Kibana](https://www.elastic.co/kibana), the visualization tool that comes with it, was used to create a dashboard on which the logs were displayed. Combining these tools, logs could easily be ingested, searched in, and visualized. 

### Brief results of the security assessment
We have split our assets into two parts: the assets concerning the users of our service, and those concerning the system itself.
The assets concerning the users are defined as:
- Usernames
- Passwords
- Emails

The three items above can contain person specific information.
Additionally, many users use the same info for multiple services, making it critical to keep this information safe if we want to attract and keep users to our service.

The one asset defined that concerns the system itself is the database, since the attackers can modify or delete it if they have access to our credentials.

#### Pentesting of our own system
Using the [OWASP ZAP](https://www.zaproxy.org/) tool to detect possible vulnerabilities in our system, only one vulnerability raised a couple of red flags. One of them was because we hadn't included an *X-Frame-Options* header in one of our HTTP responses. This header is used to protect a service against [ClickJacking](https://www.imperva.com/learn/application-security/clickjacking/) attacks, and most modern Web browsers support this header. This vulnerability was subsequently fixed in our system by using the [Secure middleware](https://github.com/unrolled/secure) by [unrolled](https://github.com/unrolled).

Another red flag was the possibility of SQL injection. We do not have any restrictions on what usernames and emails can be when a user is created. Which means that it might be possible to inject a SQL query into the database this way, even though we do not use raw SQL to query our database. To prevent this we could regulate the input of the username and email, so that the username is a maximum of 16 consecutive characters and the email has to match a regex pattern.

The last red flag was path traversal, where the tool warned that one might be able to access files by manipulating the URL. The group tried to access files that were not supposed to be public with this method, but could not succeed in any case. This flag was therefore disregarded.

Figure 7 below shows a risk matrix for the possible attacks detected by the OWASP ZAP tool.
![Risk Matrix][RiskMatrix]
<center><i> Figure 7: Risk matrix.</i> </center>

### Applied strategy for scaling and load balancing
The system was scaled on a [docker swarm cluster](https://docs.docker.com/engine/swarm/) using [terraform](https://www.terraform.io/), following [@zanderhavgaard's guide](https://github.com/itu-devops/itu-minitwit-docker-swarm-teraform). The docker swarm cluster contained 10 replicas of the MiniTwit image, spread out on multiple droplets, and could easily be scaled up if necessary. This meant that even if a droplet crashed, the uptime of the system would not be affected. 

Figure 8 below shows the droplets running on DigitalOcean.

![Docker Swarm][DockerSwarm]
<center><i> Figure 8: Droplets on DigitalOcean.</i> </center>
<br />

Unfortunately, it was not possible to keep that amount of droplets running, since the team was already over budget and could not get any further investments to the project.

After testing the docker swarm and seeing how it automatically routed requests to containers in the cluster, thereby having no need for an external loadbalancer, the droplets were destroyed.

## Lessons Learned Perspective
**Refactoring code**

Initially, the team struggled with the complete refactoring of the codebase. Without a very detailed planning process, it was decided to write in C# using [Razor pages](https://docs.microsoft.com/en-us/aspnet/core/razor-pages/?view=aspnetcore-3.1&tabs=visual-studio) as a web framework. After a huge workload, the team had a codebase that was not working. The night before the refactored code was due, a team member decided to try refactoring to Golang. He managed to get the code running, and, due to time constraints, the group decided to continue with this solution. This issue could have been avoided by analysing the assignment better and planning the refactoring from the knowledge gained.
        
**Working remotely**
    
With the very abrupt switch to forced remote work, everyone was quite anxious about how it was going to play out. Fortunately, the transition was very smooth, and no real issues came with it. In a lot of ways work has become more focused, due to the human social aspect being cut off. Not having the in-person interaction has other drawbacks, but in this project it is agreed upon by the team that it did not have an effect on the work.
    
**Migrating from an Azure to a Postgres database**

After running the database locally, a switch was made to an Azure database. The performance was great at first because of the low amount of users. However, as the userbase and data grew larger, the performance was affected significantly.

The problem arose because two things were implemented at the same time: the migration to a cloud database, and the introduction of a database abstraction layer ([Gorm](https://gorm.io/)). When the performance started deteriorating, it was assumed that Gorm was the culprit - it is, after all, one of the slowest abstraction layers for Golang (as discussed in [this issue](https://github.com/jinzhu/gorm/issues/298)). Altering the queries used to fetch the data had little to no effect, and attempts to change the abstraction layer were not successful. 

Later, it was discovered by the team that the free tier of Azure limited the performance remarkably. Using a database migration tool called [Full Convert](https://www.spectralcore.com/fullconvert), the database was successfully migrated to a Postgres database cluster hosted on DigitalOcean ([pull request 58](https://github.com/MoToSh99/MiniTwit/pull/58)). The performance suddenly improved notably, and the average response time went from 9 seconds to 0.07 seconds 

From this, we learned the following things:
- *Only one thing should be changed at a time*. The realization that the problem came from Azure would've occurred much earlier if Gorm hadn't been implemented at the same time. 
- *Even if a solution is an industry standard, you can't assume that it's the best solution for you.* Azure database solutions are used heavily in the industry, but with this system having to use the free tier, a lot of problems could've been avoided by doing the research first.

The change of performance can be observed by the purple line in figure 9 below. 

![Latest ID for our api][LatestID]
<center><i> Figure 9: The growth of latest_ids on all MiniTwit services. This service is the purple line.</i> </center>
<br />
 
The figure shows the latest id which gets updated when requests are sent to the API from the simulator. After steadily growing in the beginning, it suddenly grows significantly after migrating to Postgres.

**DevOps Style**

The style of this DevOps project differentiated it from other projects in a few ways:
* **Refactoring an old/existing project.** In other projects, systems were created from scratch with a set of requirements already defined from the get-go. Here, a simple MiniTwit service had already been made, and a translation of the code to another language was required. Also, more requirements kept getting added, which meant that sometimes existing implementations had to be changed in order to allow the new features to be implemented. 

* **CI/CD Pipelines.** Focusing on continuous delivery and having a pipeline was something new from all other projects the team has been working on so far. CI/CD pipelines had not been a focus on any other project, and therefore was very interesting to learn about and implement. It definitely accelerated the process of commiting some code to Github to getting the new version of our MiniTwit service online, while at the same time making sure that all the steps for deploying the new version were done in the correct order. CI/CD pipelines are undoubtedly something the team will consider in any of our future projects.

* **Servers.** While other projects only require an MVP that isn't actually supposed to be used, this project was under constant pressure from a simulator that used the API to send requests. This meant that during the course, errors would suddenly occur and would have to be fixed, and improvements had to be made to accommodate the amount of requests being sent. Actually having to maintain servers (like the database and the server for the service itself) was a challenge unlike anything we had encountered in other courses, and it taught us that having a running service isn't just as simple as making an implementation and then forgetting about it. 

[dependencyDiagram]: https://github.com/MoToSh99/MiniTwit/blob/dev/report/images/Dependencies.png?raw=true "Dependency Diagram"

[Pipeline]:
https://github.com/MoToSh99/MiniTwit/blob/dev/report/images/Pipeline.png?raw=true 

[Grafana]:
https://github.com/MoToSh99/MiniTwit/blob/dev/report/images/Grafana.PNG?raw=true

[RiskMatrix]: https://raw.githubusercontent.com/MoToSh99/MiniTwit/dev/report/images/riskMatrix.png

[DockerSwarm]:
https://github.com/MoToSh99/MiniTwit/blob/dev/report/images/dockerswarm.png?raw=true

[Architecture]:
https://github.com/MoToSh99/MiniTwit/blob/dev/report/images/Architecture.png?raw=true

[LatestID]:
https://github.com/MoToSh99/MiniTwit/blob/dev/report/images/Latestid.PNG?raw=true

[SonarCloud]:
https://github.com/MoToSh99/MiniTwit/blob/dev/report/images/SonarCloud.png?raw=true