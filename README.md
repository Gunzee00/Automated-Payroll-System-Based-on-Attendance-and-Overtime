# Payroll & Attendance System API

A backend system built with Golang that manages attendance, payroll, overtime, and reimbursement with JWT authentication and audit logging.

## Tech Stack
- **Language**: Golang
- **Database**: PostgreSQL
- **ORM**: GORM
- **Auth**: JWT

## Installation

git clone https://github.com/Gunzee00/Automated-Payroll-System-Based-on-Attendance-and-Overtime.git
**cd project-name
go mod tidy
Database:
This project uses two separate PostgreSQL databases. 
- Main database for the application runtime
- Test database for running unit tests
Make sure both databases are created and accessible
dealls_test (for production/development)
db_test (for unit testing)

## API usage 

- Authentication
POST:  /api/login
Authenticate user and receive a JWT token for authorization.

- Attendance Period (Admin Only)
POST: /api/admin/attendance-periods
Create a new attendance period.
Access: Admin only

- Attendance (Employee)
POST: /api/attendance/submit
Submit daily attendance. Restricted to working days (Mon–Fri).
Access: Authenticated employees

- Overtime (Employee)
POST: /api/attendance/overtime
Submit overtime hours to be included in the payroll calculation.
Access: Authenticated employees

- Reimbursement (Employee)
POST: /api/employee/reimbursements
Submit a reimbursement request with amount and description.
Access: Authenticated employees

- Payroll Processing (Admin Only)
POST /api/admin/run-payroll
Process payroll for a specific attendance period. Automatically calculates take-home pay based on attendance, overtime, and reimbursements.
Access: Admin only

- Payslip Generation (Employee)
POST /api/employee/payslips/generate
Generate a payslip for a specific attendance period. Each user can only generate once per period.
Access: Authenticated employees

GET /api/employee/payslips
Retrieve all payslips belonging to the currently authenticated user.
Access: Authenticated employees

- Payslip Summary (Admin Only)
GET /api/admin/payslips/summary
Retrieve summary of total take-home pay for all employees. Useful for payroll reporting.
Access: Admin only

## Software Software Architecture

- Routing Layer
Implemented using gorilla/mux.
All routes are defined in main.go or router.go using route prefixes (/api/...).
Middleware like JWTAuthMiddleware and AdminOnly is used to protect routes based on roles.

-  Handler Layer
Contains the business logic of each endpoint.
Located in the handlers/ directory.
Each function maps directly to a RESTful route, such as SubmitAttendance, RunPayroll, etc.

- Model Layer
All database models are defined using GORM in the models/ directory.
Includes models like User, Attendance, Payroll, Payslip, Reimbursement, etc.

- Middleware Layer
JWT verification (JWTAuthMiddleware)
Role authorization (AdminOnly)
Extracts user information and validates access before hitting handlers.

- Utility Layer
Located in utils/ directory.
Contains helper functions for JWT handling, password hashing, request ID generation, IP extraction, etc.

- Configuration Layer
Database connection and setup are handled in config/.
Supports both development and test environments with two separate databases.

- Testing Layer
Unit tests for all handlers are placed under handlers_test/.

Tests use isolated test database to avoid polluting production/development data.




