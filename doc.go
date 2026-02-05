/*

added support to IsEmptyJSON
https://golang.org/cl/241179
replaced with:
case reflect.Struct:
    return reflect.Zero(v.Type()).Interface() == v.Interface()

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
Value is a value that drivers must be able to handle. It is either nil,
a type handled by a database driver's NamedValueChecker interface,
or an instance of one of these types:

int64
float64
bool
[]byte
string
time.Time
If the driver supports cursors, a returned Value may also implement
the Rows interface in this package. This is used, for example,
when a user selects a cursor such as "select cursor(select * from my_table) from dual".
If the Rows from the select is closed, the cursor Rows will also be closed.

type Value interface{}

type Rows
Rows is an iterator over an executed query's results.

type Rows interface {
    // Columns returns the names of the columns. The number of
    // columns of the result is inferred from the length of the
    // slice. If a particular column name isn't known, an empty
    // string should be returned for that entry.
    Columns() []string

    // Close closes the rows iterator.
    Close() error

    // Next is called to populate the next row of data into
    // the provided slice. The provided slice will be the same
    // size as the Columns() are wide.
    //
    // Next should return io.EOF when there are no more rows.
    //
    // The dest should not be written to outside of Next. Care
    // should be taken when closing Rows not to modify
    // a buffer held in dest.
    Next(dest []Value) error
}

*/

package null

import "database/sql"

func ScanQuery(s sql.Scanner, name string, m map[string][]string) error {
	if temp, _ := m[name]; len(temp) > 0 {
		return s.Scan(temp[0])
	}
	return s.Scan(nil)
}

func ScanVars(s sql.Scanner, name string, m map[string]string) error {
	if temp, ok := m[name]; ok {
		return s.Scan(temp)
	}
	return s.Scan(nil)
}

type Scanners []sql.Scanner

func (self Scanners) Scan(in interface{}) (err error) {
	for i := range self {
		if err = self[i].Scan(in); err != nil {
			return
		}
	}
	return
}
