package book

type DBBook struct {
	ID          string `db:"id"`
	Name        string `db:"name"`
	Description string `db:"description"`
	Cover       string `db:"cover"`
	Narrator    string `db:"narrator"`
	Publisher   string `db:"publisher"`
	Price       int    `db:"price"`
}

// DBGroup is an intermediate struct that matches the database schema
type DBGroup struct {
	ID       string `db:"id"`
	Name     string `db:"name"`
	ShowType int    `db:"showType"`
}

func convertDBBookToProto(dbBook DBBook) *Basic {
	return &Basic{
		Id:          dbBook.ID,
		Name:        dbBook.Name,
		Description: dbBook.Description,
		Cover:       dbBook.Cover,
		Narrator:    dbBook.Narrator,
		Publisher:   dbBook.Publisher,
		Price:       uint32(dbBook.Price),
	}
}

func convertDBGroupToProto(dbGroup DBGroup, books []*Basic) *Group {
	group := &Group{
		Id:    dbGroup.ID,
		Name:  dbGroup.Name,
		Books: books,
	}

	if dbGroup.ShowType == 2 {
		group.ShowType = Group_TwoColumn
	} else {
		group.ShowType = Group_OneColumn
	}

	return group
}
