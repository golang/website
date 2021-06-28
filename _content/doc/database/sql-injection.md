<!--{
  "Title": "Avoiding SQL injection risk"
}-->

You can avoid an SQL injection risk by providing SQL parameter values as `sql`
package function arguments. Many functions in the `sql` package provide
parameters for the SQL statement and for values to be used in that statement's
parameters (others provide a parameter for a prepared statement and parameters).

Code in the following example uses the `?` symbol as a placeholder for the
`id` parameter, which is provided as a function argument:

```
// Correct format for executing an SQL statement with parameters.
rows, err := db.Query("SELECT * FROM user WHERE id = ?", id)
```

`sql` package functions that perform database operations create prepared
statements from the arguments you supply. At run time, the `sql` package turns
the SQL statement into a prepared statement and sends it along with the
parameter, which is separate.

**Note:** Parameter placeholders vary depending on the DBMS and driver
you're using. For example, [pq driver](https://pkg.go.dev/github.com/lib/pq)
for Postgres accepts a placeholder form such as `$1` instead of `?`.

You might be tempted to use a function from the `fmt` package to assemble the
SQL statement as a string with parameters included – like this:

```
// SECURITY RISK!
rows, err := db.Query(fmt.Sprintf("SELECT * FROM user WHERE id = %s", id))
```

This is not secure! When you do this, Go assembles the entire SQL statement,
replacing the `%s` format verb with the parameter value, before sending the
full statement to the DBMS. This poses an
[SQL injection](https://en.wikipedia.org/wiki/SQL_injection) risk because the
code's caller could send an unexpected SQL snippet as the `id` argument. That
snippet could complete the SQL statement in unpredictable ways that are
dangerous to your application.

For example, by passing a certain `%s` value, you might end up with something
like the following, which could return all user records in your database:

```
SELECT * FROM user WHERE id = 1 OR 1=1;
```