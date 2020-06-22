Platform - The Application
=================================
The application adopts a clean architecture. The folder structure within platform 
represent the layers of the clean architecutre. 

Entities
------------
Entities are placed in the model folder/subfolders. Entities can implement organisational
business rules. 

Application Services
--------------------------
Application services are placed in the domain folder. This represents the application 
busines rules.

Repository
------------
The repository is used for CRUD operation definitions. It is interface only

Adapter 
---------
The adapter folder provides methods to retreive repositories that are database specific.
When an application starts up it constructs an appropriate adapter and queries the 
adapter to retreive repositories.

Service
---------
The service folder represents services that can be consumed by clients, such as GraphQL,
web sockets, cli, etc.


