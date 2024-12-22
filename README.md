# DevOps project
A repository for my project for Modern Devops Practises course at FMI

# Todo List App

The classic todo list app should support creating different lists of todos. Each list should have its own name and description. In the lists, there are todos that are separate from the others. The access to lists should be managed in the database by user emails. Users should be able to log in to the application via a GitHub account. The users should be part of a specific GitHub organization in order to be able to log in.

## Implementation

### Version 1

Implement a Todo microservice that supports the following:

- **Login with GitHub**
- **Create/Modify/Delete todo lists and todos**
- **Save Todo lists and todos in the database**

### Requirements for Todo Lists

- Each todo list has an owner (the user who created the todo list).
- Owners of the todo list can read/modify todos in the list.
- Owners can add new users to the todo list.
- Only users that are part of a certain Todo list can read/modify it.
- Admins can modify every todo list.

### Requirements for GitHub Organization

The GitHub organization should have 3 teams - **readers, writers, and admins**:

- **Readers**: Participants can read todo lists they are part of.
- **Writers**:
  - Participants can create/update/delete todos in the todo lists they are part of.
  - Participants can create new Todo lists and will be part of them by default.
  - Participants can add users to Todo lists if they are the owner of the Todo list.
- **Admins**: Participants can do all of the above even without being part of a Todo list.

### Architecture Overview

- **REST Server**:
  - Implement a basic HTTP server.
  - Define Models:
    - `lists`
    - `todos`
  - Define and Implement routers for:
    - Create/Update/Delete/Get Lists
    - Create/Update/Delete/Get Todos
  - Implement service layers that are used by the routers:
    - `lists`
    - `todos`

- **Database**:
  - Define database tables, relations, primary keys, foreign keys, indices (if needed):
    - Create a local DB on `localhost:5432`.
    - Create an initial SQL migration that sets up the DB status.
    - Define GoLang DB Entities.
  - Implement converters that convert DB entities to models and models to entities.
  - Implement a repo layer that talks with the database and use it in the service layer.

- **GraphQL**:
  - Implement another server - GraphQL API facade that calls the REST API.
  - Define the GraphQL schema.
  - Call the REST APIs.

- **Security**:
  - Implement tenancy:
    - Define Users table and add "ownership" to the lists so that only the owners can modify and view the lists.
    - Use a tenant header that is propagated from the GraphQL Facade to the REST API and have the service and repo layer take the user into consideration.
  - **Auth/Authz**:
    - Define 3 GitHub organizations - read/write/admin.
    - Create a GitHub OAuth app.
    - Implement a login endpoint.
    - Implement a middleware that authenticates the user using a token.
    - Implement a middleware that authorizes the user.

- **UI**:
  - Implement a login page.
  - Implement a UI for the Lists of todos:
    - Create a new List.
    - Update existing List (name, description).
    - Delete a List.
  - Upon clicking a List, open a new page that shows all the todos for this List:
    - Add a Todo.
    - Mark todo as completed.
    - Delete todo.

- **Docker**:
  - Install Docker Desktop.
  - Create Dockerfiles for all components.

- **Kubernetes**:
  - Install Kubernetes locally using k3s for creating a cluster.
  - Define the Helm charts, which should include:
    - **REST Microservice**:
      - Deployment
      - Service
      - ConfigMap/Secret
    - **GraphQL Microservice**:
      - Deployment
      - Service
      - ConfigMap/Secret
    - **Postgres**:
      - Deployment
      - Persistent Volume Claim/PV (probably)
      - Service
      - Seed/Migration job
    - **UI**:
      - Deployment
      - Service

## Set up

- **run ./setup.sh 8000 {cluster-name}**
- **UI runs on port http://localhost:8000/**
