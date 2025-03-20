package db

/*
	What is a database transaction?
	- A single unit of work.
	- Often made up of multiple db operations.
	
	Example: Transfer 10 USD from bank account A to bank account B
	This transaction comprises 5 operations:
	1. Create a transfer record with amount = 10.
	2. Create an account entry for account A with amount = -10.
	3. Create an account entry for account B with amount = 10.
	4. Subtract 10 amount from the balance of account A.
	5. Add 10 amount to the balance of account B.

	Why do we need db transaction?
	There are 2 main reasons:
	1. To provide a reliable and consistent unit of work, even in case of system failure.
	2. To provide isolation between programs that access the database concurrently.

	In order to achieve these 2 goals, a database transaction must satify the ACID properties.
	1. Atomicity (A): Either all operations complete successfully or the transaction fails and the db is unchanged.
	2. Consistency (C): The db state must be valid after the transaction. All contraints must be satisfied.
	3. Isolation (I): Concurrent transactions must not affect each other.
	4. Durability (D): Data written by a successfull transaction must be recorded in persistent storage.
*/


/*
	In PostgreSQL, ShareLock and ExclusiveLock are two types of locks used to control concurrent access to 
	database object such as tables, rows, or other resources.
	1. ShareLock: This type of lock allows multiple transactions to read data simultaneously but prevents 
	them from modifying it.
	Use case: Ensures that no transaction modifies the data while another transaction is reading it.
	Mechanism: 
		- Multiple transaction can hold a ShareLock on the same object.
		- If a transaction holds a ShareLock on row or table, other transactions can also acquire a ShareLock,
		but an ExclusiveLock will be blocked.
	
	2. ExclusiveLock: This type of lock allows prevents all other transctions from reading or writing to the 
	blocked object.
	Use case: Ensures that only one transaction can mofify an object at a time without conflicts.
	Mechanism:
		- When a transaction holds an ExclusiveLock, all other transactions are blocked from accessing the locked
		object, including both read and write operations.
		- Two transactions cannot hold an ExclusiveLock on the same object simultaneously.
*/