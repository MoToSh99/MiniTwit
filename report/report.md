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

### Important interactions of subsystems
- Forklarer hvordan de forskellige ting er afhængige af hinanden
- Diagram 

### The current state of your systems, for example using results of static analysis and quality assessment systems
- OWASP ZAP
- SonarCloud
- Performance (Grafana)

## Process' perspective
### How do you interact as developers?
- Covid19
- Fysisk til online 
- Online messaging

### How is the team organized?
- Fixed workdays
- Self-organizing team
- One person takes an issue and tries to solve it, if not possible ask for help
    - 
### A complete description of stages and tools included in the CI/CD chains
- Diagram
- De forskellige trin i vores pipeline
- Pull request -> Checks(Deepsource(security), CodeFactor(Linting), Sonarcloud(Code complexity)) -> Code review -> Push to master -> Github Actions Release Drafter (Creates a release) -> Published Release (manually) -> Build DockerImage - Push to Dockerhub -> Pull newest docker image on remote server -> Docker compose up with new image

### Organization of your repositor(ies)
- Mono repository
- Opdelt efter features (Front-end, API, monitoring, logging etc.)
- Github.itu migration to github.com (why and when)
### Applied branching strategy
- Git-flow branching model
- simpelt diagram - hvorfor ikke en bredere model
### Applied development process and tools supporting it
- Centralised workflow
- 
### How do you monitor your systems and what precisely do you monitor?
- Grafana and prometheus
### What do you log in your systems and how do you aggregate logs?
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