// Package export provides writers that serialise drift results into
// portable file formats such as CSV and newline-delimited JSON (NDJSON).
//
// Typical usage:
//
//	e := export.New(os.Stdout, export.FormatCSV)
//	if err := e.Write("my-service", diffs); err != nil {
//		log.Fatal(err)
//	}
package export
