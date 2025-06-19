# Payroll & Attendance System API

A backend system built with Golang that manages attendance, payroll, overtime, and reimbursement with JWT authentication and audit logging.

## Tech Stack
- **Language**: Golang
- **Database**: PostgreSQL
- **ORM**: GORM
- **Auth**: JWT

## Installation

git clone https://github.com/Gunzee00/Automated-Payroll-System-Based-on-Attendance-and-Overtime.git
<br> cd project-name
<br> go mod tidy
<br> Database: This project uses two separate PostgreSQL databases. 
<br> - Main database for the application runtime
<br> - Test database for running unit tests
<br>  Make sure both databases are created and accessible
<br> dealls_test (for production/development)
<br> db_test (for unit testing)

## API usage 

- Authentication
<br> Method: POST
<br> Endpoint: /api/login
<br> Authenticate user and receive a JWT token for authorization.

- Attendance Period (Admin Only)
<br> Method: POST
<br> Endpoint: /api/admin/attendance-periods
<br> Create a new attendance period.
<br> Access: Admin only

- Attendance (Employee)
<br> Method: POST
<br> Endpoint: /api/attendance/submit
<br> Submit daily attendance. Restricted to working days (Mon–Fri).
<br> Access: Authenticated employees

- Overtime (Employee)
<br> Method:  POST
<br> Endpoint: /api/attendance/overtime
<br> Submit overtime hours to be included in the payroll calculation.
<br> Access: Authenticated employees

- Reimbursement (Employee)
<br> Endpoint: POST
<br> Endpoint: /api/employee/reimbursements
<br> Submit a reimbursement request with amount and description.
<br> Access: Authenticated employees

- Payroll Processing (Admin Only)
<br> Endpoint: POST
<br> Endpoint: /api/admin/run-payroll
<br> Process payroll for a specific attendance period. Automatically calculates take-home pay based on attendance, overtime, and reimbursements.
<br> Access: Admin only

- Payslip Generation (Employee)
<br> Endpoint: POST
<br> Endpoint: /api/employee/payslips/generate
<br> Generate a payslip for a specific attendance period. Each user can only generate once per period.
<br> Access: Authenticated employees

- <br> Endpoint: GET
<br> Endpoint: /api/employee/payslips
<br> Retrieve all payslips belonging to the currently authenticated user.
<br> Access: Authenticated employees

- Payslip Summary (Admin Only)
<br> Endpoint: GET
<br> Endpoint: /api/admin/payslips/summary
<br> Retrieve summary of total take-home pay for all employees. Useful for payroll reporting.
<br> Access: Admin only

## Software Software Architecture

- Routing Layer
<br> Implemented using gorilla/mux.
<br> All routes are defined in main.go or router.go using route prefixes (/api/...).
<br> Middleware like JWTAuthMiddleware and AdminOnly is used to protect routes based on roles.

-  Handler Layer
<br> Contains the business logic of each endpoint.
<br> Located in the handlers/ directory.
<br> Each function maps directly to a RESTful route, such as SubmitAttendance, RunPayroll, etc.

- Model Layer
<br> All database models are defined using GORM in the models/ directory.
<br> Includes models like User, Attendance, Payroll, Payslip, Reimbursement, etc.

- Middleware Layer
<br> JWT verification (JWTAuthMiddleware)
<br> Role authorization (AdminOnly)
<br> Extracts user information and validates access before hitting handlers.

- Utility Layer
<br> Located in utils/ directory.
<br> Contains helper functions for JWT handling, password hashing, request ID generation, IP extraction, etc.

- Configuration Layer
<br> Database connection and setup are handled in config.
<br> Supports both development and test environments with two separate databases.

- Testing Layer
<br> Unit tests for all handlers are placed under handlers_test.
<br>  Tests use isolated test database to avoid polluting production/development data.




