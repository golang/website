<!--{
  "Title": "Managing connections"
}-->

For the vast majority of programs, you needn't adjust the `sql.DB` connection
pool defaults. But for some advanced programs, you might need to tune the
connection pool parameters or work with connections explicitly. This topic
explains how.

The [`sql.DB`](https://pkg.go.dev/database/sql#DB) database handle is safe for
concurrent use by multiple goroutines
(meaning the handle is what other languages might call “thread-safe”). Some
other database access libraries are based on connections that can only be used
for one operation at a time. To bridge that gap, each `sql.DB` manages a pool
of active connections to the underlying database, creating new ones as needed
for parallelism in your Go program. 

The connection pool is suitable for most data access needs. When you call an
`sql.DB` `Query` or `Exec` method, the `sql.DB` implementation retrieves an
available connection from the pool or, if needed, creates one. The package
returns the connection to the pool when it's no longer needed. This supports a
high level of parallelism for database access.

### Setting connection pool properties {#connection_pool_properties}

You can set properties that guide how the `sql` package manages a connection
pool. To get statistics about the effects of these properties, use
[`DB.Stats`](https://pkg.go.dev/database/sql#DB.Stats).

#### Setting the maximum number of open connections {#max_open_connections}

[`DB.SetMaxOpenConns`](https://pkg.go.dev/database/sql#DB.SetMaxOpenConns)
imposes a limit on the number of open connections. Past this limit, new
database operations will wait for an existing operation to finish, at which
time `sql.DB` will create another connection. By default, `sql.DB` creates a
new connection any time all the existing connections are in use when a
connection is needed.

Keep in mind that setting a limit makes database usage similar to acquiring a
lock or semaphore, with the result that your application can deadlock waiting
for a new database connection.

#### Setting the maximum number of idle connections {#max_idle_connections}

[`DB.SetMaxIdleConns`](https://pkg.go.dev/database/sql#DB.SetMaxIdleConns)
changes the limit on the maximum number of idle connections `sql.DB`
maintains.

When an SQL operation finishes on a given database connection, it is not
typically shut down immediately: the application may need one again soon, and
keeping the open connection around avoids having to reconnect to the database
for the next operation. By default an `sql.DB` keeps two idle connections at
any given moment. Raising the limit can avoid frequent reconnects in programs
with significant parallelism.

#### Setting the maximum amount a time a connection can be idle {#max_idle_time}

[`DB.SetConnMaxIdleTime`](https://pkg.go.dev/database/sql#DB.SetConnMaxIdleTime)
sets the maximum length of time a connection can be idle before it is closed.
This causes the `sql.DB` to close connections that have been idle for longer
than the given duration.

By default, when an idle connection is added to the connection pool, it
remains there until it is needed again. When using `DB.SetMaxIdleConns` to
increase the number of allowed idle connections during bursts of parallel
activity, also using `DB.SetConnMaxIdleTime` can arrange to release those
connections later when the system is quiet.

#### Setting the maximum lifetime of connections {#max_connection_lifetime}

Using [`DB.SetConnMaxLifetime`](https://pkg.go.dev/database/sql#DB.SetConnMaxLifetime)
sets the maximum length of time a connection can be held open before it is
closed.

By default, a connection can be used and reused for an arbitrarily long amount
of time, subject to the limits described above. In some systems, such as those
using a load-balanced database server, it can be helpful to ensure that the
application never uses a particular connection for too long without reconnecting.

### Using dedicated connections {#dedicated_connections}

The `database/sql` package includes functions you can use when a database may
assign implicit meaning to a sequence of operations executed on a particular
connection. 

The most common example is transactions, which typically start with a `BEGIN`
command, end with a `COMMIT` or `ROLLBACK` command, and include all the
commands issued on the connection between those commands in the overall
transaction. For this use case, use the `sql` package’s transaction support.
See [Executing transactions](/doc/database/execute-transactions).

For other use cases where a sequence of individual operations must all execute
on the same connection, the `sql` package provides dedicated connections.
[`DB.Conn`](https://pkg.go.dev/database/sql#DB.Conn) obtains a dedicated
connection, an [`sql.Conn`](https://pkg.go.dev/database/sql#Conn). The
`sql.Conn` has methods `BeginTx`, `ExecContext`, `PingContext`,
`PrepareContext`, `QueryContext`, and `QueryRowContext` that behave like the
equivalent methods on DB but only use the dedicated connection. When finished
with the dedicated connection, your code must release it using `Conn.Close`.
