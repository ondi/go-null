/*
do not use multiple types in type swtich case
https://github.com/golang/go/issues/38675

https://golang.org/pkg/database/sql/#Scanner

type Scanner
Scanner is an interface used by Scan.

type Scanner interface {
    // Scan assigns a value from a database driver.
    //
    // The src value will be of one of the following types:
    //
    //    int64
    //    float64
    //    bool
    //    []byte
    //    string
    //    time.Time
    //    nil - for NULL values
    //
    // An error should be returned if the value cannot be stored
    // without loss of information.
    //
    // Reference types such as []byte are only valid until the next call to Scan
    // and should not be retained. Their underlying memory is owned by the driver.
    // If retention is necessary, copy their values before the next call to Scan.
    Scan(src interface{}) error
}

https://golang.org/pkg/database/sql/driver/#Valuer

type Valuer
Valuer is the interface providing the Value method.

Types implementing Valuer interface are able to convert themselves to a driver Value.

type Valuer interface {
    // Value returns a driver Value.
    // Value must not panic.
    Value() (Value, error)
}

type Value
Value is a value that drivers must be able to handle. It is either nil, a type handled by a database driver's NamedValueChecker interface, or an instance of one of these types:

int64
float64
bool
[]byte
string
time.Time
If the driver supports cursors, a returned Value may also implement the Rows interface in this package. This is used, for example, when a user selects a cursor such as "select cursor(select * from my_table) from dual". If the Rows from the select is closed, the cursor Rows will also be closed.

type Value interface{}

*/

package null
