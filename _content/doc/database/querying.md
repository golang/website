<!--{
  "Title": "Querying for data"
}-->

When executing an SQL statement that returns data, use one of the `Query`
methods provided in the `database/sql` package. Each of these returns a `Row`
or `Rows` whose data you can copy to variables using the `Scan` method.
You'd use these methods to, for example, execute `SELECT` statements.

When executing a statement that doesn’t return data, you can use an `Exec` or
`ExecContext` method instead. For more, see
[Executing statements that don't return data](/doc/database/change-data).

The `database/sql` package provides two ways to execute a query for results.

*   **Querying for a single row** – `QueryRow` returns at most a single `Row`
    from the database. For more, see [Querying for a single row](#single_row).
*   **Querying for multiple rows** – `Query` returns all matching rows as a
    `Rows` struct your code can loop over. For more, see
    [Querying for multiple rows](#multiple_rows).

If your code will be executing the same SQL statement repeatedly, consider
using a prepared statement. For more, see
[Using prepared statements](/doc/database/prepared-statements).

**Caution:** Don't use string formatting functions such as `fmt.Sprintf` to
assemble an SQL statement! You could introduce an SQL injection risk. For more,
see [Avoiding SQL injection risk](/doc/database/sql-injection).

### Querying for a single row {#single_row}

`QueryRow` retrieves at most a single database row, such as when you want to
look up data by a unique ID. If multiple rows are returned by the query, the
`Scan` method discards all but the first.

`QueryRowContext` works like `QueryRow` but with a `context.Context` argument.
For more, see [Canceling in-progress operations](/doc/database/cancel-operations).

The following example uses a query to find out if there's enough inventory to
support a purchase. The SQL statement returns `true` if there's enough, `false`
if not. [`Row.Scan`](https://pkg.go.dev/database/sql#Row.Scan) copies the
boolean return value into the `enough` variable through a pointer.

```
func canPurchase(id int, quantity int) (bool, error) {
	var enough bool
	// Query for a value based on a single row.
	if err := db.QueryRow("SELECT (quantity >= ?) from album where id = ?",
		quantity, id).Scan(&enough); err != nil {
		if err == sql.ErrNoRows {
			return false, fmt.Errorf("canPurchase %d: unknown album", id)
		}
		return false, fmt.Errorf("canPurchase %d: %v", id)
	}
	return enough, nil
}
```

**Note:** Parameter placeholders in prepared statements vary depending on the
DBMS and driver you're using. For example, the
[pq driver](https://pkg.go.dev/github.com/lib/pq) for Postgres requires a
placeholder like `$1` instead of `?`.

#### Handling errors {#single_row_errors}

`QueryRow` itself returns no error. Instead, `Scan` reports any error from the
combined lookup and scan. It returns
[`sql.ErrNoRows`](https://pkg.go.dev/database/sql#ErrNoRows) when the query
finds no rows.

#### Functions for returning a single row {#single_row_functions}

<table id="single-row-functions-list" class="DocTable">
  <thead>
    <tr class="DocTable-head">
      <th class="DocTable-cell" width="20%">Function</th>
      <th class="DocTable-cell">Description</th>
    </tr>
  </thead>
  <tbody>
    <tr class="DocTable-row">
      <td class="DocTable-cell">
        <code><a href="https://pkg.go.dev/database/sql#DB.QueryRow">DB.QueryRow</a></code><br />
        <code><a href="https://pkg.go.dev/database/sql#DB.QueryRowContext">DB.QueryRowContext</a></code>
      </td>
      <td class="DocTable-cell">Run a single-row query in isolation.</td>
    </tr>
    <tr class="DocTable-row">
      <td class="DocTable-cell">
        <code><a href="https://pkg.go.dev/database/sql#Tx.QueryRow">Tx.QueryRow</a></code><br />
        <code><a href="https://pkg.go.dev/database/sql#Tx.QueryRowContext">Tx.QueryRowContext</a></code>
      </td>
      <td class="DocTable-cell">Run a single-row query inside a larger transaction. For more, see
        <a href="/doc/database/execute-transactions">Executing transactions</a>.
      </td>
    </tr>
    <tr class="DocTable-row">
      <td class="DocTable-cell">
        <code><a href="https://pkg.go.dev/database/sql#Stmt.QueryRow">Stmt.QueryRow</a></code><br />
        <code><a href="https://pkg.go.dev/database/sql#Stmt.QueryRowContext">Stmt.QueryRowContext</a></code>
      </td>
      <td class="DocTable-cell">Run a single-row query using an already-prepared statement. For more,
        see <a href="/doc/database/prepared-statements">Using prepared statements</a>.
      </td>
    </tr>
    <tr class="DocTable-row">
        <td class="DocTable-cell">
  <code><a href="https://pkg.go.dev/database/sql#Conn.QueryRowContext">Conn.QueryRowContext</a></code>
      </td>
      <td class="DocTable-cell">For use with reserved connections. For more, see
        <a href="/doc/database/manage-connections">Managing connections</a>.
      </td>
    </tr>
  </tbody>
</table>

### Querying for multiple rows {#multiple_rows}

You can query for multiple rows using `Query` or `QueryContext`, which return
a `Rows` representing the query results. Your code iterates over the returned
rows using [`Rows.Next`](https://pkg.go.dev/database/sql#Rows.Next). Each
iteration calls `Scan` to copy column values into variables. 

`QueryContext` works like `Query` but with a `context.Context` argument. For
more, see [Canceling in-progress operations](/doc/database/cancel-operations).

The following example executes a query to return the albums by a specified
artist. The albums are returned in an `sql.Rows`. The code uses
[`Rows.Scan`](https://pkg.go.dev/database/sql#Rows.Scan) to copy column values
into variables represented by pointers.

```
func albumsByArtist(artist string) ([]Album, error) {
	rows, err := db.Query("SELECT * FROM album WHERE artist = ?", artist)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// An album slice to hold data from returned rows.
	var albums []Album

	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var alb Album
		if err := rows.Scan(&alb.ID, &alb.Title, &alb.Artist,
			&alb.Price, &alb.Quantity); err != nil {
			return albums, err
		}
		albums = append(albums, album)
	}
	if err = rows.Err(); err != nil {
		return albums, err
	}
	return albums, nil
}
```

Note the deferred call to [`rows.Close`](https://pkg.go.dev/database/sql#Rows.Close).
This releases any resources held by the rows no matter how the function
returns. Looping all the way through the rows also closes it implicitly,
but it is better to use `defer` to make sure `rows` is closed no matter what.

**Note:** Parameter placeholders in prepared statements vary depending on
the DBMS and driver you're using. For example, the
[pq driver](https://pkg.go.dev/github.com/lib/pq) for Postgres requires a
placeholder like `$1` instead of `?`.

#### Handling errors {#multiple_rows_errors}

Be sure to check for an error from `sql.Rows` after looping over query results.
If the query failed, this is how your code finds out.

#### Functions for returning multiple rows {#multiple_rows_functions}

<table id="multiple-row-functions-list" class="DocTable">
  <thead>
    <tr class="DocTable-head">
      <th class="DocTable-cell" width="20%">Function</th>
      <th class="DocTable-cell">Description</th>
    </tr>
  </thead>
  <tbody>
    <tr class="DocTable-row">
      <td class="DocTable-cell">
        <code><a href="https://pkg.go.dev/database/sql#DB.Query">DB.Query</a></code><br />
        <code><a href="https://pkg.go.dev/database/sql#DB.QueryContext">DB.QueryContext</a></code>
      </td>
      <td class="DocTable-cell">Run a query in isolation.</td>
    </tr>
    <tr class="DocTable-row">
      <td class="DocTable-cell">
        <code><a href="https://pkg.go.dev/database/sql#Tx.Query">Tx.Query</a></code><br />
        <code><a href="https://pkg.go.dev/database/sql#Tx.QueryContext">Tx.QueryContext</a></code>
      </td>
      <td class="DocTable-cell">Run a query inside a larger transaction. For more, see
        <a href="/doc/database/execute-transactions">Executing transactions</a>.
      </td>
    </tr>
    <tr class="DocTable-row">
      <td class="DocTable-cell">
        <code><a href="https://pkg.go.dev/database/sql#Stmt.Query">Stmt.Query</a></code><br />
        <code><a href="https://pkg.go.dev/database/sql#Stmt.QueryContext">Stmt.QueryContext</a></code>
      </td>
      <td class="DocTable-cell">Run a query using an already-prepared statement. For more, see
        <a href="/doc/database/prepared-statements">Using prepared
          statements</a>.
    </td>
    </tr>
    <tr class="DocTable-row">
      <td class="DocTable-cell">
        <code><a href="https://pkg.go.dev/database/sql#Conn.QueryContext">Conn.QueryContext</a></code>
      </td>
      <td class="DocTable-cell">For use with reserved connections. For more, see
        <a href="/doc/database/manage-connections">Managing connections</a>.
      </td>
    </tr>
  </tbody>
</table>

### Handling nullable column values {#nullable_columns}

The `database/sql` package provides several special types you can use as
arguments for the `Scan` function when a column's value might be null. Each
includes a `Valid` field that reports whether the value is non-null, and a
field holding the value if so.

Code in the following example queries for a customer name. If the name value
is null, the code substitutes another value for use in the application.

```
var s sql.NullString
err := db.QueryRow("SELECT name FROM customer WHERE id = ?", id).Scan(&s)
if err != nil {
	log.Fatal(err)
}

// Find customer name, using placeholder if not present.
name := "Valued Customer"
if s.Valid {
	name = s.String
}
```

See more about each type in the `sql` package reference:

*    [`NullBool`](https://pkg.go.dev/database/sql#NullBool)
*    [`NullFloat64`](https://pkg.go.dev/database/sql#NullFloat64)
*    [`NullInt32`](https://pkg.go.dev/database/sql#NullInt32)
*    [`NullInt64`](https://pkg.go.dev/database/sql#NullInt64)
*    [`NullString`](https://pkg.go.dev/database/sql#NullString)
*    [`NullTime`](https://pkg.go.dev/database/sql#NullTime)

### Getting data from columns {#column_data}

When looping over the rows returned by a query, you use `Scan` to copy a row’s
column values into Go values, as described in the
[`Rows.Scan`](https://pkg.go.dev/database/sql#Rows.Scan) reference.

There is a base set of data conversions supported by all drivers, such as
converting SQL `INT` to Go `int`. Some drivers extend this set of conversions;
see each individual driver's documentation for details.

As you might expect, `Scan` will convert from column types to Go types that
are similar. For example, `Scan` will convert from SQL `CHAR`, `VARCHAR`, and
`TEXT` to Go `string`. However, `Scan` will also perform a conversion to
another Go type that is a good fit for the column value. For example, if the
column is a `VARCHAR` that will always contain a number, you can specify a
numeric Go type, such as `int`, to receive the value, and `Scan` will convert
it using `strconv.Atoi` for you.

For more detail about conversions made by the `Scan` function, see the [`Rows.Scan`](https://pkg.go.dev/database/sql#Rows.Scan) reference.

### Handling multiple result sets {#multiple_result_sets}

When your database operation might return multiple result sets, you can
retrieve those by using
[`Rows.NextResultSet`](https://pkg.go.dev/database/sql#Rows.NextResultSet).
This can be useful, for example, when you're sending SQL that separately queries
multiple tables, returning a result set for each.

`Rows.NextResultSet` prepares the next result set so that a call to
`Rows.Next` retrieves the first row from that next set. It returns a boolean
indicating whether there is a next result set at all.

Code in the following example uses `DB.Query` to execute two SQL statements.
The first result set is from the first query in the procedure, retrieving all
of the rows in the `album` table. The next result set is from the second query,
retrieving rows from the `song` table.

```
rows, err := db.Query("SELECT * from album; SELECT * from song;")
if err != nil {
	log.Fatal(err)
}
defer rows.Close()

// Loop through the first result set.
for rows.Next() {
	// Handle result set.
}

// Advance to next result set.
rows.NextResultSet()

// Loop through the second result set.
for rows.Next() {
	// Handle second set.
}

// Check for any error in either result set.
if err := rows.Err(); err != nil {
	log.Fatal(err)
}
```
