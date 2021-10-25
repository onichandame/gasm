# Go CRUD

A pattern for CRUD generators.

# Concetps

A full-stack CRUD app contains 3 layers of logic:

1. Data Storage
2. Business Logic
3. Input/Output Interface

The basic data flow for requests is: interface -> business logic -> data storage -> business logic -> interface. It clearly resembles a onion model.

For most CRUD apps, the business logic does not change much. The challenge is the data models and interfaces that changes often. But the business logic layer in between sometimes needs to be updated due to the changing data/interface. This change is considered an excess work and the main pain point to be solved by this tool.

In order to model the data flow above, 5 entities of interest are involved.

## Query

Query is the input interface for read-many requests. It includes filtering, ordering, pagination and field selection. It is defined by `core.Query`

## Entity

Entity is the data model in the database. It's representation can vary based on the ORM chosen.

## DTO

DTO is the output interface for all the requests expected to return the data retrieved from the data storage. It can be represented by a struct with json tags.

For create/update requests, separate DTOs can also be defined to represent the input interface. As in many scenarios the input/output interfaces are asymmetric.

## Query Service

Query Service is an interface with all the methods to interact with the data storage. It varies with the ORM chosen.

## Assembler

Assembler is an interface with all the methods required to convert between DTO and Entity. A simple assembler is provided by `core.DefaultAssembler`. Based on the varying interfaces/data models, custom assemblers can be made to meet the demands.

# Implementation

As discussed above, the business logic of CRUD apps can be defined by the conversion policies between the Entity and the DTO, the Assembler. On the other side, the data retrieval/save are defined by the query service. Hence the architecture can be split into 2 parts: HTTP adaptor and database adaptor.

## HTTP Adaptor

This adaptor is responsible for coordinating all the components introduced so far, along with making http handlers. All components should be fed into this adaptor to make a full-blown app.

## Database Adaptor

This adaptor should implement a query service for the consumption by the HTTP adaptor.
