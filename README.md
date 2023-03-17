Go Example

A brief description of the project.
Prerequisites

Before installing this project, you will need to install Go on your system. Please visit the official Go website for installation instructions.
Installation

Clone the repository:


$ git clone https://github.com/example/go-example.git

Building the Project

To build the project, navigate to the project directory and run the go build command:
```
$ cd go-example
$ go build
```
This will create an executable file named go-example in the project directory.
Running the Project

To run the project, from project directory run the go-example command:

```
$ ./go-example
```
This will start the application and it will begin running.

Once the application is running, you can use the following curl commands to add items to the "schedule":

```
curl -X POST "http://127.0.0.1:8080/api/v1/schedule/add" -d '{"id": "bar", "start_time": "2023-03-17T22:11:17.99999Z"}'

curl -X POST "http://127.0.0.1:8080/api/v1/schedule/add" -d '{"id": "bar", "start_time": "2023-03-17T22:11:17.99999Z"}'
```
Note that you may need to modify the start_time field to match your desired schedule time.
License

This project is licensed under the MIT License.
