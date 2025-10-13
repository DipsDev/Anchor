<img width="150" height="150" alt="anchorflow_final" src="https://github.com/user-attachments/assets/d3fa4f16-b879-433a-9efa-2ecf7c83401f" />


 # ⚓ Anchor

A powerful environment tool to streamline development with ease.

Anchor automates the entire CI/CD pipeline, building, configuring,
and deploying your application into a dedicated developer or staging environment.
The goal is maximum efficiency, letting developers validate changes instantly without managing complex scripts or infrastructure.


## Key Features

* Single-Command Deployments: Create and trigger your dev pipeline with a single command.
* Simple Configuration: All config is stored inside one `Anchorfile` file.


## Configuration Reference

Anchor uses a single configuration file paired with a powerful CLI tool, in order to simply your development process.

### Configuration File

Create a file named `Anchorfile` in the root of your project. This file specifies **where** and **how** to deploy.
It uses a [HCLv2](https://github.com/hashicorp/hcl) syntax, which provides a robust way to configure environments.

### Environments

Anchor configuration is built using a root block named `environment`, which defines a series of `processes`.
Each environment is defined using a name and an optional description.
```
environment "my-great-env" {
    description = "This environment does a lot of things!"
    
    ... ... {
        ...
    }
}

environment "my-second-env" {
    ...
}
```

### Processes

Within the `environment` blocks, you define individual processes using blocks.
These processes fall into two categories, determined by their execution behavior:

* A `service` is a long-running process (e.g, a web server or database) that should continue
operating until explicitly stopped.

* A `task` is a short-lived process (e.g, a migration script) that completes and exits once its command is finished.

#### Tasks
Tasks may receive the following attributes:

| Attribute  | Type         | Required | Description                                                               |
|------------|--------------|----------|---------------------------------------------------------------------------|
| command    | string       | true     | What command should run when used.                                        |
| depends_on | list(string) | false    | A list of names of other processes that must complete successfully first. |


#### Services
Services may receive the following attributes:

| Attribute    | Type         | Required | Description                                                                                                                                         |
|--------------|--------------|----------|-----------------------------------------------------------------------------------------------------------------------------------------------------|
| engine       | string       | true     | What engine should the command use (e.g: docker, shell). For more information about a specific engine, check out [Anchor Engines](docs/engines.md). |
| depends_on   | list(string) | false    | See Tasks Above                                                                                                                                     |
| health_check | block        | true     | Defines the criteria for confirming that a service is operational.                                                                                  |



##### health_check
```
service "generic-service" {
    ...
    health_check {
        type = "http" # perform an http get request
        target = "http://localhost:3000/api/health_check"
        timeout = "60s" # fail after 60s
    }
}
```

### Complete Example

This example features a fully-written Anchorfile. The file specifies a dev environment for fullstack website.
```hcl
environment "dev" {
    description = "Starts backend, frontend and a mysql database"

    service "db" {
        engine = "docker"
        image = "mysql:latest"
        
        health_check {
            type = "tcp"
            target = "localhost:3306"
            timeout = "15s"
        }
    }
    
    service "backend" {
        engine = "shell"
        command = "cd ./backend && npm run dev"
        depends_on = ["db"]
        
        health_check {
            type = "http"
            target = "http://localhost:8080/api/v1/health"
            # use the default timeout
        }
    }
    
    service "frontend" {
        engine = "shell"
        command = "cd ./frontend && npm run dev"
        depends_on = ["db"]
        
        health_check {
            type = "http"
            target = "http://localhost:3000/"
            timeout = "15s"
        }
    }
}
```

## Running
Once configuration is complete, the powerful CLI can be used to execute the environments.\
To launch an environment's services, execute the command: `anchor apply <environment>`.