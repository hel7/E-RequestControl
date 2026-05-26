## System Overview

The primary needs of the administrator include:

* Viewing, deleting, and editing requests;
* Viewing and deleting notifications;
* Managing system users: creating, editing, deleting, and viewing user accounts;
* Data management: exporting data, creating backups of database settings and data, and restoring backups.

The primary needs of the user include:

* Creating, viewing, editing, and deleting requests;
* Viewing and deleting notifications.

---

## System Architecture

The server-side architecture of the software is built according to the principles of the three-tier model. It consists of:

* HTTP request handling layer (Handler);
* Business logic layer (Service);
* Data access layer (Repository).

This architectural approach ensures clear separation of responsibilities, simplifies maintenance and scalability, and improves system stability.

---

## Backend Implementation

The server-side was developed using the Go programming language, which is known for its high performance, maintainability, and support for concurrent execution.

The main responsibilities of the backend server include:

* Processing client requests;
* Implementing business logic for request management;
* Interacting with the database;
* Providing reliable responses for each request.

---

## API Testing

To verify the correct functionality of the REST API, testing was performed using Postman.

A set of test scenarios was created to validate:

* Correct handling of HTTP requests;
* Data format validation;
* System behavior in error scenarios;
* Compliance of HTTP status codes with REST standards.

This approach helped identify potential logical errors at early development stages and ensured that the API behavior matched the expected functionality.
