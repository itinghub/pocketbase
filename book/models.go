package book

import (
	"github.com/pocketbase/pocketbase/models"
)

// type DBBook struct {
// 	ID          string `db:"id"`
// 	Name        string `db:"name"`
// 	Description string `db:"description"`
// 	Cover       string `db:"cover"`
// 	Narrator    string `db:"narrator"`
// 	Publisher   string `db:"publisher"`
// 	Price       int    `db:"price"`
// }

// // DBGroup is an intermediate struct that matches the database schema
// type DBGroup struct {
// 	models.Record
// 	ID       string `db:"id"`
// 	Name     string `db:"name"`
// 	ShowType int    `db:"showType"`
// 	Books   []DBBook `db:"books"`
// }

func pb_toBookBasic(record *models.Record) *Basic {
	return &Basic{
		Id:          record.GetId(),
		Name:        record.GetString("name"),
		Description: record.GetString("description"),
		Cover:       record.GetString("cover"),
		Narrator:    record.GetString("narrator"),
		Publisher:   record.GetString("publisher"),
		Price:       uint32(record.GetInt("price")),
	}
}
func pb_toBookGroup(groupRecord *models.Record, expanedBooks []*models.Record) *Group {
	group := &Group{
		Id:   groupRecord.GetId(),
		Name: groupRecord.GetString("name"),
	}

	switch groupRecord.GetInt("showType") {
	case 2:
		group.ShowType = Group_TwoColumn
	default:
		group.ShowType = Group_OneColumn
	}

	protoBooks := make([]*Basic, len(expanedBooks))
	for i, book := range expanedBooks {
		protoBooks[i] = pb_toBookBasic(book)
	}
	// group.Expand =Books = protoBooks
	return group
}

func convertToGroupResult(groupRecords []*models.Record) GroupListResp {
	protoGroups := make([]*Group, len(groupRecords))
	for i, record := range groupRecords {
		protoGroups[i] = pb_toBookGroup(record, record.ExpandedAll("books"))
	}

	return GroupListResp{
		Items: protoGroups,
	}
}
