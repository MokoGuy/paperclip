package domain

type Document struct {
	ID                   int    `json:"id"`
	Title                string `json:"title"`
	CorrespondentID      *int   `json:"-"`
	CorrespondentName    string `json:"-"`
	DocumentTypeID       *int   `json:"-"`
	DocumentTypeName     string `json:"-"`
	Created              string `json:"created"`
	Added                string `json:"added"`
	Modified             string `json:"-"`
	ArchiveSerialNumber  *int   `json:"archive_serial_number,omitempty"`
	OriginalFileName     string `json:"-"`
	PageCount            int    `json:"page_count"`
	TagIDs               []int  `json:"-"`
	TagNames             []string `json:"-"`
}
