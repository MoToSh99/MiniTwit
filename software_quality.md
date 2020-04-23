# Software Quality in General

**Identify and list qualities of your MiniTwit system from the four perspectives (except of the transcendental) from the Kitchenham paper.**

* User view 
	* Reliability: Our service is distributed to third party hosts, which minimizes downtime. The responsibility for this is spread out because of an SLA from the third party services that we can rely on. 
	* Usability: Our service can be accessed from both a front-end website as well as an API.
* Manufacturing view
	* Our rework costs have been high as a result of a wrong database service being chosen the first time around, meaning that the product wasn't constructed "right the first time".
* Product view
	* We've implemented internal quality indicators in the form of measuring resonse times (Grafana) and logging all uses of the service, which has helped us improve the external product behavior. 
* Value-based view
	* Since our service is free-to-use, it's difficult to measure its value.
	* Because the project hasn't been developed with focus on it being sold, this view hasn't been considered during development.

**Did you focus on any perspective or any qualities, perhaps even without being aware of it? If yes, list these.**

* We've focused on the product view when we've implemented monitoring tools to improve the service.
* When switching databases, we focused on the performance of the service.

**Rank the identified qualities per perspective by decreasing importance to you and provide an argument for why you choose certain as the most important.**

* Most important: Performance, since our performance from the get-go was poor and service that could actually be accessed was our primary focus.
* Least important: Monitoring, because we already knew that the performance of the service was lacking, and monitoring tools didn't provide solutions to the existing problem.

**Think about and discuss with your group fellows, how you can measure the qualities that you ranked the most important. That is, try to define a set of metrics that would allow to measure these (multiple metrics per quality can be possible).**

* Performance is easily measured by e.g. accessing the site, looking at our monitoring dashboard, and assessing the amount of errors registered in our logs.

# How maintainable are your systems?

**How can you identify and measure maintainability of your MiniTwit systems?**

* Maintainability in our system can be measured by lines of code, coupling of our classes, number of nested function calls, and complexity of the code. 

**Is your MiniTwit system maintainable?**

* Due to our codebase being fairly simple, well defined, and loosely coupled, we have kept the complexity of our system down, increasing our maintainability. Because of these characteristics, we feel that our maintainability is high. 

**Collect a set of characteristics that make your system maintainable. Try to include more than just the source code.**

* Loose coupling, extensive monitoring, error logging, and simple source code.

# Do you have Technical Debt in your systems?

**What is Technical Debt for you?**

* Using the newest versions of dependencies, regular updating, and conformity to the newest industry standards all help lowering technical debt. 

**Describe how you could identify and measure it?**

* Monitoring the number of issues posted on GitHub.
* Using a package manager to regularly update dependencies. 
* Using an external service to analyse the code. 

**Result of running SonarQube**


![Result](https://i.imgur.com/gB9Ix6t.png)



**Technical debt according to the documentation of SonarQube**

* "Effort to fix all Code Smells". Maintainability of the system, and how long it will take to fix Code Smells. 

**Does that correspond to your understanding of technical debt?**

* No. Our understanding concerned updating dependencies and using industry standards for the system, while the documentation defined technical debt as the effort to fix bad code, like complexity, naming conventions and class density (Code Smells). 
