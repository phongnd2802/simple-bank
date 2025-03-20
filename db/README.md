# Deeply Understand Transaction Isolation Levels

## Read Phenomena

### 1. Dirty Read
A **dirty read** happens when a transaction reads data written by another concurrent transaction that has not been committed yet. This is problematic because we don't know if that other transaction will eventually be committed or rolled back. If a rollback occurs, we might end up using incorrect data.

### 2. Non-repeatable Read
A **non-repeatable read** occurs when a transaction reads the same record twice and sees different values because another transaction modified and committed the row after the first read.

### 3. Phantom Read
A **phantom read** is a similar phenomenon but affects queries that search for multiple rows instead of just one. In this case, when the same query is re-executed, a different set of rows is returned due to changes made by other recently committed transactions. These changes could include:
- Inserting new rows
- Deleting existing rows that satisfy the search condition of the current transaction's query

### 4. Serialization Anomaly
A **serialization anomaly** occurs when the result of a group of concurrently committed transactions cannot be achieved if we try to run them sequentially in any order without overlapping each other.

## 4 Standard Isolation Levels

### 1. Read Uncommitted
Transactions at this level can see data written by other uncommitted transactions, thus allowing the dirty read phenomenon to happen.

### 2. Read Committed
Transactions can only see data that has been committed by other transactions. Because of this, dirty reads are no longer possible.

### 3. Repeatable Read
This level ensures that the same `SELECT` query will always return the same result, no matter how many times it is executed, even if some other concurrent transactions have committed new changes that satisfy the query.

### 4. Serializable
Concurrent transactions running at this level are guaranteed to produce the same result as if they were executed sequentially in some order, one after another, without overlapping. Essentially, this means that there exists at least one way to order these concurrent transactions so that if we run them one by one, the final result will be the same.

