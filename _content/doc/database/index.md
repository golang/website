<!--{
  "Title": "Accessing relational databases",
  "Breadcrumb": true
}-->

Using Go, you can incorporate a wide variety of databases and data access
approaches into your applications. Topics in this section describe how to use
the standard library's [`database/sql`](https://pkg.go.dev/database/sql)
package to access relational databases.

For an introductory tutorial to data access with Go, please see
[Tutorial: Accessing a relational database](/doc/tutorial/database-access).

Go supports other data access technologies as well, including ORM libraries
for higher-level access to relational databases, and also non-relational
NoSQL data stores.

*   **Object-relational mapping (ORM) libraries.** While the `database/sql`
    package includes functions for lower-level data access logic, you can
    also use Go to access data stores at a higher abstraction level. For more
    about two popular object-relational mapping (ORM) libraries for Go, see
    [GORM](https://gorm.io/index.html) ([package reference](https://pkg.go.dev/gorm.io/gorm))
    and [ent](https://entgo.io/) ([package reference](https://pkg.go.dev/entgo.io/ent)).
*   **NoSQL data stores.** The Go community has developed drivers for the
    majority of NoSQL data stores, including [MongoDB](https://docs.mongodb.com/drivers/go/)
    and [Couchbase](https://docs.couchbase.com/go-sdk/current/hello-world/overview.html).
    You can search [pkg.go.dev](https://pkg.go.dev/) for more.

### Supported database management systems {#supported_dbms}

Go supports all of the most common relational database management systems,
including MySQL, Oracle, Postgres, SQL Server, SQLite, and more.

You'll find a complete list of drivers at the
[SQLDrivers](https://github.com/golang/go/wiki/SQLDrivers) page.

### Functions to execute queries or make database changes {#functions}

The `database/sql` package includes functions specifically designed for the
kind of database operation you're executing. For example, while you can use
`Query` or `QueryRow` to execute queries, `QueryRow` is designed for the case
when you're expecting only a single row, omitting the overhead of returning
an `sql.Rows` that includes only one row. You can use the `Exec` function
to make database changes with SQL statements such as `INSERT`, `UPDATE`, or
`DELETE`.

For more, see the following:

*   [Executing SQL statements that don't return data](/doc/database/change-data)
*   [Querying for data](/doc/database/querying)

### Transactions {#transactions}

Through `sql.Tx`, you can write code to execute database operations in a
transaction. In a transaction, multiple operations can be performed together
and conclude with a final commit, to apply all the changes in one atomic
step, or a rollback, to discard them.

For more about transactions, see [Executing transactions](/doc/database/execute-transactions).

### Query cancellation {#query_cancellation}

You can use `context.Context` when you want the ability to cancel a database
operation, such as when the client's connection closes or the operation runs
longer than you want it to.

For any database operation, you can use a `database/sql` package function
that takes `Context` as an argument. Using the `Context`, you can specify a
timeout or deadline for the operation. You can also use the `Context` to
propagate a cancellation request through your application to the function
executing an SQL statement, ensuring that resources are freed up if they're
no longer needed.

For more, see [Canceling in-progress operations](/doc/database/cancel-operations).

### Managed connection pool {#connection_pool}

When you use the `sql.DB` database handle, you're connecting with a built-in
connection pool that creates and disposes of connections according to your
code's needs. A handle through `sql.DB` is the most common way to do
database access with Go. For more, see
[Opening a database handle](/doc/database/open-handle).

The `database/sql` package manages the connection pool for you. However, for
more advanced needs, you can set connection pool properties as described in
[Setting connection pool properties](/doc/database/manage-connections#connection_pool_properties).

For those operations in which you need a single reserved connection, the
`database/sql` package provides [`sql.Conn`](https://pkg.go.dev/database/sql#Conn).
`Conn` is especially useful when a transaction with `sql.Tx` would be a
poor choice.

For example, your code might need to:

*   Make schema changes through a DDL, including logic that contains its
    own transaction semantics. Mixing `sql` package transaction functions with
    SQL transaction statements is a poor practice, as described in
    [Executing transactions](/doc/database/execute-transactions).
*   Perform query locking operations that create temporary tables.

For more, see [Using dedicated connections](/doc/database/manage-connections#dedicated_connections).
