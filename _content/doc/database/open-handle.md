<!--{
  "Title": "Opening a database handle",
  "Breadcrumb": true
}-->

The [`database/sql`](https://pkg.go.dev/database/sql) package simplifies
database access by reducing the need
for you to manage connections. Unlike many data access APIs, with
`database/sql` you don't explicitly open a connection, do work, then close
the connection. Instead, your code opens a database handle that represents
a connection pool, then executes data access operations with the handle,
calling a `Close` method only when needed to free resources, such as those
held by retrieved rows or a prepared statement.

In other words, it's the database handle, represented by an
[`sql.DB`](https://pkg.go.dev/database/sql#DB), that
handles connections, opening and closing them on your code's behalf. As your
code uses the handle to execute database operations, those operations have
concurrent access to the database. For more, see
[Managing connections](/doc/database/manage-connections).

**Note:** You can also reserve a database connection. For more
information, see
[Using dedicated connections](/doc/database/manage-connections#dedicated_connections).

In addition to the APIs available in the `database/sql` package, the Go
community has developed drivers for all of the most common (and many uncommon)
database management systems (DBMSes).

When opening a database handle, you follow these high-level steps:

1. Locate a driver.

    A driver translates requests and responses between your Go code and the
    database. For more, see [Locating and importing a database driver](#database_driver).

2. Open a database handle.

    After you've imported the driver, you can open a handle for a specific
    database. For more, see [Opening a database handle](#opening_handle).

3. Confirm a connection.

    Once you've opened a database handle, your code can check that a
    connection is available. For more, see [Confirming a connection](#confirm_connection).

Your code typically won’t explicitly open or close database connections -- that's
done by the database handle. However, your code should free resources it
obtains along the way, such as an `sql.Rows` containing query results. For
more, see [Freeing resources](#free_resources).

### Locating and importing a database driver {#database_driver}

You'll need a database driver that supports the DBMS you're using. To locate
a driver for your database, see [SQLDrivers](https://github.com/golang/go/wiki/SQLDrivers).

To make the driver available to your code, you import it as you would
another Go package. Here's an example:

```
import "github.com/go-sql-driver/mysql"
```

Note that if you're not calling any functions directly from the driver
package –- such as when it's being used implicitly by the `sql` package --
you'll need to use a blank import, which prefixes the import path with an
underscore:


```
import _ "github.com/go-sql-driver/mysql"
```

**Note:** As a best practice, avoid using the database driver's own API
for database operations. Instead, use functions in the `database/sql`
package. This will help keep your code loosely coupled with the DBMS,
making it easier to switch to a different DBMS if you need to.

### Opening a database handle {#opening_handle}

An `sql.DB` database handle provides the ability to read from and write to a
database, either individually or in a transaction.

You can get a database handle by calling either `sql.Open` (which takes a
connection string) or `sql.OpenDB` (which takes a `driver.Connector`). Both
return a pointer to an [`sql.DB`](https://pkg.go.dev/database/sql#DB).

**Note:** Be sure to keep your database credentials out of your Go source.
For more, see [Storing database credentials](#store_credentials).

#### Opening with a connection string {#open_connection_string}

Use the [`sql.Open` function](https://pkg.go.dev/database/sql#Open) when you
want to connect using a connection string. The format for the string will vary
depending on the driver you're using. 

Here's an example for MySQL:

```
db, err = sql.Open("mysql", "username:password@tcp(127.0.0.1:3306)/jazzrecords")
if err != nil {
	log.Fatal(err)
}
```

However, you'll likely find that capturing connection properties in a more
structured way gives you code that's more readable. The details will vary by
driver.

For example, you could replace the preceding example with the following, which
uses the MySQL driver's [`Config`](https://pkg.go.dev/github.com/go-sql-driver/mysql#Config)
to specify properties and its
[`FormatDSN method`](https://pkg.go.dev/github.com/go-sql-driver/mysql#Config.FormatDSN)
to build a connection string.

```
// Specify connection properties.
cfg := mysql.Config{
	User:   username,
	Passwd: password,
	Net:    "tcp",
	Addr:   "127.0.0.1:3306",
	DBName: "jazzrecords",
}

// Get a database handle.
db, err = sql.Open("mysql", cfg.FormatDSN())
if err != nil {
	log.Fatal(err)
}
```

#### Opening with a Connector {#open_connector}

Use the [`sql.OpenDB function`](https://pkg.go.dev/database/sql#OpenDB) when
you want to take advantage of driver-specific connection features that aren't
available in a connection string. Each driver supports its own set of
connection properties, often providing ways to customize the connection request
specific to the DBMS.

Adapting the preceding `sql.Open` example to use `sql.OpenDB`, you could
create a handle with code such as the following:

```
// Specify connection properties.
cfg := mysql.Config{
	User:   username,
	Passwd: password,
	Net:    "tcp",
	Addr:   "127.0.0.1:3306",
	DBName: "jazzrecords",
}

// Get a driver-specific connector.
connector, err := mysql.NewConnector(&cfg)
if err != nil {
	log.Fatal(err)
}

// Get a database handle.
db = sql.OpenDB(connector)
```

#### Handling errors {#handle_errors}

Your code should check for an error from attempting to create a handle, such
as with `sql.Open`. This won't be a connection error. Instead, you'll get an
error if `sql.Open` was unable to initialize the handle. This could happen,
for example, if it's unable to parse the DSN you specified.

### Confirming a connection {#confirm_connection}

When you open a database handle, the `sql` package may not create a new
database connection itself right away. Instead, it may create the connection
when your code needs it. If you won't be using the database right away and
want to confirm that a connection could be established, call
[`Ping`](https://pkg.go.dev/database/sql#DB.Ping) or
[`PingContext`](https://pkg.go.dev/database/sql#DB.PingContext).

Code in the following example pings the database to confirm a connection.

```
db, err = sql.Open("mysql", connString)

// Confirm a successful connection.
if err := db.Ping(); err != nil {
	log.Fatal(err)
}
```

### Storing database credentials {#store_credentials}

Avoid storing database credentials in your Go source, which could expose the
contents of your database to others. Instead, find a way to store them in a
location outside your code but available to it. For example, consider a
secret keeper app that stores credentials and provides an API your code can
use to retrieve credentials for authenticating with your DBMS.

One popular approach is to store the secrets in the environment before the
program starts, perhaps loaded from a secret manager, and then your Go program
can read them using [`os.Getenv`](https://pkg.go.dev/os#Getenv):

```
username := os.Getenv("DB_USER")
password := os.Getenv("DB_PASS")
```

This approach also lets you set the environment variables yourself for local
testing. 

### Freeing resources {#free_resources}

Although you don't manage or close connections explicitly with the
`database/sql` package, your code should free resources it has obtained when
they're no longer needed. Those can include resources held by an `sql.Rows`
representing data returned from a query or an `sql.Stmt` representing a
prepared statement.

Typically, you close resources by deferring a call to a `Close` function so
that resources are released before the enclosing function exits.

Code in the following example defers `Close` to free the resource held by
[`sql.Rows`](https://pkg.go.dev/database/sql#Rows).

```
rows, err := db.Query("SELECT * FROM album WHERE artist = ?", artist)
if err != nil {
	log.Fatal(err)
}
defer rows.Close()

// Loop through returned rows.
```
