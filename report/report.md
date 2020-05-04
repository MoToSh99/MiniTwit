# DevOps report

## System's perspective
### Design of your ITU-MiniTwit systems
- Looking at the old code
- Make it work
- Database ORM
### Architecture of your ITU-MiniTwit systems
- Diagram of the system 
    - Api is one system, frontend part is one system,.. 


### All dependencies of your ITU-MiniTwit systems on all levels of abstraction and development stages
- Digram over de forskellige containers
- Forskellige trin i vores program/programmer


![alt text][dependencyDiagram]

### Important interactions of subsystems
- Forklarer hvordan de forskellige ting er afhængige af hinanden
- Diagram 

### The current state of your systems, for example using results of static analysis and quality assessment systems
- OWASP ZAP
- SonarCloud
- Performance (Grafana)

## Process' perspective
### How do you interact as developers?
Initially the group interacted in person, communicating before and after lectures at school. Any communication between these sessions, happened through messenger. 

Due to COVID-19, the approach was changed. Everyone was forced to work remotely. The same strategy was kept, instead of being physically in the same room, a discord server was used. Giving the opporunity to discuss in the group, while the lectures were being given.

Generally if additional communication was required, someone in the team called for a discord meeting at a time that everyone was available.

Most of the development happened independantly, but in the case of more people working together, screen sharing was used through discord. 

### How is the team organized? - Emil
#### Fixed workdays
The had twoo weekly meetup days, where the work to be done and the work already done would be discussed. Inbetween the individual team memebers would have task/issues that they would due. If a task was to difficult or the team member needed the input from the team the task could be dicussed and worked on during these weekly meetings . 
#### Self-organizing team
Our team has been working as a self organized team. Each week, after the lecture the team would distribute the tasks to be done between them. Then if anyone ran into problems they would contact the rest of the team and either ask for help or hand the issue over to another team member. The responsibility of reaching out to the rest of the team lays with the  individual team member. 
### A complete description of stages and tools included in the CI/CD chains - TM
- Diagram
- De forskellige trin i vores pipeline
- Pull request -> Checks(Deepsource(security), CodeFactor(Linting), Sonarcloud(Code complexity)) -> Code review -> Push to dev -> Push to master -> Github Actions Release Drafter (Creates a release) -> Published Release (manually) -> Build DockerImage - Push to Dockerhub -> Pull newest docker image on remote server -> Docker compose up with new image

DeepSource : GO (Security)
CodeFactor (Linting)
SonarCloud (Clean Code)

### Organization of your repositor(ies)
- Mono repository
- Opdelt efter features (Front-end, API, monitoring, logging etc.)
- Github.itu migration to github.com (why and when)
### Applied branching strategy
The git-flow branching model was chosen, due to it being tried and tested in the group. Everyone had positive previous experience with this strategy for similar projects. Considering the small team size and modular work delegation everyone agreed to use this model to avoid over-complicating the development process.

![image alt][GithubFlow]
### Applied development process and tools supporting it
- Centralised workflow
- 
### How do you monitor your systems and what precisely do you monitor? - TM
- Grafana and prometheus
### What do you log in your systems and how do you aggregate logs? - TM
- All API calls
### Brief results of the security assessment
- https://github.com/MoToSh99/MiniTwit/blob/dev/software_quality.md
### Applied strategy for scaling and load balancing
- Docker swarm terraform (Se guiden)
    - Kostede kassen
## Lessons Learned Perspective
### Describe the biggest issues, how you solved them, and which are major lessons learned
- Refactoring code
    - From python to C# to Golang
- Working remotely has not been a problem
    - Discord is your friend
- Database migration
    - Local DB to azure 
    - From Azure (MSsql) to Digital Ocean (Postgres)



- Noter 
- Tjek session for extra



[dependencyDiagram]: https://github.com/MoToSh99/MiniTwit/blob/dev/report/images/Dependencies.png?raw=true "Dependency Diagram"

[GithubFlow]: https://raw.githubusercontent.com/MoToSh99/MiniTwit/dev/report/images/GithubFlow.png