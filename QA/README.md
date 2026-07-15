# QA Testing Report — Request Manager API

## Project Overview

Manual API testing of a diploma project built with **Go (Gin framework)** — a request/ticket management system with JWT-based authentication, role-based access control (user/admin), tickets CRUD, and a notification system.

**Testing approach:** Black-box API testing via Postman, covering functional, negative, boundary, and security (access control) test cases.

## Summary

| Metric | Count |
|---|---|
| Test cases executed | 13 |
| Bugs found | 4 |
| Critical severity | 2 |
| Medium severity | 2 |
| Modules covered | Tickets (CRUD), Authentication, Notifications |

**Key finding:** A critical broken-authentication vulnerability (`BUG-003`) allows any unauthenticated user to create an administrator account.

---

## Test Cases

| ID | Title | Priority | Result |
|---|---|---|---|
| IDTC-001 | Creation of ticket without optional "description" field | High | Pass |
| IDTC-002 | Creation of ticket with both "title" and "description" fields | High | Pass |
| IDTC-003 | Creation of tickets by different valid users | High | Pass |
| IDTC-004 | Creation of ticket without required "title" field | High | Pass |
| IDTC-005 | Ticket creation without authorization token | High | Pass |
| IDTC-006 | Access to dashboard without an active session | High | Pass |
| IDTC-007 | Unauthorized ticket update — User1 → User2's ticket (IDOR) | High | **Fail → BUG-001** |
| IDTC-008 | Unauthorized ticket update — User2 → User1's ticket (IDOR) | High | **Fail → BUG-001** |
| IDTC-009 | Unauthorized ticket deletion — User1 → User2's ticket | High | **Fail → BUG-002** |
| IDTC-010 | Unauthorized ticket deletion — User2 → User1's ticket | High | **Fail → BUG-002** |
| IDTC-011 | Mass assignment via `RoleID` field in `POST /api/auth/register` | Medium | Pass |
| IDTC-012 | `GET /api/notifications` returns only the logged-in user's notifications | Medium | Pass |

Detailed step-by-step test case documentation (Preconditions / Steps / Expected / Actual) is available in the https://docs.google.com/spreadsheets/d/1jDJLtdq8zAVaHIGheazsf-gNH7O5NOrBNuLB6yRx8RI/edit?usp=sharing

---

## Bug Reports

### BUG-001 — Broken Access Control on Ticket Update (IDOR)
| Field | Details |
|---|---|
| **Severity** | Critical |
| **Priority** | High |
| **Environment** | Local, API level (Postman). Not exposed via UI. |
| **Steps to Reproduce** | 1. Log in as User1, create a ticket, note the ticket ID.<br>2. Log in as User2.<br>3. Send `PUT /api/tickets/{User1_ticket_id}` using User2's token.<br>4. Send the request. |
| **Expected Result** | `403 Forbidden` — a user should not be able to modify tickets they don't own. |
| **Actual Result** | `200 OK` — ticket successfully updated by a non-owner. |

### BUG-002 — Incorrect Status Code on Unauthorized Ticket Deletion
| Field | Details |
|---|---|
| **Severity** | Medium |
| **Priority** | Medium |
| **Environment** | Local, API level (Postman) |
| **Steps to Reproduce** | 1. Log in as User1, create a ticket, note the ticket ID.<br>2. Log in as User2.<br>3. Send `DELETE /api/tickets/{User1_ticket_id}` using User2's token. |
| **Expected Result** | `403 Forbidden` or `404 Not Found`, with a clear error message. |
| **Actual Result** | `500 Internal Server Error`, message: `"ticket not found or does not belong to user"`. Access control works, but the error is surfaced as a server crash rather than a handled response. |

### BUG-003 — Unauthenticated Admin Account Creation (Critical)
| Field | Details |
|---|---|
| **Severity** | Critical |
| **Priority** | High |
| **Environment** | Local, API level (Postman) |
| **Steps to Reproduce** | 1. Send `POST /api/auth/registerAdmin` with a valid user body.<br>2. Do **not** include an `Authorization` header.<br>3. Send the request. |
| **Expected Result** | `401 Unauthorized` — this endpoint should require an active admin session to create new admin accounts. |
| **Actual Result** | `200 OK` — an admin account is created with **no authentication required**. |
| **Notes** | Root cause: `registerAdmin` is registered in the public `/auth` route group, before the `userIdentity`/`adminRequired` middleware is applied. Suggested fix: remove this endpoint entirely and create admins exclusively through `/api/admin/users` (which already requires admin auth), or add `userIdentity` + `adminRequired` middleware directly to this route. |

### BUG-004 — Misleading Success Response on Notification Deletion
| Field | Details |
|---|---|
| **Severity** | Medium |
| **Priority** | Medium |
| **Environment** | Local, API level (Postman) |
| **Steps to Reproduce** | 1. Log in as User A, get a notification ID belonging to User A.<br>2. Log in as User B.<br>3. Send `DELETE /api/notifications/{User_A_notification_id}` using User B's token.<br>4. Check the response, then verify the database state directly. |
| **Expected Result** | `403 Forbidden` or `404 Not Found`; response should reflect the actual outcome. |
| **Actual Result** | `200 OK`, message `"marked as read"` — but the record is **not** modified in the database. The API reports success even though nothing happened. |

---

## Tools Used
- **Postman** — manual API request testing
- **Docker logs** — verifying backend behavior and response codes
- **Direct database inspection** — confirming actual data state vs. API response claims

## Tested By
Bohdan — QA Manual, self-directed testing of a personal diploma project.
