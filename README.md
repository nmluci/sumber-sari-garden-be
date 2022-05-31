# Sumber Sari Garden Backend

## Table of Contents
- [Sumber Sari Garden Backend](#sumber-sari-garden-backend)
  - [Table of Contents](#table-of-contents)
  - [Preface](#preface)
  - [Design Pattern](#design-pattern)
    - [Domain Table](#domain-table)
  - [Database Viewpoints](#database-viewpoints)
  - [Information System Viewpoint](#information-system-viewpoint)
  - [Queries Ran](#queries-ran)
  - [Notes](#notes)
  

## Preface
Herein lies the source code of Sumber Sari Garden's Online Shop's Backend, which was coded in Go for one of several course completions' requirement such as Database Practicum, Information System, and System Design and Analysis. Some technology used in this Object{} is Go Language, SQL, and Markdown. 

## Design Pattern
DDD (Domain Driven Design) Pattern

### Domain Table
Domain | Folder
-------|--------------
auth   | Authentication & Users
constant | Constant used in projects globally
dto | Requests & Responses Structs
models | Entity Structs (Repository-approachs)
ping | Ping Service (debugging purpose)
product | Product, Coupon, and Category
usercart | Carts, Checkout, Verification

## Database Viewpoints
- [x] DDL (db/migrations)
- [x] INSERT (product, usercart)
- [x] UPDATE (product, usercart)
- [x] DELETE (product, usercart)
- [x] SELECT w/ Logical Op & Like (usercart)
- [x] SELECT w/ Order By (usercart)
- [x] SELECT w/ Alias (auth, product, usercart)
- [x] SELECT w/ Between (usercart)
- [x] SELECT w/ JOIN (product, usercart)
- [x] SELECT w/ UNION (usercart)
- [x] SELECT w/ IN (usercart)
- [x] Aggregation Func (usercart)
- [x] Having (usercart)
- [ ] VIEW (not possible, too complex)
- [ ] Stored Procedure and Cursor (already handled by backend logics obviously...)
- [x] Trigger (product, usercart)

## Information System Viewpoint
Category | Implementations
-------- | ---------------``
SCM      | Products
TPS      | Order, Checkout

## Queries Ran
[Click Me!](/docs/queries.md)
## Notes
- compiled program (as per 2022-05-21 00:27:00 AM), is in RELEASE section (FOR WINDOWS)
- for other os, compiles it by yourself please :)